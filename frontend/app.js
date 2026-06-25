const ORDER_CONTACT_URL = "https://t.me/aroma_type_test_bot";
const DEFAULT_TAG_IDS = [
  "psych_drive",
  "psych_focus",
  "psych_aesthetic",
  "psych_power",
  "fresh",
  "clean",
  "daily",
  "office",
  "light",
  "calm",
  "romantic",
  "soft",
  "warm",
  "bright",
  "energy",
  "party",
  "noticeable",
  "mystery",
  "deep",
  "night",
  "reliable",
  "cozy",
  "date",
  "trail",
  "morning",
];

const state = {
  authMode: "login",
  token: localStorage.getItem("aroma_token") || "",
  user: readJSON("aroma_user"),
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

function readJSON(key) {
  try {
    return JSON.parse(localStorage.getItem(key) || "null");
  } catch {
    return null;
  }
}

function setAuth(result) {
  state.token = result.accessToken;
  state.user = result.user;
  localStorage.setItem("aroma_token", result.accessToken);
  localStorage.setItem("aroma_user", JSON.stringify(result.user));
}

function clearAuth() {
  state.token = "";
  state.user = null;
  localStorage.removeItem("aroma_token");
  localStorage.removeItem("aroma_user");
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
    auth: renderAuth,
    home: renderHome,
    intro: renderIntro,
    quiz: renderQuiz,
    loading: renderLoading,
    profile: renderProfile,
    results: renderResults,
    admin: renderAdmin,
  };

  (screens[currentRoute] || renderHome)();
}

function phone(content, className = "") {
  app.innerHTML = `<main class="phone ${className}"><div class="screen-transition">${content}</div></main>`;
}

function renderAuth() {
  phone(`
    <section class="screen auth-screen">
      <div class="brand-row">
        <div>
          <p class="subtitle">Найдите аромат, который подходит именно вам.</p>
          <span class="brand-big">Aroma Type<span class="spark">✦</span></span>
        </div>
      </div>

      <div class="auth-panel">
        <div class="segmented" aria-label="Auth mode">
          <button class="segment ${state.authMode === "login" ? "active" : ""}" data-action="auth-mode" data-mode="login" type="button">Вход</button>
          <button class="segment ${state.authMode === "register" ? "active" : ""}" data-action="auth-mode" data-mode="register" type="button">Регистрация</button>
        </div>

        <form class="form" data-form="auth" style="margin-top: 18px;">
          <label class="field">
            <span>Email</span>
            <input class="input" name="email" type="email" autocomplete="email" placeholder="name@example.com" required />
          </label>
          <label class="field">
            <span>Пароль</span>
            <input class="input" name="password" type="password" autocomplete="${state.authMode === "login" ? "current-password" : "new-password"}" placeholder="минимум 6 символов" required minlength="6" />
          </label>
          <button class="btn" type="submit">${state.authMode === "login" ? "Войти" : "Создать аккаунт"}</button>
        </form>

      </div>
    </section>
  `);
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
        <p class="small-copy result-count">Подобрано ${pluralize(totalItems, ["аромат", "аромата", "ароматов"])} для вас</p>
        <button class="btn" data-action="show-results" type="button">Показать мои ароматы</button>
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
        <div class="step-label">${totalItems}<br>${pluralize(totalItems, ["вариант", "варианта", "вариантов"]).replace(/^\d+\s/, "")}</div>
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
        ${visibleItems.length ? visibleItems.map(renderRecommendationCard).join("") : `<p class="small-copy">В этой категории пока нет ароматов.</p>`}
      </div>

      <div class="bottom-actions">
        <p class="small-copy" style="text-align: center; margin-bottom: 12px;">Набор из 5 миниатюр · Доставка включена</p>
        <button class="btn" data-action="order-set" type="button">Заказать сет пробников</button>
        <button class="btn btn-secondary" data-action="restart-quiz" type="button">Пройти тест заново</button>
      </div>
    </section>
  `);
}

function renderRecommendationCard(item) {
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

function renderAdmin() {
  const isAdmin = state.user && state.user.role === "admin";

  phone(`
    <section class="screen admin-screen">
      <div class="brand-row">
        <button class="top-back" data-action="${state.token ? "go-home" : "auth-back"}" type="button" aria-label="Назад">‹</button>
        <span class="brand">Aroma Type<span class="spark">✦</span></span>
      </div>

      <h1 class="admin-title">Добавить продукт</h1>

      ${isAdmin ? renderAdminForm() : renderAdminLogin()}
    </section>
  `);
}

function renderAdminLogin() {
  return `
    <form class="form" data-form="admin-login">
      <label class="field">
        <span>Email</span>
        <input class="input" name="email" type="email" value="admin@aromatype.local" required />
      </label>
      <label class="field">
        <span>Пароль</span>
        <input class="input" name="password" type="password" placeholder="Пароль администратора" required />
      </label>
      <button class="btn" type="submit">Войти как админ</button>
    </form>
  `;
}

function renderAdminForm() {
  return `
    <form class="admin-grid" data-form="create-fragrance">
      <label class="field">
        <span>Название</span>
        <input class="input" name="name" placeholder="Введите название товара" required />
      </label>
      <label class="field">
        <span>Бренд</span>
        <input class="input" name="brand" placeholder="Например, Juliette Has A Gun" required />
      </label>
      <div class="two-col">
        <label class="field">
          <span>Стоимость</span>
          <input class="input" name="price" inputmode="decimal" placeholder="8393" required />
        </label>
        <label class="field">
          <span>Объемы</span>
          <input class="input" name="volumes" placeholder="50:8393, 100:12990" required />
        </label>
      </div>
      <label class="field">
        <span>Фото</span>
        <input class="input" name="photo" type="file" accept="image/jpeg,image/png,image/webp" />
        <input name="imageUrl" type="hidden" />
        <img id="admin-photo-preview" class="upload-preview hidden" alt="" />
      </label>
      <label class="field">
        <span>Основные аккорды</span>
        <input class="input" name="mainAccords" placeholder="сладкий, ванильный, фруктовый" required />
      </label>
      <label class="field">
        <span>Верхние ноты</span>
        <input class="input" name="topNotes" placeholder="клубника, бергамот" />
      </label>
      <label class="field">
        <span>Средние ноты</span>
        <input class="input" name="middleNotes" placeholder="мороженое, жасмин" />
      </label>
      <label class="field">
        <span>Базовые ноты</span>
        <input class="input" name="baseNotes" placeholder="ваниль, мускус" />
      </label>
      <label class="field">
        <span>Описание</span>
        <textarea class="textarea" name="description" placeholder="Введите описание" required></textarea>
      </label>
      <label class="field">
        <span>Теги для подбора</span>
        <input class="input" name="tagIds" list="tag-suggestions" placeholder="psych_drive, fresh, clean" />
        <datalist id="tag-suggestions">
          ${DEFAULT_TAG_IDS.map((tag) => `<option value="${tag}"></option>`).join("")}
        </datalist>
      </label>
      <button class="btn" type="submit">Сохранить товар</button>
    </form>
  `;
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

  const formType = form.dataset.form;
  if (formType === "auth") {
    await submitAuth(form);
  }
  if (formType === "admin-login") {
    await submitAdminLogin(form);
  }
  if (formType === "create-fragrance") {
    await submitCreateFragrance(form);
  }
}

async function submitAuth(form) {
  const data = new FormData(form);
  const path = state.authMode === "login" ? "/api/auth/login" : "/api/auth/register";

  try {
    const result = await api(path, {
      method: "POST",
      body: {
        email: data.get("email"),
        password: data.get("password"),
      },
    });
    setAuth(result);
    showToast(state.authMode === "login" ? "Вы вошли" : "Аккаунт создан");
    navigate("home");
  } catch (error) {
    showToast(error.message);
  }
}

async function submitAdminLogin(form) {
  const data = new FormData(form);

  try {
    const result = await api("/api/auth/login", {
      method: "POST",
      body: {
        email: data.get("email"),
        password: data.get("password"),
      },
    });
    setAuth(result);
    if (result.user.role !== "admin") {
      showToast("У этого аккаунта нет прав администратора");
      return;
    }
    renderAdmin();
  } catch (error) {
    showToast(error.message);
  }
}

async function submitCreateFragrance(form) {
  const data = new FormData(form);
  let imageUrl = data.get("imageUrl") || "";
  const photo = data.get("photo");

  try {
    if (photo && photo.size > 0) {
      const upload = new FormData();
      upload.append("photo", photo);
      const uploaded = await api("/api/admin/uploads/fragrance-photo", {
        method: "POST",
        formData: upload,
        auth: true,
      });
      imageUrl = uploaded.imageUrl;
    }

    const payload = {
      name: data.get("name"),
      brand: data.get("brand"),
      imageUrl,
      price: numberValue(data.get("price")),
      volumeOptions: parseVolumeOptions(data.get("volumes")),
      description: data.get("description"),
      topNotes: splitList(data.get("topNotes")),
      middleNotes: splitList(data.get("middleNotes")),
      baseNotes: splitList(data.get("baseNotes")),
      mainAccords: splitList(data.get("mainAccords")),
      tagIds: splitList(data.get("tagIds")),
      isActive: true,
    };

    const product = await api("/api/admin/fragrances", {
      method: "POST",
      body: payload,
      auth: true,
    });

    showToast("Товар сохранён");
    form.reset();
    state.selectedProduct = product;
    navigate(`product/${product.id}`);
  } catch (error) {
    showToast(error.message);
  }
}

function handleClick(event) {
  const target = event.target.closest("[data-action]");
  if (!target) return;

  const action = target.dataset.action;

  if (action === "auth-mode") {
    state.authMode = target.dataset.mode;
    renderAuth();
  }
  if (action === "api-settings") {
    openApiSettings();
  }
  if (action === "open-admin") {
    navigate("admin");
  }
  if (action === "auth-back") {
    navigate("auth");
  }
  if (action === "logout") {
    clearAuth();
    navigate("auth");
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
  const input = event.target;
  if (input.name !== "photo" || !input.files || !input.files[0]) return;

  const preview = document.getElementById("admin-photo-preview");
  if (!preview) return;

  preview.src = URL.createObjectURL(input.files[0]);
  preview.classList.remove("hidden");
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

  if (options.auth) {
    headers.set("Authorization", `Bearer ${state.token}`);
  }

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

function splitList(value) {
  return String(value || "")
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean);
}

function parseVolumeOptions(value) {
  return splitList(value).map((item) => {
    const [volume, price] = item.split(":").map((part) => part.trim());
    return {
      volumeMl: parseInt(volume, 10),
      price: numberValue(price),
    };
  }).filter((item) => item.volumeMl > 0 && item.price >= 0);
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
