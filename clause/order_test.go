package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestOrderBy(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Order{Expr: "$-.age ASC, $-.name DESC"}},
			gqlWant: "ORDER BY $-.age ASC, $-.name DESC",
		},
		{
			clauses: []clause.Interface{clause.Order{Expr: "$var.dst DESC"}},
			gqlWant: "ORDER BY $var.dst DESC",
		},
		{
			clauses: []clause.Interface{clause.Order{Expr: "$-.age ASC, $-.name DESC"}, clause.Order{Expr: "$var.dst DESC"}},
			gqlWant: "ORDER BY $var.dst DESC",
		},
		{
			clauses: []clause.Interface{clause.Order{Expr: ""}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
