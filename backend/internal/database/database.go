package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	role TEXT NOT NULL CHECK (role IN ('user', 'admin')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS fragrances (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	brand TEXT NOT NULL,
	gender TEXT NOT NULL CHECK (gender IN ('male', 'female', 'unisex')),
	image_url TEXT NOT NULL DEFAULT '',
	volume NUMERIC(8, 2) NOT NULL DEFAULT 0,
	price NUMERIC(10, 2) NOT NULL DEFAULT 0,
	stock_status TEXT NOT NULL CHECK (stock_status IN ('in_stock', 'out_of_stock')),
	description TEXT NOT NULL DEFAULT '',
	top_notes JSONB NOT NULL DEFAULT '[]'::JSONB,
	middle_notes JSONB NOT NULL DEFAULT '[]'::JSONB,
	base_notes JSONB NOT NULL DEFAULT '[]'::JSONB,
	main_accords JSONB NOT NULL DEFAULT '[]'::JSONB,
	longevity INTEGER NOT NULL DEFAULT 3 CHECK (longevity BETWEEN 1 AND 5),
	projection INTEGER NOT NULL DEFAULT 3 CHECK (projection BETWEEN 1 AND 5),
	visibility INTEGER NOT NULL DEFAULT 3 CHECK (visibility BETWEEN 1 AND 5),
	versatility INTEGER NOT NULL DEFAULT 3 CHECK (versatility BETWEEN 1 AND 5),
	seasons JSONB NOT NULL DEFAULT '[]'::JSONB,
	time_of_day JSONB NOT NULL DEFAULT '[]'::JSONB,
	situations JSONB NOT NULL DEFAULT '[]'::JSONB,
	matching_profiles JSONB NOT NULL DEFAULT '[]'::JSONB,
	why_recommended TEXT NOT NULL DEFAULT '',
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tags (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	type TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS fragrance_tags (
	fragrance_id TEXT NOT NULL REFERENCES fragrances(id) ON DELETE CASCADE,
	tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
	weight INTEGER NOT NULL DEFAULT 1,
	PRIMARY KEY (fragrance_id, tag_id)
);

CREATE TABLE IF NOT EXISTS questions (
	id TEXT PRIMARY KEY,
	text TEXT NOT NULL,
	type TEXT NOT NULL CHECK (type IN ('single_choice', 'multiple_choice')),
	sort_order INTEGER NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS answer_options (
	id TEXT PRIMARY KEY,
	question_id TEXT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
	text TEXT NOT NULL,
	value TEXT NOT NULL,
	sort_order INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS answer_option_tags (
	answer_option_id TEXT NOT NULL REFERENCES answer_options(id) ON DELETE CASCADE,
	tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
	weight INTEGER NOT NULL DEFAULT 1,
	PRIMARY KEY (answer_option_id, tag_id)
);

CREATE TABLE IF NOT EXISTS questionnaire_sessions (
	id TEXT PRIMARY KEY,
	user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
	status TEXT NOT NULL CHECK (status IN ('started', 'completed')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	completed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS questionnaire_answers (
	id TEXT PRIMARY KEY,
	session_id TEXT NOT NULL REFERENCES questionnaire_sessions(id) ON DELETE CASCADE,
	question_id TEXT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
	answer_option_id TEXT NOT NULL REFERENCES answer_options(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS aroma_profiles (
	id TEXT PRIMARY KEY,
	session_id TEXT NOT NULL REFERENCES questionnaire_sessions(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS recommendations (
	id TEXT PRIMARY KEY,
	session_id TEXT NOT NULL REFERENCES questionnaire_sessions(id) ON DELETE CASCADE,
	algorithm_version TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS recommendation_items (
	id TEXT PRIMARY KEY,
	recommendation_id TEXT NOT NULL REFERENCES recommendations(id) ON DELETE CASCADE,
	fragrance_id TEXT NOT NULL REFERENCES fragrances(id) ON DELETE CASCADE,
	score INTEGER NOT NULL,
	reason TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
	id TEXT PRIMARY KEY,
	user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
	fragrance_id TEXT NOT NULL REFERENCES fragrances(id) ON DELETE RESTRICT,
	status TEXT NOT NULL DEFAULT 'new',
	contact_name TEXT NOT NULL DEFAULT '',
	contact_phone TEXT NOT NULL DEFAULT '',
	contact_telegram TEXT NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`)
	if err != nil {
		return err
	}

	return seedMVPData(ctx, pool)
}

func seedMVPData(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
INSERT INTO tags (id, name, type) VALUES
	('calm', 'спокойный', 'character'),
	('confidence', 'уверенный', 'character'),
	('bright', 'яркий', 'character'),
	('energy', 'энергичный', 'character'),
	('mystery', 'загадочный', 'character'),
	('romantic', 'романтичный', 'character'),
	('soft', 'мягкий', 'feeling'),
	('office', 'офисный', 'usage'),
	('date', 'для свиданий', 'usage'),
	('party', 'для вечеринок', 'usage'),
	('daily', 'повседневный', 'usage'),
	('clean', 'чистый', 'feeling'),
	('fresh', 'свежий', 'feeling'),
	('warm', 'тёплый', 'feeling'),
	('cozy', 'уютный', 'feeling'),
	('deep', 'глубокий', 'feeling'),
	('light', 'лёгкий', 'feeling'),
	('noticeable', 'заметный', 'usage'),
	('trail', 'шлейфовый', 'usage'),
	('morning', 'утренний', 'usage'),
	('day', 'дневной', 'usage'),
	('evening', 'вечерний', 'usage'),
	('night', 'ночной', 'usage'),
	('experimental', 'необычный', 'style'),
	('reliable', 'надежный', 'style')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, type = EXCLUDED.type;

INSERT INTO questions (id, text, type, sort_order, is_active) VALUES
	('q1', 'Если описать вас в трёх словах через аромат, что это будет?', 'single_choice', 1, TRUE),
	('q2', 'В каких жизненных моментах вы хотите ощущать аромат?', 'single_choice', 2, TRUE),
	('q3', 'Представьте, что аромат — это цвет. Какой вы выберете?', 'single_choice', 3, TRUE),
	('q4', 'Насколько вы хотите, чтобы люди замечали ваш аромат?', 'single_choice', 4, TRUE),
	('q5', 'Если аромат был бы временем суток, какое это будет время?', 'single_choice', 5, TRUE),
	('q6', 'Какой эффект на людей вам нравится оставлять через аромат?', 'single_choice', 6, TRUE),
	('q7', 'Вы любите, когда аромат необычный или предпочитаете знакомое?', 'single_choice', 7, TRUE)
ON CONFLICT (id) DO UPDATE SET text = EXCLUDED.text, type = EXCLUDED.type, sort_order = EXCLUDED.sort_order, is_active = EXCLUDED.is_active;

INSERT INTO answer_options (id, question_id, text, value, sort_order) VALUES
	('q1_a1', 'q1', 'Спокойный и уверенный', 'calm_confident', 1),
	('q1_a2', 'q1', 'Яркий и энергичный', 'bright_energy', 2),
	('q1_a3', 'q1', 'Тайный и загадочный', 'secret_mystery', 3),
	('q1_a4', 'q1', 'Романтичный и мягкий', 'romantic_soft', 4),
	('q2_a1', 'q2', 'На работе или учебе', 'work_study', 1),
	('q2_a2', 'q2', 'На встречах и свиданиях', 'dates_meetings', 2),
	('q2_a3', 'q2', 'На вечеринках и фестивалях', 'parties_festivals', 3),
	('q2_a4', 'q2', 'В повседневной жизни, как часть образа', 'daily_style', 4),
	('q3_a1', 'q3', 'Светлый и прозрачный (легкий, свежий)', 'light_fresh', 1),
	('q3_a2', 'q3', 'Тёплый и уютный (мягкий, согревающий)', 'warm_cozy', 2),
	('q3_a3', 'q3', 'Яркий и насыщенный (выразительный, динамичный)', 'bright_dynamic', 3),
	('q3_a4', 'q3', 'Темный и глубокий (интригующий, загадочный)', 'dark_deep', 4),
	('q4_a1', 'q4', 'Почти не ощущают', 'barely_noticeable', 1),
	('q4_a2', 'q4', 'Замечают время от времени', 'sometimes_noticeable', 2),
	('q4_a3', 'q4', 'Часто обращают внимание', 'often_noticeable', 3),
	('q4_a4', 'q4', 'Шлейф, который невозможно не заметить', 'strong_trail', 4),
	('q5_a1', 'q5', 'Утро, свежо и бодро', 'morning_fresh', 1),
	('q5_a2', 'q5', 'День, комфортно и спокойно', 'day_calm', 2),
	('q5_a3', 'q5', 'Вечер, тепло и уютно', 'evening_warm', 3),
	('q5_a4', 'q5', 'Ночь, таинственно и чувственно', 'night_mystery', 4),
	('q6_a1', 'q6', 'Согревающее впечатление, мягкость', 'warm_soft_effect', 1),
	('q6_a2', 'q6', 'Энергия и драйв', 'energy_drive', 2),
	('q6_a3', 'q6', 'Очарование и притяжение', 'charm_attraction', 3),
	('q6_a4', 'q6', 'Тайна и интрига', 'mystery_intrigue', 4),
	('q7_a1', 'q7', 'Интересно пробовать новое', 'try_new', 1),
	('q7_a2', 'q7', 'Предпочитаю узнаваемое и надежное', 'reliable_known', 2),
	('q7_a3', 'q7', 'Иногда экспериментирую', 'sometimes_experiment', 3),
	('q7_a4', 'q7', 'Люблю неожиданное, смелое', 'bold_unexpected', 4)
ON CONFLICT (id) DO UPDATE SET question_id = EXCLUDED.question_id, text = EXCLUDED.text, value = EXCLUDED.value, sort_order = EXCLUDED.sort_order;

INSERT INTO answer_option_tags (answer_option_id, tag_id, weight) VALUES
	('q1_a1', 'calm', 3), ('q1_a1', 'confidence', 2), ('q1_a2', 'bright', 3), ('q1_a2', 'energy', 3),
	('q1_a3', 'mystery', 3), ('q1_a3', 'deep', 2), ('q1_a4', 'romantic', 3), ('q1_a4', 'soft', 3),
	('q2_a1', 'office', 3), ('q2_a1', 'day', 2), ('q2_a2', 'date', 3), ('q2_a2', 'romantic', 2),
	('q2_a3', 'party', 3), ('q2_a3', 'bright', 2), ('q2_a4', 'daily', 3), ('q2_a4', 'reliable', 2),
	('q3_a1', 'light', 3), ('q3_a1', 'fresh', 3), ('q3_a1', 'clean', 2), ('q3_a2', 'warm', 3),
	('q3_a2', 'cozy', 3), ('q3_a2', 'soft', 2), ('q3_a3', 'bright', 3), ('q3_a3', 'energy', 2),
	('q3_a4', 'deep', 3), ('q3_a4', 'mystery', 2), ('q4_a1', 'light', 3), ('q4_a1', 'calm', 2),
	('q4_a2', 'daily', 2), ('q4_a2', 'reliable', 2), ('q4_a3', 'noticeable', 3), ('q4_a3', 'confidence', 2),
	('q4_a4', 'trail', 3), ('q4_a4', 'deep', 2), ('q5_a1', 'morning', 3), ('q5_a1', 'fresh', 2),
	('q5_a2', 'day', 3), ('q5_a2', 'calm', 2), ('q5_a3', 'evening', 3), ('q5_a3', 'warm', 2),
	('q5_a4', 'night', 3), ('q5_a4', 'mystery', 2), ('q6_a1', 'warm', 3), ('q6_a1', 'soft', 2),
	('q6_a2', 'energy', 3), ('q6_a2', 'bright', 2), ('q6_a3', 'romantic', 3), ('q6_a3', 'date', 2),
	('q6_a4', 'mystery', 3), ('q6_a4', 'deep', 2), ('q7_a1', 'experimental', 3), ('q7_a1', 'bright', 1),
	('q7_a2', 'reliable', 3), ('q7_a2', 'daily', 1), ('q7_a3', 'experimental', 2), ('q7_a3', 'reliable', 1),
	('q7_a4', 'experimental', 3), ('q7_a4', 'mystery', 1)
ON CONFLICT (answer_option_id, tag_id) DO UPDATE SET weight = EXCLUDED.weight;

INSERT INTO fragrances (
	id, name, brand, gender, image_url, volume, price, stock_status, description,
	top_notes, middle_notes, base_notes, main_accords, longevity, projection, visibility, versatility,
	seasons, time_of_day, situations, matching_profiles, why_recommended, is_active
) VALUES
	('fresh-office', 'Fresh Office', 'AromaType', 'unisex', '', 2, 490, 'in_stock',
	 'Чистый и свежий аромат для учебы, офиса и спокойного повседневного образа.',
	 '["бергамот", "лимон"]', '["лаванда"]', '["мускус"]', '["свежий", "чистый"]', 3, 2, 2, 5,
	 '["весна", "лето"]', '["утро", "день"]', '["офис", "учеба", "повседневность"]',
	 '["Спокойный минималист"]', 'Подходит для спокойного и чистого образа без лишней громкости.', TRUE),
	('warm-date', 'Warm Date', 'AromaType', 'unisex', '', 2, 540, 'in_stock',
	 'Мягкий теплый аромат для встреч, свиданий и уютного вечернего настроения.',
	 '["кардамон"]', '["жасмин"]', '["ваниль", "мускус"]', '["тёплый", "уютный"]', 4, 3, 3, 4,
	 '["осень", "зима"]', '["вечер"]', '["свидание", "встреча"]',
	 '["Мягкая романтика"]', 'Подходит для теплого, мягкого и притягательного впечатления.', TRUE),
	('bright-party', 'Bright Party', 'AromaType', 'unisex', '', 2, 590, 'in_stock',
	 'Яркий энергичный аромат для вечеринок, фестивалей и заметного образа.',
	 '["грейпфрут"]', '["имбирь"]', '["амброксан"]', '["яркий", "энергичный"]', 4, 4, 4, 3,
	 '["весна", "лето"]', '["день", "вечер"]', '["вечеринка", "фестиваль"]',
	 '["Яркая энергия"]', 'Подходит для выразительного и динамичного образа.', TRUE),
	('mystic-night', 'Mystic Night', 'AromaType', 'unisex', '', 2, 650, 'in_stock',
	 'Глубокий и загадочный аромат для ночного настроения и смелого впечатления.',
	 '["черный перец"]', '["ладан"]', '["ветивер", "амброксан"]', '["глубокий", "загадочный"]', 5, 4, 5, 2,
	 '["осень", "зима"]', '["ночь"]', '["вечер", "особый случай"]',
	 '["Таинственный акцент"]', 'Подходит для глубокого, интригующего и необычного образа.', TRUE),
	('daily-soft', 'Daily Soft', 'AromaType', 'unisex', '', 2, 450, 'in_stock',
	 'Надежный мягкий аромат на каждый день, когда хочется комфорта и универсальности.',
	 '["мандарин"]', '["пудровые ноты"]', '["белый мускус"]', '["мягкий", "повседневный"]', 3, 2, 2, 5,
	 '["весна", "осень"]', '["день"]', '["повседневность", "учеба", "прогулка"]',
	 '["Сбалансированный профиль"]', 'Подходит для узнаваемого, надежного и комфортного повседневного выбора.', TRUE)
ON CONFLICT (id) DO UPDATE SET
	name = EXCLUDED.name,
	brand = EXCLUDED.brand,
	gender = EXCLUDED.gender,
	image_url = EXCLUDED.image_url,
	volume = EXCLUDED.volume,
	price = EXCLUDED.price,
	stock_status = EXCLUDED.stock_status,
	description = EXCLUDED.description,
	top_notes = EXCLUDED.top_notes,
	middle_notes = EXCLUDED.middle_notes,
	base_notes = EXCLUDED.base_notes,
	main_accords = EXCLUDED.main_accords,
	longevity = EXCLUDED.longevity,
	projection = EXCLUDED.projection,
	visibility = EXCLUDED.visibility,
	versatility = EXCLUDED.versatility,
	seasons = EXCLUDED.seasons,
	time_of_day = EXCLUDED.time_of_day,
	situations = EXCLUDED.situations,
	matching_profiles = EXCLUDED.matching_profiles,
	why_recommended = EXCLUDED.why_recommended,
	is_active = EXCLUDED.is_active,
	updated_at = now();

INSERT INTO fragrance_tags (fragrance_id, tag_id, weight) VALUES
	('fresh-office', 'fresh', 3), ('fresh-office', 'clean', 3), ('fresh-office', 'office', 3), ('fresh-office', 'calm', 2), ('fresh-office', 'day', 2), ('fresh-office', 'light', 2),
	('warm-date', 'warm', 3), ('warm-date', 'cozy', 3), ('warm-date', 'date', 3), ('warm-date', 'romantic', 2), ('warm-date', 'soft', 2), ('warm-date', 'evening', 2),
	('bright-party', 'bright', 3), ('bright-party', 'energy', 3), ('bright-party', 'party', 3), ('bright-party', 'noticeable', 2), ('bright-party', 'experimental', 1),
	('mystic-night', 'mystery', 3), ('mystic-night', 'deep', 3), ('mystic-night', 'night', 3), ('mystic-night', 'trail', 2), ('mystic-night', 'experimental', 2),
	('daily-soft', 'daily', 3), ('daily-soft', 'soft', 3), ('daily-soft', 'reliable', 3), ('daily-soft', 'calm', 2), ('daily-soft', 'light', 2)
ON CONFLICT (fragrance_id, tag_id) DO UPDATE SET weight = EXCLUDED.weight;
`)
	return err
}
