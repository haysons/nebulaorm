package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestDeleteEdge(t *testing.T) {
	edge1 := edgeTest{
		SrcID: "101",
		DstID: "102",
		Rank:  1,
	}
	edge2 := edgeTest{
		SrcID: "201",
		DstID: "202",
	}
	edge3 := edgeTest{
		SrcID: "301",
		DstID: "302",
		Rank:  2,
	}
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.DeleteEdge{EdgeTypeName: "serve", Edges: `"player100" -> "team204"@0`}},
			gqlWant: `DELETE EDGE serve "player100" -> "team204"@0`,
		},
		{
			clauses: []clause.Interface{clause.DeleteEdge{EdgeTypeName: "serve", Edges: []string{`"player100" -> "team204"`, `"player101" -> "team204"@1`}}},
			gqlWant: `DELETE EDGE serve "player100" -> "team204", "player101" -> "team204"@1`,
		},
		{
			clauses: []clause.Interface{clause.DeleteEdge{EdgeTypeName: "serve", Edges: `"player100" -> "team204"@0`}, clause.DeleteEdge{EdgeTypeName: "follow", Edges: `$-.src -> $-.dst @ $-.rank`}},
			gqlWant: `DELETE EDGE follow $-.src -> $-.dst @ $-.rank`,
		},
		{
			clauses: []clause.Interface{clause.DeleteEdge{EdgeTypeName: "edge_test", Edges: []*edgeTest{&edge1, &edge2, &edge3}}},
			gqlWant: `DELETE EDGE edge_test "101"->"102"@1, "201"->"202", "301"->"302"@2`,
		},
		{
			clauses: []clause.Interface{clause.DeleteEdge{EdgeTypeName: "edge_test", Edges: []edgeTest{edge1, edge2, edge3}}},
			gqlWant: `DELETE EDGE edge_test "101"->"102"@1, "201"->"202", "301"->"302"@2`,
		},
		{
			clauses: []clause.Interface{clause.DeleteEdge{EdgeTypeName: "edge_test"}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.DeleteEdge{EdgeTypeName: "edge_test", Edges: map[string]string{}}},
			errWant: clause.ErrInvalidClauseParams,
		},
		{
			clauses: []clause.Interface{clause.DeleteEdge{Edges: []edgeTest{edge1, edge2, edge3}}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}

type edgeTest struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
	Name  string `norm:"prop:name"`
	Age   int    `norm:"prop:age"`
}

func (e edgeTest) EdgeTypeName() string {
	return "edge_test"
}
