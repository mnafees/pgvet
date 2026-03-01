-- Should trigger: DELETE without WHERE deletes all rows
DELETE FROM users;

-- Should NOT trigger: DELETE with WHERE
DELETE FROM users WHERE id = 1;

-- Should NOT trigger: DELETE with USING and WHERE
DELETE FROM orders USING expired WHERE orders.id = expired.id;
