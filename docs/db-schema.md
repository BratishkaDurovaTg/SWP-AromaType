# AromaType Database Schema

## fragrances

- id
- name
- brand
- image_url
- price
- volume_options
- description
- top_notes
- middle_notes
- base_notes
- main_accords
- psychotype
- psychotype_scores
- is_active
- created_at
- updated_at

## tags

- id
- name
- type
- created_at

## fragrance_tags

- fragrance_id
- tag_id
- weight

## questions

- id
- text
- type
- sort_order
- is_active

## answer_options

- id
- question_id
- text
- value
- sort_order

## answer_option_tags

- answer_option_id
- tag_id
- weight

## Recommendation Logic

- The questionnaire contains 8 single-choice questions.
- Each answer option maps to one of four psychotype tags through `answer_option_tags`:
  `psych_drive`, `psych_focus`, `psych_aesthetic`, or `psych_power`.
- In the current questionnaire, all `A` answers map to `psych_drive`,
  all `B` answers map to `psych_focus`, all `C` answers map to
  `psych_aesthetic`, and all `D` answers map to `psych_power`.
- Each fragrance stores one primary psychotype in `psychotype`:
  `drive`, `focus`, `aesthetic`, `power`, or `balanced`.
- Each fragrance also stores a flexible 0-100 score vector in `psychotype_scores`:
  `drive`, `focus`, `aesthetic`, and `power`.
- The backend builds a user psychotype profile from selected answer options.
- Each active fragrance receives a score by comparing the user profile with the
  fragrance score vector.
- The API returns only the top 5 highest-scoring fragrances with short
  explanations.
