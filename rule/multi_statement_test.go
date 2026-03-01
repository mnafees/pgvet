package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestMultiStatement(t *testing.T) {
	r := &rule.MultiStatement{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags two statements",
			sql:   "DELETE FROM a WHERE id = 1; DELETE FROM b WHERE id = 1;",
			wantN: 1,
		},
		{
			name:  "flags three statements",
			sql:   "SELECT 1; SELECT 2; SELECT 3;",
			wantN: 2,
		},
		{
			name:  "allows single statement",
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
			diags := r.CheckMulti(result.Stmts, tt.sql)
			if len(diags) != tt.wantN {
				t.Errorf("got %d diagnostics, want %d: %+v", len(diags), tt.wantN, diags)
			}
		})
	}
}
