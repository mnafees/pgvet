package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestDistinctOnOrder(t *testing.T) {
	r := &rule.DistinctOnOrder{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags DISTINCT ON without ORDER BY",
			sql:   "SELECT DISTINCT ON (dept_id) dept_id, name FROM employees",
			wantN: 1,
		},
		{
			name:  "allows DISTINCT ON with matching ORDER BY",
			sql:   "SELECT DISTINCT ON (dept_id) dept_id, name FROM employees ORDER BY dept_id, salary DESC",
			wantN: 0,
		},
		{
			name:  "flags DISTINCT ON with non-matching ORDER BY",
			sql:   "SELECT DISTINCT ON (dept_id) dept_id, name FROM employees ORDER BY name",
			wantN: 1,
		},
		{
			name:  "allows plain DISTINCT",
			sql:   "SELECT DISTINCT status FROM orders",
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
