package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
)

type InsertWithoutColumns struct{}

func (r *InsertWithoutColumns) Name() string { return "insert-without-columns" }
func (r *InsertWithoutColumns) Description() string {
	return "INSERT without column list depends on column order — list columns explicitly"
}

func (r *InsertWithoutColumns) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	ins := stmt.Stmt.GetInsertStmt()
	if ins == nil {
		return nil
	}
	if len(ins.Cols) > 0 {
		return nil
	}
	// DEFAULT VALUES has no SelectStmt — that's fine.
	if ins.SelectStmt == nil {
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
