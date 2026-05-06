-- migrations/000001_create_auth_table.down.sql
-- Откат создания таблицы пользователей
DROP TABLE IF EXISTS auth; 
DROP TABLE IF EXISTS sessions;