package rule_test

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

func TestUpdateWithoutWhere(t *testing.T) {
	r := &rule.UpdateWithoutWhere{}

	tests := []struct {
		name  string
		sql   string
		wantN int
	}{
		{
			name:  "flags UPDATE without WHERE",
			sql:   "UPDATE users SET active = false",
			wantN: 1,
		},
		{
			name:  "allows UPDATE with WHERE",
			sql:   "UPDATE users SET active = false WHERE id = 1",
			wantN: 0,
		},
		{
			name:  "allows UPDATE with FROM and WHERE",
			sql:   "UPDATE orders SET status = 'shipped' FROM shipments WHERE orders.id = shipments.order_id",
			wantN: 0,
		},
		{
			name:  "ignores non-UPDATE",
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
