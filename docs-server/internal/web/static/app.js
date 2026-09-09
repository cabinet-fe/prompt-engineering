// docs-server 文档 UI：首屏库列表、按 path 层级折叠的目录树与文档查看入口。
// marked 与 highlight.js 是 P4 渲染视图的 vendored 依赖，这里先完成同源模块直引。
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
  for (const slug of libraries) {
    const item = document.createElement("li");
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = slug;
    button.addEventListener("click", () => selectLibrary(slug));
    item.appendChild(button);
    libraryList.appendChild(item);
  }
}

async function selectLibrary(slug) {
  state.library = slug;
  state.activePath = null;
  for (const button of libraryList.querySelectorAll("button")) {
    button.classList.toggle("active", button.textContent === slug);
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

// openDocument 进入查看区域：本阶段先占位，Markdown 渲染留给下一阶段。
function openDocument(doc) {
  state.activePath = doc.path;
  for (const button of docTree.querySelectorAll("button[data-path]")) {
    button.classList.toggle("active", button.dataset.path === doc.path);
  }
  showPlaceholder(`「${doc.title}」（${doc.path}）已选中，文档渲染即将提供。`);
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
