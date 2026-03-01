-- Should trigger: char(n) pads with spaces
CREATE TABLE t (code char(10));

-- Should NOT trigger: varchar is fine
CREATE TABLE t2 (name varchar(100));

-- Should NOT trigger: text is fine
CREATE TABLE t3 (name text);

-- Should trigger: CAST to char
SELECT x::char(5) FROM t;
