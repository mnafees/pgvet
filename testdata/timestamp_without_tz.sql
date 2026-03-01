-- Should trigger: timestamp without time zone
CREATE TABLE t (created_at timestamp);

-- Should NOT trigger: timestamptz is correct
CREATE TABLE t2 (created_at timestamptz);

-- Should NOT trigger: timestamp with time zone is correct
CREATE TABLE t3 (created_at timestamp with time zone);

-- Should trigger: CAST to timestamp without time zone
SELECT now()::timestamp;
