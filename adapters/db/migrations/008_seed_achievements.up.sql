INSERT INTO achievements (
    code,
    name,
    achievement_group,
    position
)
VALUES
    (
        'first_purchase',
        'First Purchase',
        'shopping',
        1
    ),
    (
        'three_purchases',
        'Three Purchases',
        'shopping',
        2
    ),
    (
        'five_purchases',
        'Five Purchases',
        'shopping',
        3
    ),
    (
        'ten_purchases',
        'Ten Purchases',
        'shopping',
        4
    ),
    (
        'twenty_purchases',
        'Twenty Purchases',
        'shopping',
        5
    )
ON CONFLICT (code) DO NOTHING;
