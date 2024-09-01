package statement

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestStatement(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				stmt := New()
				stmt.SetClausesBuild([]string{clause.FromName, clause.OrderName})
				stmt.AddClause(&clause.From{VID: "team1"})
				stmt.AddClause(&clause.Order{Expr: "$-.id"})
				return stmt
			},
			want: `FROM "team1" ORDER BY $-.id;`,
		},
		{
			stmt: func() *Statement {
				stmt := New()
				stmt.SetClausesBuild([]string{clause.GoName, clause.OrderName})
				stmt.SetPartType(PartTypeFetch)
				stmt.AddClause(&clause.Go{StepStart: 1, StepEnd: 2})
				stmt.AddClause(&clause.Order{Expr: "$-.id"})
				return stmt
			},
			want: `GO 1 TO 2 STEPS ORDER BY $-.id;`,
		},
		{
			stmt: func() *Statement {
				stmt := New()
				stmt.Go(1).Pipe().Pipe().GroupBy("$-.id")
				return stmt
			},
			want: `GO 1 STEPS | GROUP BY $-.id;`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got an unexpected error: %v", err)
				}
				return
			}
			if ngql != tt.want {
				t.Errorf("NGQL = %v, want %v", ngql, tt.want)
			}
		})
	}
}
