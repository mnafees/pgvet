package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestBanCharType(t *testing.T) {
	r := &rule.BanCharType{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags char(n) in CREATE TABLE",
			sql:   "CREATE TABLE t (code char(10))",
			wantN: 1,
		},
		{
			name:  "flags character(n) in CREATE TABLE",
			sql:   "CREATE TABLE t (code character(10))",
			wantN: 1,
		},
		{
			name:  "allows varchar in CREATE TABLE",
			sql:   "CREATE TABLE t (name varchar(100))",
			wantN: 0,
		},
		{
			name:  "allows text in CREATE TABLE",
			sql:   "CREATE TABLE t (name text)",
			wantN: 0,
		},
		{
			name:  "flags CAST to char",
			sql:   "SELECT x::char(5) FROM t",
			wantN: 1,
		},
		{
			name:  "allows CAST to text",
			sql:   "SELECT x::text FROM t",
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
