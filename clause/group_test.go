package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestGroupBy(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Group{Expr: "$-.Name"}},
			gqlWant: "GROUP BY $-.Name",
		},
		{
			clauses: []clause.Interface{clause.Group{Expr: "$-.player"}},
			gqlWant: "GROUP BY $-.player",
		},
		{
			clauses: []clause.Interface{clause.Group{Expr: "age"}, clause.Group{Expr: "name"}},
			gqlWant: "GROUP BY name",
		},
		{
			clauses: []clause.Interface{clause.Group{}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
