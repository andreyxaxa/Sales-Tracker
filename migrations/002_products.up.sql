CREATE TABLE IF NOT EXISTS products
(
    id          BIGSERIAL      PRIMARY KEY,
    category_id BIGINT         NOT NULL REFERENCES categories(id),
    name        VARCHAR(255)   NOT NULL,
    base_price  NUMERIC(12, 2) NOT NULL CHECK (base_price > 0)
);