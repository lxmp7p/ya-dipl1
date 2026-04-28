-- migrations/000001_create_auth_table.up.sql
-- Создание таблицы пользователей
CREATE TABLE auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE sessions (
	session_id TEXT PRIMARY KEY,
	login TEXT NOT NULL
);