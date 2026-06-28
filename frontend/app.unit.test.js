import { fireEvent, screen } from "@testing-library/dom";
import { describe, expect, it, vi } from "vitest";
import { loadApp } from "./testUtils.js";

describe("Aroma Type frontend components", () => {
  it("renders the home screen and opens the intro screen", async () => {
    await loadApp(vi.fn());

    expect(screen.getByRole("button", { name: /Начать подбор/i })).toBeInTheDocument();
    expect(screen.getByText(/Персональные рекомендации/i)).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: /Начать подбор/i }));

    expect(screen.getByText(/Расскажите/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Продолжить/i })).toBeInTheDocument();
  });
});
