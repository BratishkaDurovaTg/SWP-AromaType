import "@testing-library/jest-dom/vitest";
import { vi } from "vitest";

export function jsonResponse(data, status = 200) {
  return {
    ok: status >= 200 && status < 300,
    status,
    headers: new Headers({
      "content-type": "application/json"
    }),
    json: async () => data,
    text: async () => JSON.stringify(data)
  };
}

function createMemoryStorage() {
  const values = new Map();

  return {
    clear: vi.fn(() => values.clear()),
    getItem: vi.fn((key) => values.get(String(key)) ?? null),
    key: vi.fn((index) => Array.from(values.keys())[index] ?? null),
    removeItem: vi.fn((key) => values.delete(String(key))),
    setItem: vi.fn((key, value) => values.set(String(key), String(value))),
    get length() {
      return values.size;
    }
  };
}

export async function loadApp(fetchMock = vi.fn()) {
  vi.resetModules();

  document.body.innerHTML = `
    <div id="app" class="app-shell"></div>
    <div id="toast" class="toast" role="status" aria-live="polite"></div>
  `;

  const storage = createMemoryStorage();
  Object.defineProperty(window, "localStorage", {
    configurable: true,
    value: storage
  });
  Object.defineProperty(globalThis, "localStorage", {
    configurable: true,
    value: storage
  });
  window.localStorage.setItem("aroma_api_base", "");
  window.location.hash = "";

  window.Telegram = {
    WebApp: {
      ready: vi.fn(),
      expand: vi.fn(),
      setHeaderColor: vi.fn(),
      setBackgroundColor: vi.fn(),
      openTelegramLink: vi.fn()
    }
  };

  globalThis.fetch = fetchMock;

  await import("../app.js");
}
