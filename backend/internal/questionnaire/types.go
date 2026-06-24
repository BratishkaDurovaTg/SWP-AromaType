package questionnaire

type Question struct {
	ID        string         `json:"id"`
	Text      string         `json:"text"`
	Type      string         `json:"type"`
	SortOrder int            `json:"sortOrder"`
	Options   []AnswerOption `json:"options"`
}

type AnswerOption struct {
	ID         string `json:"id"`
	QuestionID string `json:"questionId"`
	Text       string `json:"text"`
	Value      string `json:"value"`
	SortOrder  int    `json:"sortOrder"`
}

type RecommendationRequest struct {
	AnswerOptionIDs []string `json:"answerOptionIds"`
}

type RecommendationResponse struct {
	Profile    Profile              `json:"profile"`
	Items      []RecommendationItem `json:"items"`
	TotalItems int                  `json:"totalItems"`
}

type Profile struct {
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Tags            []string      `json:"tags"`
	ProfileBars     []ScoreMetric `json:"profileBars"`
	CharacterTraits []ScoreMetric `json:"characterTraits"`
	KeyNotes        []string      `json:"keyNotes"`
}

type ScoreMetric struct {
	Label   string `json:"label"`
	Percent int    `json:"percent"`
}

type RecommendationItem struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Brand        string   `json:"brand"`
	ImageURL     string   `json:"imageUrl"`
	Price        string   `json:"price"`
	MainAccords  []string `json:"mainAccords"`
	KeyNotes     []string `json:"keyNotes"`
	MatchPercent int      `json:"matchPercent"`
	Score        int      `json:"score"`
	Reason       string   `json:"reason"`
}

type Fragrance struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Brand         string         `json:"brand"`
	ImageURL      string         `json:"imageUrl"`
	Price         string         `json:"price"`
	VolumeOptions []VolumeOption `json:"volumeOptions"`
	Description   string         `json:"description"`
	TopNotes      []string       `json:"topNotes"`
	MiddleNotes   []string       `json:"middleNotes"`
	BaseNotes     []string       `json:"baseNotes"`
	MainAccords   []string       `json:"mainAccords"`
	IsActive      bool           `json:"isActive"`
}

type VolumeOption struct {
	VolumeML int     `json:"volumeMl"`
	Price    float64 `json:"price"`
}

type CreateFragranceRequest struct {
	Name          string         `json:"name"`
	Brand         string         `json:"brand"`
	ImageURL      string         `json:"imageUrl"`
	Price         float64        `json:"price"`
	VolumeOptions []VolumeOption `json:"volumeOptions"`
	Description   string         `json:"description"`
	TopNotes      []string       `json:"topNotes"`
	MiddleNotes   []string       `json:"middleNotes"`
	BaseNotes     []string       `json:"baseNotes"`
	MainAccords   []string       `json:"mainAccords"`
	TagIDs        []string       `json:"tagIds"`
	IsActive      *bool          `json:"isActive"`
}
