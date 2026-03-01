-- Should trigger: GROUP BY ordinal position
SELECT status, count(*) FROM users GROUP BY 1;

-- Should trigger: multiple ordinals
SELECT a, b, count(*) FROM t GROUP BY 1, 2;

-- Should NOT trigger: GROUP BY column name
SELECT status, count(*) FROM users GROUP BY status;

-- Should NOT trigger: no GROUP BY
SELECT id FROM users;
