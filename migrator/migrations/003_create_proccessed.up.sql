CREATE TABLE proccessed_events (
  event_id UUID PRIMARY KEY NOT NULL,
  proccessed_at TIMESTAMPTZ DEFAULT NOW ()
)