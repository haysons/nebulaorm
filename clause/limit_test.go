package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestLimit(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Limit{Limit: 1}},
			gqlWant: "LIMIT 1",
		},
		{
			clauses: []clause.Interface{clause.Limit{Limit: 5}, clause.Limit{Limit: 5, Offset: 3}, clause.Limit{Limit: 1}},
			gqlWant: "LIMIT 1",
		},
		{
			clauses: []clause.Interface{clause.Limit{Limit: 5, Offset: 3}},
			gqlWant: "LIMIT 3, 5",
		},
		{
			clauses: []clause.Interface{clause.Limit{Limit: 0, Offset: 3}},
			gqlWant: "LIMIT 3, 0",
		},
		{
			clauses: []clause.Interface{clause.Limit{Limit: -1, Offset: 3}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
