// 调用记录虚拟表格视图：仅渲染可见区间行，查询结果单行省略、鼠标悬停浮层完整展示。

const CALL_ACTION_LABELS = { search: "检索", get: "取文档", toc: "目录" };
const ROW_HEIGHT = 40;
const OVERSCAN = 6;
const COLUMNS = ["时间", "库", "动作", "查询词", "查询结果", "IP"];

// renderCallsView 在 viewer 内挂载虚拟表格视图，返回切换视图时必须调用的清理函数。
export function renderCallsView(viewer, calls, total, onRefresh) {
  const wrap = document.createElement("div");
  wrap.className = "calls";

  const headerBar = document.createElement("div");
  headerBar.className = "calls-header";

  const titleGroup = document.createElement("div");
  const heading = document.createElement("h1");
  heading.textContent = "调用记录";
  const meta = document.createElement("p");
  meta.className = "calls-meta";
  meta.textContent =
    total === 0
      ? "最多保留最近 200 条调用记录。"
      : `共 ${total} 条（最多保留最近 200 条）`;
  titleGroup.append(heading, meta);

  const refreshBtn = document.createElement("button");
  refreshBtn.type = "button";
  refreshBtn.className = "calls-refresh";
  refreshBtn.textContent = "刷新";
  refreshBtn.addEventListener("click", () => onRefresh());

  headerBar.append(titleGroup, refreshBtn);
  wrap.appendChild(headerBar);

  if (calls.length === 0) {
    const p = document.createElement("p");
    p.className = "placeholder";
    p.textContent = "还没有任何检索调用。";
    wrap.appendChild(p);
    viewer.replaceChildren(wrap);
    return () => {};
  }

  const viewport = document.createElement("div");
  viewport.className = "calls-vtable-viewport";

  const inner = document.createElement("div");
  inner.className = "calls-vtable-inner";
  inner.setAttribute("role", "table");
  inner.setAttribute("aria-label", "调用记录");

  const tableHeader = document.createElement("div");
  tableHeader.className = "calls-vtable-header";
  tableHeader.setAttribute("role", "row");
  for (const label of COLUMNS) {
    const th = document.createElement("div");
    th.className = "calls-vtable-th";
    th.setAttribute("role", "columnheader");
    th.textContent = label;
    tableHeader.appendChild(th);
  }

  const body = document.createElement("div");
  body.className = "calls-vtable-body";
  body.style.height = `${calls.length * ROW_HEIGHT}px`;

  const slice = document.createElement("div");
  slice.className = "calls-vtable-slice";
  slice.setAttribute("role", "rowgroup");

  body.appendChild(slice);
  inner.append(tableHeader, body);
  viewport.appendChild(inner);
  wrap.appendChild(viewport);
  viewer.replaceChildren(wrap);

  const popover = document.createElement("div");
  popover.className = "calls-popover";
  popover.hidden = true;
  document.body.appendChild(popover);

  let disposed = false;
  let rafId = 0;
  let lastStart = -1;
  let lastEnd = -1;
  let hideTimer = 0;
  let activeCell = null;

  function hidePopoverNow() {
    if (hideTimer) {
      clearTimeout(hideTimer);
      hideTimer = 0;
    }
    popover.hidden = true;
    activeCell = null;
  }

  function scheduleHidePopover() {
    if (hideTimer) clearTimeout(hideTimer);
    hideTimer = setTimeout(() => {
      hideTimer = 0;
      popover.hidden = true;
      activeCell = null;
    }, 120);
  }

  function cancelHidePopover() {
    if (hideTimer) {
      clearTimeout(hideTimer);
      hideTimer = 0;
    }
  }

  function renderWindow(force = false) {
    if (disposed) return;
    const headerHeight = tableHeader.offsetHeight || 38;
    const viewHeight = Math.max(ROW_HEIGHT, viewport.clientHeight - headerHeight);
    const scrollTop = Math.max(0, viewport.scrollTop);
    const firstVisible = Math.floor(scrollTop / ROW_HEIGHT);
    const visibleCount = Math.ceil(viewHeight / ROW_HEIGHT);
    const start = Math.max(0, firstVisible - OVERSCAN);
    const end = Math.min(calls.length, firstVisible + visibleCount + OVERSCAN);

    if (!force && start === lastStart && end === lastEnd) return;
    lastStart = start;
    lastEnd = end;

    slice.style.transform = `translateY(${start * ROW_HEIGHT}px)`;
    const frag = document.createDocumentFragment();
    for (let i = start; i < end; i++) {
      frag.appendChild(buildRow(calls[i]));
    }
    slice.replaceChildren(frag);
  }

  function scheduleRender() {
    if (disposed || rafId) return;
    rafId = requestAnimationFrame(() => {
      rafId = 0;
      renderWindow(false);
    });
  }

  function onScroll() {
    hidePopoverNow();
    scheduleRender();
  }

  function onSliceMouseOver(e) {
    const target = e.target instanceof Element ? e.target.closest("[data-full-text]") : null;
    if (!target || !slice.contains(target)) return;
    cancelHidePopover();
    if (activeCell === target && !popover.hidden) return;
    activeCell = target;
    popover.textContent = target.dataset.fullText;
    popover.hidden = false;
    positionPopover(target, popover);
  }

  function onSliceMouseOut(e) {
    const target = e.target instanceof Element ? e.target.closest("[data-full-text]") : null;
    if (!target) return;
    const related = e.relatedTarget instanceof Node ? e.relatedTarget : null;
    if (related && (target.contains(related) || popover.contains(related))) return;
    scheduleHidePopover();
  }

  function onPopoverMouseEnter() {
    cancelHidePopover();
  }

  function onPopoverMouseLeave(e) {
    const related = e.relatedTarget instanceof Node ? e.relatedTarget : null;
    if (related && activeCell && activeCell.contains(related)) return;
    scheduleHidePopover();
  }

  viewport.addEventListener("scroll", onScroll, { passive: true });
  slice.addEventListener("mouseover", onSliceMouseOver);
  slice.addEventListener("mouseout", onSliceMouseOut);
  popover.addEventListener("mouseenter", onPopoverMouseEnter);
  popover.addEventListener("mouseleave", onPopoverMouseLeave);

  let ro = null;
  if (typeof ResizeObserver === "function") {
    ro = new ResizeObserver(() => scheduleRender());
    ro.observe(viewport);
  }

  renderWindow(true);

  return () => {
    disposed = true;
    if (rafId) {
      cancelAnimationFrame(rafId);
      rafId = 0;
    }
    hidePopoverNow();
    ro?.disconnect();
    viewport.removeEventListener("scroll", onScroll);
    slice.removeEventListener("mouseover", onSliceMouseOver);
    slice.removeEventListener("mouseout", onSliceMouseOut);
    popover.removeEventListener("mouseenter", onPopoverMouseEnter);
    popover.removeEventListener("mouseleave", onPopoverMouseLeave);
    popover.remove();
  };
}

function buildRow(call) {
  const row = document.createElement("div");
  row.className = "calls-vtable-row";
  row.setAttribute("role", "row");
  row.style.height = `${ROW_HEIGHT}px`;

  row.append(
    createCell(formatTime(call.created_at)),
    createCell(call.library || "（跨库）"),
    createCell(CALL_ACTION_LABELS[call.action] || call.action),
    createCell(call.query || "—", call.query || ""),
    createCell(formatResultPreview(call.result), call.result || "", "calls-result-cell"),
    createCell(call.ip || "—"),
  );
  return row;
}

function createCell(text, fullText = "", extraClass = "") {
  const cell = document.createElement("div");
  cell.className = "calls-vtable-td";
  if (extraClass) cell.classList.add(extraClass);
  cell.setAttribute("role", "cell");
  cell.textContent = text;
  if (fullText) {
    cell.classList.add("has-popover");
    cell.dataset.fullText = fullText;
  }
  return cell;
}

function formatResultPreview(result) {
  if (!result) return "—";
  return result.replace(/\s*\r?\n\s*/g, "  ");
}

function formatTime(iso) {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return iso;
  return date.toLocaleString("zh-CN", { hour12: false });
}

function positionPopover(cell, popover) {
  const rect = cell.getBoundingClientRect();
  const pad = 12;
  const gap = 6;
  const maxW = Math.min(560, Math.max(240, window.innerWidth - pad * 2));
  popover.style.maxWidth = `${maxW}px`;
  popover.style.left = "0px";
  popover.style.top = "0px";

  const pw = popover.offsetWidth;
  const ph = popover.offsetHeight;
  let left = Math.min(rect.left, window.innerWidth - pw - pad);
  left = Math.max(pad, left);

  const spaceBelow = window.innerHeight - rect.bottom - gap - pad;
  const spaceAbove = rect.top - gap - pad;
  let top;
  if (ph <= spaceBelow || spaceBelow >= spaceAbove) {
    top = Math.min(rect.bottom + gap, window.innerHeight - ph - pad);
  } else {
    top = Math.max(pad, rect.top - gap - ph);
  }
  top = Math.max(pad, top);

  popover.style.left = `${Math.round(left)}px`;
  popover.style.top = `${Math.round(top)}px`;
}
