package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestDeleteWithoutWhere(t *testing.T) {
	r := &rule.DeleteWithoutWhere{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags DELETE without WHERE",
			sql:   "DELETE FROM users",
			wantN: 1,
		},
		{
			name:  "allows DELETE with WHERE",
			sql:   "DELETE FROM users WHERE id = 1",
			wantN: 0,
		},
		{
			name:  "allows DELETE with USING and WHERE",
			sql:   "DELETE FROM orders USING expired WHERE orders.id = expired.id",
			wantN: 0,
		},
		{
			name:  "ignores non-DELETE",
			sql:   "SELECT 1",
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
