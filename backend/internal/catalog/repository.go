package catalog

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("fragrance not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListFragrances(ctx context.Context) ([]questionnaire.Fragrance, error) {
	rows, err := r.db.Query(ctx, `
SELECT id, name, brand, image_url, price::TEXT, volume_options, description,
	top_notes, middle_notes, base_notes, main_accords, psychotype, psychotype_scores, is_active
FROM fragrances
ORDER BY name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]questionnaire.Fragrance, 0)
	for rows.Next() {
		item, err := scanFragrance(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) GetFragrance(ctx context.Context, id string) (questionnaire.Fragrance, error) {
	row := r.db.QueryRow(ctx, `
SELECT id, name, brand, image_url, price::TEXT, volume_options, description,
	top_notes, middle_notes, base_notes, main_accords, psychotype, psychotype_scores, is_active
FROM fragrances
WHERE id = $1
`, id)

	item, err := scanFragrance(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return questionnaire.Fragrance{}, ErrNotFound
	}
	return item, err
}

func (r *Repository) UpsertFragrance(ctx context.Context, item questionnaire.Fragrance) error {
	volumeOptions, err := marshalJSON(item.VolumeOptions)
	if err != nil {
		return err
	}
	topNotes, err := marshalJSON(item.TopNotes)
	if err != nil {
		return err
	}
	middleNotes, err := marshalJSON(item.MiddleNotes)
	if err != nil {
		return err
	}
	baseNotes, err := marshalJSON(item.BaseNotes)
	if err != nil {
		return err
	}
	mainAccords, err := marshalJSON(item.MainAccords)
	if err != nil {
		return err
	}
	psychotypeScores, err := marshalJSON(item.PsychotypeScores)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
INSERT INTO fragrances (
	id, name, brand, image_url, price, volume_options, description,
	top_notes, middle_notes, base_notes, main_accords, psychotype, psychotype_scores, is_active
) VALUES (
	$1, $2, $3, $4, $5::NUMERIC, $6::JSONB, $7,
	$8::JSONB, $9::JSONB, $10::JSONB, $11::JSONB, $12, $13::JSONB, $14
)
ON CONFLICT (id) DO UPDATE SET
	name = EXCLUDED.name,
	brand = EXCLUDED.brand,
	image_url = EXCLUDED.image_url,
	price = EXCLUDED.price,
	volume_options = EXCLUDED.volume_options,
	description = EXCLUDED.description,
	top_notes = EXCLUDED.top_notes,
	middle_notes = EXCLUDED.middle_notes,
	base_notes = EXCLUDED.base_notes,
	main_accords = EXCLUDED.main_accords,
	psychotype = EXCLUDED.psychotype,
	psychotype_scores = EXCLUDED.psychotype_scores,
	is_active = EXCLUDED.is_active,
	updated_at = now()
`, item.ID, item.Name, item.Brand, item.ImageURL, item.Price, volumeOptions, item.Description,
		topNotes, middleNotes, baseNotes, mainAccords, item.Psychotype, psychotypeScores, item.IsActive)
	return err
}

type fragranceScanner interface {
	Scan(dest ...any) error
}

func scanFragrance(row fragranceScanner) (questionnaire.Fragrance, error) {
	var item questionnaire.Fragrance
	var volumeOptions, topNotes, middleNotes, baseNotes, mainAccords, psychotypeScores []byte

	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Brand,
		&item.ImageURL,
		&item.Price,
		&volumeOptions,
		&item.Description,
		&topNotes,
		&middleNotes,
		&baseNotes,
		&mainAccords,
		&item.Psychotype,
		&psychotypeScores,
		&item.IsActive,
	)
	if err != nil {
		return questionnaire.Fragrance{}, err
	}

	if err := unmarshalJSON(volumeOptions, &item.VolumeOptions); err != nil {
		return questionnaire.Fragrance{}, err
	}
	if err := unmarshalJSON(topNotes, &item.TopNotes); err != nil {
		return questionnaire.Fragrance{}, err
	}
	if err := unmarshalJSON(middleNotes, &item.MiddleNotes); err != nil {
		return questionnaire.Fragrance{}, err
	}
	if err := unmarshalJSON(baseNotes, &item.BaseNotes); err != nil {
		return questionnaire.Fragrance{}, err
	}
	if err := unmarshalJSON(mainAccords, &item.MainAccords); err != nil {
		return questionnaire.Fragrance{}, err
	}
	if err := unmarshalJSON(psychotypeScores, &item.PsychotypeScores); err != nil {
		return questionnaire.Fragrance{}, err
	}

	return item, nil
}

func marshalJSON(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func unmarshalJSON(raw []byte, destination any) error {
	if len(raw) == 0 {
		raw = []byte("null")
	}
	return json.Unmarshal(raw, destination)
}
