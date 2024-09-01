package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestUpdateVertex(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.UpdateVertex{VID: "player101", TagUpdate: &playerUpdate{"age": &clause.Expr{Str: "age + 2"}}}},
			gqlWant: `UPDATE VERTEX ON player "player101" SET age = age + 2`,
		},
		{
			clauses: []clause.Interface{clause.UpdateVertex{IsUpsert: true, VID: "player666", TagUpdate: &playerUpdate{"age": 31}}},
			gqlWant: `UPSERT VERTEX ON player "player666" SET age = 31`,
		},
		{
			clauses: []clause.Interface{clause.UpdateVertex{IsUpsert: true, VID: "player668", TagUpdate: playerUpdate{"name": "Amber", "age": &clause.Expr{Str: "age + 1"}}}},
			gqlWant: `UPSERT VERTEX ON player "player668" SET name = "Amber", age = age + 1`,
		},
		{
			clauses: []clause.Interface{clause.UpdateVertex{VID: 101, TagUpdate: &playerTag{Name: "hayson", Age: 26}}},
			gqlWant: `UPDATE VERTEX ON player 101 SET name = "hayson", age = 26`,
		},
		{
			clauses: []clause.Interface{clause.UpdateVertex{}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.UpdateVertex{VID: 101}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}

type playerUpdate map[string]interface{}

func (m playerUpdate) VertexTagName() string {
	return "player"
}

type playerTag struct {
	Name   string
	Age    int
	Gender int
}

func (m playerTag) VertexTagName() string {
	return "player"
}
