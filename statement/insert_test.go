package statement

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().InsertVertex(&t1{VID: "10"})
			},
			want: `INSERT VERTEX t1() VALUES "10":();`,
		},
		{
			stmt: func() *Statement {
				return New().InsertVertex(&t2{VID: "11", Name: "n1", Age: 12})
			},
			want: `INSERT VERTEX t2(name, age) VALUES "11":("n1", 12);`,
		},
		{
			stmt: func() *Statement {
				return New().InsertVertex([]t2{{VID: "13", Name: "n3", Age: 12}, {VID: "14", Name: "n4", Age: 8}})
			},
			want: `INSERT VERTEX t2(name, age) VALUES "13":("n3", 12), "14":("n4", 8);`,
		},
		{
			stmt: func() *Statement {
				return New().InsertVertex(v1{VID: "21", T1: t3{P1: 321}, T2: t4{P2: "hello"}})
			},
			want: `INSERT VERTEX t3(p1), t4(p2) VALUES "21":(321, "hello");`,
		},
		{
			stmt: func() *Statement {
				return New().InsertVertex(t2{VID: "1", Name: "n3", Age: 14}, true)
			},
			want: `INSERT VERTEX IF NOT EXISTS t2(name, age) VALUES "1":("n3", 14);`,
		},
		{
			stmt: func() *Statement {
				return New().InsertEdge(&e1{SrcID: "10", DstID: "11"})
			},
			want: `INSERT EDGE e1() VALUES "10"->"11":();`,
		},
		{
			stmt: func() *Statement {
				return New().InsertEdge(e1{SrcID: "10", DstID: "11", Rank: 1})
			},
			want: `INSERT EDGE e1() VALUES "10"->"11"@1:();`,
		},
		{
			stmt: func() *Statement {
				return New().InsertEdge(e2{SrcID: "11", DstID: "13", Name: "n1", Age: 1})
			},
			want: `INSERT EDGE e2(name, age) VALUES "11"->"13":("n1", 1);`,
		},
		{
			stmt: func() *Statement {
				return New().InsertEdge([]*e2{{SrcID: "12", DstID: "13", Name: "n1", Age: 1}, {SrcID: "13", DstID: "14", Name: "n2", Age: 2}})
			},
			want: `INSERT EDGE e2(name, age) VALUES "12"->"13":("n1", 1), "13"->"14":("n2", 2);`,
		},
		{
			stmt: func() *Statement {
				return New().InsertEdge(e2{SrcID: "14", DstID: "15", Rank: 1, Name: "n2", Age: 13}, true)
			},
			want: `INSERT EDGE IF NOT EXISTS e2(name, age) VALUES "14"->"15"@1:("n2", 13);`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got an unexpected error: %v", err)
				}
				return
			}
			if ngql != tt.want {
				t.Errorf("NGQL = %v, want %v", ngql, tt.want)
			}
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

type v1 struct {
	VID string `norm:"vertex_id"`
	T1  t3
	T2  t4
}

func (v *v1) VertexID() string {
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

func (t *t4) VertexTagName() string {
	return "t4"
}

type e1 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
}

func (e e1) EdgeTypeName() string {
	return "e1"
}

type e2 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
	Name  string `norm:"prop:name"`
	Age   int    `norm:"prop:age"`
}

func (e *e2) EdgeTypeName() string {
	return "e2"
}
