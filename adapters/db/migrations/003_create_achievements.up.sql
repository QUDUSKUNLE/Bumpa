CREATE TABLE achievements (
  code TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  achievement_group TEXT NOT NULL,
  position INT NOT NULL,
  UNIQUE (achievement_group, position)
);
