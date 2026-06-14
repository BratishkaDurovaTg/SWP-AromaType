# AromaType Database Schema

## users

- id
- email
- password_hash
- role
- created_at
- updated_at

## fragrances

- id
- name
- brand
- gender
- image_url
- volume
- price
- stock_status
- description
- top_notes
- middle_notes
- base_notes
- main_accords
- longevity
- projection
- visibility
- versatility
- seasons
- time_of_day
- situations
- matching_profiles
- why_recommended
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

## questionnaire_sessions

- id
- user_id
- status
- created_at
- completed_at

## questionnaire_answers

- id
- session_id
- question_id
- answer_option_id
- created_at

## aroma_profiles

- id
- session_id
- name
- description
- created_at

## recommendations

- id
- session_id
- algorithm_version
- created_at

## recommendation_items

- id
- recommendation_id
- fragrance_id
- score
- reason

## orders

- id
- user_id
- fragrance_id
- status
- contact_name
- contact_phone
- contact_telegram
- comment
- created_at
- updated_at

## Recommendation Logic

- Answer options are connected to tags through `answer_option_tags`.
- Fragrances are connected to tags through `fragrance_tags`.
- The backend builds a user tag profile from selected answer options.
- Each active fragrance receives a score based on matching tag weights.
- The API returns the highest-scoring fragrances with short explanations.
