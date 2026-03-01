package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestTimestampWithoutTimezone(t *testing.T) {
	r := &rule.TimestampWithoutTimezone{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags timestamp in CREATE TABLE",
			sql:   "CREATE TABLE t (created_at timestamp)",
			wantN: 1,
		},
		{
			name:  "flags timestamp without time zone",
			sql:   "CREATE TABLE t (created_at timestamp without time zone)",
			wantN: 1,
		},
		{
			name:  "allows timestamptz in CREATE TABLE",
			sql:   "CREATE TABLE t (created_at timestamptz)",
			wantN: 0,
		},
		{
			name:  "allows timestamp with time zone",
			sql:   "CREATE TABLE t (created_at timestamp with time zone)",
			wantN: 0,
		},
		{
			name:  "flags CAST to timestamp",
			sql:   "SELECT now()::timestamp",
			wantN: 1,
		},
		{
			name:  "allows CAST to timestamptz",
			sql:   "SELECT now()::timestamptz",
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
