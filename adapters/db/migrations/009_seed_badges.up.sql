INSERT INTO badges (
    code,
    name,
    required_achievements
)
VALUES
    ('bronze', 'Bronze', 1),
    ('silver', 'Silver', 3),
    ('gold', 'Gold', 5),
    ('platinum', 'Platinum', 10),
    ('diamond', 'Diamond', 20)
ON CONFLICT (code) DO NOTHING;
