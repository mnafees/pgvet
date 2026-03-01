package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestOffsetWithoutLimit(t *testing.T) {
	r := &rule.OffsetWithoutLimit{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags OFFSET without LIMIT",
			sql:   "SELECT * FROM users OFFSET 10",
			wantN: 1,
		},
		{
			name:  "allows OFFSET with LIMIT",
			sql:   "SELECT * FROM users LIMIT 10 OFFSET 5",
			wantN: 0,
		},
		{
			name:  "allows LIMIT without OFFSET",
			sql:   "SELECT * FROM users LIMIT 10",
			wantN: 0,
		},
		{
			name:  "allows no OFFSET or LIMIT",
			sql:   "SELECT * FROM users",
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
