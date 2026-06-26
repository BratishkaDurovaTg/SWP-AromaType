package catalogbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/catalog"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
	"github.com/google/uuid"
)

const maxTelegramPhotoSize = 8 << 20

type Config struct {
	Token     string
	Password  string
	UploadDir string
}

type Bot struct {
	cfg      Config
	logger   *slog.Logger
	repo     *catalog.Repository
	telegram *telegramClient
	sessions map[int64]*session
}

type session struct {
	authenticated bool
	action        string
	targetID      string
	step          int
	draft         questionnaire.Fragrance
}

func New(cfg Config, logger *slog.Logger, repo *catalog.Repository) *Bot {
	return &Bot{
		cfg:      cfg,
		logger:   logger,
		repo:     repo,
		telegram: newTelegramClient(cfg.Token),
		sessions: make(map[int64]*session),
	}
}

func (b *Bot) Run(ctx context.Context) error {
	if b.cfg.Token == "" {
		return errors.New("CATALOG_BOT_TOKEN must not be empty")
	}
	if b.cfg.Password == "" {
		return errors.New("CATALOG_BOT_PASSWORD must not be empty")
	}
	if b.cfg.UploadDir == "" {
		return errors.New("UPLOAD_DIR must not be empty")
	}
	if err := os.MkdirAll(b.cfg.UploadDir, 0o755); err != nil {
		return err
	}

	offset := 0
	b.logger.Info("catalog bot started")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		updates, err := b.telegram.getUpdates(ctx, offset)
		if err != nil {
			b.logger.Error("failed to get telegram updates", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			if update.Message == nil {
				continue
			}
			if err := b.handleMessage(ctx, *update.Message); err != nil {
				b.logger.Error("failed to handle message", "error", err)
				_ = b.telegram.sendMessage(ctx, update.Message.Chat.ID, "Ошибка: "+err.Error())
			}
		}
	}
}

func (b *Bot) handleMessage(ctx context.Context, msg message) error {
	chatID := msg.Chat.ID
	s := b.session(chatID)
	text := strings.TrimSpace(msg.Text)

	if text == "/start" {
		if s.authenticated {
			return b.sendMenu(ctx, chatID)
		}
		s.action = "password"
		return b.telegram.sendMessage(ctx, chatID, "Введите пароль админки каталога.")
	}

	if !s.authenticated {
		if s.action != "password" {
			s.action = "password"
			return b.telegram.sendMessage(ctx, chatID, "Сначала авторизация. Отправьте /start и введите пароль.")
		}
		if text != b.cfg.Password {
			return b.telegram.sendMessage(ctx, chatID, "Неверный пароль. Попробуйте еще раз.")
		}
		s.authenticated = true
		s.action = ""
		return b.sendMenu(ctx, chatID)
	}

	if len(msg.Photo) > 0 && s.action == "photo" {
		return b.handlePhoto(ctx, chatID, s, msg.Photo)
	}
	if len(msg.Photo) > 0 && s.action == "add" && s.step == addStepPhoto {
		return b.handleAddPhoto(ctx, chatID, s, msg.Photo)
	}

	if text == "" {
		return b.telegram.sendMessage(ctx, chatID, "Отправьте команду или текстовое значение.")
	}
	if text == "/cancel" {
		b.sessions[chatID] = &session{authenticated: true}
		return b.telegram.sendMessage(ctx, chatID, "Действие отменено.")
	}

	if s.action == "add" {
		return b.handleAddStep(ctx, chatID, s, text)
	}

	return b.handleCommand(ctx, chatID, s, text)
}

func (b *Bot) handleCommand(ctx context.Context, chatID int64, s *session, text string) error {
	command, rest, _ := strings.Cut(text, " ")
	switch command {
	case "/help", "help":
		return b.sendMenu(ctx, chatID)
	case "/list":
		return b.list(ctx, chatID)
	case "/view":
		return b.view(ctx, chatID, strings.TrimSpace(rest))
	case "/add":
		s.action = "add"
		s.step = addStepID
		s.draft = questionnaire.Fragrance{
			IsActive:         true,
			Psychotype:       "balanced",
			PsychotypeScores: questionnaire.PsychotypeScores{Drive: 50, Focus: 50, Aesthetic: 50, Power: 50},
		}
		return b.telegram.sendMessage(ctx, chatID, addPrompt(s.step))
	case "/edit":
		id := strings.TrimSpace(rest)
		if id == "" {
			return b.telegram.sendMessage(ctx, chatID, "Формат: /edit fragrance-id")
		}
		item, err := b.repo.GetFragrance(ctx, id)
		if err != nil {
			return err
		}
		return b.telegram.sendMessage(ctx, chatID, formatFragrance(item)+"\n\n"+editHelp(id))
	case "/set":
		return b.setField(ctx, chatID, strings.TrimSpace(rest))
	case "/photo":
		id := strings.TrimSpace(rest)
		if id == "" {
			return b.telegram.sendMessage(ctx, chatID, "Формат: /photo fragrance-id")
		}
		if _, err := b.repo.GetFragrance(ctx, id); err != nil {
			return err
		}
		s.action = "photo"
		s.targetID = id
		return b.telegram.sendMessage(ctx, chatID, "Отправьте фото для товара "+id+".")
	case "/toggle":
		return b.toggle(ctx, chatID, strings.TrimSpace(rest))
	default:
		return b.telegram.sendMessage(ctx, chatID, "Не понял команду.\n\n"+menuText())
	}
}

func (b *Bot) session(chatID int64) *session {
	s, ok := b.sessions[chatID]
	if !ok {
		s = &session{}
		b.sessions[chatID] = s
	}
	return s
}

func (b *Bot) sendMenu(ctx context.Context, chatID int64) error {
	return b.telegram.sendMessage(ctx, chatID, "Доступ открыт.\n\n"+menuText())
}

func menuText() string {
	return strings.TrimSpace(`Команды каталога:
/list — список товаров
/view id — полная карточка
/add — добавить товар пошагово
/edit id — показать команды редактирования
/set id field value — изменить поле
/photo id — заменить фото
/toggle id — включить/выключить товар
/cancel — отменить текущее действие

Поля для /set:
name, brand, price, volumes, description, top, middle, base, accords, psychotype, scores, active, image_url`)
}

func editHelp(id string) string {
	return fmt.Sprintf(`Примеры:
/set %[1]s name Miami Shake
/set %[1]s price 8393
/set %[1]s volumes 50:8393, 100:12990
/set %[1]s top клубника, бергамот
/set %[1]s psychotype aesthetic
/set %[1]s scores drive:20, focus:35, aesthetic:90, power:25
/set %[1]s active yes
/photo %[1]s`, id)
}

func (b *Bot) list(ctx context.Context, chatID int64) error {
	items, err := b.repo.ListFragrances(ctx)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return b.telegram.sendMessage(ctx, chatID, "Каталог пуст.")
	}

	var chunks []string
	var builder strings.Builder
	for index, item := range items {
		line := fmt.Sprintf("%d. %s — %s / %s / %s / active=%t\n", index+1, item.ID, item.Brand, item.Name, item.Psychotype, item.IsActive)
		if builder.Len()+len(line) > 3500 {
			chunks = append(chunks, builder.String())
			builder.Reset()
		}
		builder.WriteString(line)
	}
	if builder.Len() > 0 {
		chunks = append(chunks, builder.String())
	}
	for _, chunk := range chunks {
		if err := b.telegram.sendMessage(ctx, chatID, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) view(ctx context.Context, chatID int64, id string) error {
	if id == "" {
		return b.telegram.sendMessage(ctx, chatID, "Формат: /view fragrance-id")
	}
	item, err := b.repo.GetFragrance(ctx, id)
	if err != nil {
		return err
	}
	return b.telegram.sendMessage(ctx, chatID, formatFragrance(item))
}

func (b *Bot) setField(ctx context.Context, chatID int64, rest string) error {
	id, rest, ok := strings.Cut(rest, " ")
	if !ok {
		return b.telegram.sendMessage(ctx, chatID, "Формат: /set id field value")
	}
	field, value, ok := strings.Cut(strings.TrimSpace(rest), " ")
	if !ok {
		return b.telegram.sendMessage(ctx, chatID, "Формат: /set id field value")
	}

	item, err := b.repo.GetFragrance(ctx, id)
	if err != nil {
		return err
	}
	if err := applyField(&item, field, value); err != nil {
		return b.telegram.sendMessage(ctx, chatID, fieldError(field))
	}
	if err := b.repo.UpsertFragrance(ctx, item); err != nil {
		return err
	}
	return b.telegram.sendMessage(ctx, chatID, "Сохранено.\n\n"+formatFragrance(item))
}

func (b *Bot) toggle(ctx context.Context, chatID int64, id string) error {
	if id == "" {
		return b.telegram.sendMessage(ctx, chatID, "Формат: /toggle fragrance-id")
	}
	item, err := b.repo.GetFragrance(ctx, id)
	if err != nil {
		return err
	}
	item.IsActive = !item.IsActive
	if err := b.repo.UpsertFragrance(ctx, item); err != nil {
		return err
	}
	return b.telegram.sendMessage(ctx, chatID, fmt.Sprintf("Готово: %s active=%t", item.ID, item.IsActive))
}

func applyField(item *questionnaire.Fragrance, field string, value string) error {
	field = strings.ToLower(strings.TrimSpace(field))
	value = strings.TrimSpace(value)
	switch field {
	case "name", "название":
		item.Name = value
	case "brand", "бренд":
		item.Brand = value
	case "price", "цена":
		price, err := parsePrice(value)
		if err != nil {
			return err
		}
		item.Price = price
	case "volumes", "volume", "объем", "объемы":
		volumes, err := parseVolumes(value)
		if err != nil {
			return err
		}
		item.VolumeOptions = volumes
	case "description", "desc", "описание":
		item.Description = value
	case "top", "top_notes", "верх":
		item.TopNotes = splitList(value)
	case "middle", "middle_notes", "середина":
		item.MiddleNotes = splitList(value)
	case "base", "base_notes", "база":
		item.BaseNotes = splitList(value)
	case "accords", "main_accords", "аккорды":
		item.MainAccords = splitList(value)
	case "psychotype", "type", "тип":
		psychotype, err := normalizePsychotype(value)
		if err != nil {
			return err
		}
		item.Psychotype = psychotype
	case "scores", "score", "баллы":
		scores, err := parseScores(value)
		if err != nil {
			return err
		}
		item.PsychotypeScores = scores
	case "active", "is_active":
		active, err := parseActive(value)
		if err != nil {
			return err
		}
		item.IsActive = active
	case "image_url", "image", "photo", "фото":
		item.ImageURL = value
	default:
		return errInvalidValue
	}
	return validateFragrance(*item)
}

func fieldError(field string) string {
	return "Не смог сохранить поле " + field + ". Проверь формат. /edit id покажет примеры."
}

func validateFragrance(item questionnaire.Fragrance) error {
	if _, err := validateID(item.ID); err != nil {
		return err
	}
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Brand) == "" {
		return errInvalidValue
	}
	if _, err := parsePrice(item.Price); err != nil {
		return err
	}
	if _, err := normalizePsychotype(item.Psychotype); err != nil {
		return err
	}
	return nil
}

func formatFragrance(item questionnaire.Fragrance) string {
	return fmt.Sprintf(`%s
%s — %s
Цена: %s
Объемы: %s
Active: %t
Психотип: %s
Scores: drive=%d focus=%d aesthetic=%d power=%d
Фото: %s
Аккорды: %s
Верхние: %s
Средние: %s
Базовые: %s
Описание: %s`,
		item.ID,
		item.Brand,
		item.Name,
		item.Price,
		formatVolumes(item.VolumeOptions),
		item.IsActive,
		item.Psychotype,
		item.PsychotypeScores.Drive,
		item.PsychotypeScores.Focus,
		item.PsychotypeScores.Aesthetic,
		item.PsychotypeScores.Power,
		emptyDash(item.ImageURL),
		strings.Join(item.MainAccords, ", "),
		strings.Join(item.TopNotes, ", "),
		strings.Join(item.MiddleNotes, ", "),
		strings.Join(item.BaseNotes, ", "),
		emptyDash(item.Description),
	)
}

func formatVolumes(values []questionnaire.VolumeOption) string {
	if len(values) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, fmt.Sprintf("%d:%s", value.VolumeML, strconv.FormatFloat(value.Price, 'f', -1, 64)))
	}
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}

func emptyDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func (b *Bot) handlePhoto(ctx context.Context, chatID int64, s *session, photos []photoSize) error {
	item, err := b.repo.GetFragrance(ctx, s.targetID)
	if err != nil {
		return err
	}
	imageURL, err := b.saveTelegramPhoto(ctx, photos)
	if err != nil {
		return err
	}
	item.ImageURL = imageURL
	if err := b.repo.UpsertFragrance(ctx, item); err != nil {
		return err
	}
	s.action = ""
	s.targetID = ""
	return b.telegram.sendMessage(ctx, chatID, "Фото обновлено: "+imageURL)
}

func (b *Bot) saveTelegramPhoto(ctx context.Context, photos []photoSize) (string, error) {
	if len(photos) == 0 {
		return "", errInvalidValue
	}
	best := photos[0]
	for _, photo := range photos[1:] {
		if photo.FileSize > best.FileSize || photo.Width*photo.Height > best.Width*best.Height {
			best = photo
		}
	}

	file, err := b.telegram.getFile(ctx, best.FileID)
	if err != nil {
		return "", err
	}
	bytes, err := b.telegram.downloadFile(ctx, file.FilePath)
	if err != nil {
		return "", err
	}
	if len(bytes) > maxTelegramPhotoSize {
		return "", fmt.Errorf("photo is too large")
	}

	extension := strings.ToLower(filepath.Ext(file.FilePath))
	if extension == "" {
		extension = ".jpg"
	}
	fileName := uuid.NewString() + extension
	path := filepath.Join(b.cfg.UploadDir, fileName)
	if err := os.WriteFile(path, bytes, 0o644); err != nil {
		return "", err
	}
	return "/uploads/" + fileName, nil
}
