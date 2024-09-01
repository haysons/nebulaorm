package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestUpdateEdge(t *testing.T) {
	e11 := &edge1{
		SrcID: "player100",
		DstID: "team204",
	}
	e21 := &edge2{
		SrcID: "player100",
		DstID: "team204",
		Rank:  2,
	}
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: e11, PropsUpdate: map[string]interface{}{"start_year": clause.Expr{Str: "start_year + 1"}}}},
			gqlWant: `UPDATE EDGE ON e1 "player100"->"team204" SET start_year = start_year + 1`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: e11, PropsUpdate: map[string]*clause.Expr{"start_year": {Str: "start_year + 1"}}}},
			gqlWant: `UPDATE EDGE ON e1 "player100"->"team204" SET start_year = start_year + 1`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: `e1 "player100" -> "team204"`, PropsUpdate: map[string]interface{}{"start_year": clause.Expr{Str: "start_year + 1"}}}},
			gqlWant: `UPDATE EDGE ON e1 "player100" -> "team204" SET start_year = start_year + 1`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: e21, PropsUpdate: &edge2{Rank: 3, Name: "hayson", Age: 26}}},
			gqlWant: `UPDATE EDGE ON e2 "player100"->"team204"@2 SET name = "hayson", age = 26`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: e21, PropsUpdate: map[string]interface{}{"name": "hayson", "age": clause.Expr{Str: "age + 1"}}}},
			gqlWant: `UPDATE EDGE ON e2 "player100"->"team204"@2 SET name = "hayson", age = age + 1`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{IsUpsert: true, Edge: e21, PropsUpdate: edge2{Rank: 3, Name: "hayson", Age: 26}}},
			gqlWant: `UPSERT EDGE ON e2 "player100"->"team204"@2 SET name = "hayson", age = 26`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: e21, PropsUpdate: &edge2{Rank: 3, Name: "hayson"}}},
			gqlWant: `UPDATE EDGE ON e2 "player100"->"team204"@2 SET name = "hayson"`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: e21, PropsUpdate: map[string]string{"name": "hayson"}}},
			gqlWant: `UPDATE EDGE ON e2 "player100"->"team204"@2 SET name = "hayson"`,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: ""}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.UpdateEdge{Edge: `e1 "player100" -> "team204"`}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
