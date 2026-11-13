/* 提示词工程门户 · 交互 */
(() => {
  "use strict";

  /* ---------- 跑马灯：克隆轨道内容实现无缝循环 ---------- */
  const tickerInner = document.querySelector(".ticker-inner");
  if (tickerInner) tickerInner.appendChild(tickerInner.firstElementChild.cloneNode(true));

  /* ---------- 场景选读：推荐做法 ---------- */
  const PICKS = {
    first: {
      cmd: "setup",
      desc: "新仓库只需运行一次。通过简短沟通确认项目类型，生成架构与规范文档、模块地图及辅助脚本，并在根 AGENTS.md 写入索引与「收尾」约定——后续无论走工作流还是直接对话，交付前都会自动更新文档并触发独立审查。"
    },
    vague: {
      cmd: "explore → blueprint → rush",
      desc: "需求含糊或边界不清时使用。explore 只沟通梳理、不写代码，逐一澄清歧义并在你确认无误后才推进到 blueprint 拆解和 rush 实现，避免方向跑偏。"
    },
    rush: {
      cmd: "rush 审批单支持批量转办",
      desc: "全自动编排器：可在流程起点或 explore / blueprint / 阶段实现与 review 之后的任意位置介入，读取当前进度断点续跑。无依赖阶段并行开发，每阶段完成即审查，不通过自动返工，全部通过后自动归档，适合大中型功能。"
    },
    steps: {
      cmd: "blueprint → rush <feature>",
      desc: "需求明确但想先确认技术规格、验收标准与阶段切分时，先跑 blueprint 生成 spec 与阶段任务，确认无误后再交给 rush 自动并行实现、审查与归档。"
    },
    small: {
      cmd: "直接对话（如：把 UserCard 的头像换成懒加载）",
      desc: "日常开发、改 Bug 或小功能无需走规格拆解，直接对话编码即可。setup 生成的根 AGENTS.md「收尾」约定会在交付前自动更新受影响文档，非琐碎改动还会自动派独立子代理做 git 审查。"
    },
    review: {
      cmd: "review  /  review main  /  review <需求名> P2",
      desc: "只做代码审查，不修改代码。在独立子代理中比对功能契约、代码规范与逻辑正确性，通过后协助提交，未通过则列出具体阻断项。"
    },
    sync: {
      cmd: "sync-docs src/modules/approval",
      desc: "日常改动会按根 AGENTS.md「收尾」当轮更新文档；当发现已有文档漏改、需要批量对齐，或需新增自定义文档 / 子包代码地图时调用 sync-docs。"
    },
    accept: {
      cmd: "acceptance",
      desc: "为项目配置自动化验收标准。先自动扫描现有的测试与构建命令，针对性生成端到端测试或接口验收规范，作为后续开发与审查的硬性指标。"
    }
  };
  const pickChips = [...document.querySelectorAll(".pick-row")];
  const selectPick = (key) => {
    const p = PICKS[key];
    if (!p) return;
    pickChips.forEach((c) => c.classList.toggle("active", c.dataset.pick === key));
    document.getElementById("pick-cmd").textContent = p.cmd;
    document.getElementById("pick-desc").textContent = p.desc;
  };
  pickChips.forEach((c) => c.addEventListener("click", () => selectPick(c.dataset.pick)));

  /* ---------- 流程：节点详情 ---------- */
  const FLOWS = {
    setup: {
      title: "STEP 1 · setup · 初始化底座",
      cmd: "setup",
      desc: "新仓库只需运行一次。通过简短沟通确认项目类型，生成架构与规范文档、模块地图及辅助脚本，并在根 AGENTS.md 写入文档索引与「收尾」约定。日常对话与工程技能都基于这套底座运行。"
    },
    explore: {
      title: "STEP 2 · explore · 需求探索",
      cmd: "explore 审批单支持批量转办",
      desc: "只沟通梳理、不改写代码。逐条分析并澄清需求中的不明确之处，向你提问确认。直到所有关键疑惑消除并经你确认，才进入下一阶段。"
    },
    blueprint: {
      title: "STEP 3 · blueprint · 规格与任务切分",
      cmd: "blueprint approval-batch-transfer",
      desc: "一次产出技术设计规格（spec.md）与可并行的阶段任务（tasks/Pn.md），制定可判定的验收标准，并标明是否涉及架构级变更。"
    },
    rush: {
      title: "STEP 4 · rush · 并行实现与编排",
      cmd: "rush approval-batch-transfer",
      desc: "按阶段依赖自动调度实现子代理并行编码，遵守代码规范并当轮同步文档，每完成一个阶段即刻送交独立审查，全部通过后自动归档。"
    },
    review: {
      title: "STEP 5 · review · 独立审查与归档",
      cmd: "review approval-batch-transfer P2",
      desc: "在独立子代理中审查代码与文档同步情况，主会话只接收审查结果。通过后自动生成本地提交，全部阶段通过时自动清理 cooking 临时目录。"
    }
  };
  const flowNodes = [...document.querySelectorAll(".flow-node")];
  const selectFlow = (key) => {
    const f = FLOWS[key];
    if (!f) return;
    flowNodes.forEach((n) => n.classList.toggle("active", n.dataset.flow === key));
    document.getElementById("flow-title").textContent = f.title;
    document.getElementById("flow-cmd").textContent = f.cmd;
    document.getElementById("flow-desc").textContent = f.desc;
  };
  flowNodes.forEach((n) => n.addEventListener("click", () => selectFlow(n.dataset.flow)));

  /* ---------- 技能筛选 ---------- */
  const filterChips = [...document.querySelectorAll(".filter-chip")];
  const skillCards = [...document.querySelectorAll(".skill-card")];
  filterChips.forEach((chip) => {
    chip.addEventListener("click", () => {
      filterChips.forEach((c) => c.classList.toggle("active", c === chip));
      const f = chip.dataset.filter;
      skillCards.forEach((card) => card.classList.toggle("hide", !(f === "all" || card.dataset.cat === f)));
    });
  });

  /* ---------- 技能详情弹层 ---------- */
  const modal = document.getElementById("skill-modal");
  const modalSheet = modal.querySelector(".modal-sheet");
  const modalCat = document.getElementById("modal-cat");
  const modalTitle = document.getElementById("modal-title");
  const modalBody = document.getElementById("modal-body");
  let lastFocus = null;

  const openModal = (card) => {
    lastFocus = card;
    modalTitle.textContent = card.querySelector(".skill-name").textContent;
    modalCat.textContent = card.querySelector(".skill-cat").textContent;
    modalSheet.style.setProperty("--cat", getComputedStyle(card).getPropertyValue("--cat"));
    modalBody.innerHTML = card.querySelector(".skill-detail").innerHTML;
    modal.hidden = false;
    document.body.style.overflow = "hidden";
    modal.querySelector(".modal-close").focus();
  };
  const closeModal = () => {
    modal.hidden = true;
    document.body.style.overflow = "";
    if (lastFocus) lastFocus.focus();
  };
  skillCards.forEach((card) => card.addEventListener("click", () => openModal(card)));
  modal.addEventListener("click", (e) => { if (e.target.closest("[data-close]")) closeModal(); });
  addEventListener("keydown", (e) => { if (e.key === "Escape" && !modal.hidden) closeModal(); });

  /* ---------- 复制 + Toast ---------- */
  const toast = document.getElementById("toast");
  let toastTimer = 0;
  const showToast = (msg) => {
    toast.textContent = msg;
    toast.classList.add("show");
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => toast.classList.remove("show"), 1600);
  };
  const copyText = async (text) => {
    try {
      await navigator.clipboard.writeText(text);
      showToast("已复制到剪贴板！");
    } catch {
      const ta = document.createElement("textarea");
      ta.value = text;
      ta.style.position = "fixed";
      ta.style.opacity = "0";
      document.body.appendChild(ta);
      ta.select();
      try { document.execCommand("copy"); showToast("已复制到剪贴板！"); }
      catch { showToast("复制失败，请手动复制"); }
      ta.remove();
    }
  };
  document.querySelectorAll(".copy-btn").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      e.stopPropagation();
      copyText(btn.dataset.copy);
    });
  });

  /* ---------- 导航当前区块 ---------- */
  const navLinks = [...document.querySelectorAll(".navbar-links a")];
  const sections = navLinks
    .map((a) => document.querySelector(a.getAttribute("href")))
    .filter(Boolean);
  const navIO = new IntersectionObserver((entries) => {
    for (const e of entries) {
      if (e.isIntersecting) {
        navLinks.forEach((a) => a.classList.toggle("active", a.getAttribute("href") === "#" + e.target.id));
      }
    }
  }, { rootMargin: "-38% 0px -55% 0px" });
  sections.forEach((s) => navIO.observe(s));
})();
