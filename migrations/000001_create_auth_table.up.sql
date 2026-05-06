-- migrations/000001_create_auth_table.up.sql
-- Создание таблицы пользователей
CREATE TABLE auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE sessions (
	session_id TEXT PRIMARY KEY,
	login TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE balance (
    user_id UUID PRIMARY KEY,
	balance DECIMAL(10,2) NOT NULL DEFAULT 0,
	withdrawn DECIMAL(10,2) NOT NULL DEFAULT 0
);

ALTER TABLE sessions 
ADD CONSTRAINT fk_sessions_user 
FOREIGN KEY (user_id) REFERENCES auth(id);

ALTER TABLE balance 
ADD CONSTRAINT fk_sessions_user 
FOREIGN KEY (user_id) REFERENCES auth(id);