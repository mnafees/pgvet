# pgvet

A static analysis tool for PostgreSQL SQL files, powered by the real PostgreSQL parser via [pg_query_go](https://github.com/pganalyze/pg_query_go).

pgvet parses your `.sql` files using the same parser that runs inside PostgreSQL itself and checks for common anti-patterns and correctness issues — no running database required.

## Install

```bash
go install github.com/mnafees/pgvet@latest
```

Note: The first build takes ~3 minutes due to CGO compilation of the embedded PostgreSQL parser. Subsequent builds are fast.

## Usage

```bash
# Check a file
pgvet queries.sql

# Check a directory recursively
pgvet sql/

# Check multiple paths
pgvet queries/ migrations/ views.sql

# Read from stdin
echo "SELECT * FROM users" | pgvet

# JSON output (for CI integration)
pgvet --format json sql/

# Run only specific rules
pgvet --rules not-in-subquery,select-star sql/

# Exclude specific rules
pgvet --exclude select-star sql/
```

Exit codes: `0` = no issues, `1` = issues found, `2` = usage/parse error.

## Rules

### Default rules

These rules run by default:

| Rule | Severity | Description |
|------|----------|-------------|
| `select-star` | warning | `SELECT *` in the outermost query is fragile — list columns explicitly |
| `limit-without-order` | warning | `LIMIT` without `ORDER BY` produces non-deterministic results (exempts `LIMIT 1`) |
| `not-in-subquery` | error | `NOT IN (SELECT ...)` is broken when the subquery can return NULLs — use `NOT EXISTS` instead |
| `for-update-no-skip` | warning | `FOR UPDATE` without `SKIP LOCKED` or `NOWAIT` can cause lock contention |
| `distinct-on-order` | warning | `DISTINCT ON` without a matching leading `ORDER BY` produces non-deterministic results |
| `null-comparison` | error | `= NULL` or `<> NULL` always yields NULL — use `IS NULL` or `IS NOT NULL` |
| `update-without-where` | warning | `UPDATE` without `WHERE` updates every row in the table |
| `delete-without-where` | warning | `DELETE` without `WHERE` deletes every row in the table |
| `insert-without-columns` | warning | `INSERT` without column list depends on column order — list columns explicitly |
| `ban-char-type` | warning | `char(n)` pads with spaces — use `text` or `varchar` instead |
| `timestamp-without-timezone` | warning | `timestamp` without time zone loses timezone context — use `timestamptz` instead |
| `order-by-ordinal` | warning | `ORDER BY` ordinal position is fragile — use column names or expressions |
| `group-by-ordinal` | warning | `GROUP BY` ordinal position is fragile — use column names or expressions |
| `like-starts-with-wildcard` | warning | `LIKE`/`ILIKE` pattern starting with `%` prevents index usage |
| `offset-without-limit` | warning | `OFFSET` without `LIMIT` returns all remaining rows — likely a mistake |

### Opt-in rules

These rules must be explicitly enabled with `--rules`:

| Rule | Severity | Description |
|------|----------|-------------|
| `multi-statement` | error | Multiple statements in a single query block — CTEs from the first statement are not visible to subsequent ones |

## Output formats

### Text (default)

```
queries.sql:3:8: warning: [select-star] SELECT * in outermost query is fragile — list columns explicitly
queries.sql:7:30: error: [not-in-subquery] NOT IN (SELECT ...) is broken when the subquery can return NULLs — use NOT EXISTS instead
```

### JSON

```json
[
  {
    "rule": "select-star",
    "message": "SELECT * in outermost query is fragile — list columns explicitly",
    "file": "queries.sql",
    "line": 3,
    "col": 8,
    "severity": "warning"
  }
]
```

## Writing custom rules

pgvet has a simple rule interface:

```go
type Rule interface {
    Name() string
    Description() string
    Check(stmt *pg_query.RawStmt, sql string) []Diagnostic
}
```

Each rule receives a single parsed statement and the original SQL text. The `walker` package provides a generic AST traversal helper so rules don't need their own recursion logic.

## License

MIT
