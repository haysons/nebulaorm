package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestDeleteVertex(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: "team1"}},
			gqlWant: `DELETE VERTEX "team1"`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: "team1", WithEdge: true}},
			gqlWant: `DELETE VERTEX "team1" WITH EDGE`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: clause.Expr{Str: "$-.id"}, WithEdge: true}},
			gqlWant: `DELETE VERTEX $-.id WITH EDGE`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: &clause.Expr{Str: "$-.id"}, WithEdge: true}},
			gqlWant: `DELETE VERTEX $-.id WITH EDGE`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: []int{1, 2, 3}}},
			gqlWant: `DELETE VERTEX 1, 2, 3`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: []string{"1", "2", "3"}, WithEdge: true}},
			gqlWant: `DELETE VERTEX "1", "2", "3" WITH EDGE`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: []clause.Expr{{Str: "$-.id"}}, WithEdge: true}},
			gqlWant: `DELETE VERTEX $-.id WITH EDGE`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{VID: []*clause.Expr{{Str: "$-.id"}}, WithEdge: true}},
			gqlWant: `DELETE VERTEX $-.id WITH EDGE`,
		},
		{
			clauses: []clause.Interface{clause.DeleteVertex{}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
