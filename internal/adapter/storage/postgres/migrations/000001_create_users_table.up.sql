CREATE TYPE user_role_enum AS ENUM ('admin', 'client', 'delivery', 'warehouse');

CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    username   VARCHAR(255) CHECK ( length(username) >= 8 ) NOT NUlL UNIQUE,
    email      VARCHAR(255) CHECK ( length(email) >= 12 )    NOT NULL UNIQUE,
    password   VARCHAR(255)                                NOT NULL,
    role       user_role_enum                              NOT NULL DEFAULT ('client'),
    created_at TIMESTAMP                                   NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP                                   NOT NULL DEFAULT (now())
);