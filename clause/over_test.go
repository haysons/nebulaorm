package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestOver(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Over{EdgeTypeList: []string{"serve"}}},
			gqlWant: "OVER serve",
		},
		{
			clauses: []clause.Interface{clause.Over{EdgeTypeList: []string{"serve", "follow"}, Direction: clause.OverDirectBidirect}, clause.Over{EdgeTypeList: []string{"teach"}, Direction: clause.OverDirectReversely}},
			gqlWant: "OVER serve, follow, teach REVERSELY",
		},
		{
			clauses: []clause.Interface{clause.Over{EdgeTypeList: []string{""}, Direction: clause.OverDirectBidirect}, clause.Over{EdgeTypeList: []string{"serve", "follow"}}},
			gqlWant: "OVER serve, follow BIDIRECT",
		},
		{
			clauses: []clause.Interface{clause.Over{EdgeTypeList: []string{""}}, clause.Over{EdgeTypeList: []string{""}}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.Over{EdgeTypeList: []string{"serve"}, Direction: "test"}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
