DELETE FROM badges
WHERE code IN (
    'bronze',
    'silver',
    'gold',
    'platinum',
    'diamond'
);
