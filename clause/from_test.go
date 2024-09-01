package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestFrom(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.From{VID: "player100"}},
			gqlWant: `FROM "player100"`,
		},
		{
			clauses: []clause.Interface{clause.From{VID: []string{"player100", "player102"}}},
			gqlWant: `FROM "player100", "player102"`,
		},
		{
			clauses: []clause.Interface{clause.From{VID: []int{1, 2, 3}}},
			gqlWant: `FROM 1, 2, 3`,
		},
		{
			clauses: []clause.Interface{clause.From{VID: clause.Expr{Str: "$-.id"}}},
			gqlWant: `FROM $-.id`,
		},
		{
			clauses: []clause.Interface{clause.From{}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
