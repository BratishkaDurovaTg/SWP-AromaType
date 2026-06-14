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
	Profile Profile              `json:"profile"`
	Items   []RecommendationItem `json:"items"`
}

type Profile struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type RecommendationItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Gender      string `json:"gender"`
	ImageURL    string `json:"imageUrl"`
	Price       string `json:"price"`
	StockStatus string `json:"stockStatus"`
	Score       int    `json:"score"`
	Reason      string `json:"reason"`
}
