package questionnaire

import (
	"context"
	"encoding/json"
	"errors"

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
	f.image_url,
	f.price::TEXT,
	f.top_notes,
	f.middle_notes,
	f.base_notes,
	f.main_accords,
	ft.tag_id,
	t.name AS tag_name,
	ft.weight
FROM fragrances f
JOIN fragrance_tags ft ON ft.fragrance_id = f.id
JOIN tags t ON t.id = ft.tag_id
WHERE f.is_active = TRUE
ORDER BY f.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[FragranceTagRow])
}

func (r *Repository) GetFragranceByID(ctx context.Context, id string) (Fragrance, error) {
	var fragrance Fragrance
	var topNotes, middleNotes, baseNotes, mainAccords []byte
	var volumeOptions []byte

	err := r.db.QueryRow(ctx, `
SELECT
	id,
	name,
	brand,
	image_url,
	price::TEXT,
	volume_options,
	description,
	top_notes,
	middle_notes,
	base_notes,
	main_accords,
	is_active
FROM fragrances
WHERE id = $1 AND is_active = TRUE
`, id).Scan(
		&fragrance.ID,
		&fragrance.Name,
		&fragrance.Brand,
		&fragrance.ImageURL,
		&fragrance.Price,
		&volumeOptions,
		&fragrance.Description,
		&topNotes,
		&middleNotes,
		&baseNotes,
		&mainAccords,
		&fragrance.IsActive,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Fragrance{}, ErrFragranceNotFound
	}
	if err != nil {
		return Fragrance{}, err
	}

	if err := decodeStringArray(topNotes, &fragrance.TopNotes); err != nil {
		return Fragrance{}, err
	}
	if err := decodeStringArray(middleNotes, &fragrance.MiddleNotes); err != nil {
		return Fragrance{}, err
	}
	if err := decodeStringArray(baseNotes, &fragrance.BaseNotes); err != nil {
		return Fragrance{}, err
	}
	if err := decodeStringArray(mainAccords, &fragrance.MainAccords); err != nil {
		return Fragrance{}, err
	}
	if err := decodeVolumeOptions(volumeOptions, &fragrance.VolumeOptions); err != nil {
		return Fragrance{}, err
	}

	return fragrance, nil
}

func (r *Repository) CreateFragrance(ctx context.Context, fragrance Fragrance, tagIDs []string) (Fragrance, error) {
	topNotes, err := encodeStringArray(fragrance.TopNotes)
	if err != nil {
		return Fragrance{}, err
	}
	middleNotes, err := encodeStringArray(fragrance.MiddleNotes)
	if err != nil {
		return Fragrance{}, err
	}
	baseNotes, err := encodeStringArray(fragrance.BaseNotes)
	if err != nil {
		return Fragrance{}, err
	}
	mainAccords, err := encodeStringArray(fragrance.MainAccords)
	if err != nil {
		return Fragrance{}, err
	}
	volumeOptions, err := encodeVolumeOptions(fragrance.VolumeOptions)
	if err != nil {
		return Fragrance{}, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Fragrance{}, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
INSERT INTO fragrances (
	id, name, brand, image_url, price, volume_options, description,
	top_notes, middle_notes, base_notes, main_accords, is_active
) VALUES (
	$1, $2, $3, $4, $5, $6::JSONB, $7,
	$8::JSONB, $9::JSONB, $10::JSONB, $11::JSONB, $12
)
RETURNING price::TEXT
`,
		fragrance.ID,
		fragrance.Name,
		fragrance.Brand,
		fragrance.ImageURL,
		fragrance.Price,
		volumeOptions,
		fragrance.Description,
		topNotes,
		middleNotes,
		baseNotes,
		mainAccords,
		fragrance.IsActive,
	).Scan(&fragrance.Price)
	if err != nil {
		return Fragrance{}, err
	}

	for _, tagID := range tagIDs {
		if _, err := tx.Exec(ctx, `
INSERT INTO fragrance_tags (fragrance_id, tag_id, weight)
VALUES ($1, $2, 1)
ON CONFLICT (fragrance_id, tag_id) DO UPDATE SET weight = EXCLUDED.weight
`, fragrance.ID, tagID); err != nil {
			return Fragrance{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Fragrance{}, err
	}

	return fragrance, nil
}

func encodeStringArray(values []string) (string, error) {
	if values == nil {
		values = []string{}
	}
	bytes, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func decodeStringArray(raw []byte, destination *[]string) error {
	if len(raw) == 0 {
		*destination = []string{}
		return nil
	}
	return json.Unmarshal(raw, destination)
}

func encodeVolumeOptions(values []VolumeOption) (string, error) {
	if values == nil {
		values = []VolumeOption{}
	}
	bytes, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func decodeVolumeOptions(raw []byte, destination *[]VolumeOption) error {
	if len(raw) == 0 {
		*destination = []VolumeOption{}
		return nil
	}
	return json.Unmarshal(raw, destination)
}

type FragranceTagRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Brand       string `db:"brand"`
	ImageURL    string `db:"image_url"`
	Price       string `db:"price"`
	TopNotes    []byte `db:"top_notes"`
	MiddleNotes []byte `db:"middle_notes"`
	BaseNotes   []byte `db:"base_notes"`
	MainAccords []byte `db:"main_accords"`
	TagID       string `db:"tag_id"`
	TagName     string `db:"tag_name"`
	Weight      int    `db:"weight"`
}
