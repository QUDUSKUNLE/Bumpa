CREATE TABLE badges (
  code TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  required_achievements INT NOT NULL
);

CREATE TABLE user_badges (
  user_id UUID NOT NULL REFERENCES users(id),
  badge_code TEXT NOT NULL REFERENCES badges(code),
  unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, badge_code)
);
