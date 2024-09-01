package statement

import (
	"github.com/haysons/nebulaorm/clause"
)

// UpdateVertex generate update vertex clause
// the vid parameter is used to specify the vertex id that needs to be updated, and tagUpdate specifies that the tag needs
// to be updated as well as the value of the updated attribute. if the update of an attribute is performed through a structure,
// only the non-zero exported fields in the structure will be updated, the
//
//	type t2 struct {
//		VID  string `norm:"vertex_id"`
//		Name string `norm:"prop:name"`
//		Age  int    `norm:"prop:age"`
//	}
//
//	func (t t2) VertexTagName() string {
//		return "t2"
//	}
//
// UPDATE VERTEX ON t2 "10" SET name = "hayson
// stmt.UpdateVertex("10", &t2{Name: "hayson"})
//
// if you want to update a zero-valued field, you can proactively specify the field to be updated
//
// UPDATE VERTEX ON t2 "10" SET name = "hayson", age = 0
// stmt.UpdateVertex("10", &t2{Name: "hayson"}, clause.WithPropNames([]string{"name", "age"}))
//
// you can also use map[string]interface{} to update, especially if you need to update a field as an expression, in which case
// you have to manually specify the name of the tag that needs to be updated.
//
// UPDATE VERTEX ON player "player101" SET age = age + 2 WHEN name == "Tony Parker" YIELD name AS Name, age AS Age
// stmt.UpdateVertex("player101", map[string]interface{}{"age": clause.Expr{Str: "age + 2"}}, clause.WithTagName("player")).
// When("name == ?", "Tony Parker").Yield("name AS Name, age AS Age")
//
// other uses can be found in./update_test
func (stmt *Statement) UpdateVertex(vid interface{}, tagUpdate interface{}, opts ...clause.Option) *Statement {
	updateOpts := new(clause.Options)
	for _, opt := range opts {
		opt(updateOpts)
	}
	stmt.AddClause(&clause.UpdateVertex{
		IsUpsert:  false,
		VID:       vid,
		TagUpdate: tagUpdate,
		Opts:      *updateOpts,
	})
	stmt.SetPartType(PartTypeUpdateVertex)
	return stmt
}

// UpsertVertex generate upsert vertex clause
// specific usage reference UpdateVertex
func (stmt *Statement) UpsertVertex(vid interface{}, tagUpdate interface{}, opts ...clause.Option) *Statement {
	updateOpts := new(clause.Options)
	for _, opt := range opts {
		opt(updateOpts)
	}
	stmt.AddClause(&clause.UpdateVertex{
		IsUpsert:  true,
		VID:       vid,
		TagUpdate: tagUpdate,
		Opts:      *updateOpts,
	})
	stmt.SetPartType(PartTypeUpdateVertex)
	return stmt
}

// UpdateEdge generate update edge clause
// The edge parameter is the edge that needs to be updated, you can either use a string to represent the edge that needs to
// be updated, or you can pass in an edge variable directly, which reduces the annoying string splicing. propsUpdate says
// The list of properties that need to be updated, the properties can be updated via a struct, this will only update the non-zero
// exported fields in the struct, if you want to update the zero-valued fields you can use clause.WithPropNames
// If you want to update a zero-valued field, you can use clause.WithPropNames to specify the field to be updated.
//
//	type e2 struct {
//		SrcID string `norm:"edge_src_id"`
//		DstID string `norm:"edge_dst_id"`
//		Rank  int    `norm:"edge_rank"`
//		Name  string `norm:"prop:name"`
//		Age   int    `norm:"prop:age"`
//	}
//
//	func (e *e2) EdgeTypeName() string {
//		return "e2"
//	}
//
// UPDATE EDGE ON e2 "player100" -> "team204"@0 SET name = "hayson", age = 26`
// stmt.UpdateEdge(`e2 "player100" -> "team204"@0`, &e2{Name: "hayson", Age: 26})
//
// UPDATE EDGE ON e2 "player100"->"team204" SET age = 26
// stmt.UpdateEdge(e2{SrcID: "player100", DstID: "team204"}, &e2{Age: 26})
//
// UPDATE EDGE ON e2 "player100"->"team204"@1 SET name = ""
// stmt.UpdateEdge(e2{SrcID: "player100", DstID: "team204", Rank: 1}, &e2{Age: 26}, clause.WithPropNames([]string{"name"}))
//
// it is also possible to specify the properties to be updated using map[string]interface{}
//
// UPDATE EDGE ON e2 "player100"->"team204" SET start_year = start_year + 1 WHEN end_year > 2010 YIELD start_year, end_year
// stmt.UpdateEdge(e2{SrcID: "player100", DstID: "team204"}, map[string]interface{}{"start_year": clause.Expr{Str: "start_year + 1"}}).
// When("end_year > ?", 2010).Yield("start_year, end_year")
func (stmt *Statement) UpdateEdge(edge interface{}, propsUpdate interface{}, opts ...clause.Option) *Statement {
	updateOpts := new(clause.Options)
	for _, opt := range opts {
		opt(updateOpts)
	}
	stmt.AddClause(&clause.UpdateEdge{
		IsUpsert:    false,
		Edge:        edge,
		PropsUpdate: propsUpdate,
		Opts:        *updateOpts,
	})
	stmt.SetPartType(PartTypeUpdateEdge)
	return stmt
}

// UpsertEdge generate upsert edge clause
// specific usage reference UpdateEdge
func (stmt *Statement) UpsertEdge(edge interface{}, propsUpdate interface{}, opts ...clause.Option) *Statement {
	updateOpts := new(clause.Options)
	for _, opt := range opts {
		opt(updateOpts)
	}
	stmt.AddClause(&clause.UpdateEdge{
		IsUpsert:    true,
		Edge:        edge,
		PropsUpdate: propsUpdate,
		Opts:        *updateOpts,
	})
	stmt.SetPartType(PartTypeUpdateEdge)
	return stmt
}

// When mainly used to generate when clause in update type statements
// specific usage reference Where
func (stmt *Statement) When(query string, args ...interface{}) *Statement {
	if query == "" {
		return stmt
	}
	stmt.AddClause(&clause.When{
		Conditions: []clause.Condition{stmt.buildCondition(clause.OperatorAnd, query, args...)},
	})
	return stmt
}
