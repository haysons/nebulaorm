package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestYield(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Yield{ExprList: []string{"properties(vertex).name, properties(vertex).age"}}},
			gqlWant: "YIELD properties(vertex).name, properties(vertex).age",
		},
		{
			clauses: []clause.Interface{clause.Yield{ExprList: []string{"AVG($-.Age) as Avg_age"}}, clause.Yield{ExprList: []string{"count(*) as Num_friends"}}},
			gqlWant: "YIELD AVG($-.Age) as Avg_age, count(*) as Num_friends",
		},
		{
			clauses: []clause.Interface{clause.Yield{ExprList: []string{"properties(vertex).age as v"}}, clause.Yield{Distinct: true}},
			gqlWant: "YIELD DISTINCT properties(vertex).age as v",
		},
		{
			clauses: []clause.Interface{clause.Yield{ExprList: []string{"properties(vertex).age as v"}}, clause.Yield{Distinct: true}},
			gqlWant: "YIELD DISTINCT properties(vertex).age as v",
		},
		{
			clauses: []clause.Interface{clause.Yield{}, clause.Yield{Distinct: true}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
