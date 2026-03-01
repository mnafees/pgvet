-- Should trigger: LIKE with leading %
SELECT * FROM users WHERE name LIKE '%test';

-- Should trigger: ILIKE with leading %
SELECT * FROM users WHERE name ILIKE '%test';

-- Should NOT trigger: LIKE with trailing % only
SELECT * FROM users WHERE name LIKE 'test%';

-- Should NOT trigger: normal comparison
SELECT * FROM users WHERE name = 'test';
