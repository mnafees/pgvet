package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestNotInSubquery(t *testing.T) {
	r := &rule.NotInSubquery{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags NOT IN (SELECT ...)",
			sql:   "SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM banned)",
			wantN: 1,
		},
		{
			name:  "allows IN (SELECT ...)",
			sql:   "SELECT * FROM users WHERE id IN (SELECT user_id FROM active)",
			wantN: 0,
		},
		{
			name:  "allows NOT IN with literal list",
			sql:   "SELECT * FROM users WHERE status NOT IN ('banned', 'deleted')",
			wantN: 0,
		},
		{
			name:  "no subquery at all",
			sql:   "SELECT id FROM users WHERE active = true",
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
