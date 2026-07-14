CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    subtotal NUMERIC(10,2) NOT NULL,
    order_total NUMERIC(10,2) NOT NULL,
    shipping_cost NUMERIC(10,2) NOT NULL,

    payment_id UUID,
    waybill UUID,

    order_status TEXT NOT NULL DEFAULT 'PENDING_PAYMENT',
    payment_status TEXT NOT NULL DEFAULT 'PENDING',

    first_name TEXT NOT NULL,
    last_name TEXT,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    address1 TEXT NOT NULL,
    address2 TEXT,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    zip TEXT NOT NULL,

    idempotency_key UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_idempotency_key ON orders (idempotency_key);
