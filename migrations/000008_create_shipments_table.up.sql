CREATE TABLE IF NOT EXISTS shipments (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
  idempotency_key UUID,

  status TEXT NOT NULL DEFAULT 'AWAITING_PICKUP',
  carrier TEXT NOT NULL,
  shipping_cost NUMERIC(10,2) NOT NULL,
  tracking_number UUID UNIQUE NOT NULL,

  shipped_at TIMESTAMPTZ,
  delivered_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipments_idempotency_key ON shipments(idempotency_key);
