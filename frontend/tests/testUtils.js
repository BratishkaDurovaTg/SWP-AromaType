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

export async function loadApp(fetchMock = vi.fn()) {
  vi.resetModules();

  document.body.innerHTML = `
    <div id="app" class="app-shell"></div>
    <div id="toast" class="toast" role="status" aria-live="polite"></div>
  `;

  localStorage.clear();
  localStorage.setItem("aroma_api_base", "");
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
