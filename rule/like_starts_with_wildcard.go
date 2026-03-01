package rule

import (
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type LikeStartsWithWildcard struct{}

func (r *LikeStartsWithWildcard) Name() string { return "like-starts-with-wildcard" }
func (r *LikeStartsWithWildcard) Description() string {
	return "LIKE/ILIKE pattern starting with % prevents index usage"
}

func (r *LikeStartsWithWildcard) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		ae := node.GetAExpr()
		if ae == nil {
			return true
		}
		if ae.Kind != pg_query.A_Expr_Kind_AEXPR_LIKE && ae.Kind != pg_query.A_Expr_Kind_AEXPR_ILIKE {
			return true
		}

		if ae.Rexpr == nil {
			return true
		}
		ac := ae.Rexpr.GetAConst()
		if ac == nil {
			return true
		}
		sv := ac.GetSval()
		if sv == nil {
			return true
		}
		if strings.HasPrefix(sv.Sval, "%") {
			line, col := offsetToLineCol(sql, int(ae.Location))
			diags = append(diags, Diagnostic{
				Rule:     r.Name(),
				Message:  r.Description(),
				Line:     line,
				Col:      col,
				Severity: SeverityWarning,
			})
		}
		return true
	})

	return diags
}
