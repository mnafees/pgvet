-- Should trigger: DISTINCT ON without ORDER BY
SELECT DISTINCT ON (department_id) department_id, name
FROM employees;

-- Should NOT trigger: DISTINCT ON with matching ORDER BY
SELECT DISTINCT ON (department_id) department_id, name, salary
FROM employees
ORDER BY department_id, salary DESC;

-- Should trigger: DISTINCT ON with non-matching ORDER BY
SELECT DISTINCT ON (department_id) department_id, name
FROM employees
ORDER BY name;

-- Should NOT trigger: plain DISTINCT (not DISTINCT ON)
SELECT DISTINCT status FROM orders;
