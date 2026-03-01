package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestForUpdateNoSkip(t *testing.T) {
	r := &rule.ForUpdateNoSkip{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags FOR UPDATE without SKIP LOCKED",
			sql:   "SELECT id FROM items ORDER BY id FOR UPDATE",
			wantN: 1,
		},
		{
			name:  "allows FOR UPDATE SKIP LOCKED",
			sql:   "SELECT id FROM items ORDER BY id FOR UPDATE SKIP LOCKED",
			wantN: 0,
		},
		{
			name:  "allows FOR UPDATE NOWAIT",
			sql:   "SELECT id FROM items ORDER BY id FOR UPDATE NOWAIT",
			wantN: 0,
		},
		{
			name:  "allows FOR SHARE (not exclusive)",
			sql:   "SELECT id FROM items WHERE id = 1 FOR SHARE",
			wantN: 0,
		},
		{
			name:  "no locking clause",
			sql:   "SELECT id FROM items",
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
