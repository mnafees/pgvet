package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
)

type SelectStar struct{}

func (r *SelectStar) Name() string        { return "select-star" }
func (r *SelectStar) Description() string { return "SELECT * in outermost query is fragile — list columns explicitly" }

func (r *SelectStar) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	sel := stmt.Stmt.GetSelectStmt()
	if sel == nil {
		return nil
	}
	return r.checkOutermost(sel, sql, stmt.StmtLocation)
}

func (r *SelectStar) checkOutermost(sel *pg_query.SelectStmt, sql string, baseOffset int32) []Diagnostic {
	// For UNION/INTERSECT/EXCEPT, check both sides.
	if sel.Op != pg_query.SetOperation_SETOP_NONE {
		var diags []Diagnostic
		if sel.Larg != nil {
			diags = append(diags, r.checkOutermost(sel.Larg, sql, baseOffset)...)
		}
		if sel.Rarg != nil {
			diags = append(diags, r.checkOutermost(sel.Rarg, sql, baseOffset)...)
		}
		return diags
	}

	var diags []Diagnostic
	for _, target := range sel.TargetList {
		rt := target.GetResTarget()
		if rt == nil {
			continue
		}
		cr := rt.Val.GetColumnRef()
		if cr == nil {
			continue
		}
		for _, field := range cr.Fields {
			if field.GetAStar() != nil {
				line, col := offsetToLineCol(sql, int(cr.Location))
				diags = append(diags, Diagnostic{
					Rule:     r.Name(),
					Message:  r.Description(),
					Line:     line,
					Col:      col,
					Severity: SeverityWarning,
				})
			}
		}
	}
	return diags
}
