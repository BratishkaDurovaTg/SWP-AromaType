# Dynamic View — AromaType

The dynamic view describes runtime interactions, sequence flows, and data paths
for the main user scenarios.

## 1. User Completes Questionnaire and Receives Recommendations

This is the core workflow — the user answers 8 psychotype questions and receives
up to 5 fragrance recommendations.

```
User              Mini App              API Server           PostgreSQL
 │                   │                     │                    │
 │  1. Open Mini App │                     │                    │
 │──────────────────>│                     │                    │
 │                   │  2. GET /api/questions                    │
 │                   │────────────────────>│                    │
 │                   │                     │  3. SELECT FROM    │
 │                   │                     │     questions      │
 │                   │                     │     JOIN           │
 │                   │                     │     answer_options │
 │                   │                     │───────────────────>│
 │                   │                     │  4. rows           │
 │                   │                     │<───────────────────│
 │                   │  5. Questions JSON  │                    │
 │                   │<────────────────────│                    │
 │  6. Render quiz   │                     │                    │
 │<──────────────────│                     │                    │
 │                   │                     │                    │
 │  7. Select answer │                     │                    │
 │     8 times       │                     │                    │
 │═══════════════════>│                     │                    │
 │                   │                     │                    │
 │  8. Submit        │                     │                    │
 │──────────────────>│                     │                    │
 │                   │  9. POST /api/                              │
 │                   │     recommendations                        │
 │                   │     body: {answerOptionIds: [...]}          │
 │                   │────────────────────>│                    │
 │                   │                     │                    │
 │                   │                     │  10. For each answer: │
 │                   │                     │    SELECT option    │
 │                   │                     │    tag weights      │
 │                   │                     │───────────────────>│
 │                   │                     │  11. weights        │
 │                   │                     │<───────────────────│
 │                   │                     │                    │
 │                   │                     │  12. SELECT active  │
 │                   │                     │     fragrance tags  │
 │                   │                     │───────────────────>│
 │                   │                     │  13. rows           │
 │                   │                     │<───────────────────│
 │                   │                     │                    │
 │                   │                     │  14. Compute scores │
 │                   │                     │      ┌────────────┐ │
 │                   │                     │      │ Dot product │ │
 │                   │                     │      │ user vs     │ │
 │                   │                     │      │ fragrance   │ │
 │                   │                     │      │ psychotype  │ │
 │                   │                     │      │ + tag match │ │
 │                   │                     │      │ score       │ │
 │                   │                     │      └────────────┘ │
 │                   │                     │                    │
 │                   │                     │  15. Sort, limit 5,│
 │                   │                     │      normalize %   │
 │                   │                     │                    │
 │                   │  16. RecommendationsJSON                │
 │                   │<────────────────────│                    │
 │                   │                     │                    │
 │  17. Show profile,│                     │                    │
 │      results list │                     │                    │
 │<──────────────────│                     │                    │
```

### Recommendation Algorithm

```
For each user answer:
  answer ──> answer_option.tags ──> accumulate tag weights

User psychotype scores:
  drive    = Σ weight where tag = psych_drive
  focus    = Σ weight where tag = psych_focus
  aesthetic= Σ weight where tag = psych_aesthetic
  power    = Σ weight where tag = psych_power

For each active fragrance:
  psychotype_match = dot(user_scores, fragrance.psychotype_scores)
  tag_match_score  = Σ user_weight * fragrance_weight (non-psych tags)
  total_score      = psychotype_match + tag_match_score

Sort by total_score DESC → take top 5 → normalise to 70-99%
```

## 2. User Views Product Details

```
User              Mini App              API Server           PostgreSQL
 │                   │                     │                    │
 │  1. Tap product   │                     │                    │
 │──────────────────>│                     │                    │
 │                   │  2. GET /api/                              │
 │                   │     fragrances/{id}                       │
 │                   │────────────────────>│                    │
 │                   │                     │  3. SELECT with   │
 │                   │                     │     JSONB decoding│
 │                   │                     │───────────────────>│
 │                   │                     │  4. fragrance row │
 │                   │                     │<───────────────────│
 │                   │  5. Product JSON    │                    │
 │                   │<────────────────────│                    │
 │  6. Show product  │                     │                    │
 │     card with     │                     │                    │
 │     volume select │                     │                    │
 │<──────────────────│                     │                    │
```

## 3. Admin Adds a Fragrance via Catalog Bot

```
Admin              Catalog Bot             Database           Telegram API
 │                     │                      │                    │
 │  1. /add            │                      │                    │
 │────────────────────>│                      │                    │
 │                     │                      │                    │
 │  2. Prompt: ID      │                      │                    │
 │<────────────────────│                      │                    │
 │  3. "fresh-rose"    │                      │                    │
 │────────────────────>│                      │                    │
 │  ... (16 steps) ... │                      │                    │
 │                     │                      │                    │
 │  4. Prompt: photo   │                      │                    │
 │<────────────────────│                      │                    │
 │  5. Send photo      │                      │                    │
 │────────────────────>│                      │                    │
 │                     │                      │                    │
 │                     │  6. getFile/down-    │                    │
 │                     │     loadFile         │                    │
 │                     │──────────────────────────────────────────>│
 │                     │  7. file bytes       │                    │
 │                     │<──────────────────────────────────────────│
 │                     │                      │                    │
 │                     │  8. INSERT INTO      │                    │
 │                     │     fragrances       │                    │
 │                     │─────────────────────>│                    │
 │                     │  9. OK               │                    │
 │                     │<─────────────────────│                    │
 │                     │                      │                    │
 │  10. "Fragrance     │                      │                    │
 │      created!"      │                      │                    │
 │<────────────────────│                      │                    │
```

## 4. User Places an Order

```
User              Mini App                Telegram API
 │                   │                        │
 │  1. Add to cart   │                        │
 │──────────────────>│                        │
 │                   │  2. Persist to         │
 │                   │     localStorage       │
 │                   │     (aroma_cart_v1)    │
 │                   │                        │
 │  3. Go to cart    │                        │
 │──────────────────>│                        │
 │                   │  4. Show cart items    │
 │                   │     + totals           │
 │<──────────────────│                        │
 │                   │                        │
 │  5. Go to checkout│                        │
 │──────────────────>│                        │
 │                   │                        │
 │  6. Fill          │                        │
 │     recipient     │                        │
 │     form, address │                        │
 │═══════════════════>│                        │
 │                   │                        │
 │  7. Submit order  │                        │
 │──────────────────>│  8. Open Telegram      │
 │                   │     chat with config-  │
 │                   │     ured contact       │
 │                   │     (tg://)            │
 │                   │───────────────────────>│
 │                   │                        │
 │  9. Order details │                        │
 │     sent to       │                        │
 │     seller        │                        │
```

The order is delivered via a Telegram URI — no order management is stored in the
backend. Cart and order data persist only in the user's browser `localStorage`.
