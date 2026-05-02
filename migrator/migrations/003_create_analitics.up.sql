CREATE TABLE order_events_analytics (
  id UUID PRIMARY KEY,
  event_type VARCHAR(50),
  payload JSONB,
  received_at TIMESTAMPTZ DEFAULT NOW()
);
