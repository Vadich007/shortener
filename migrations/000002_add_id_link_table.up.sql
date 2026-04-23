-- Добавление колонки id как BIGSERIAL (автоинкремент)
ALTER TABLE links ADD COLUMN id BIGSERIAL NOT NULL;

-- id - первичным ключом
ALTER TABLE links ADD PRIMARY KEY (id);

-- Создание индекса на поле original_url
CREATE UNIQUE INDEX IF NOT EXISTS idx_original_url_unique ON links (original_url);