package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestOrderByOrdinal(t *testing.T) {
	r := &rule.OrderByOrdinal{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags ORDER BY 1",
			sql:   "SELECT id, name FROM users ORDER BY 1",
			wantN: 1,
		},
		{
			name:  "flags ORDER BY 1, 2",
			sql:   "SELECT id, name FROM users ORDER BY 1, 2",
			wantN: 2,
		},
		{
			name:  "allows ORDER BY column name",
			sql:   "SELECT id, name FROM users ORDER BY id",
			wantN: 0,
		},
		{
			name:  "allows no ORDER BY",
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
