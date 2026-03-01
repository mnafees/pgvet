-- Should trigger: OFFSET without LIMIT
SELECT * FROM users OFFSET 10;

-- Should NOT trigger: OFFSET with LIMIT
SELECT * FROM users LIMIT 10 OFFSET 5;

-- Should NOT trigger: LIMIT without OFFSET
SELECT * FROM users LIMIT 10;

-- Should NOT trigger: no OFFSET or LIMIT
SELECT * FROM users;
