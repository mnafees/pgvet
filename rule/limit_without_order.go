package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type LimitWithoutOrder struct{}

func (r *LimitWithoutOrder) Name() string        { return "limit-without-order" }
func (r *LimitWithoutOrder) Description() string { return "LIMIT without ORDER BY produces non-deterministic results" }

func (r *LimitWithoutOrder) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		sel := node.GetSelectStmt()
		if sel == nil {
			return true
		}

		if sel.LimitCount == nil {
			return true
		}
		if len(sel.SortClause) > 0 {
			return true
		}

		// Exempt LIMIT 1 — commonly used for existence checks.
		if isLimitOne(sel.LimitCount) {
			return true
		}

		loc := sel.LimitCount
		line, col := offsetToLineCol(sql, int(locOf(loc)))
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

func isLimitOne(node *pg_query.Node) bool {
	c := node.GetAConst()
	if c == nil {
		return false
	}
	ival := c.GetIval()
	return ival != nil && ival.Ival == 1
}

func locOf(node *pg_query.Node) int32 {
	switch {
	case node.GetAConst() != nil:
		return node.GetAConst().Location
	case node.GetTypeCast() != nil:
		return node.GetTypeCast().Location
	case node.GetColumnRef() != nil:
		return node.GetColumnRef().Location
	case node.GetFuncCall() != nil:
		return node.GetFuncCall().Location
	default:
		return 0
	}
}
