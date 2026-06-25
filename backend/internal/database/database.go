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
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS recommendation_items;
DROP TABLE IF EXISTS recommendations;
DROP TABLE IF EXISTS aroma_profiles;
DROP TABLE IF EXISTS questionnaire_answers;
DROP TABLE IF EXISTS questionnaire_sessions;

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
	image_url TEXT NOT NULL DEFAULT '',
	price NUMERIC(10, 2) NOT NULL DEFAULT 0,
	volume_options JSONB NOT NULL DEFAULT '[]'::JSONB,
	description TEXT NOT NULL DEFAULT '',
	top_notes JSONB NOT NULL DEFAULT '[]'::JSONB,
	middle_notes JSONB NOT NULL DEFAULT '[]'::JSONB,
	base_notes JSONB NOT NULL DEFAULT '[]'::JSONB,
	main_accords JSONB NOT NULL DEFAULT '[]'::JSONB,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE fragrances ADD COLUMN IF NOT EXISTS volume_options JSONB NOT NULL DEFAULT '[]'::JSONB;
ALTER TABLE fragrances DROP COLUMN IF EXISTS gender;
ALTER TABLE fragrances DROP COLUMN IF EXISTS volume;
ALTER TABLE fragrances DROP COLUMN IF EXISTS stock_status;
ALTER TABLE fragrances DROP COLUMN IF EXISTS longevity;
ALTER TABLE fragrances DROP COLUMN IF EXISTS projection;
ALTER TABLE fragrances DROP COLUMN IF EXISTS visibility;
ALTER TABLE fragrances DROP COLUMN IF EXISTS versatility;
ALTER TABLE fragrances DROP COLUMN IF EXISTS seasons;
ALTER TABLE fragrances DROP COLUMN IF EXISTS time_of_day;
ALTER TABLE fragrances DROP COLUMN IF EXISTS situations;
ALTER TABLE fragrances DROP COLUMN IF EXISTS matching_profiles;
ALTER TABLE fragrances DROP COLUMN IF EXISTS why_recommended;

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
	('reliable', 'надежный', 'style'),
	('psych_drive', 'Драйв / Экстраверсия', 'psychotype'),
	('psych_focus', 'Интеллект / Фокус', 'psychotype'),
	('psych_aesthetic', 'Эстетика / Гедонизм', 'psychotype'),
	('psych_power', 'Власть / Доминанта', 'psychotype')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, type = EXCLUDED.type;

INSERT INTO questions (id, text, type, sort_order, is_active) VALUES
	('q1', 'Вы входите в незнакомую компанию или на встречу. Ваше действие?', 'single_choice', 1, TRUE),
	('q2', 'Где вы чувствуете максимальный прилив энергии?', 'single_choice', 2, TRUE),
	('q3', 'Ваши планы резко рухнули из-за форс-мажора. Реакция?', 'single_choice', 3, TRUE),
	('q4', 'Что ваш образ должен транслировать окружающим?', 'single_choice', 4, TRUE),
	('q5', 'Вы беретесь за сложный новый проект. Что поможет победить?', 'single_choice', 5, TRUE),
	('q6', 'Материал, который для вас означает «качество»:', 'single_choice', 6, TRUE),
	('q7', 'Эффект от аромата, который вам нужен уже через 15 минут утром:', 'single_choice', 7, TRUE),
	('q8', 'Какое описание аромата вызывает наибольший отклик?', 'single_choice', 8, TRUE)
ON CONFLICT (id) DO UPDATE SET text = EXCLUDED.text, type = EXCLUDED.type, sort_order = EXCLUDED.sort_order, is_active = EXCLUDED.is_active;

INSERT INTO answer_options (id, question_id, text, value, sort_order) VALUES
	('q1_a1', 'q1', 'Легко привлеку внимание, я в центре событий.', 'drive_attention', 1),
	('q1_a2', 'q1', 'Займу позицию наблюдателя со стороны.', 'focus_observer', 2),
	('q1_a3', 'q1', 'Найду одного собеседника для спокойной беседы.', 'aesthetic_private_talk', 3),
	('q1_a4', 'q1', 'Возьму инициативу и направлю разговор.', 'power_initiative', 4),
	('q2_a1', 'q2', 'В движении: фестивали, тусовки, путешествия.', 'drive_motion', 1),
	('q2_a2', 'q2', 'В тишине: за сложной и глубокой работой.', 'focus_deep_work', 2),
	('q2_a3', 'q2', 'В красивом месте: концептуальный ресторан, выставка.', 'aesthetic_beautiful_place', 3),
	('q2_a4', 'q2', 'В условиях вызова: важные переговоры или спорт.', 'power_challenge', 4),
	('q3_a1', 'q3', 'Легко импровизирую на ходу.', 'drive_improvise', 1),
	('q3_a2', 'q3', 'Спокойно анализирую новые вводные данные.', 'focus_analyze', 2),
	('q3_a3', 'q3', 'Беру паузу, чтобы вернуть внутренний баланс.', 'aesthetic_balance', 3),
	('q3_a4', 'q3', 'Беру контроль на себя и перестраиваю ситуацию.', 'power_control', 4),
	('q4_a1', 'q4', 'Харизму, открытость и оптимизм.', 'drive_charisma', 1),
	('q4_a2', 'q4', 'Независимость, интеллект и загадку.', 'focus_intellect', 2),
	('q4_a3', 'q4', 'Безупречный вкус, стиль и ухоженность.', 'aesthetic_taste', 3),
	('q4_a4', 'q4', 'Силу, авторитет и уверенность лидера.', 'power_authority', 4),
	('q5_a1', 'q5', 'Способность быстро адаптироваться и общаться.', 'drive_adapt', 1),
	('q5_a2', 'q5', 'Глубокий анализ и алгоритмы.', 'focus_algorithms', 2),
	('q5_a3', 'q5', 'Прошлый опыт и перфекционизм в деталях.', 'aesthetic_perfection', 3),
	('q5_a4', 'q5', 'Железная дисциплина и вера в свой результат.', 'power_discipline', 4),
	('q6_a1', 'q6', 'Технологичная ткань, легкий лен или хлопок.', 'drive_light_fabric', 1),
	('q6_a2', 'q6', 'Плотное сукно, тяжелый твид или шерсть.', 'focus_dense_wool', 2),
	('q6_a3', 'q6', 'Мягкий кашемир премиум-качества или шелк.', 'aesthetic_cashmere', 3),
	('q6_a4', 'q6', 'Толстая дорогая кожа или грубая замша.', 'power_leather', 4),
	('q7_a1', 'q7', 'Заряд энергии: мгновенно проснуться и взбодриться.', 'drive_energy_boost', 1),
	('q7_a2', 'q7', 'Режим фокуса: ментальный кокон для работы.', 'focus_cocoon', 2),
	('q7_a3', 'q7', 'Гедонизм: комфорт и эстетическое удовольствие.', 'aesthetic_hedonism', 3),
	('q7_a4', 'q7', 'Эффект брони: чувство силы и неуязвимости.', 'power_armor', 4),
	('q8_a1', 'q8', 'Мокрый асфальт, цитрусы, мята и морской бриз.', 'drive_wet_asphalt', 1),
	('q8_a2', 'q8', 'Старая библиотека, дорогой алкоголь, ладан и специи.', 'focus_library', 2),
	('q8_a3', 'q8', 'Уходовый крем, новая замша и мягкая ваниль.', 'aesthetic_cream_suede', 3),
	('q8_a4', 'q8', 'Салон нового авто, дорогой табак, костер и смола.', 'power_tobacco_resin', 4)
ON CONFLICT (id) DO UPDATE SET question_id = EXCLUDED.question_id, text = EXCLUDED.text, value = EXCLUDED.value, sort_order = EXCLUDED.sort_order;

DELETE FROM answer_option_tags
WHERE answer_option_id IN (
	SELECT id FROM answer_options WHERE question_id IN ('q1', 'q2', 'q3', 'q4', 'q5', 'q6', 'q7', 'q8')
);

INSERT INTO answer_option_tags (answer_option_id, tag_id, weight) VALUES
	('q1_a1', 'psych_drive', 3), ('q1_a2', 'psych_focus', 3), ('q1_a3', 'psych_aesthetic', 3), ('q1_a4', 'psych_power', 3),
	('q2_a1', 'psych_drive', 3), ('q2_a2', 'psych_focus', 3), ('q2_a3', 'psych_aesthetic', 3), ('q2_a4', 'psych_power', 3),
	('q3_a1', 'psych_drive', 3), ('q3_a2', 'psych_focus', 3), ('q3_a3', 'psych_aesthetic', 3), ('q3_a4', 'psych_power', 3),
	('q4_a1', 'psych_drive', 3), ('q4_a2', 'psych_focus', 3), ('q4_a3', 'psych_aesthetic', 3), ('q4_a4', 'psych_power', 3),
	('q5_a1', 'psych_drive', 3), ('q5_a2', 'psych_focus', 3), ('q5_a3', 'psych_aesthetic', 3), ('q5_a4', 'psych_power', 3),
	('q6_a1', 'psych_drive', 3), ('q6_a2', 'psych_focus', 3), ('q6_a3', 'psych_aesthetic', 3), ('q6_a4', 'psych_power', 3),
	('q7_a1', 'psych_drive', 3), ('q7_a2', 'psych_focus', 3), ('q7_a3', 'psych_aesthetic', 3), ('q7_a4', 'psych_power', 3),
	('q8_a1', 'psych_drive', 3), ('q8_a2', 'psych_focus', 3), ('q8_a3', 'psych_aesthetic', 3), ('q8_a4', 'psych_power', 3)
ON CONFLICT (answer_option_id, tag_id) DO UPDATE SET weight = EXCLUDED.weight;

INSERT INTO fragrances (
	id, name, brand, image_url, price, volume_options, description,
	top_notes, middle_notes, base_notes, main_accords, is_active
) VALUES
	('fresh-office', 'Fresh Office', 'AromaType', '', 490, '[{"volumeMl":50,"price":490},{"volumeMl":100,"price":890}]',
	 'Чистый и свежий аромат для учебы, офиса и спокойного повседневного образа.',
	 '["бергамот", "лимон"]', '["лаванда"]', '["мускус"]', '["свежий", "чистый"]', TRUE),
	('warm-date', 'Warm Date', 'AromaType', '', 540, '[{"volumeMl":50,"price":540},{"volumeMl":100,"price":980}]',
	 'Мягкий теплый аромат для встреч, свиданий и уютного вечернего настроения.',
	 '["кардамон"]', '["жасмин"]', '["ваниль", "мускус"]', '["тёплый", "уютный"]', TRUE),
	('bright-party', 'Bright Party', 'AromaType', '', 590, '[{"volumeMl":50,"price":590},{"volumeMl":100,"price":1090}]',
	 'Яркий энергичный аромат для вечеринок, фестивалей и заметного образа.',
	 '["грейпфрут"]', '["имбирь"]', '["амброксан"]', '["яркий", "энергичный"]', TRUE),
	('mystic-night', 'Mystic Night', 'AromaType', '', 650, '[{"volumeMl":50,"price":650},{"volumeMl":100,"price":1190}]',
	 'Глубокий и загадочный аромат для ночного настроения и смелого впечатления.',
	 '["черный перец"]', '["ладан"]', '["ветивер", "амброксан"]', '["глубокий", "загадочный"]', TRUE),
	('daily-soft', 'Daily Soft', 'AromaType', '', 450, '[{"volumeMl":50,"price":450},{"volumeMl":100,"price":790}]',
	 'Надежный мягкий аромат на каждый день, когда хочется комфорта и универсальности.',
	 '["мандарин"]', '["пудровые ноты"]', '["белый мускус"]', '["мягкий", "повседневный"]', TRUE)
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
	is_active = EXCLUDED.is_active,
	updated_at = now();

INSERT INTO fragrance_tags (fragrance_id, tag_id, weight) VALUES
	('fresh-office', 'fresh', 3), ('fresh-office', 'clean', 3), ('fresh-office', 'office', 3), ('fresh-office', 'calm', 2), ('fresh-office', 'day', 2), ('fresh-office', 'light', 2), ('fresh-office', 'psych_focus', 3), ('fresh-office', 'psych_drive', 1),
	('warm-date', 'warm', 3), ('warm-date', 'cozy', 3), ('warm-date', 'date', 3), ('warm-date', 'romantic', 2), ('warm-date', 'soft', 2), ('warm-date', 'evening', 2), ('warm-date', 'psych_aesthetic', 3),
	('bright-party', 'bright', 3), ('bright-party', 'energy', 3), ('bright-party', 'party', 3), ('bright-party', 'noticeable', 2), ('bright-party', 'experimental', 1), ('bright-party', 'psych_drive', 3),
	('mystic-night', 'mystery', 3), ('mystic-night', 'deep', 3), ('mystic-night', 'night', 3), ('mystic-night', 'trail', 2), ('mystic-night', 'experimental', 2), ('mystic-night', 'psych_focus', 2), ('mystic-night', 'psych_power', 3),
	('daily-soft', 'daily', 3), ('daily-soft', 'soft', 3), ('daily-soft', 'reliable', 3), ('daily-soft', 'calm', 2), ('daily-soft', 'light', 2), ('daily-soft', 'psych_aesthetic', 2), ('daily-soft', 'psych_focus', 1)
ON CONFLICT (fragrance_id, tag_id) DO UPDATE SET weight = EXCLUDED.weight;
`)
	return err
}
