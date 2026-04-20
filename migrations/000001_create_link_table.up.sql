-- Создание таблицу для хранения URL
CREATE TABLE IF NOT EXISTS links (
    shorted_url TEXT NOT NULL,
    original_url TEXT NOT NULL
);

-- Создание индекса для ускорения поиска
CREATE INDEX idx_shorted_url ON links(shorted_url);