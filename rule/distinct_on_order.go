package rule

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type DistinctOnOrder struct{}

func (r *DistinctOnOrder) Name() string { return "distinct-on-order" }
func (r *DistinctOnOrder) Description() string {
	return "DISTINCT ON without a matching leading ORDER BY produces non-deterministic results"
}

func (r *DistinctOnOrder) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		sel := node.GetSelectStmt()
		if sel == nil {
			return true
		}

		if len(sel.DistinctClause) == 0 {
			return true
		}

		// A plain DISTINCT (no ON) has an empty list but DistinctClause is non-nil.
		// DISTINCT ON has actual expression nodes. Check if these are real expressions
		// vs. just a marker for plain DISTINCT.
		hasDistinctOn := false
		for _, d := range sel.DistinctClause {
			if d.GetColumnRef() != nil || d.GetFuncCall() != nil || d.GetAExpr() != nil || d.GetAConst() != nil {
				hasDistinctOn = true
				break
			}
		}
		if !hasDistinctOn {
			return true
		}

		// No ORDER BY at all — always flag.
		if len(sel.SortClause) == 0 {
			loc := distinctOnLocation(sel)
			line, col := offsetToLineCol(sql, loc)
			diags = append(diags, Diagnostic{
				Rule:     r.Name(),
				Message:  "DISTINCT ON without any ORDER BY — row selection is arbitrary",
				Line:     line,
				Col:      col,
				Severity: SeverityWarning,
			})
			return true
		}

		// Check that leading ORDER BY columns match DISTINCT ON columns.
		for i, distNode := range sel.DistinctClause {
			if i >= len(sel.SortClause) {
				loc := distinctOnLocation(sel)
				line, col := offsetToLineCol(sql, loc)
				diags = append(diags, Diagnostic{
					Rule:     r.Name(),
					Message:  fmt.Sprintf("DISTINCT ON has %d columns but ORDER BY only has %d leading columns", len(sel.DistinctClause), len(sel.SortClause)),
					Line:     line,
					Col:      col,
					Severity: SeverityWarning,
				})
				break
			}

			sortNode := sel.SortClause[i].GetSortBy()
			if sortNode == nil {
				continue
			}

			if !nodesEqual(distNode, sortNode.Node) {
				loc := distinctOnLocation(sel)
				line, col := offsetToLineCol(sql, loc)
				diags = append(diags, Diagnostic{
					Rule:     r.Name(),
					Message:  fmt.Sprintf("DISTINCT ON column %d does not match ORDER BY column %d", i+1, i+1),
					Line:     line,
					Col:      col,
					Severity: SeverityWarning,
				})
				break
			}
		}

		return true
	})

	return diags
}

func distinctOnLocation(sel *pg_query.SelectStmt) int {
	if len(sel.DistinctClause) > 0 {
		if cr := sel.DistinctClause[0].GetColumnRef(); cr != nil {
			return int(cr.Location)
		}
	}
	return 0
}

// nodesEqual compares two nodes for structural equality by checking column references.
func nodesEqual(a, b *pg_query.Node) bool {
	if a == nil || b == nil {
		return a == b
	}

	// Compare column references.
	crA := a.GetColumnRef()
	crB := b.GetColumnRef()
	if crA != nil && crB != nil {
		if len(crA.Fields) != len(crB.Fields) {
			return false
		}
		for i := range crA.Fields {
			sA := crA.Fields[i].GetString_()
			sB := crB.Fields[i].GetString_()
			if sA == nil || sB == nil {
				return sA == nil && sB == nil
			}
			if sA.Sval != sB.Sval {
				return false
			}
		}
		return true
	}

	return false
}
