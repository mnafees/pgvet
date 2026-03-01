-- Should trigger: ORDER BY ordinal position
SELECT id, name FROM users ORDER BY 1;

-- Should trigger: multiple ordinals
SELECT id, name FROM users ORDER BY 1, 2;

-- Should NOT trigger: ORDER BY column name
SELECT id, name FROM users ORDER BY id;

-- Should NOT trigger: no ORDER BY
SELECT id, name FROM users;
