-- Should trigger: LIMIT without ORDER BY
SELECT id, name FROM users LIMIT 10;

-- Should NOT trigger: LIMIT 1 (existence check exemption)
SELECT id FROM users WHERE email = 'test@test.com' LIMIT 1;

-- Should NOT trigger: LIMIT with ORDER BY
SELECT id, name FROM users ORDER BY id ASC LIMIT 10;

-- Should trigger: LIMIT 5 without ORDER BY
SELECT * FROM events WHERE type = 'click' LIMIT 5;
