package statement

import "github.com/haysons/nebulaorm/clause"

// DeleteVertex generate delete vertex clause
//
// DELETE VERTEX "team1"
// stmt.DeleteVertex("team1")
//
// DELETE VERTEX "team1", "team2" WITH EDGE
// stmt.DeleteVertex([]string{"team1", "team2"}, true)
//
// DELETE VERTEX $-.id
// stmt.DeleteVertex(clause.Expr{Str: "$-.id"})
func (stmt *Statement) DeleteVertex(vid interface{}, withEdge ...bool) *Statement {
	var withEdgeOpt bool
	if len(withEdge) > 0 {
		withEdgeOpt = withEdge[0]
	}
	stmt.AddClause(&clause.DeleteVertex{
		VID:      vid,
		WithEdge: withEdgeOpt,
	})
	stmt.SetPartType(PartTypeDeleteVertex)
	return stmt
}

// DeleteEdge generate delete edge clause
//
// DELETE EDGE serve "player100" -> "team204"@0
// stmt.DeleteEdge("serve", `"player100" -> "team204"@0`)
//
// DELETE EDGE serve "player100" -> "team204"@0, "player100" -> "team204"@1
// stmt.DeleteEdge("serve", []string{`"player100" -> "team204"@0`, `"player100" -> "team204"@1`})
//
// manually concatenating an expression for an edge can be tedious, and you can pass in a structure representing
// the edge (implementing the resolver.EdgeTypeName interface) to indicate the list of edges to be removed
//
//	type edgeServe struct {
//		 SrcID string `norm:"edge_src_id"`
//		 DstID string `norm:"edge_dst_id"`
//		 Rank  int    `norm:"edge_rank"`
//	}
//
//	func (e edgeServe) EdgeTypeName() string {
//		 return "serve"
//	}
//
// DELETE EDGE serve "player100"->"team204"@1
// stmt.DeleteEdge("serve", &edgeServe{SrcID: "player100", DstID: "team204", Rank: 1})
//
// DELETE EDGE serve "player100"->"team204", "player101"->"team204"@1
// stmt.DeleteEdge("serve", []edgeServe{{SrcID: "player100", DstID: "team204"}, {SrcID: "player101", DstID: "team204", Rank: 1}})
func (stmt *Statement) DeleteEdge(edgeTypeName string, edges interface{}) *Statement {
	stmt.AddClause(&clause.DeleteEdge{
		EdgeTypeName: edgeTypeName,
		Edges:        edges,
	})
	stmt.SetPartType(PartTypeDeleteEdge)
	return stmt
}
