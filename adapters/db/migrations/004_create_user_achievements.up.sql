CREATE TABLE user_achievements (
  user_id UUID NOT NULL REFERENCES users(id),
  achievement_code TEXT NOT NULL REFERENCES achievements(code),
  unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, achievement_code)
);
