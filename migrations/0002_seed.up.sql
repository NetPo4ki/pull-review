INSERT INTO teams (team_name) VALUES
    ('backend'),
    ('payments')
ON CONFLICT DO NOTHING;

INSERT INTO users (user_id, username, team_name, is_active) VALUES
    ('u1', 'Egor', 'backend', true),
    ('u2', 'Katya',   'backend',  true),
    ('u3', 'Kirill', 'payments', true),
    ('u4', 'Anya',  'payments', false)
ON CONFLICT (user_id) DO UPDATE
SET username = EXCLUDED.username,
    team_name = EXCLUDED.team_name,
    is_active = EXCLUDED.is_active;