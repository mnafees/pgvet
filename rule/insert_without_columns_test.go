package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestInsertWithoutColumns(t *testing.T) {
	r := &rule.InsertWithoutColumns{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags INSERT VALUES without columns",
			sql:   "INSERT INTO users VALUES (1, 'alice')",
			wantN: 1,
		},
		{
			name:  "allows INSERT VALUES with columns",
			sql:   "INSERT INTO users (id, name) VALUES (1, 'alice')",
			wantN: 0,
		},
		{
			name:  "flags INSERT SELECT without columns",
			sql:   "INSERT INTO users SELECT * FROM temp",
			wantN: 1,
		},
		{
			name:  "allows INSERT SELECT with columns",
			sql:   "INSERT INTO users (id, name) SELECT id, name FROM temp",
			wantN: 0,
		},
		{
			name:  "allows DEFAULT VALUES",
			sql:   "INSERT INTO users DEFAULT VALUES",
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
