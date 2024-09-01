package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Lookup{TypeName: "player"}},
			gqlWant: "LOOKUP ON player",
		},
		{
			clauses: []clause.Interface{clause.Lookup{TypeName: "follow"}},
			gqlWant: "LOOKUP ON follow",
		},
		{
			clauses: []clause.Interface{clause.Lookup{TypeName: "player"}, clause.Lookup{TypeName: "follow"}},
			gqlWant: "LOOKUP ON follow",
		},
		{
			clauses: []clause.Interface{clause.Lookup{TypeName: ""}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
