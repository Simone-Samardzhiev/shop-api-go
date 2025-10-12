CREATE TABLE category_section
(
    id   UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE categories
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    section_id UUID REFERENCES category_section (id) ON DELETE CASCADE,
    UNIQUE (name, section_id)
);

CREATE TABLE subcategories
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    category_id UUID REFERENCES categories (id) ON DELETE CASCADE,
    UNIQUE (name, category_id)
);

CREATE TABLE products
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL UNIQUE CHECK ( length(name) > 8 ),
    description TEXT           NOT NULL CHECK ( length(description) > 24 ),
    price       NUMERIC(10, 2) NOT NULL CHECK ( price > 0 ),
    rating      NUMERIC(2, 1)  NOT NULL CHECK ( rating >= 0 AND rating <= 5 ),
    count       INT            NOT NULL CHECK ( count >= 0 ),
    image_url   TEXT           NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP      NOT NULL DEFAULT now()
);

CREATE TABLE product_subcategories
(
    product_id     UUID REFERENCES products (id) ON DELETE CASCADE,
    subcategory_id UUID REFERENCES subcategories (id) ON DELETE RESTRICT,
    PRIMARY KEY (product_id, subcategory_id)
);

