package questionnaire

import (
	"context"
	"errors"
	"sort"
	"strings"
)

var ErrNoAnswers = errors.New("at least one answer option is required")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
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
					Gender:      row.Gender,
					ImageURL:    row.ImageURL,
					Price:       row.Price,
					StockStatus: row.StockStatus,
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
		Profile: buildProfile(userTagWeights, tagNames),
		Items:   items,
	}, nil
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
	description := "Вам подходят ароматы с понятным характером, умеренной заметностью и несколькими сценариями использования."

	profileScores := map[string]int{
		"mystery": tagWeights["mystery"] + tagWeights["deep"] + tagWeights["night"],
		"bright":  tagWeights["bright"] + tagWeights["energy"] + tagWeights["party"] + tagWeights["noticeable"] + tagWeights["trail"],
		"romance": tagWeights["romantic"] + tagWeights["soft"] + tagWeights["warm"] + tagWeights["cozy"] + tagWeights["date"],
		"calm":    tagWeights["calm"] + tagWeights["clean"] + tagWeights["fresh"] + tagWeights["office"] + tagWeights["daily"] + tagWeights["reliable"] + tagWeights["light"],
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
	case "mystery":
		profileName = "Таинственный акцент"
		description = "Вам ближе глубокие, необычные и интригующие ароматы, которые создают запоминающийся образ."
	case "bright":
		profileName = "Яркая энергия"
		description = "Вам подходят выразительные, динамичные ароматы, которые заметны и поддерживают активный образ."
	case "romance":
		profileName = "Мягкая романтика"
		description = "Вам подходят мягкие, теплые и притягательные ароматы для близкого общения и спокойного впечатления."
	case "calm":
		profileName = "Спокойный минималист"
		description = "Вам подходят чистые, спокойные и ненавязчивые ароматы для повседневности, офиса и учебы."
	}

	if len(names) > 5 {
		names = names[:5]
	}

	return Profile{Name: profileName, Description: description, Tags: names}
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

func hasTag(tagWeights map[string]int, tagID string) bool {
	return tagWeights[tagID] > 0
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
