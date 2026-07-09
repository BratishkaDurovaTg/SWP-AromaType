import { fireEvent, screen, waitFor } from "@testing-library/dom";
import { describe, expect, it, vi } from "vitest";
import { jsonResponse, loadApp } from "./testUtils.js";

const questionsFixture = [
  {
    id: "q1",
    text: "Что вы хотите выразить через аромат?",
    type: "single_choice",
    options: [
      {
        id: "a-focus",
        text: "Собранность и фокус"
      }
    ]
  }
];

const recommendationsFixture = {
  profile: {
    name: "Интеллект и фокус",
    description: "Вам подходят собранные, глубокие и немного загадочные ароматы.",
    profileBars: [
      { label: "Драйв", percent: 30 },
      { label: "Фокус", percent: 96 },
      { label: "Эстетика", percent: 55 },
      { label: "Доминанта", percent: 42 }
    ],
    characterTraits: [
      { label: "Внутренняя собранность", percent: 96 }
    ],
    keyNotes: ["ветивер", "ладан", "черный перец"]
  },
  items: [
    {
      id: "p1",
      brand: "AromaType",
      name: "Mystic Night",
      gender: "male",
      imageUrl: "",
      reason: "Совпадает с вашим профилем по тегам: Власть / Доминанта.",
      keyNotes: ["Черный перец", "Ладан", "Ветивер"],
      mainAccords: ["deep"],
      psychotype: "power",
      psychotypeScores: { drive: 10, focus: 25, aesthetic: 10, power: 55 },
      matchPercent: 99,
      score: 1250,
      price: "8393"
    },
    {
      id: "p2",
      brand: "AromaType",
      name: "Fresh Office",
      gender: "unisex",
      imageUrl: "",
      reason: "Совпадает с вашим профилем по тегам: Интеллект / Фокус.",
      keyNotes: ["Бергамот", "Лимон", "Лаванда"],
      mainAccords: ["fresh"],
      psychotype: "focus",
      psychotypeScores: { drive: 10, focus: 70, aesthetic: 15, power: 5 },
      matchPercent: 95,
      score: 1100,
      price: "7200"
    }
  ],
  totalItems: 2
};

const productFixture = {
  id: "p1",
  brand: "AromaType",
  name: "Mystic Night",
  gender: "female",
  imageUrl: "",
  price: "8393",
  volumeOptions: [
    { volumeMl: 3, price: 8393 },
    { volumeMl: 5, price: 11990 },
    { volumeMl: 10, price: 18990 }
  ],
  description: "Глубокий и загадочный аромат.",
  topNotes: ["Черный перец"],
  middleNotes: ["Ладан"],
  baseNotes: ["Ветивер"],
  mainAccords: ["Глубокий", "Загадочный"],
  psychotype: "power",
  psychotypeScores: { drive: 10, focus: 25, aesthetic: 10, power: 55 },
  isActive: true
};

describe("Aroma Type frontend API integration", () => {
  it("loads questions, sends selected answers, and renders recommendation results", async () => {
    const fetchMock = vi.fn(async (url, options = {}) => {
      if (url === "/api/questions") {
        return jsonResponse(questionsFixture);
      }

      if (url === "/api/recommendations") {
        expect(options.method).toBe("POST");
        expect(JSON.parse(options.body)).toEqual({
          answerOptionIds: ["a-focus"]
        });

        return jsonResponse(recommendationsFixture);
      }

      if (url === "/api/fragrances/p1") {
        return jsonResponse(productFixture);
      }

      throw new Error(`Unexpected API call: ${url}`);
    });

    await loadApp(fetchMock);

    fireEvent.click(screen.getByRole("button", { name: /Начать подбор/i }));
    fireEvent.click(screen.getByRole("button", { name: /Продолжить/i }));

    await waitFor(() => {
      expect(screen.getByText(/Что вы хотите выразить/i)).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole("button", { name: /Собранность и фокус/i }));
    fireEvent.click(screen.getByRole("button", { name: /Показать результат/i }));

    await waitFor(() => {
      expect(screen.getByText("Интеллект и фокус")).toBeInTheDocument();
    }, { timeout: 2500 });

    expect(fetchMock).toHaveBeenCalledWith(
      "/api/questions",
      expect.objectContaining({
        method: "GET"
      })
    );

    expect(fetchMock).toHaveBeenCalledWith(
      "/api/recommendations",
      expect.objectContaining({
        method: "POST"
      })
    );

    const profileImage = document.querySelector(".profile-hero-photo");
    expect(profileImage).toBeInTheDocument();
    expect(profileImage).toHaveAttribute("src", "./assets/profile-focus.png");

    fireEvent.click(screen.getByRole("button", { name: /Показать мои ароматы/i }));

    expect(screen.getAllByText("Mystic Night").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Fresh Office").length).toBeGreaterThan(0);

    const cardNumbers = [...document.querySelectorAll(".card-number")].map((node) =>
      node.textContent.trim()
    );
    expect(cardNumbers).toEqual(["01", "02"]);

    const addButtons = document.querySelectorAll(".card-add-btn");
    expect(addButtons).toHaveLength(2);

    fireEvent.click(screen.getAllByText("Mystic Night")[0]);

    await waitFor(() => {
      expect(screen.getByText("для женщин")).toBeInTheDocument();
    });

    expect(screen.getByText("8 393 ₽")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "5" }));
    expect(screen.getByText("11 990 ₽")).toBeInTheDocument();

    fireEvent.click(screen.getAllByRole("button", { name: /Добавить/i })[0]);
    expect(screen.getByText("Аромат добавлен в корзину")).toBeInTheDocument();
  });
});
