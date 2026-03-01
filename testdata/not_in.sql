-- Should trigger: NOT IN (SELECT ...)
SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM banned_users);

-- Should NOT trigger: regular IN (SELECT ...)
SELECT * FROM users WHERE id IN (SELECT user_id FROM active_users);

-- Should NOT trigger: NOT IN with literal list
SELECT * FROM users WHERE status NOT IN ('banned', 'deleted');

-- Should trigger: NOT IN subquery in complex expression
SELECT u.id FROM users u
WHERE u.department_id NOT IN (SELECT id FROM departments WHERE active = false);
