-- Удаление PRIMARY KEY и колонку id
ALTER TABLE links DROP COLUMN id CASCADE;

-- Удаление уникального индекса
DROP INDEX IF EXISTS idx_original_url_unique;