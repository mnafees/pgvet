package rule

import (
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type ForUpdateNoSkip struct{}

func (r *ForUpdateNoSkip) Name() string { return "for-update-no-skip" }
func (r *ForUpdateNoSkip) Description() string {
	return "FOR UPDATE without SKIP LOCKED or NOWAIT can cause lock contention"
}

func (r *ForUpdateNoSkip) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		sel := node.GetSelectStmt()
		if sel == nil {
			return true
		}

		for _, lockNode := range sel.LockingClause {
			lc := lockNode.GetLockingClause()
			if lc == nil {
				continue
			}

			// Only flag FOR UPDATE and FOR NO KEY UPDATE (the exclusive lock modes).
			if lc.Strength != pg_query.LockClauseStrength_LCS_FORUPDATE &&
				lc.Strength != pg_query.LockClauseStrength_LCS_FORNOKEYUPDATE {
				continue
			}

			// Check if SKIP LOCKED or NOWAIT is specified.
			if lc.WaitPolicy == pg_query.LockWaitPolicy_LockWaitSkip ||
				lc.WaitPolicy == pg_query.LockWaitPolicy_LockWaitError {
				continue
			}

			// LockingClause has no Location field; find "FOR UPDATE" in the SQL.
			loc := strings.Index(strings.ToUpper(sql), "FOR UPDATE")
			if loc < 0 {
				loc = 0
			}
			line, col := offsetToLineCol(sql, loc)
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
