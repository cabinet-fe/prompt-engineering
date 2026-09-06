/* Prompt Engineering 门户 · 交互 */
(() => {
  "use strict";
  const reduced = matchMedia("(prefers-reduced-motion: reduce)").matches;

  /* ---------- 滚动进场 ---------- */
  const io = new IntersectionObserver((entries) => {
    for (const e of entries) {
      if (e.isIntersecting) {
        e.target.classList.add("in");
        io.unobserve(e.target);
      }
    }
  }, { threshold: 0.15, rootMargin: "0px 0px -6% 0px" });
  document.querySelectorAll(".reveal").forEach((el) => io.observe(el));

  /* ---------- 场景指引：推荐做法 ---------- */
  const PICKS = {
    first: {
      cmd: "setup",
      desc: "新仓库只需运行一次。通过简短沟通确认项目类型，生成技术规范、模块地图并配置辅助脚本，建立统一的开发底座。后续其他技能都以此为依据。"
    },
    vague: {
      cmd: "explore → to-spec → to-tasks → implement → archive",
      desc: "需求含糊或边界不清时使用。explore 只沟通梳理、不写代码，逐一澄清歧义并在你确认无误后才推进到后续环节，避免方向跑偏。"
    },
    rush: {
      cmd: "rush 审批单支持批量转办",
      desc: "全自动任务编排器：自动完成需求规格制定、任务拆解、并行编码与代码审查。做完一个阶段立即审查，不通过自动返工修复，适合大中型功能。"
    },
    steps: {
      cmd: "to-spec → to-tasks → implement → review → archive",
      desc: "需求明确时可跳过前期探讨，直接进入技术规格拆解与分步实现。每一步你都可以手动检查验收标准与实际代码，节奏完全由你掌控。"
    },
    small: {
      cmd: "implement 把 UserCard 的头像换成懒加载",
      desc: "改个 Bug、加个字段等日常小改动，无需复杂的规格拆解，直接描述需求即可编码。改完自动调用代码审查子代理，审查通过后自动生成本地提交。"
    },
    review: {
      cmd: "review  /  review main  /  review <需求名> P2",
      desc: "只做代码审查，不修改代码。在独立子代理中比对功能契约、代码规范与逻辑正确性，通过后协助提交，未通过则列出具体阻断项。"
    },
    sync: {
      cmd: "sync-docs src/modules/approval",
      desc: "手动修改代码导致现有文档与代码不一致时调用。精准修正受影响的架构说明与规范，未受影响的内容绝不乱动，且不新建多余文件。"
    },
    accept: {
      cmd: "acceptance",
      desc: "为项目配置自动化验收标准。先自动扫描现有的测试与构建命令，针对性生成端到端测试或接口验收规范，作为后续开发与审查的硬性指标。"
    }
  };
  const pickNote = document.getElementById("pick-note");
  const pickCmd = document.getElementById("pick-cmd");
  const pickDesc = document.getElementById("pick-desc");
  const pickChips = [...document.querySelectorAll(".pick-chip")];
  const selectPick = (key) => {
    const p = PICKS[key];
    if (!p) return;
    pickChips.forEach((c) => c.classList.toggle("active", c.dataset.pick === key));
    pickCmd.textContent = p.cmd;
    pickDesc.textContent = p.desc;
    pickNote.classList.remove("pop");
    void pickNote.offsetWidth; // 重触发动画
    pickNote.classList.add("pop");
  };
  pickChips.forEach((c) => c.addEventListener("click", () => selectPick(c.dataset.pick)));
  selectPick("first");

  /* ---------- 流程：节点详情 ---------- */
  const FLOWS = {
    setup: {
      title: "🧱 setup · 初始化底座",
      cmd: "setup",
      desc: "新仓库运行一次。通过简要沟通识别项目类型，生成架构规范、模块地图与辅助脚本。其他所有工程技能都会基于这套底座运行。"
    },
    explore: {
      title: "🧭 explore · 需求探索",
      cmd: "explore 审批单支持批量转办",
      desc: "只沟通梳理、不改写代码。逐条分析并澄清需求中的不明确之处，向你提问确认。直到所有关键疑惑消除并经你确认，才进入下一阶段。"
    },
    "to-spec": {
      title: "📜 to-spec · 制定规格",
      cmd: "to-spec approval-batch-transfer",
      desc: "把明确需求转化为技术设计规格。制定客观可判定的验收标准，并精确划定本次改动涉及的文件范围。"
    },
    "to-tasks": {
      title: "✂️ to-tasks · 拆解阶段任务",
      cmd: "to-tasks approval-batch-transfer",
      desc: "将大功能细化拆解为多个独立阶段任务（P1、P2 等），理清任务依赖关系。无前后依赖的任务可直接并行开发。"
    },
    implement: {
      title: "🔨 implement · 编码实现",
      cmd: "implement approval-batch-transfer P1",
      desc: "按阶段任务编写代码（或针对小需求直接修改）。严格遵循既有代码风格，完成必要测试并同步修正受影响的文档，随后自动发起审查。"
    },
    review: {
      title: "🔍 review · 独立审查",
      cmd: "review approval-batch-transfer P2",
      desc: "在独立子代理中审查代码，主会话只接收审查结果。比对功能实现与代码规范，通过后自动在本地提交，不通过则指出具体问题并返工。"
    },
    archive: {
      title: "🗃️ archive · 交付与归档",
      cmd: "archive approval-batch-transfer",
      desc: "所有阶段均完成并审查通过后，核对项目文档状态，清理开发过程中的临时记录（cooking 目录），完成交付。"
    }
  };
  const flowNote = document.getElementById("flow-note");
  const flowTitle = document.getElementById("flow-title");
  const flowCmd = document.getElementById("flow-cmd");
  const flowDesc = document.getElementById("flow-desc");
  const flowNodes = [...document.querySelectorAll(".flow-node")];
  const selectFlow = (key) => {
    const f = FLOWS[key];
    if (!f) return;
    flowNodes.forEach((n) => n.classList.toggle("active", n.dataset.flow === key));
    flowTitle.textContent = f.title;
    flowCmd.textContent = f.cmd;
    flowDesc.textContent = f.desc;
    flowNote.classList.remove("pop");
    void flowNote.offsetWidth;
    flowNote.classList.add("pop");
  };
  flowNodes.forEach((n) => n.addEventListener("click", () => selectFlow(n.dataset.flow)));
  selectFlow("setup");

  /* ---------- 技能筛选 ---------- */
  const filterChips = [...document.querySelectorAll(".filter-chip")];
  const skillCards = [...document.querySelectorAll(".skill-card")];
  filterChips.forEach((chip) => {
    chip.addEventListener("click", () => {
      filterChips.forEach((c) => c.classList.toggle("active", c === chip));
      const f = chip.dataset.filter;
      skillCards.forEach((card) => {
        const show = f === "all" || card.dataset.cat === f;
        card.classList.toggle("hide", !show);
        if (show) card.classList.add("in");
      });
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
  const navLinks = [...document.querySelectorAll(".nav-links a")];
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

  /* ---------- 涂鸦视差 ---------- */
  if (!reduced) {
    const px = [...document.querySelectorAll("[data-parallax]")];
    let ticking = false;
    const applyParallax = () => {
      ticking = false;
      const y = scrollY;
      for (const el of px) {
        el.style.translate = `0 ${(y * parseFloat(el.dataset.parallax)).toFixed(1)}px`;
      }
    };
    addEventListener("scroll", () => {
      if (!ticking) {
        ticking = true;
        requestAnimationFrame(applyParallax);
      }
    }, { passive: true });
  }
})();
