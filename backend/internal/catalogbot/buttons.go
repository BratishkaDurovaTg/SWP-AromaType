package catalogbot

import (
	"fmt"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
)

const listButtonLimit = 40

type editableField struct {
	key   string
	label string
}

var editableFields = []editableField{
	{key: "name", label: "Название"},
	{key: "brand", label: "Бренд"},
	{key: "price", label: "Цена"},
	{key: "volumes", label: "Объемы"},
	{key: "description", label: "Описание"},
	{key: "accords", label: "Аккорды"},
	{key: "top", label: "Верхние ноты"},
	{key: "middle", label: "Средние ноты"},
	{key: "base", label: "Базовые ноты"},
	{key: "psychotype", label: "Психотип"},
	{key: "scores", label: "Scores"},
	{key: "active", label: "Активность"},
	{key: "image_url", label: "URL фото"},
}

func mainMenuKeyboard() replyKeyboardMarkup {
	return replyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]keyboardButton{
			{{Text: "Каталог"}, {Text: "Добавить товар"}},
			{{Text: "Помощь"}, {Text: "Отмена"}},
		},
	}
}

func listKeyboard(items []questionnaire.Fragrance) inlineKeyboardMarkup {
	rows := make([][]inlineKeyboardButton, 0)
	limit := len(items)
	if limit > listButtonLimit {
		limit = listButtonLimit
	}

	for index := 0; index < limit; index++ {
		item := items[index]
		text := fmt.Sprintf("%d. %s", index+1, item.Name)
		if item.Brand != "" {
			text = fmt.Sprintf("%d. %s / %s", index+1, item.Brand, item.Name)
		}
		rows = append(rows, []inlineKeyboardButton{{
			Text:         text,
			CallbackData: "view:" + item.ID,
		}})
	}

	rows = append(rows, []inlineKeyboardButton{
		{Text: "Добавить товар", CallbackData: "add"},
		{Text: "Помощь", CallbackData: "help"},
	})
	return inlineKeyboardMarkup{InlineKeyboard: rows}
}

func itemKeyboard(id string) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{InlineKeyboard: [][]inlineKeyboardButton{
		{
			{Text: "Редактировать", CallbackData: "edit:" + id},
			{Text: "Фото", CallbackData: "photo:" + id},
		},
		{
			{Text: "Вкл/выкл", CallbackData: "toggle:" + id},
			{Text: "Каталог", CallbackData: "list"},
		},
	}}
}

func editFieldKeyboard(id string) inlineKeyboardMarkup {
	rows := make([][]inlineKeyboardButton, 0, len(editableFields)/2+2)
	for index := 0; index < len(editableFields); index += 2 {
		row := []inlineKeyboardButton{fieldButton(id, editableFields[index])}
		if index+1 < len(editableFields) {
			row = append(row, fieldButton(id, editableFields[index+1]))
		}
		rows = append(rows, row)
	}
	rows = append(rows, []inlineKeyboardButton{
		{Text: "Фото", CallbackData: "photo:" + id},
		{Text: "Назад к товару", CallbackData: "view:" + id},
	})
	return inlineKeyboardMarkup{InlineKeyboard: rows}
}

func fieldButton(id string, field editableField) inlineKeyboardButton {
	return inlineKeyboardButton{
		Text:         field.label,
		CallbackData: "field:" + id + ":" + field.key,
	}
}
