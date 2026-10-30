// docs-server 文档 UI：首屏库列表、按 path 层级折叠的目录树与 Markdown 阅读视图。
// 渲染用的 marked 与 highlight.js 均为 vendored 单文件 ESM，直引同源路径。
import { marked } from "/vendor/marked.esm.js";
import hljs from "/vendor/highlight.esm.js";

const state = { library: null, activePath: null };

const libraryList = document.getElementById("library-list");
const docTree = document.getElementById("doc-tree");
const viewer = document.getElementById("viewer");
const callsButton = document.getElementById("calls-button");

loadLibraries();
callsButton.addEventListener("click", openCalls);

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
  callsButton.classList.remove("active");
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
  callsButton.classList.remove("active");
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

// 调用记录动作的展示文案；每页条数与服务端 /api/v1/calls 缺省 limit 一致。
const callActionLabels = { search: "检索", get: "取文档", toc: "目录" };
const callsPageSize = 50;
const callsState = { page: 1 };

// openCalls 在查看区分页展示调用记录（服务端仅保留 7 天）：时间 / 库 / 动作 / 查询词 / IP。
async function openCalls(page = callsState.page) {
  callsButton.classList.add("active");
  for (const button of docTree.querySelectorAll("button[data-path]")) {
    button.classList.remove("active");
  }
  let body;
  try {
    const params = new URLSearchParams({
      limit: callsPageSize,
      offset: (page - 1) * callsPageSize,
    });
    const res = await fetch(`/api/v1/calls?${params}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    body = await res.json();
  } catch (err) {
    showNotice(viewer, `调用记录加载失败：${err.message}`);
    return;
  }
  // 记录被滚动清理后当前页可能已空，退回第一页。
  if (body.calls.length === 0 && page > 1) {
    openCalls(1);
    return;
  }
  callsState.page = page;
  renderCalls(body.calls, body.total);
}

// renderCalls 渲染当前页调用记录表与分页导航：时间转本地时区，库为空表示跨库检索。
function renderCalls(calls, total) {
  const wrap = document.createElement("div");
  wrap.className = "calls";

  const heading = document.createElement("h1");
  heading.textContent = "调用记录";
  const start = (callsState.page - 1) * callsPageSize;
  const meta = document.createElement("p");
  meta.className = "calls-meta";
  meta.textContent =
    total === 0
      ? "服务端仅保留最近 7 天的调用记录。"
      : `第 ${start + 1}–${start + calls.length} 条 / 共 ${total} 条（服务端仅保留 7 天）`;
  wrap.append(heading, meta);

  if (calls.length === 0) {
    const p = document.createElement("p");
    p.className = "placeholder";
    p.textContent = "近 7 天还没有任何检索调用。";
    wrap.appendChild(p);
    viewer.replaceChildren(wrap);
    return;
  }

  const table = document.createElement("table");
  table.className = "calls-table";
  const thead = document.createElement("thead");
  const headRow = document.createElement("tr");
  for (const label of ["时间", "库", "动作", "查询词", "IP"]) {
    const th = document.createElement("th");
    th.textContent = label;
    headRow.appendChild(th);
  }
  thead.appendChild(headRow);
  const tbody = document.createElement("tbody");
  for (const call of calls) {
    const row = document.createElement("tr");
    appendCell(row, formatTime(call.created_at));
    appendCell(row, call.library || "（跨库）");
    appendCell(row, callActionLabels[call.action] || call.action);
    appendCell(row, call.query || "—", "wrap");
    appendCell(row, call.ip);
    tbody.appendChild(row);
  }
  table.append(thead, tbody);
  wrap.appendChild(table);

  // 分页导航：上一页 / 页码 / 下一页。
  const pages = Math.max(1, Math.ceil(total / callsPageSize));
  const nav = document.createElement("div");
  nav.className = "calls-nav";
  nav.append(
    callsNavButton("上一页", callsState.page > 1, () => openCalls(callsState.page - 1)),
    callsNavLabel(`第 ${callsState.page} / ${pages} 页`),
    callsNavButton("下一页", callsState.page < pages, () => openCalls(callsState.page + 1)),
  );
  wrap.appendChild(nav);
  viewer.replaceChildren(wrap);
  viewer.scrollTop = 0;
}

// callsNavButton 构造分页按钮，不可用时置 disabled。
function callsNavButton(label, enabled, onClick) {
  const button = document.createElement("button");
  button.type = "button";
  button.textContent = label;
  button.disabled = !enabled;
  if (enabled) {
    button.addEventListener("click", onClick);
  }
  return button;
}

// callsNavLabel 构造分页页码文本。
function callsNavLabel(text) {
  const span = document.createElement("span");
  span.className = "calls-nav-label";
  span.textContent = text;
  return span;
}

// appendCell 往表格行追加一个 td，className 可选。
function appendCell(row, text, className) {
  const td = document.createElement("td");
  td.textContent = text;
  if (className) td.className = className;
  row.appendChild(td);
}

// formatTime 把服务端 UTC 时间（RFC3339）转为本地时区显示。
function formatTime(iso) {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return iso;
  return date.toLocaleString("zh-CN", { hour12: false });
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
