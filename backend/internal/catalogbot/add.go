package catalogbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
)

const (
	addStepID = iota
	addStepName
	addStepBrand
	addStepPrice
	addStepVolumes
	addStepDescription
	addStepTop
	addStepMiddle
	addStepBase
	addStepAccords
	addStepPsychotype
	addStepScores
	addStepActive
	addStepPhoto
)

func addPrompt(step int) string {
	switch step {
	case addStepID:
		return "ID товара латиницей, например miami-shake."
	case addStepName:
		return "Название товара."
	case addStepBrand:
		return "Бренд."
	case addStepPrice:
		return "Цена, например 8393."
	case addStepVolumes:
		return "Объемы в формате 50:8393, 100:12990. Если не нужно, отправьте -"
	case addStepDescription:
		return "Описание товара."
	case addStepTop:
		return "Верхние ноты через запятую. Если нет, отправьте -"
	case addStepMiddle:
		return "Средние ноты через запятую. Если нет, отправьте -"
	case addStepBase:
		return "Базовые ноты через запятую. Если нет, отправьте -"
	case addStepAccords:
		return "Основные аккорды через запятую."
	case addStepPsychotype:
		return "Психотип: drive, focus, aesthetic, power или balanced."
	case addStepScores:
		return "Scores: drive:20, focus:35, aesthetic:90, power:25"
	case addStepActive:
		return "Активен? yes/no"
	case addStepPhoto:
		return "Отправьте фото товара или напишите skip."
	default:
		return "Следующее значение."
	}
}

func (b *Bot) handleAddStep(ctx context.Context, chatID int64, s *session, text string) error {
	text = strings.TrimSpace(text)
	switch s.step {
	case addStepID:
		id, err := validateID(text)
		if err != nil {
			return b.telegram.sendMessage(ctx, chatID, "ID должен быть латиницей: miami-shake, citrus_day. Попробуйте еще раз.")
		}
		if _, err := b.repo.GetFragrance(ctx, id); err == nil {
			return b.telegram.sendMessage(ctx, chatID, "Такой ID уже есть. Введите другой ID.")
		}
		s.draft.ID = id
	case addStepName:
		s.draft.Name = text
	case addStepBrand:
		s.draft.Brand = text
	case addStepPrice:
		price, err := parsePrice(text)
		if err != nil {
			return b.telegram.sendMessage(ctx, chatID, "Цена должна быть числом, например 8393.")
		}
		s.draft.Price = price
	case addStepVolumes:
		volumes, err := parseVolumes(text)
		if err != nil {
			return b.telegram.sendMessage(ctx, chatID, "Формат объемов: 50:8393, 100:12990 или -")
		}
		s.draft.VolumeOptions = volumes
	case addStepDescription:
		s.draft.Description = text
	case addStepTop:
		s.draft.TopNotes = optionalList(text)
	case addStepMiddle:
		s.draft.MiddleNotes = optionalList(text)
	case addStepBase:
		s.draft.BaseNotes = optionalList(text)
	case addStepAccords:
		s.draft.MainAccords = optionalList(text)
	case addStepPsychotype:
		psychotype, err := normalizePsychotype(text)
		if err != nil {
			return b.telegram.sendMessage(ctx, chatID, "Психотип: drive, focus, aesthetic, power или balanced.")
		}
		s.draft.Psychotype = psychotype
	case addStepScores:
		scores, err := parseScores(text)
		if err != nil {
			return b.telegram.sendMessage(ctx, chatID, "Формат scores: drive:20, focus:35, aesthetic:90, power:25")
		}
		s.draft.PsychotypeScores = scores
	case addStepActive:
		active, err := parseActive(text)
		if err != nil {
			return b.telegram.sendMessage(ctx, chatID, "Ответьте yes или no.")
		}
		s.draft.IsActive = active
	case addStepPhoto:
		if strings.EqualFold(text, "skip") || text == "-" {
			return b.finishAdd(ctx, chatID, s)
		}
		return b.telegram.sendMessage(ctx, chatID, "На этом шаге отправьте фото или skip.")
	}

	s.step++
	return b.telegram.sendMessageMarkup(ctx, chatID, addPrompt(s.step), mainMenuKeyboard())
}

func (b *Bot) handleAddPhoto(ctx context.Context, chatID int64, s *session, photos []photoSize) error {
	imageURL, err := b.saveTelegramPhoto(ctx, photos)
	if err != nil {
		return err
	}
	s.draft.ImageURL = imageURL
	return b.finishAdd(ctx, chatID, s)
}

func (b *Bot) finishAdd(ctx context.Context, chatID int64, s *session) error {
	if err := validateFragrance(s.draft); err != nil {
		return b.telegram.sendMessage(ctx, chatID, "Товар не сохранен: не хватает обязательных данных. /cancel и начните заново.")
	}
	if s.draft.MainAccords == nil {
		s.draft.MainAccords = []string{}
	}
	if s.draft.TopNotes == nil {
		s.draft.TopNotes = []string{}
	}
	if s.draft.MiddleNotes == nil {
		s.draft.MiddleNotes = []string{}
	}
	if s.draft.BaseNotes == nil {
		s.draft.BaseNotes = []string{}
	}
	if s.draft.VolumeOptions == nil {
		s.draft.VolumeOptions = []questionnaire.VolumeOption{}
	}
	if err := b.repo.UpsertFragrance(ctx, s.draft); err != nil {
		return err
	}
	saved := s.draft
	s.action = ""
	s.step = 0
	s.draft = questionnaire.Fragrance{}
	return b.telegram.sendMessageMarkup(ctx, chatID, fmt.Sprintf("Товар создан.\n\n%s", formatFragrance(saved)), itemKeyboard(saved.ID))
}

func optionalList(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" || strings.EqualFold(value, "skip") {
		return []string{}
	}
	return splitList(value)
}
