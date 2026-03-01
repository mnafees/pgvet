package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestSelectStar(t *testing.T) {
	r := &rule.SelectStar{}

	tests := []struct {
		name    string
		sql     string
		wantN   int
	}{
		{
			name:  "flags SELECT *",
			sql:   "SELECT * FROM users WHERE id = 1",
			wantN: 1,
		},
		{
			name:  "allows explicit columns",
			sql:   "SELECT id, name FROM users",
			wantN: 0,
		},
		{
			name:  "flags SELECT * with CTE at outer level",
			sql:   "WITH f AS (SELECT id FROM users) SELECT * FROM f",
			wantN: 1,
		},
		{
			name:  "allows SELECT * only inside CTE",
			sql:   "WITH f AS (SELECT * FROM users) SELECT id, name FROM f",
			wantN: 0,
		},
		{
			name:  "flags UNION with SELECT *",
			sql:   "SELECT * FROM users UNION ALL SELECT * FROM admins",
			wantN: 2,
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
