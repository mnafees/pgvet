package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestLimitWithoutOrder(t *testing.T) {
	r := &rule.LimitWithoutOrder{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags LIMIT without ORDER BY",
			sql:   "SELECT id FROM users LIMIT 10",
			wantN: 1,
		},
		{
			name:  "exempts LIMIT 1",
			sql:   "SELECT id FROM users WHERE email = 'x' LIMIT 1",
			wantN: 0,
		},
		{
			name:  "allows LIMIT with ORDER BY",
			sql:   "SELECT id FROM users ORDER BY id LIMIT 10",
			wantN: 0,
		},
		{
			name:  "no LIMIT at all",
			sql:   "SELECT id FROM users",
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
