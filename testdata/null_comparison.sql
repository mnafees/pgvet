-- Should trigger: comparing with = NULL always yields NULL
SELECT * FROM users WHERE id = NULL;

-- Should trigger: comparing with <> NULL always yields NULL
SELECT * FROM users WHERE status <> NULL;

-- Should NOT trigger: IS NULL is the correct way
SELECT * FROM users WHERE id IS NULL;

-- Should NOT trigger: IS NOT NULL is the correct way
SELECT * FROM users WHERE id IS NOT NULL;
