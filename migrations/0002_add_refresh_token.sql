-- +goose Up
ALTER TABLE sessions
    ADD COLUMN refresh_token TEXT UNIQUE;

-- Обновляем существующие сессии, устанавливая уникальные refresh_token
UPDATE sessions
SET refresh_token = gen_random_uuid()::text
WHERE refresh_token IS NULL;

-- Делаем поле обязательным после обновления существующих записей
ALTER TABLE sessions
    ALTER COLUMN refresh_token SET NOT NULL;

-- +goose Down
ALTER TABLE sessions
    DROP COLUMN refresh_token; 