CREATE TABLE orders (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL,
  status VARCHAR(20) NOT NULL CHECK (status IN ('created', 'paid', 'shipped')),
  amount INTEGER NOT NULL CHECK (amount > 0),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

