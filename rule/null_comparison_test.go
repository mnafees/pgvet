package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestNullComparison(t *testing.T) {
	r := &rule.NullComparison{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags = NULL",
			sql:   "SELECT * FROM users WHERE id = NULL",
			wantN: 1,
		},
		{
			name:  "flags <> NULL",
			sql:   "SELECT * FROM users WHERE id <> NULL",
			wantN: 1,
		},
		{
			name:  "flags NULL = col",
			sql:   "SELECT * FROM users WHERE NULL = id",
			wantN: 1,
		},
		{
			name:  "allows IS NULL",
			sql:   "SELECT * FROM users WHERE id IS NULL",
			wantN: 0,
		},
		{
			name:  "allows IS NOT NULL",
			sql:   "SELECT * FROM users WHERE id IS NOT NULL",
			wantN: 0,
		},
		{
			name:  "allows normal comparison",
			sql:   "SELECT * FROM users WHERE id = 1",
			wantN: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pg_query.Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			diags := r.Check(result.Stmts[0], tt.sql)
			if len(diags) != tt.wantN {
				t.Errorf("got %d diagnostics, want %d: %+v", len(diags), tt.wantN, diags)
			}
		})
	}
}
