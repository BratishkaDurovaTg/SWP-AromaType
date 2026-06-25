package questionnaire

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrNoAnswers         = errors.New("at least one answer option is required")
	ErrFragranceNotFound = errors.New("fragrance not found")
	ErrInvalidFragrance  = errors.New("invalid fragrance")
)

type Service struct {
	repo fragranceRepository
}

type fragranceRepository interface {
	GetQuestions(ctx context.Context) ([]Question, error)
	GetOptionTagWeights(ctx context.Context, optionID string) (map[string]int, error)
	GetTagNames(ctx context.Context, tagIDs []string) (map[string]string, error)
	GetActiveFragranceTags(ctx context.Context) ([]FragranceTagRow, error)
	GetFragranceByID(ctx context.Context, id string) (Fragrance, error)
	CreateFragrance(ctx context.Context, fragrance Fragrance, tagIDs []string) (Fragrance, error)
}

func NewService(repo fragranceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetQuestions(ctx context.Context) ([]Question, error) {
	return s.repo.GetQuestions(ctx)
}

func (s *Service) Recommend(ctx context.Context, answerOptionIDs []string) (RecommendationResponse, error) {
	answerOptionIDs = uniqueNonEmpty(answerOptionIDs)
	if len(answerOptionIDs) == 0 {
		return RecommendationResponse{}, ErrNoAnswers
	}

	userTagWeights := make(map[string]int)
	for _, optionID := range answerOptionIDs {
		weights, err := s.repo.GetOptionTagWeights(ctx, optionID)
		if err != nil {
			return RecommendationResponse{}, err
		}
		for tagID, weight := range weights {
			userTagWeights[tagID] += weight
		}
	}

	tagIDs := mapKeys(userTagWeights)
	tagNames, err := s.repo.GetTagNames(ctx, tagIDs)
	if err != nil {
		return RecommendationResponse{}, err
	}

	fragranceRows, err := s.repo.GetActiveFragranceTags(ctx)
	if err != nil {
		return RecommendationResponse{}, err
	}

	fragrances := make(map[string]*fragranceScore)
	for _, row := range fragranceRows {
		item, ok := fragrances[row.ID]
		if !ok {
			item = &fragranceScore{
				RecommendationItem: RecommendationItem{
					ID:          row.ID,
					Name:        row.Name,
					Brand:       row.Brand,
					ImageURL:    row.ImageURL,
					Price:       row.Price,
					MainAccords: decodeStrings(row.MainAccords),
					KeyNotes:    firstStrings(append(append(decodeStrings(row.TopNotes), decodeStrings(row.MiddleNotes)...), decodeStrings(row.BaseNotes)...), 3),
				},
				matchedTags: make(map[string]int),
			}
			fragrances[row.ID] = item
		}

		if userWeight, ok := userTagWeights[row.TagID]; ok {
			item.Score += userWeight * row.Weight
			item.matchedTags[row.TagName] += userWeight * row.Weight
		}
	}

	items := make([]RecommendationItem, 0, len(fragrances))
	for _, fragrance := range fragrances {
		if fragrance.Score <= 0 {
			continue
		}
		fragrance.Reason = buildReason(fragrance.matchedTags)
		items = append(items, fragrance.RecommendationItem)
	}

	topScore := 0
	for _, item := range items {
		if item.Score > topScore {
			topScore = item.Score
		}
	}
	for index := range items {
		items[index].MatchPercent = scorePercent(items[index].Score, topScore)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].Name < items[j].Name
		}
		return items[i].Score > items[j].Score
	})

	if len(items) > 5 {
		items = items[:5]
	}

	return RecommendationResponse{
		Profile:    buildProfile(userTagWeights, tagNames),
		Items:      items,
		TotalItems: len(items),
	}, nil
}

func (s *Service) GetFragrance(ctx context.Context, id string) (Fragrance, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Fragrance{}, ErrFragranceNotFound
	}
	return s.repo.GetFragranceByID(ctx, id)
}

func (s *Service) CreateFragrance(ctx context.Context, payload CreateFragranceRequest) (Fragrance, error) {
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Brand = strings.TrimSpace(payload.Brand)
	payload.ImageURL = strings.TrimSpace(payload.ImageURL)
	payload.Description = strings.TrimSpace(payload.Description)
	payload.VolumeOptions = cleanVolumeOptions(payload.VolumeOptions)

	if payload.Name == "" ||
		payload.Brand == "" ||
		payload.Price < 0 {
		return Fragrance{}, ErrInvalidFragrance
	}

	isActive := true
	if payload.IsActive != nil {
		isActive = *payload.IsActive
	}

	fragrance := Fragrance{
		ID:            uuid.NewString(),
		Name:          payload.Name,
		Brand:         payload.Brand,
		ImageURL:      payload.ImageURL,
		Price:         strconv.FormatFloat(payload.Price, 'f', -1, 64),
		VolumeOptions: payload.VolumeOptions,
		Description:   payload.Description,
		TopNotes:      cleanStrings(payload.TopNotes),
		MiddleNotes:   cleanStrings(payload.MiddleNotes),
		BaseNotes:     cleanStrings(payload.BaseNotes),
		MainAccords:   cleanStrings(payload.MainAccords),
		IsActive:      isActive,
	}

	return s.repo.CreateFragrance(ctx, fragrance, uniqueNonEmpty(payload.TagIDs))
}

type fragranceScore struct {
	RecommendationItem
	matchedTags map[string]int
}

func buildProfile(tagWeights map[string]int, tagNames map[string]string) Profile {
	type weightedTag struct {
		id     string
		name   string
		weight int
	}

	tags := make([]weightedTag, 0, len(tagWeights))
	for id, weight := range tagWeights {
		tags = append(tags, weightedTag{id: id, name: tagNames[id], weight: weight})
	}

	sort.Slice(tags, func(i, j int) bool {
		if tags[i].weight == tags[j].weight {
			return tags[i].name < tags[j].name
		}
		return tags[i].weight > tags[j].weight
	})

	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag.name != "" {
			names = append(names, tag.name)
		}
	}

	profileName := "Сбалансированный профиль"
	description := "Ваши ответы распределились между несколькими типами. Вам подходят ароматы с понятным характером, балансом комфорта и выразительности."

	profileScores := map[string]int{
		"drive":     tagWeights["psych_drive"],
		"focus":     tagWeights["psych_focus"],
		"aesthetic": tagWeights["psych_aesthetic"],
		"power":     tagWeights["psych_power"],
	}

	dominantProfile := "balanced"
	dominantScore := 0
	for profile, score := range profileScores {
		if score > dominantScore {
			dominantProfile = profile
			dominantScore = score
		}
	}

	switch dominantProfile {
	case "drive":
		profileName = "Драйв и экстраверсия"
		description = "Вам близки ароматы с движением, свежим импульсом и открытой харизмой. Они помогают быстро включиться, звучать живо и не терять энергию."
	case "focus":
		profileName = "Интеллект и фокус"
		description = "Вам подходят собранные, глубокие и немного загадочные ароматы. Они создают ощущение концентрации, дистанции и внутренней опоры."
	case "aesthetic":
		profileName = "Эстетика и гедонизм"
		description = "Вам важны тактильность, ухоженность и удовольствие от деталей. Лучше всего раскрываются мягкие, чистые и красиво собранные композиции."
	case "power":
		profileName = "Власть и доминанта"
		description = "Вам подходят плотные, уверенные и статусные ароматы. Они дают эффект собранности, силы и заметного личного присутствия."
	}

	if len(names) > 5 {
		names = names[:5]
	}

	totalPsychScore := profileScores["drive"] + profileScores["focus"] + profileScores["aesthetic"] + profileScores["power"]

	return Profile{
		Name:        profileName,
		Description: description,
		Tags:        names,
		ProfileBars: []ScoreMetric{
			{Label: "Драйв", Percent: profilePercent(profileScores["drive"], totalPsychScore)},
			{Label: "Фокус", Percent: profilePercent(profileScores["focus"], totalPsychScore)},
			{Label: "Эстетика", Percent: profilePercent(profileScores["aesthetic"], totalPsychScore)},
			{Label: "Доминанта", Percent: profilePercent(profileScores["power"], totalPsychScore)},
		},
		CharacterTraits: []ScoreMetric{
			{Label: "Социальный импульс", Percent: clampPercent(profileScores["drive"]*9 + profileScores["power"]*4)},
			{Label: "Внутренняя собранность", Percent: clampPercent(profileScores["focus"]*9 + profileScores["power"]*5)},
			{Label: "Чувство вкуса", Percent: clampPercent(profileScores["aesthetic"]*9 + profileScores["focus"]*3)},
		},
		KeyNotes: firstStrings(names, 5),
	}
}

func buildReason(matchedTags map[string]int) string {
	tags := make([]string, 0, len(matchedTags))
	for tag := range matchedTags {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	if len(tags) > 3 {
		tags = tags[:3]
	}
	if len(tags) == 0 {
		return "Аромат попал в подборку по общему совпадению с вашими ответами."
	}

	return "Совпадает с вашим профилем по тегам: " + strings.Join(tags, ", ") + "."
}

func uniqueNonEmpty(values []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func mapKeys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func cleanStrings(values []string) []string {
	return uniqueNonEmpty(values)
}

func cleanVolumeOptions(values []VolumeOption) []VolumeOption {
	result := make([]VolumeOption, 0, len(values))
	seen := make(map[int]struct{})
	for _, value := range values {
		if value.VolumeML <= 0 || value.Price < 0 {
			continue
		}
		if _, ok := seen[value.VolumeML]; ok {
			continue
		}
		seen[value.VolumeML] = struct{}{}
		result = append(result, value)
	}
	return result
}

func decodeStrings(raw []byte) []string {
	var values []string
	if err := decodeStringArray(raw, &values); err != nil {
		return []string{}
	}
	return values
}

func firstStrings(values []string, limit int) []string {
	values = cleanStrings(values)
	if len(values) > limit {
		return values[:limit]
	}
	return values
}

func scorePercent(score, topScore int) int {
	if score <= 0 || topScore <= 0 {
		return 0
	}
	percent := 70 + (score * 30 / topScore)
	if percent > 99 {
		return 99
	}
	return percent
}

func profilePercent(score, totalScore int) int {
	if score <= 0 || totalScore <= 0 {
		return 0
	}
	percent := score * 100 / totalScore
	if percent > 100 {
		return 100
	}
	return percent
}

func clampPercent(value int) int {
	if value < 12 {
		return 12
	}
	if value > 96 {
		return 96
	}
	return value
}
