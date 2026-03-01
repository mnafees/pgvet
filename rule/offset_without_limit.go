package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type OffsetWithoutLimit struct{}

func (r *OffsetWithoutLimit) Name() string { return "offset-without-limit" }
func (r *OffsetWithoutLimit) Description() string {
	return "OFFSET without LIMIT returns all remaining rows — likely a mistake"
}

func (r *OffsetWithoutLimit) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		sel := node.GetSelectStmt()
		if sel == nil {
			return true
		}
		if sel.LimitOffset == nil {
			return true
		}
		if sel.LimitCount != nil {
			return true
		}

		line, col := offsetToLineCol(sql, int(locOf(sel.LimitOffset)))
		diags = append(diags, Diagnostic{
			Rule:     r.Name(),
			Message:  r.Description(),
			Line:     line,
			Col:      col,
			Severity: SeverityWarning,
		})
		return true
	})

	return diags
}
