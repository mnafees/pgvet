package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
)

type DeleteWithoutWhere struct{}

func (r *DeleteWithoutWhere) Name() string { return "delete-without-where" }
func (r *DeleteWithoutWhere) Description() string {
	return "DELETE without WHERE deletes every row in the table"
}

func (r *DeleteWithoutWhere) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	d := stmt.Stmt.GetDeleteStmt()
	if d == nil {
		return nil
	}
	if d.WhereClause != nil {
		return nil
	}

	line, col := offsetToLineCol(sql, int(stmt.StmtLocation))
	return []Diagnostic{{
		Rule:     r.Name(),
		Message:  r.Description(),
		Line:     line,
		Col:      col,
		Severity: SeverityWarning,
	}}
}
