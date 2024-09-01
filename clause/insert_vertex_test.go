package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"reflect"
	"testing"
)

func TestInsertVertex(t *testing.T) {
	t11 := &t1{VID: "10"}
	t21 := t2{VID: "11", Name: "n1", Age: 12}
	t22 := t2{VID: "12", Name: "n2", Age: 18}
	v31 := v3{
		VID: "21",
		T1:  &t3{P1: 321},
		T2:  t4{P2: "hello"},
	}
	v32 := &v3{
		VID: "22",
		T1:  &t3{P1: 456},
		T2:  t4{P2: "world"},
	}
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.InsertVertex{Vertexes: reflect.ValueOf(t11)}},
			gqlWant: `INSERT VERTEX t1() VALUES "10":()`,
		},
		{
			clauses: []clause.Interface{clause.InsertVertex{Vertexes: reflect.ValueOf(t21)}},
			gqlWant: `INSERT VERTEX t2(name, age) VALUES "11":("n1", 12)`,
		},
		{
			clauses: []clause.Interface{clause.InsertVertex{Vertexes: reflect.ValueOf([]*t2{&t21, &t22})}},
			gqlWant: `INSERT VERTEX t2(name, age) VALUES "11":("n1", 12), "12":("n2", 18)`,
		},
		{
			clauses: []clause.Interface{clause.InsertVertex{Vertexes: reflect.ValueOf(v31)}},
			gqlWant: `INSERT VERTEX t3(p1), t4(p2) VALUES "21":(321, "hello")`,
		},
		{
			clauses: []clause.Interface{clause.InsertVertex{Vertexes: reflect.ValueOf([]v3{v31, *v32})}},
			gqlWant: `INSERT VERTEX t3(p1), t4(p2) VALUES "21":(321, "hello"), "22":(456, "world")`,
		},
		{
			clauses: []clause.Interface{clause.InsertVertex{IfNotExist: true, Vertexes: reflect.ValueOf([]v3{v31, *v32})}},
			gqlWant: `INSERT VERTEX IF NOT EXISTS t3(p1), t4(p2) VALUES "21":(321, "hello"), "22":(456, "world")`,
		},
		{
			clauses: []clause.Interface{clause.InsertVertex{IfNotExist: true}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}

type t1 struct {
	VID string `norm:"vertex_id"`
}

func (t t1) VertexID() string {
	return t.VID
}

func (t t1) VertexTagName() string {
	return "t1"
}

type t2 struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

func (t t2) VertexID() string {
	return t.VID
}

func (t t2) VertexTagName() string {
	return "t2"
}

type v3 struct {
	VID string `norm:"vertex_id"`
	T1  *t3
	T2  t4
}

func (v *v3) VertexID() string {
	return v.VID
}

type t3 struct {
	P1 int
}

func (t *t3) VertexTagName() string {
	return "t3"
}

type t4 struct {
	P2 string
}

func (t t4) VertexTagName() string {
	return "t4"
}
