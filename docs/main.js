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

  /* ---------- 先挑路：处方 ---------- */
  const PICKS = {
    first: {
      cmd: "setup",
      desc: "一次性。访谈后分类，写 PROJECT.md、铺 .agents/docs/ 和脚本、覆写根 AGENTS.md。跑完 node .agents/scripts/precheck.mjs 应为 PASS。其它技能发现没 setup，会停下来让你先跑，不代跑。"
    },
    vague: {
      cmd: "explore → to-spec → to-tasks → implement → archive",
      desc: "从 explore 起步，分步走。它只问不写代码：把「不问就可能做错」的歧义逐条记进 goal.md，每轮最多问 5 题，未决清零、范围经你确认才往下走。"
    },
    rush: {
      cmd: "rush 审批单想支持批量转办",
      desc: "编排器不是另一套流程——产物、模板、阶段并行规则跟单独调用完全一样。可并行的阶段同时开 implement，一返回就立刻单独 review，不通过自动返工，单阶段最多 3 轮。"
    },
    steps: {
      cmd: "to-spec → to-tasks → implement → review → archive",
      desc: "需求本来就明确的话跳过 explore，从 to-spec 起步，每一步你自己把关：spec 验收标准必须可判定，implement 一次只做一个未阻塞阶段。"
    },
    small: {
      cmd: "implement 把 UserCard 的头像换成懒加载",
      desc: "参数既不是 cooking 标识、也不是单独的 P<n> 时走直写：不读 spec/tasks，不碰 cooking，小 diff。说错的已有文档当场改。改完派 review 子代理，通过后 git-commit auto。"
    },
    review: {
      cmd: "review  /  review main  /  review <feature> P2",
      desc: "只评不改，必须在子代理里评：主会话只派发、只听结论，不读 diff。评审轴：Spec、Standards（含已有文档有没有被说错）、正确性（仅 git 路径）。"
    },
    sync: {
      cmd: "sync-docs src/modules/approval",
      desc: "未走 implement 时点名。只改已经被代码说错的已有文档；没有被说错就结束。禁止新建文件。"
    },
    accept: {
      cmd: "acceptance",
      desc: "先检索仓库已有的测试 / e2e / HTTP / 构建命令，再结合项目类别推荐手段——检索完之前不会点名具体工具。之后 to-tasks 按它追加完成标准，review 按它评。代价是 token 和工时明显增加。"
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
      title: "🧱 setup",
      cmd: "setup",
      desc: "每个目标仓库先跑一次。访谈后分类，写 PROJECT.md（代码类还有架构/规范/地图/坏味道）、复制脚本、覆写根 AGENTS.md。其它技能发现没 setup 会停下让你先跑，不代跑。"
    },
    explore: {
      title: "🧭 explore",
      cmd: "explore 审批单想支持批量转办",
      desc: "只问不写代码。把「不问就可能做错」的歧义逐条写进 goal.md 的「未决问题」，每轮挑最阻塞的问、最多 5 题；查得到的事实派子代理，不问你。未决清零、范围经你确认，才写「确认：已确认」。"
    },
    "to-spec": {
      title: "📜 to-spec",
      cmd: "to-spec approval-batch-transfer",
      desc: "把明确需求写成 spec.md。验收标准必须可判定——「体验好」这种不收；「影响文件」只写新增/删除/修改的仓库相对路径，写完要能被 spec-files.mjs parse 通过。"
    },
    "to-tasks": {
      title: "✂️ to-tasks",
      cmd: "to-tasks approval-batch-transfer",
      desc: "把 spec.md 拆成 tasks/P1.md、P2.md…。Pn 是阶段 id，不是必须串行的序号——前置为「无」的阶段一上来就能并行。"
    },
    implement: {
      title: "🔨 implement",
      cmd: "implement approval-batch-transfer P1",
      desc: "一次只做一个未阻塞阶段，做完在同一对话派 review 子代理——不自己评、不提交。说错的已有文档当场改。参数既不是标识也不是 P<n> 时走直写：不碰 cooking，小 diff，对照 SMELLS.md 收掉本次引入的坏味道。"
    },
    review: {
      title: "🔍 review",
      cmd: "review approval-batch-transfer P2",
      desc: "只评不改，必须在子代理里评。通过后由派发方 git-commit auto（本地提交，不 push）；不通过就把阻塞项列出来返工，再评一轮，该阶段下游不许开工。"
    },
    archive: {
      title: "🗃️ archive",
      cmd: "archive approval-batch-transfer",
      desc: "要求每个阶段都「实现：完成 + 评审：通过」。确认已有文档已对齐后删掉整个 cooking 目录，不蒸馏 spec。"
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
      showToast("复制好了！✎");
    } catch {
      const ta = document.createElement("textarea");
      ta.value = text;
      ta.style.position = "fixed";
      ta.style.opacity = "0";
      document.body.appendChild(ta);
      ta.select();
      try { document.execCommand("copy"); showToast("复制好了！✎"); }
      catch { showToast("复制失败，手动来吧"); }
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
