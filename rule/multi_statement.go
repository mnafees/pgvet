package rule

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v6"
)

type MultiStatement struct{}

func (r *MultiStatement) Name() string { return "multi-statement" }
func (r *MultiStatement) Description() string {
	return "Multiple statements in a single query block — CTEs from the first statement are not visible to subsequent ones"
}

// Check is special: it operates on the full parse result, not a single RawStmt.
// The analyzer calls CheckMulti instead.
func (r *MultiStatement) Check(stmt *pg_query.RawStmt, _ string) []Diagnostic {
	// Single-statement check is a no-op; see CheckMulti.
	_ = stmt
	return nil
}

// CheckMulti checks if a single SQL string parses into multiple top-level statements.
func (r *MultiStatement) CheckMulti(stmts []*pg_query.RawStmt, sql string) []Diagnostic {
	if len(stmts) <= 1 {
		return nil
	}

	// Flag the second statement onward.
	var diags []Diagnostic
	for i := 1; i < len(stmts); i++ {
		loc := int(stmts[i].StmtLocation)
		line, col := offsetToLineCol(sql, loc)
		diags = append(diags, Diagnostic{
			Rule:     r.Name(),
			Message:  fmt.Sprintf("Statement %d of %d in a single block — each statement should be separate", i+1, len(stmts)),
			Line:     line,
			Col:      col,
			Severity: SeverityError,
		})
	}
	return diags
}
