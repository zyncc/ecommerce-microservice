CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL UNIQUE REFERENCES product(id) ON DELETE CASCADE,
    small INT NOT NULL CHECK (small >= 0),
    medium INT NOT NULL CHECK (medium >= 0),
    large INT NOT NULL CHECK (large >= 0),
    extra_large INT NOT NULL CHECK (extra_large >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
