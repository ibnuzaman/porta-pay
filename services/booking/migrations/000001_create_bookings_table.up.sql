CREATE TABLE IF NOT EXISTS bookings (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL,
    route_id    BIGINT NOT NULL,
    qty         INT    NOT NULL CHECK (qty > 0),
    status      TEXT   NOT NULL CHECK (status IN ('CREATED','PAID','CONFIRMED','EXPIRED')),
    price_total BIGINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);