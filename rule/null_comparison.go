package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type NullComparison struct{}

func (r *NullComparison) Name() string { return "null-comparison" }
func (r *NullComparison) Description() string {
	return "= NULL or <> NULL always yields NULL — use IS NULL or IS NOT NULL"
}

func (r *NullComparison) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		ae := node.GetAExpr()
		if ae == nil {
			return true
		}
		if ae.Kind != pg_query.A_Expr_Kind_AEXPR_OP {
			return true
		}

		opName := ""
		for _, n := range ae.Name {
			if s := n.GetString_(); s != nil {
				opName = s.Sval
			}
		}
		if opName != "=" && opName != "<>" {
			return true
		}

		if isNullConst(ae.Lexpr) || isNullConst(ae.Rexpr) {
			line, col := offsetToLineCol(sql, int(ae.Location))
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

func isNullConst(node *pg_query.Node) bool {
	if node == nil {
		return false
	}
	ac := node.GetAConst()
	return ac != nil && ac.Isnull
}
