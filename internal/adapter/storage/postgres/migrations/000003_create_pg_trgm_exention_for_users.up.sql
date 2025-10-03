CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX users_username_trgm_key
ON users USING gin (username gin_trgm_ops)