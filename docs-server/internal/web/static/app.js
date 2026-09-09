// docs-server 文档 UI：首屏库列表、按 path 层级折叠的目录树与 Markdown 阅读视图。
// 渲染用的 marked 与 highlight.js 均为 vendored 单文件 ESM，直引同源路径。
import { marked } from "/vendor/marked.esm.js";
import hljs from "/vendor/highlight.esm.js";

const state = { library: null, activePath: null };

const libraryList = document.getElementById("library-list");
const docTree = document.getElementById("doc-tree");
const viewer = document.getElementById("viewer");

loadLibraries();

async function loadLibraries() {
  let libraries;
  try {
    const res = await fetch("/api/v1/libraries");
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    ({ libraries } = await res.json());
  } catch (err) {
    showNotice(libraryList, `库列表加载失败：${err.message}`);
    return;
  }
  libraryList.replaceChildren();
  if (libraries.length === 0) {
    showNotice(libraryList, "还没有任何库，请先推送文档。");
    return;
  }
  for (const lib of libraries) {
    const item = document.createElement("li");
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = lib.slug;
    button.dataset.slug = lib.slug;
    button.title = `${lib.documents} 篇文档`;
    button.addEventListener("click", () => selectLibrary(lib.slug));
    item.appendChild(button);
    libraryList.appendChild(item);
  }
}

async function selectLibrary(slug) {
  state.library = slug;
  state.activePath = null;
  for (const button of libraryList.querySelectorAll("button")) {
    button.classList.toggle("active", button.dataset.slug === slug);
  }

  let documents;
  try {
    const res = await fetch(`/api/v1/libraries/${encodeURIComponent(slug)}/documents`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    ({ documents } = await res.json());
  } catch (err) {
    showNotice(docTree, `文档列表加载失败：${err.message}`);
    return;
  }
  docTree.replaceChildren();
  renderTree(docTree, buildTree(documents));
  if (documents.length === 0) {
    showNotice(docTree, "该库还没有文档。");
    return;
  }
  showPlaceholder(`已选择库「${slug}」，点击目录树中的文档查看。`);
}

// buildTree 把扁平的 path 列表按「/」层级组装成嵌套节点：children 为子目录，docs 为本层文档。
function buildTree(documents) {
  const root = { children: new Map(), docs: [] };
  for (const doc of documents) {
    const parts = doc.path.split("/");
    let node = root;
    for (const dir of parts.slice(0, -1)) {
      if (!node.children.has(dir)) {
        node.children.set(dir, { children: new Map(), docs: [] });
      }
      node = node.children.get(dir);
    }
    node.docs.push({ ...doc, name: parts[parts.length - 1] });
  }
  return root;
}

// renderTree 渲染一层目录树：目录用 details/summary 原生折叠展开，文档为可点击按钮。
function renderTree(container, node) {
  const entries = [
    ...[...node.children.entries()].map(([name, child]) => ({ kind: "dir", name, child })),
    ...node.docs.map((doc) => ({ kind: "doc", name: doc.name, doc })),
  ].sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));

  const list = document.createElement("ul");
  for (const entry of entries) {
    const item = document.createElement("li");
    if (entry.kind === "dir") {
      const details = document.createElement("details");
      const summary = document.createElement("summary");
      summary.textContent = entry.name + "/";
      details.append(summary);
      renderTree(details, entry.child);
      item.appendChild(details);
    } else {
      const button = document.createElement("button");
      button.type = "button";
      button.textContent = entry.name;
      button.dataset.path = entry.doc.path;
      button.addEventListener("click", () => openDocument(entry.doc));
      item.appendChild(button);
    }
    list.appendChild(item);
  }
  container.appendChild(list);
}

// openDocument 进入查看：一次请求拿全文后整篇渲染，非流式、非分页。
async function openDocument(doc) {
  state.activePath = doc.path;
  for (const button of docTree.querySelectorAll("button[data-path]")) {
    button.classList.toggle("active", button.dataset.path === doc.path);
  }
  const url = `/api/v1/libraries/${encodeURIComponent(state.library)}/documents/${encodeDocPath(doc.path)}`;
  let body;
  try {
    const res = await fetch(url);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    ({ content: body } = await res.json());
  } catch (err) {
    showNotice(viewer, `文档加载失败：${err.message}`);
    return;
  }
  renderDocument(body);
}

// encodeDocPath 只编码路径段内的字符、保留分隔符「/」，与 docs-search 查询脚本一致。
function encodeDocPath(docPath) {
  return docPath.split("/").map(encodeURIComponent).join("/");
}

// renderDocument 把整篇 Markdown 一次性渲染进查看区：marked 裸用（内容自推送管道自控，
// 无 DOMPurify），渲染后补标题 id、代码高亮与页内 TOC。
function renderDocument(markdown) {
  const article = document.createElement("div");
  article.className = "doc-body";
  article.innerHTML = marked.parse(markdown);

  const headings = article.querySelectorAll("h1, h2, h3, h4, h5, h6");
  const used = new Set();
  for (const heading of headings) {
    heading.id = uniqueId(slug(heading.textContent), used);
  }

  const fragment = document.createDocumentFragment();
  const toc = buildToc(headings);
  if (toc) fragment.appendChild(toc);
  fragment.appendChild(article);
  viewer.replaceChildren(fragment);

  for (const block of article.querySelectorAll("pre code[class^='language-'], pre code[class*=' language-']")) {
    hljs.highlightElement(block);
  }
}

// slug 把标题文本压成锚点 id：保留字母数字（含中文），其余序列转连字符。
function slug(text) {
  return (
    text.trim().toLowerCase().replace(/[^\p{L}\p{N}]+/gu, "-").replace(/^-+|-+$/g, "") ||
    "heading"
  );
}

// uniqueId 在 base 撞号时追加数字后缀，保证同篇内 id 唯一。
function uniqueId(base, used) {
  let id = base;
  for (let n = 2; used.has(id); n++) id = `${base}-${n}`;
  used.add(id);
  return id;
}

// buildToc 用渲染后的标题生成页内 TOC：按层级缩进，点击锚点滚动定位到对应标题。
function buildToc(headings) {
  if (headings.length === 0) return null;
  const nav = document.createElement("nav");
  nav.className = "toc";
  const label = document.createElement("p");
  label.className = "toc-title";
  label.textContent = "目录";
  const list = document.createElement("ul");
  for (const heading of headings) {
    const item = document.createElement("li");
    item.className = `toc-l${heading.tagName[1]}`;
    const link = document.createElement("a");
    link.href = `#${heading.id}`;
    link.textContent = heading.textContent;
    item.appendChild(link);
    list.appendChild(item);
  }
  nav.append(label, list);
  return nav;
}

function showPlaceholder(message) {
  const p = document.createElement("p");
  p.className = "placeholder";
  p.textContent = message;
  viewer.replaceChildren(p);
}

function showNotice(container, message) {
  const p = document.createElement("p");
  p.className = "notice";
  p.textContent = message;
  container.replaceChildren(p);
}
