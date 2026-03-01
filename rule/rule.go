package rule

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
)

// Severity levels for diagnostics.
const (
	SeverityWarning = "warning"
	SeverityError   = "error"
)

// Diagnostic represents a single issue found by a rule.
type Diagnostic struct {
	Rule     string `json:"rule"`
	Message  string `json:"message"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Col      int    `json:"col"`
	Severity string `json:"severity"`
}

// Rule is the interface that all pgvet rules must implement.
type Rule interface {
	Name() string
	Description() string
	Check(stmt *pg_query.RawStmt, sql string) []Diagnostic
}
