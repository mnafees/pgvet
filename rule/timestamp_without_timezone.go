package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/walker"
)

type TimestampWithoutTimezone struct{}

func (r *TimestampWithoutTimezone) Name() string { return "timestamp-without-timezone" }
func (r *TimestampWithoutTimezone) Description() string {
	return "timestamp without time zone loses timezone context — use timestamptz instead"
}

func (r *TimestampWithoutTimezone) Check(stmt *pg_query.RawStmt, sql string) []Diagnostic {
	var diags []Diagnostic

	// Check CREATE TABLE columns.
	if cs := stmt.Stmt.GetCreateStmt(); cs != nil {
		for _, elt := range cs.TableElts {
			if cd := elt.GetColumnDef(); cd != nil {
				if isTimestampWithoutTZ(cd.TypeName) {
					line, col := offsetToLineCol(sql, int(cd.Location))
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
	}

	// Check ALTER TABLE ADD COLUMN.
	if as := stmt.Stmt.GetAlterTableStmt(); as != nil {
		for _, cmd := range as.Cmds {
			ac := cmd.GetAlterTableCmd()
			if ac == nil {
				continue
			}
			if ac.Subtype != pg_query.AlterTableType_AT_AddColumn {
				continue
			}
			if ac.Def == nil {
				continue
			}
			if cd := ac.Def.GetColumnDef(); cd != nil {
				if isTimestampWithoutTZ(cd.TypeName) {
					line, col := offsetToLineCol(sql, int(cd.Location))
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
	}

	// Check CAST / :: in queries.
	walker.Walk(stmt.Stmt, func(node *pg_query.Node) bool {
		tc := node.GetTypeCast()
		if tc == nil {
			return true
		}
		if isTimestampWithoutTZ(tc.TypeName) {
			line, col := offsetToLineCol(sql, int(tc.Location))
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

func isTimestampWithoutTZ(tn *pg_query.TypeName) bool {
	if tn == nil || len(tn.Names) == 0 {
		return false
	}
	last := tn.Names[len(tn.Names)-1]
	if s := last.GetString_(); s != nil {
		return s.Sval == "timestamp"
	}
	return false
}
