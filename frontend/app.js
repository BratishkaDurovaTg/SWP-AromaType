const ORDER_CONTACT_URL = "https://t.me/aroma_type_test_bot";
const PROFILE_HERO_IMAGES = {
  "интеллект и фокус": "./assets/profile-focus.png",
  "драйв и экстраверсия": "./assets/profile-drive.png",
  "эстетика и гедонизм": "./assets/profile-aesthetic.png",
  "власть и доминанта": "./assets/profile-dominance.png",
  "сбалансированный профиль": "./assets/profile-balanced.png",
  "сбалансированного профиля": "./assets/profile-balanced.png",
};

function normalizeProfileName(profileName) {
  return String(profileName || "")
    .trim()
    .toLowerCase()
    .replace(/ё/g, "е")
    .replace(/\s+/g, " ");
}

function getProfileHeroImage(profileName) {
  const normalized = normalizeProfileName(profileName);

  if (PROFILE_HERO_IMAGES[normalized]) {
    return PROFILE_HERO_IMAGES[normalized];
  }

  if (normalized.includes("интеллект") || normalized.includes("фокус")) {
    return "./assets/profile-focus.png";
  }

  if (normalized.includes("драйв") || normalized.includes("экстраверс")) {
    return "./assets/profile-drive.png";
  }

  if (normalized.includes("эстет") || normalized.includes("гедонизм")) {
    return "./assets/profile-aesthetic.png";
  }

  if (normalized.includes("власть") || normalized.includes("доминант")) {
    return "./assets/profile-dominance.png";
  }

  if (normalized.includes("сбаланс")) {
    return "./assets/profile-balanced.png";
  }

  return "./assets/profile-balanced.png";
}

function injectProfileHeroStyles() {
  if (document.getElementById("profile-hero-styles")) return;

  const style = document.createElement("style");
  style.id = "profile-hero-styles";
  style.textContent = `
    .profile-hero-image {
      position: relative !important;
      height: 360px !important;
      margin: 0 -16px !important;
      overflow: hidden !important;
      background: var(--cream) !important;
    }

    .profile-hero-photo {
      position: absolute;
      inset: 0;
      width: 100%;
      height: 100%;
      display: block;
      object-fit: cover;
      object-position: center top;
      z-index: 0;
    }

    .profile-hero-image::before {
      content: none !important;
      display: none !important;
    }

    .profile-hero-image::after {
      content: "";
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      height: 150px;
      background: linear-gradient(
        180deg,
        rgba(236, 229, 211, 0) 0%,
        var(--cream) 100%
      );
      pointer-events: none;
      z-index: 1;
    }

    .profile-eyebrow {
      position: absolute !important;
      left: 18px !important;
      bottom: 20px !important;
      z-index: 2 !important;
      min-width: 160px !important;
      height: 44px !important;
      border: 1.5px solid rgba(255, 255, 255, 0.75) !important;
      border-radius: 500px !important;
      display: inline-flex !important;
      align-items: center !important;
      justify-content: center !important;
      padding: 0 18px !important;
      background: rgba(249, 244, 241, 0.42) !important;
      color: var(--ink) !important;
      font-family: var(--display) !important;
      font-size: 12px !important;
      backdrop-filter: blur(10px);
    }
  `;
  document.head.appendChild(style);
}

const state = {
  apiBase: resolveApiBase(),
  questions: [],
  selectedAnswers: new Map(),
  currentQuestionIndex: 0,
  recommendations: null,
  activeFilter: "all",
  selectedProduct: null,
  selectedVolumeIndex: 0,
};

const app = document.getElementById("app");
const toast = document.getElementById("toast");

injectProfileHeroStyles();
initTelegram();
window.addEventListener("hashchange", render);
document.addEventListener("submit", handleSubmit);
document.addEventListener("click", handleClick);
document.addEventListener("change", handleChange);
render();

function initTelegram() {
  const tg = window.Telegram && window.Telegram.WebApp;
  if (!tg) return;
  tg.ready();
  tg.expand();
  tg.setHeaderColor("#e9e1cf");
  tg.setBackgroundColor("#e9e1cf");
}

function resolveApiBase() {
  const saved = localStorage.getItem("aroma_api_base");
  if (saved !== null) return saved.trim().replace(/\/$/, "");

  const isLocal = ["localhost", "127.0.0.1", ""].includes(window.location.hostname);
  if (isLocal || window.location.protocol === "file:") {
    return "http://localhost:8080";
  }

  return "";
}

function route() {
  const hash = window.location.hash.replace(/^#/, "");
  return hash || "home";
}

function navigate(nextRoute) {
  window.location.hash = nextRoute;
  if (route() === nextRoute) render();
}

function render() {
  const currentRoute = route();

  if (currentRoute.startsWith("product/")) {
    renderProduct(currentRoute.replace("product/", ""));
    return;
  }

  const screens = {
    home: renderHome,
    intro: renderIntro,
    quiz: renderQuiz,
    loading: renderLoading,
    profile: renderProfile,
    results: renderResults,
  };

  (screens[currentRoute] || renderHome)();
}

function phone(content, className = "") {
  app.innerHTML = `<main class="phone ${className}"><div class="screen-transition">${content}</div></main>`;
}

function renderHome() {
  phone(`
    <section class="screen start-screen">
      <div class="start-heading">
        <p>Найдите аромат, который<br>подходит именно вам.</p>
        <h1>Aroma Type<span>✦</span></h1>
      </div>

      <div class="start-photo" aria-hidden="true"></div>

      <div class="start-bottom">
        <button class="start-cta" data-action="start-intro" type="button">Начать подбор</button>
        <p>Персональные рекомендации на<br>основе вашего стиля и образа жизни.</p>
      </div>
    </section>
  `);
}

function renderIntro() {
  phone(`
    <section class="screen intro-screen">
      <div class="intro-time">3 минуты</div>
      <h1 class="intro-title">Расскажите<br>немного о<br>себе</h1>
      <div class="intro-line"></div>
      <p class="intro-copy">Мы подберем ваш<br><span>персональный аромат</span><br>за 8 вопросов</p>
      <div class="intro-spark" aria-hidden="true"></div>

      <div class="bottom-actions">
        <button class="btn btn-brown" data-action="start-quiz" type="button">Продолжить</button>
      </div>
    </section>
  `);
}

async function renderQuiz() {
  phone(`
    <section class="screen soft-screen">
      <div class="state-block">
        <div>
          <div class="spinner"></div>
          <p>Загружаем вопросы</p>
        </div>
      </div>
    </section>
  `);

  try {
    if (state.questions.length === 0) {
      state.questions = await api("/api/questions");
    }
    renderQuestion();
  } catch (error) {
    renderState("Не удалось загрузить вопросы", error.message, "start-quiz");
  }
}

function renderQuestion() {
  const total = state.questions.length;
  const question = state.questions[state.currentQuestionIndex];
  if (!question) {
    submitRecommendations();
    return;
  }

  const selected = state.selectedAnswers.get(question.id) || new Set();
  const isMultiple = question.type === "multiple_choice";

  phone(`
    <section class="screen question-screen screen-with-footer">
      <div class="brand-row question-header">
        <button class="top-back" data-action="quiz-back" type="button" aria-label="Назад">‹</button>
        <span class="brand">Aroma Type<span class="spark">✦</span></span>
        <div class="step-label">Шаг<br>${state.currentQuestionIndex + 1} / ${total}</div>
      </div>
      <div class="hairline"></div>

      <h1 class="question-title">${highlightQuestion(question.text)}</h1>

      <div class="question-options">
        ${question.options.map((option) => `
          <button class="option-btn ${selected.has(option.id) ? "active" : ""}" data-action="select-answer" data-question-id="${escapeAttr(question.id)}" data-option-id="${escapeAttr(option.id)}" data-multiple="${isMultiple}" type="button">
            ${escapeHTML(option.text)}
          </button>
        `).join("")}
      </div>

      <div class="bottom-actions">
        <div class="quiz-footer">
          <button class="btn" data-action="quiz-next" type="button">${state.currentQuestionIndex + 1 === total ? "Показать результат" : "Далее"}</button>
        </div>
      </div>
    </section>
  `);
}

async function submitRecommendations() {
  const answerOptionIds = Array.from(state.selectedAnswers.values()).flatMap((set) => Array.from(set));

  navigate("loading");

  try {
    const [recommendations] = await Promise.all([
      api("/api/recommendations", {
        method: "POST",
        body: { answerOptionIds },
      }),
      delay(1300),
    ]);
    state.recommendations = recommendations;
    state.activeFilter = "all";
    navigate("profile");
  } catch (error) {
    renderState("Не удалось собрать подборку", error.message, "quiz");
  }
}

function renderLoading() {
  phone(`
    <section class="screen loading-screen">
      <div class="loading-topline"><span></span><b>Анализ</b><span></span></div>
      <h1>Подбираем<br>рекомендации...</h1>
      <div class="loading-spinner" aria-hidden="true"></div>
    </section>
  `);
}

function renderProfile() {
  if (!state.recommendations) {
    navigate("quiz");
    return;
  }

  const { profile, totalItems } = state.recommendations;
  const heroImage = getProfileHeroImage(profile.name);

  phone(`
    <section class="screen profile-screen screen-with-footer">
      <div class="brand-row profile-brand">
        <span class="brand">Aroma Type<span class="spark">✦</span></span>
      </div>

      <div class="profile-hero-image">
        <img
          class="profile-hero-photo"
          src="${escapeAttr(heroImage)}"
          alt="${escapeAttr(profile.name)}"
          loading="lazy"
        />
        <span class="eyebrow profile-eyebrow">Парфюмерный тип</span>
      </div>

      <h1 class="profile-title">${escapeHTML(profile.name)}</h1>
      <p class="profile-description">${escapeHTML(profile.description)}</p>

      <div class="divider"></div>
      <h2 class="section-title">Профиль аромата</h2>
      ${renderProfileBars(profile.profileBars || [])}

      <h2 class="section-title">Черты характера</h2>
      ${renderMetrics(profile.characterTraits || [])}

      <div class="divider"></div>
      <h2 class="section-title">Ключевые ноты</h2>
      <div class="tag-row">${(profile.keyNotes || []).map(renderTag).join("")}</div>

      <div class="bottom-actions">
        <p class="small-copy result-count">Подобрано ${pluralize(totalItems, ["аромат", "аромата", "ароматов"])} для вас</p>
        <button class="btn" data-action="show-results" type="button">Показать мои ароматы</button>
        <button class="btn btn-secondary" data-action="restart-quiz" type="button">Пройти тест заново</button>
      </div>
    </section>
  `);
}

  const { profile, items, totalItems } = state.recommendations;

  phone(`
    <section class="screen profile-screen screen-with-footer">
      <div class="brand-row profile-brand">
        <span class="brand">Aroma Type<span class="spark">✦</span></span>
      </div>

      <div class="profile-hero-image" aria-hidden="true">
        <span class="eyebrow profile-eyebrow">Парфюмерный тип</span>
      </div>

      <h1 class="profile-title">${escapeHTML(profile.name)}</h1>
      <p class="profile-description">${escapeHTML(profile.description)}</p>

      <div class="divider"></div>
      <h2 class="section-title">Профиль аромата</h2>
      ${renderProfileBars(profile.profileBars || [])}

      <h2 class="section-title">Черты характера</h2>
      ${renderMetrics(profile.characterTraits || [])}

      <div class="divider"></div>
      <h2 class="section-title">Ключевые ноты</h2>
      <div class="tag-row">${(profile.keyNotes || []).map(renderTag).join("")}</div>

      <div class="bottom-actions">
  <p class="small-copy" style="width: 100%; text-align: center; margin-bottom: 12px;">Набор из 5 миниатюр</p>
  <button class="btn" data-action="order-set" type="button">Заказать сет пробников</button>
  <button class="btn btn-secondary" data-action="restart-quiz" type="button">Пройти тест заново</button>
      </div>
    </section>
  `);
}

function renderResults() {
  if (!state.recommendations) {
    navigate("quiz");
    return;
  }

  const { profile, items, totalItems } = state.recommendations;
  const filters = unique(["all", ...items.flatMap((item) => item.mainAccords || [])]).slice(0, 4);
  const visibleItems = state.activeFilter === "all"
    ? items
    : items.filter((item) => (item.mainAccords || []).includes(state.activeFilter));

  phone(`
    <section class="screen podbor-screen screen-with-footer">
      <div class="brand-row podbor-header">
  <button class="top-back" data-action="back-profile" type="button" aria-label="Назад">‹</button>
  <span class="brand">Aroma Type<span class="spark">✦</span></span>
  <button class="top-cart" data-action="order-set" type="button" aria-label="Корзина">
    <img src="./assets/shopping-bag.png" alt="" class="top-cart-icon" />
  </button>
      </div>
      <div class="hairline"></div>

      <p class="small-copy podbor-note">На основе вашего профиля · ${escapeHTML(profile.name)}</p>
      <h1 class="podbor-title">Ваша<br>подборка</h1>
      <p class="podbor-copy">Каждый аромат выбран под ваш тип — чистый, лёгкий, сдержанный.</p>
      <div class="divider"></div>

      <div class="filters">
        ${filters.map((filter) => `<button class="chip ${state.activeFilter === filter ? "active" : ""}" data-action="filter" data-filter="${escapeAttr(filter)}" type="button">${filter === "all" ? "Все" : escapeHTML(capitalize(filter))}</button>`).join("")}
      </div>

      <div class="card-list">
        ${visibleItems.length ? visibleItems.map((item, index) => renderRecommendationCard(item, index)).join("") : `<p class="small-copy">В этой категории пока нет ароматов.</p>`}
      </div>

      <div class="bottom-actions">
        <p class="small-copy" style="text-align: center; margin-bottom: 12px;">Набор из 5 миниатюр · Доставка включена</p>
        <button class="btn" data-action="order-set" type="button">В корзину</button>
        <button class="btn btn-secondary" data-action="restart-quiz" type="button">Пройти тест заново</button>
      </div>
    </section>
  `);
}

function renderRecommendationCard(item, index = 0) {
  const cardNumber = String(index + 1).padStart(2, "0");

  return `
    <article class="fragrance-card" data-action="open-product" data-product-id="${escapeAttr(item.id)}">
      ${renderImage(item.imageUrl, item.name, "product-image")}
      <div>
        <div class="card-brand">${escapeHTML(item.brand)}</div>
        <h3 class="card-title">${escapeHTML(item.name)}</h3>
        <p class="card-reason">${escapeHTML(item.reason || "")}</p>
        <div class="tag-row">${(item.keyNotes || []).slice(0, 3).map(renderTag).join("")}</div>
        <div class="match-row">
          <div>
            <span>Совпадение</span>
            <div class="match-line"><div class="match-fill" style="width: ${safePercent(item.matchPercent)}%;"></div></div>
          </div>
          <span>${safePercent(item.matchPercent)}%</span>
        </div>
      </div>

      <button
        class="card-add-btn"
        data-action="add-to-cart"
        data-product-id="${escapeAttr(item.id)}"
        type="button"
        aria-label="Добавить в корзину"
      >+</button>

      <span class="card-number">${cardNumber}</span>
    </article>
  `;
}

async function renderProduct(productId) {
  phone(`
    <section class="screen">
      <div class="state-block">
        <div>
          <div class="spinner"></div>
          <p>Открываем аромат</p>
        </div>
      </div>
    </section>
  `);

  try {
    state.selectedProduct = await api(`/api/fragrances/${encodeURIComponent(productId)}`);
    state.selectedVolumeIndex = 0;
    renderProductLoaded();
  } catch (error) {
    renderState("Не удалось открыть аромат", error.message, "results");
  }
}

function renderProductLoaded() {
  const product = state.selectedProduct;
  const volume = (product.volumeOptions || [])[state.selectedVolumeIndex];
  const price = volume ? formatPrice(volume.price) : formatPrice(product.price);

  phone(`
    <section class="screen product-screen screen-with-footer">
      <div class="brand-row">
        <button class="top-back" data-action="back-results" type="button" aria-label="Назад">‹</button>
        <span class="brand">Aroma Type<span class="spark">✦</span></span>
        <span></span>
      </div>

      <div class="product-hero">
        <div class="product-brand">${escapeHTML(product.brand)}</div>
        <h1 class="product-title">${escapeHTML(product.name)}</h1>
        <div class="product-gender">унисекс</div>
        ${renderImage(product.imageUrl, product.name, "product-detail-image")}
      </div>

      <div class="buy-row">
        <div class="volume-box">
          <div class="volume-title">Объем / мл</div>
          <div class="volume-options">
            ${(product.volumeOptions || []).map((option, index) => `
              <button class="volume-option ${index === state.selectedVolumeIndex ? "active" : ""}" data-action="select-volume" data-volume-index="${index}" type="button">${option.volumeMl}</button>
            `).join("") || `<button class="volume-option active" type="button">50</button>`}
          </div>
        </div>
        <div>
          <div class="price">${price}</div>
          <button class="btn btn-small" data-action="order-product" type="button">Заказать</button>
        </div>
      </div>

      <p class="product-description">${escapeHTML(product.description || "Описание аромата скоро появится.")}</p>

      <div class="divider"></div>
      <h2 class="card-title">Аккорды</h2>
      <div class="tag-row">${(product.mainAccords || []).map(renderTag).join("")}</div>

      <div class="divider"></div>
      <h2 class="card-title">Пирамида парфюма</h2>
      <div class="pyramid">
        ${renderNoteLevel("Верхние ноты", product.topNotes, "◇")}
        ${renderNoteLevel("Средние ноты", product.middleNotes, "◌")}
        ${renderNoteLevel("Базовые ноты", product.baseNotes, "♢")}
      </div>

      <div class="bottom-actions">
        <button class="btn" data-action="order-product" type="button">Заказать аромат</button>
      </div>
    </section>
  `);
}

function renderState(title, message, action) {
  phone(`
    <section class="screen">
      <div class="state-block">
        <div>
          <h1 class="headline-tight">${escapeHTML(title)}</h1>
          <p class="body-copy">${escapeHTML(message || "Попробуйте ещё раз.")}</p>
          <button class="btn" data-action="${escapeAttr(action)}" type="button">Повторить</button>
        </div>
      </div>
    </section>
  `);
}

function renderProfileBars(bars) {
  if (!bars.length) return "";
  return `
    <div class="profile-bars">
      ${bars.map((bar) => `
        <div class="profile-bar">
          <div class="profile-bar-line"><div class="profile-bar-fill" style="width: ${safePercent(bar.percent)}%;"></div></div>
          <div class="profile-bar-label">${escapeHTML(bar.label)}</div>
        </div>
      `).join("")}
    </div>
  `;
}

function renderMetrics(metrics) {
  if (!metrics.length) return "";
  return `
    <div class="metric-list">
      ${metrics.map((metric, index) => `
        <div class="metric">
          <span class="metric-icon">${index + 1}</span>
          <span class="metric-label">${escapeHTML(metric.label)}</span>
          <span class="metric-percent">${safePercent(metric.percent)}%</span>
          <div class="metric-line"><div class="metric-fill" style="width: ${safePercent(metric.percent)}%;"></div></div>
        </div>
      `).join("")}
    </div>
  `;
}

function renderNoteLevel(label, notes, symbol) {
  const value = (notes || []).join(", ") || "Не указано";
  return `
    <div>
      <div class="note-level">${escapeHTML(label)}</div>
      <div class="note-symbol">${symbol}</div>
      <div class="note-value">${escapeHTML(value)}</div>
    </div>
  `;
}

function renderTag(value) {
  return `<span class="tag">${escapeHTML(capitalize(value))}</span>`;
}

function highlightQuestion(text) {
  const escaped = escapeHTML(text);
  const replacements = [
    ["через аромат", "<span>через аромат</span>"],
    ["незнакомую компанию", "<span>незнакомую компанию</span>"],
    ["максимальный прилив энергии", "<span>максимальный прилив энергии</span>"],
    ["форс-мажора", "<span>форс-мажора</span>"],
    ["транслировать окружающим", "<span>транслировать окружающим</span>"],
    ["сложный новый проект", "<span>сложный новый проект</span>"],
    ["означает «качество»", "<span>означает «качество»</span>"],
    ["уже через 15 минут", "<span>уже через 15 минут</span>"],
    ["наибольший отклик", "<span>наибольший отклик</span>"],
  ];

  return replacements.reduce((result, [source, target]) => result.replace(source, target), escaped);
}

function renderImage(url, alt, className) {
  const src = imageURL(url);
  if (!src) {
    return `<div class="${className} placeholder">${escapeHTML(alt)}</div>`;
  }
  return `<img class="${className}" src="${escapeAttr(src)}" alt="${escapeAttr(alt)}" loading="lazy" />`;
}

async function handleSubmit(event) {
  const form = event.target.closest("form");
  if (!form) return;
  event.preventDefault();
}

function handleClick(event) {
  const target = event.target.closest("[data-action]");
  if (!target) return;

  const action = target.dataset.action;

  if (action === "api-settings") {
    openApiSettings();
  }
  if (action === "start-intro") {
    navigate("intro");
  }
  if (action === "start-quiz") {
    state.selectedAnswers = new Map();
    state.currentQuestionIndex = 0;
    state.recommendations = null;
    navigate("quiz");
  }
  if (action === "select-answer") {
    selectAnswer(target);
  }
  if (action === "quiz-next") {
    goNextQuestion();
  }
  if (action === "quiz-back") {
    goPreviousQuestion();
  }
  if (action === "restart-quiz") {
    state.selectedAnswers = new Map();
    state.currentQuestionIndex = 0;
    state.recommendations = null;
    navigate("quiz");
  }
  if (action === "show-results") {
    navigate("results");
  }
  if (action === "back-profile") {
    navigate("profile");
  }
  if (action === "go-home") {
    navigate("home");
  }
  if (action === "filter") {
    state.activeFilter = target.dataset.filter;
    renderResults();
  }
  if (action === "open-product") {
    navigate(`product/${target.dataset.productId}`);
  }
  if (action === "back-results") {
    navigate("results");
  }
  if (action === "select-volume") {
    state.selectedVolumeIndex = Number(target.dataset.volumeIndex);
    renderProductLoaded();
  }
  if (action === "order-set" || action === "order-product") {
    openOrderContact();
  }
}

function handleChange(event) {
  void event;
}

function selectAnswer(target) {
  const questionId = target.dataset.questionId;
  const optionId = target.dataset.optionId;
  const isMultiple = target.dataset.multiple === "true";
  const selected = state.selectedAnswers.get(questionId) || new Set();

  if (isMultiple) {
    if (selected.has(optionId)) selected.delete(optionId);
    else selected.add(optionId);
  } else {
    selected.clear();
    selected.add(optionId);
  }

  state.selectedAnswers.set(questionId, selected);

  const options = target.closest(".question-options");
  if (!options) return;

  if (!isMultiple) {
    options.querySelectorAll(".option-btn").forEach((button) => {
      button.classList.remove("active");
    });
  }
  target.classList.toggle("active", selected.has(optionId));
}

function goNextQuestion() {
  const question = state.questions[state.currentQuestionIndex];
  const selected = state.selectedAnswers.get(question.id);

  if (!selected || selected.size === 0) {
    showToast("Выберите вариант ответа");
    return;
  }

  state.currentQuestionIndex += 1;
  if (state.currentQuestionIndex >= state.questions.length) {
    submitRecommendations();
    return;
  }
  renderQuestion();
}

function goPreviousQuestion() {
  if (state.currentQuestionIndex === 0) {
    navigate("intro");
    return;
  }
  state.currentQuestionIndex -= 1;
  renderQuestion();
}

function openApiSettings() {
  const nextValue = window.prompt("API base URL", state.apiBase || window.location.origin);
  if (nextValue === null) return;

  state.apiBase = nextValue.trim().replace(/\/$/, "");
  localStorage.setItem("aroma_api_base", state.apiBase);
  showToast("API URL сохранён");
}

function openOrderContact() {
  const product = state.selectedProduct;
  const text = product
    ? `Здравствуйте! Хочу заказать аромат ${product.brand} ${product.name}.`
    : "Здравствуйте! Хочу заказать сет пробников Aroma Type.";
  const separator = ORDER_CONTACT_URL.includes("?") ? "&" : "?";
  const url = `${ORDER_CONTACT_URL}${separator}text=${encodeURIComponent(text)}`;

  const tg = window.Telegram && window.Telegram.WebApp;
  if (tg) {
    tg.openTelegramLink(url);
    return;
  }
  window.open(url, "_blank", "noopener,noreferrer");
}

async function api(path, options = {}) {
  const headers = new Headers(options.headers || {});
  const request = {
    method: options.method || "GET",
    headers,
  };

  if (options.formData) {
    request.body = options.formData;
  } else if (options.body !== undefined) {
    headers.set("Content-Type", "application/json");
    request.body = JSON.stringify(options.body);
  }

  const response = await fetch(`${state.apiBase}${path}`, request);
  const contentType = response.headers.get("content-type") || "";
  const payload = contentType.includes("application/json")
    ? await response.json()
    : await response.text();

  if (!response.ok) {
    const message = typeof payload === "object" && payload.message
      ? payload.message
      : "Запрос не выполнен";
    throw new Error(message);
  }

  return payload;
}

function numberValue(value) {
  const parsed = Number(String(value || "0").replace(",", ".").replace(/[^\d.]/g, ""));
  return Number.isFinite(parsed) ? parsed : 0;
}

function formatPrice(value) {
  const number = numberValue(value);
  if (!number) return "Цена по запросу";
  return `${new Intl.NumberFormat("ru-RU").format(number)} ₽`;
}

function imageURL(value) {
  if (!value) return "";
  if (/^(https?:|data:|blob:)/.test(value)) return value;
  return `${state.apiBase}${value}`;
}

function unique(values) {
  return Array.from(new Set(values.filter(Boolean)));
}

function capitalize(value) {
  const text = String(value || "");
  return text ? text[0].toUpperCase() + text.slice(1) : "";
}

function safePercent(value) {
  const number = Number(value);
  if (!Number.isFinite(number)) return 0;
  return Math.max(0, Math.min(100, Math.round(number)));
}

function pluralize(count, forms) {
  const absolute = Math.abs(Number(count)) % 100;
  const lastDigit = absolute % 10;
  let form = forms[2];

  if (absolute > 10 && absolute < 20) {
    form = forms[2];
  } else if (lastDigit > 1 && lastDigit < 5) {
    form = forms[1];
  } else if (lastDigit === 1) {
    form = forms[0];
  }

  return `${count} ${form}`;
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function escapeHTML(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function escapeAttr(value) {
  return escapeHTML(value);
}

let toastTimer = null;
function showToast(message) {
  toast.textContent = message;
  toast.classList.add("show");
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => toast.classList.remove("show"), 3200);
}
