-- Should trigger: SELECT * at outermost level
SELECT * FROM users WHERE id = 1;

-- Should NOT trigger: SELECT * inside CTE, explicit columns at outer level
WITH all_users AS (SELECT * FROM users)
SELECT id, name FROM all_users;

-- Should trigger: SELECT * at outer level even with CTE
WITH filtered AS (SELECT id FROM users WHERE active = true)
SELECT * FROM filtered;

-- Should NOT trigger: explicit column list
SELECT id, name, email FROM users;
