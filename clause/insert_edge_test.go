package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"reflect"
	"testing"
)

func TestInsertEdge(t *testing.T) {
	e11 := &edge1{
		SrcID: "10",
		DstID: "11",
	}
	e12 := &edge1{
		SrcID: "10",
		DstID: "11",
		Rank:  1,
	}
	e21 := edge2{
		SrcID: "11",
		DstID: "13",
		Name:  "n1",
		Age:   1,
	}
	e22 := edge2{
		SrcID: "12",
		DstID: "13",
		Name:  "n1",
		Age:   1,
	}
	e23 := edge2{
		SrcID: "13",
		DstID: "14",
		Name:  "n2",
		Age:   2,
	}
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.InsertEdge{Edges: reflect.ValueOf(e11)}},
			gqlWant: `INSERT EDGE e1() VALUES "10"->"11":()`,
		},
		{
			clauses: []clause.Interface{clause.InsertEdge{Edges: reflect.ValueOf(e12)}},
			gqlWant: `INSERT EDGE e1() VALUES "10"->"11"@1:()`,
		},
		{
			clauses: []clause.Interface{clause.InsertEdge{Edges: reflect.ValueOf(e21)}},
			gqlWant: `INSERT EDGE e2(name, age) VALUES "11"->"13":("n1", 1)`,
		},
		{
			clauses: []clause.Interface{clause.InsertEdge{Edges: reflect.ValueOf([]edge2{e22, e23})}},
			gqlWant: `INSERT EDGE e2(name, age) VALUES "12"->"13":("n1", 1), "13"->"14":("n2", 2)`,
		},
		{
			clauses: []clause.Interface{clause.InsertEdge{IfNotExist: true, Edges: reflect.ValueOf([]edge2{e22, e23})}},
			gqlWant: `INSERT EDGE IF NOT EXISTS e2(name, age) VALUES "12"->"13":("n1", 1), "13"->"14":("n2", 2)`,
		},
		{
			clauses: []clause.Interface{clause.InsertEdge{IfNotExist: true, Edges: reflect.ValueOf([]*edge2{&e22, &e23})}},
			gqlWant: `INSERT EDGE IF NOT EXISTS e2(name, age) VALUES "12"->"13":("n1", 1), "13"->"14":("n2", 2)`,
		},
		{
			clauses: []clause.Interface{clause.InsertEdge{IfNotExist: true}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}

type edge1 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
}

func (e edge1) EdgeTypeName() string {
	return "e1"
}

type edge2 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
	Name  string `norm:"prop:name"`
	Age   int    `norm:"prop:age"`
}

func (e edge2) EdgeTypeName() string {
	return "e2"
}
