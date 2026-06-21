package questionnaire

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetQuestions(ctx context.Context) ([]Question, error) {
	rows, err := r.db.Query(ctx, `
SELECT q.id, q.text, q.type, q.sort_order, o.id, o.text, o.value, o.sort_order
FROM questions q
LEFT JOIN answer_options o ON o.question_id = q.id
WHERE q.is_active = TRUE
ORDER BY q.sort_order, o.sort_order
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionsByID := make(map[string]*Question)
	orderedIDs := make([]string, 0)

	for rows.Next() {
		var question Question
		var option AnswerOption
		if err := rows.Scan(
			&question.ID,
			&question.Text,
			&question.Type,
			&question.SortOrder,
			&option.ID,
			&option.Text,
			&option.Value,
			&option.SortOrder,
		); err != nil {
			return nil, err
		}

		existing, ok := questionsByID[question.ID]
		if !ok {
			question.Options = []AnswerOption{}
			questionsByID[question.ID] = &question
			orderedIDs = append(orderedIDs, question.ID)
			existing = &question
		}

		if option.ID != "" {
			option.QuestionID = question.ID
			existing.Options = append(existing.Options, option)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	questions := make([]Question, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		questions = append(questions, *questionsByID[id])
	}

	return questions, nil
}

func (r *Repository) GetOptionTagWeights(ctx context.Context, optionID string) (map[string]int, error) {
	rows, err := r.db.Query(ctx, `
SELECT tag_id, weight
FROM answer_option_tags
WHERE answer_option_id = $1
`, optionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weights := make(map[string]int)
	for rows.Next() {
		var tagID string
		var weight int
		if err := rows.Scan(&tagID, &weight); err != nil {
			return nil, err
		}
		weights[tagID] += weight
	}

	return weights, rows.Err()
}

func (r *Repository) GetTagNames(ctx context.Context, tagIDs []string) (map[string]string, error) {
	rows, err := r.db.Query(ctx, `
SELECT id, name
FROM tags
WHERE id = ANY($1)
`, tagIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := make(map[string]string)
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		names[id] = name
	}

	return names, rows.Err()
}

func (r *Repository) GetActiveFragranceTags(ctx context.Context) ([]FragranceTagRow, error) {
	rows, err := r.db.Query(ctx, `
SELECT
	f.id,
	f.name,
	f.brand,
	f.gender,
	f.image_url,
	f.price::TEXT,
	f.stock_status,
	ft.tag_id,
	t.name AS tag_name,
	ft.weight
FROM fragrances f
JOIN fragrance_tags ft ON ft.fragrance_id = f.id
JOIN tags t ON t.id = ft.tag_id
WHERE f.is_active = TRUE AND f.stock_status = 'in_stock'
ORDER BY f.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[FragranceTagRow])
}

type FragranceTagRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Brand       string `db:"brand"`
	Gender      string `db:"gender"`
	ImageURL    string `db:"image_url"`
	Price       string `db:"price"`
	StockStatus string `db:"stock_status"`
	TagID       string `db:"tag_id"`
	TagName     string `db:"tag_name"`
	Weight      int    `db:"weight"`
}
