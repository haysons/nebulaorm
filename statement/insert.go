package statement

import (
	"github.com/haysons/nebulaorm/clause"
	"reflect"
)

// InsertVertex generate delete vertex clause
// vertexes must be a vertex that can be parsed, that is, the vertex itself needs to implement the VertexIDStr or VertexIDInt64
// interface, and the vertex's tag needs to implement the VertexTagNamer interface, in most cases if the node has only one tag,
// it only needs the structure to implement the VertexIDStr(VertexIDInt64) and VertexTagNamer interface at the same time. If a
// vertex has multiple tags, the tag is used as the attribute interface of the structure.
// Note: currently, the structure does not support embedded
//
//	type t2 struct {
//		VID  string `norm:"vertex_id"`
//		Name string `norm:"prop:name"`
//		Age  int    `norm:"prop:age"`
//	}
//
//	func (t t2) VertexID() string {
//		return t.VID
//	}
//
//	func (t t2) VertexTagName() string {
//		return "t2"
//	}
//
// INSERT VERTEX t2(name, age) VALUES "11":("n1", 12)
// stmt.InsertVertex(&t2{VID: "11", Name: "n1", Age: 12})
//
// INSERT VERTEX t2(name, age) VALUES "13":("n3", 12), "14":("n4", 8)
// stmt.InsertVertex([]t2{{VID: "13", Name: "n3", Age: 12}, {VID: "14", Name: "n4", Age: 8}})
//
//	type v1 struct {
//		VID string `norm:"vertex_id"`
//		T1  t3
//		T2  t4
//	}
//
//	func (v *v1) VertexID() string {
//		return v.VID
//	}
//
//	type t3 struct {
//		P1 int
//	}
//
//	func (t *t3) VertexTagName() string {
//		return "t3"
//	}
//
//	type t4 struct {
//		P2 string
//	}
//
//	func (t *t4) VertexTagName() string {
//		return "t4"
//	}
//
// INSERT VERTEX t3(p1), t4(p2) VALUES "21":(321, "hello"),
// stmt.InsertVertex(v1{VID: "21", T1: t3{P1: 321}, T2: t4{P2: "hello"}})
// more usage reference: ./insert_test.go
func (stmt *Statement) InsertVertex(vertexes interface{}, ifNotExist ...bool) *Statement {
	var notExistOpt bool
	if len(ifNotExist) > 0 {
		notExistOpt = ifNotExist[0]
	}
	stmt.AddClause(&clause.InsertVertex{
		IfNotExist: notExistOpt,
		Vertexes:   reflect.ValueOf(vertexes),
	})
	stmt.SetPartType(PartTypeInsertVertex)
	return stmt
}

// InsertEdge generate delete edge clause
// edges must be parseable, meaning they need to implement the EdgeTypeNamer interface and have their src_id, dst_id,
// and rank specified in the struct field's tag
// Note: currently, the structure does not support embedded
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
// INSERT EDGE e2(name, age) VALUES "11"->"13":("n1", 1)
// stmt.InsertEdge(e2{SrcID: "11", DstID: "13", Name: "n1", Age: 1})
//
// INSERT EDGE e2(name, age) VALUES "12"->"13":("n1", 1), "13"->"14":("n2", 2)
// stmt.InsertEdge([]*e2{{SrcID: "12", DstID: "13", Name: "n1", Age: 1}, {SrcID: "13", DstID: "14", Name: "n2", Age: 2}})
//
// INSERT EDGE IF NOT EXISTS e2(name, age) VALUES "14"->"15"@1:("n2", 13)
// stmt.InsertEdge(e2{SrcID: "14", DstID: "15", Rank: 1, Name: "n2", Age: 13}, true)
func (stmt *Statement) InsertEdge(edges interface{}, ifNotExist ...bool) *Statement {
	var notExistOpt bool
	if len(ifNotExist) > 0 {
		notExistOpt = ifNotExist[0]
	}
	stmt.AddClause(&clause.InsertEdge{
		IfNotExist: notExistOpt,
		Edges:      reflect.ValueOf(edges),
	})
	stmt.SetPartType(PartTypeInsertEdge)
	return stmt
}
