CREATE TABLE IF NOT EXISTS sales
(
    id             BIGSERIAL     PRIMARY KEY,
    product_id     BIGINT        NOT NULL REFERENCES products(id),
    quantity       INT           NOT NULL CHECK(quantity > 0),
    unit_price     NUMERIC       NOT NULL CHECK(unit_price > 0),
    payment_method VARCHAR(10)   NOT NULL CHECK(payment_method IN ('cash', 'card', 'online')),
    sold_at        TIMESTAMPTZ   NOT NULL,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);