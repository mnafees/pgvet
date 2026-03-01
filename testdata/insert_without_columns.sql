-- Should trigger: INSERT without column list
INSERT INTO users VALUES (1, 'alice');

-- Should NOT trigger: INSERT with column list
INSERT INTO users (id, name) VALUES (1, 'alice');

-- Should trigger: INSERT SELECT without column list
INSERT INTO users SELECT * FROM temp;

-- Should NOT trigger: DEFAULT VALUES
INSERT INTO users DEFAULT VALUES;
