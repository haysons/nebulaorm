package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestFetch(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Fetch{Names: []string{"player"}, VID: "player100"}},
			gqlWant: `FETCH PROP ON player "player100"`,
		},
		{
			clauses: []clause.Interface{clause.Fetch{Names: []string{"player"}, VID: []string{"player100", "player101", "player102"}}},
			gqlWant: `FETCH PROP ON player "player100", "player101", "player102"`,
		},
		{
			clauses: []clause.Interface{clause.Fetch{Names: []string{"player"}, VID: []int{1, 2, 3}}},
			gqlWant: `FETCH PROP ON player 1, 2, 3`,
		},
		{
			clauses: []clause.Interface{clause.Fetch{Names: []string{"serve"}, VID: clause.Expr{Str: `"player100" -> "team204"`}}},
			gqlWant: `FETCH PROP ON serve "player100" -> "team204"`,
		},
		{
			clauses: []clause.Interface{clause.Fetch{Names: []string{"serve"}, VID: []*clause.Expr{{Str: `"player100" -> "team204"`}, {Str: `"player133" -> "team202"`}}}},
			gqlWant: `FETCH PROP ON serve "player100" -> "team204", "player133" -> "team202"`,
		},
		{
			clauses: []clause.Interface{clause.Fetch{Names: []string{"player"}}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.Fetch{VID: []interface{}{"player100"}}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
