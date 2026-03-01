package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
)

type UpdateWithoutWhere struct{}

func (r *UpdateWithoutWhere) Name() string { return "update-without-where" }
func (r *UpdateWithoutWhere) Description() string {
	return "UPDATE without WHERE updates every row in the table"
}

func (r *UpdateWithoutWhere) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	u := stmt.Stmt.GetUpdateStmt()
	if u == nil {
		return nil
	}
	if u.WhereClause != nil {
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
