package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type NotInSubquery struct{}

func (r *NotInSubquery) Name() string { return "not-in-subquery" }
func (r *NotInSubquery) Description() string {
	return "NOT IN (SELECT ...) is broken when the subquery can return NULLs — use NOT EXISTS instead"
}

func (r *NotInSubquery) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		// PostgreSQL parses "NOT IN (SELECT ...)" as:
		//   BoolExpr(NOT_EXPR, args: [SubLink(ANY_SUBLINK, ...)])
		be := node.GetBoolExpr()
		if be == nil {
			return true
		}
		if be.Boolop != pg_query.BoolExprType_NOT_EXPR {
			return true
		}

		for _, arg := range be.Args {
			sl := arg.GetSubLink()
			if sl == nil {
				continue
			}
			if sl.SubLinkType != pg_query.SubLinkType_ANY_SUBLINK {
				continue
			}

			line, col := offsetToLineCol(sql, int(be.Location))
			diags = append(diags, Diagnostic{
				Rule:     r.Name(),
				Message:  r.Description(),
				Line:     line,
				Col:      col,
				Severity: SeverityError,
			})
		}

		return true
	})

	return diags
}
