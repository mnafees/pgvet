package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type OrderByOrdinal struct{}

func (r *OrderByOrdinal) Name() string { return "order-by-ordinal" }
func (r *OrderByOrdinal) Description() string {
	return "ORDER BY ordinal position is fragile — use column names or expressions"
}

func (r *OrderByOrdinal) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		sel := node.GetSelectStmt()
		if sel == nil {
			return true
		}

		for _, sc := range sel.SortClause {
			sb := sc.GetSortBy()
			if sb == nil || sb.Node == nil {
				continue
			}
			ac := sb.Node.GetAConst()
			if ac == nil {
				continue
			}
			if ac.GetIval() != nil {
				line, col := offsetToLineCol(sql, int(ac.Location))
				diags = append(diags, Diagnostic{
					Rule:     r.Name(),
					Message:  r.Description(),
					Line:     line,
					Col:      col,
					Severity: SeverityWarning,
				})
			}
		}
		return true
	})

	return diags
}
