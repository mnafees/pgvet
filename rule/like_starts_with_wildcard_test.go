package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestLikeStartsWithWildcard(t *testing.T) {
	r := &rule.LikeStartsWithWildcard{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags LIKE with leading %",
			sql:   "SELECT * FROM users WHERE name LIKE '%test'",
			wantN: 1,
		},
		{
			name:  "flags ILIKE with leading %",
			sql:   "SELECT * FROM users WHERE name ILIKE '%test'",
			wantN: 1,
		},
		{
			name:  "flags LIKE with leading % and trailing %",
			sql:   "SELECT * FROM users WHERE name LIKE '%test%'",
			wantN: 1,
		},
		{
			name:  "allows LIKE with trailing % only",
			sql:   "SELECT * FROM users WHERE name LIKE 'test%'",
			wantN: 0,
		},
		{
			name:  "allows exact LIKE",
			sql:   "SELECT * FROM users WHERE name LIKE 'test'",
			wantN: 0,
		},
		{
			name:  "allows non-LIKE comparison",
			sql:   "SELECT * FROM users WHERE name = 'test'",
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
