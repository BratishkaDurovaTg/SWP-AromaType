package catalogbot

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
)

var (
	errInvalidValue = errors.New("invalid value")
	idPattern       = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{1,80}$`)
	scoresPattern   = regexp.MustCompile(`^\d{1,3},\d{1,3},\d{1,3},\d{1,3}$`)
)

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		key := strings.ToLower(part)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, part)
	}
	return result
}

func parsePrice(value string) (string, error) {
	value = strings.TrimSpace(strings.ReplaceAll(value, ",", "."))
	if value == "" {
		return "", errInvalidValue
	}
	price, err := strconv.ParseFloat(value, 64)
	if err != nil || price < 0 {
		return "", errInvalidValue
	}
	return strconv.FormatFloat(price, 'f', -1, 64), nil
}

func parseVolumes(value string) ([]questionnaire.VolumeOption, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" || strings.EqualFold(value, "skip") {
		return []questionnaire.VolumeOption{}, nil
	}

	parts := splitList(value)
	result := make([]questionnaire.VolumeOption, 0, len(parts))
	seen := make(map[int]struct{}, len(parts))
	for _, part := range parts {
		volume, price, ok := strings.Cut(part, ":")
		if !ok {
			return nil, fmt.Errorf("%w: volume must use 3:8393 format", errInvalidValue)
		}
		volumeML, err := strconv.Atoi(strings.TrimSpace(volume))
		if err != nil || !isSupportedVolume(volumeML) {
			return nil, errInvalidValue
		}
		if _, ok := seen[volumeML]; ok {
			return nil, errInvalidValue
		}
		seen[volumeML] = struct{}{}
		priceValue, err := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(price), ",", "."), 64)
		if err != nil || priceValue < 0 {
			return nil, errInvalidValue
		}
		result = append(result, questionnaire.VolumeOption{VolumeML: volumeML, Price: priceValue})
	}
	return result, nil
}

func parseScores(value string) (questionnaire.PsychotypeScores, error) {
	var scores questionnaire.PsychotypeScores
	value = strings.TrimSpace(value)
	if value == "" {
		return scores, errInvalidValue
	}

	if !strings.Contains(value, ":") {
		return parseOrderedScores(value)
	}

	for _, part := range splitList(value) {
		key, rawScore, ok := strings.Cut(part, ":")
		if !ok {
			return scores, fmt.Errorf("%w: score must use drive:50 format", errInvalidValue)
		}
		score, err := strconv.Atoi(strings.TrimSpace(rawScore))
		if err != nil || score < 0 || score > 100 {
			return scores, errInvalidValue
		}

		switch strings.ToLower(strings.TrimSpace(key)) {
		case "drive", "a":
			scores.Drive = score
		case "focus", "b":
			scores.Focus = score
		case "aesthetic", "c":
			scores.Aesthetic = score
		case "power", "d":
			scores.Power = score
		default:
			return scores, fmt.Errorf("%w: unknown score key %q", errInvalidValue, key)
		}
	}
	if scores.Drive+scores.Focus+scores.Aesthetic+scores.Power != 100 {
		return scores, fmt.Errorf("%w: scores sum must be 100", errInvalidValue)
	}
	return scores, nil
}

func parseOrderedScores(value string) (questionnaire.PsychotypeScores, error) {
	var scores questionnaire.PsychotypeScores
	if !scoresPattern.MatchString(value) {
		return scores, fmt.Errorf("%w: scores must use 20,20,40,20 format", errInvalidValue)
	}

	parts := strings.Split(value, ",")
	values := make([]int, len(parts))
	for index, part := range parts {
		score, err := strconv.Atoi(part)
		if err != nil || score < 0 || score > 100 {
			return scores, errInvalidValue
		}
		values[index] = score
	}

	scores = questionnaire.PsychotypeScores{
		Drive:     values[0],
		Focus:     values[1],
		Aesthetic: values[2],
		Power:     values[3],
	}
	if scores.Drive+scores.Focus+scores.Aesthetic+scores.Power != 100 {
		return scores, fmt.Errorf("%w: scores sum must be 100", errInvalidValue)
	}
	return scores, nil
}

func parseActive(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes", "y", "true", "1", "on", "active", "да", "д", "вкл":
		return true, nil
	case "no", "n", "false", "0", "off", "inactive", "нет", "н", "выкл":
		return false, nil
	default:
		return false, errInvalidValue
	}
}

func normalizePsychotype(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "drive", "focus", "aesthetic", "power", "balanced":
		return value, nil
	case "a":
		return "drive", nil
	case "b":
		return "focus", nil
	case "c":
		return "aesthetic", nil
	case "d":
		return "power", nil
	default:
		return "", errInvalidValue
	}
}

func parseGender(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "male", "m", "man", "men", "м", "муж", "мужской":
		return questionnaire.GenderMale, nil
	case "female", "f", "woman", "women", "ж", "жен", "женский":
		return questionnaire.GenderFemale, nil
	case "unisex", "u", "унисекс", "универсальный":
		return questionnaire.GenderUnisex, nil
	default:
		return "", errInvalidValue
	}
}

func isSupportedVolume(volumeML int) bool {
	switch volumeML {
	case 3, 5, 10:
		return true
	default:
		return false
	}
}

func setVolumePrice(values []questionnaire.VolumeOption, volumeML int, price float64) []questionnaire.VolumeOption {
	next := make([]questionnaire.VolumeOption, 0, len(values)+1)
	updated := false
	for _, value := range values {
		if value.VolumeML == volumeML {
			next = append(next, questionnaire.VolumeOption{VolumeML: volumeML, Price: price})
			updated = true
			continue
		}
		next = append(next, value)
	}
	if !updated {
		next = append(next, questionnaire.VolumeOption{VolumeML: volumeML, Price: price})
	}
	return next
}

func volumePrice(values []questionnaire.VolumeOption, volumeML int) (float64, bool) {
	for _, value := range values {
		if value.VolumeML == volumeML {
			return value.Price, true
		}
	}
	return 0, false
}

func validateID(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if !idPattern.MatchString(value) {
		return "", errInvalidValue
	}
	return value, nil
}
