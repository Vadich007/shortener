-- Добавление колонки id как BIGSERIAL (автоинкремент)
ALTER TABLE links ADD COLUMN id BIGSERIAL NOT NULL;

-- Деление её первичным ключом
ALTER TABLE links ADD PRIMARY KEY (id);