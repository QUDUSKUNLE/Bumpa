CREATE TABLE payments (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  badge_code TEXT NOT NULL REFERENCES badges(code),
  amount_kobo BIGINT NOT NULL,
  provider_reference TEXT,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, badge_code)
);
