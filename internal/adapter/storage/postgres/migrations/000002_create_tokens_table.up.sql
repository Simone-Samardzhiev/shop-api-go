CREATE TYPE token_type_enum AS ENUM ('access_token', 'refresh_token');

CREATE TABLE tokens
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES users (id),
    token_type token_type_enum NOT NULL,
    expires    TIMESTAMP       NOT NULL
);