package statement

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().UpdateVertex("10", &t2{Name: "hayson", Age: 26})
			},
			want: `UPDATE VERTEX ON t2 "10" SET name = "hayson", age = 26;`,
		},
		{
			stmt: func() *Statement {
				return New().UpdateVertex("10", &t2{Name: "hayson"})
			},
			want: `UPDATE VERTEX ON t2 "10" SET name = "hayson";`,
		},
		{
			stmt: func() *Statement {
				return New().UpdateVertex("10", &t2{Name: "hayson"}, clause.WithPropNames([]string{"name", "age"}))
			},
			want: `UPDATE VERTEX ON t2 "10" SET name = "hayson", age = 0;`,
		},
		{
			stmt: func() *Statement {
				return New().UpdateVertex("player101", map[string]interface{}{"age": clause.Expr{Str: "age + 2"}}, clause.WithTagName("player")).
					When("name == ?", "Tony Parker").Yield("name AS Name, age AS Age")
			},
			want: `UPDATE VERTEX ON player "player101" SET age = age + 2 WHEN name == "Tony Parker" YIELD name AS Name, age AS Age;`,
		},
		{
			stmt: func() *Statement {
				return New().UpsertVertex("player666", map[string]int{"age": 30}, clause.WithTagName("player")).
					When("name == ?", "Joe").Yield("name AS Name, age AS Age")
			},
			want: `UPSERT VERTEX ON player "player666" SET age = 30 WHEN name == "Joe" YIELD name AS Name, age AS Age;`,
		},
		{
			stmt: func() *Statement {
				return New().UpsertVertex("player101", map[string]clause.Expr{"age": {Str: "age + 2"}}, clause.WithTagName("player")).
					When("name == ?", "Tony Parker").Yield("name AS Name, age AS Age")
			},
			want: `UPSERT VERTEX ON player "player101" SET age = age + 2 WHEN name == "Tony Parker" YIELD name AS Name, age AS Age;`,
		},
		{
			stmt: func() *Statement {
				return New().UpdateEdge(`e2 "player100" -> "team204"@0`, &e2{Name: "hayson", Age: 26})
			},
			want: `UPDATE EDGE ON e2 "player100" -> "team204"@0 SET name = "hayson", age = 26;`,
		},
		{
			stmt: func() *Statement {
				return New().UpdateEdge(e2{SrcID: "player100", DstID: "team204"}, &e2{Age: 26})
			},
			want: `UPDATE EDGE ON e2 "player100"->"team204" SET age = 26;`,
		},
		{
			stmt: func() *Statement {
				return New().UpdateEdge(e2{SrcID: "player100", DstID: "team204", Rank: 1}, &e2{Age: 26}, clause.WithPropNames([]string{"name"}))
			},
			want: `UPDATE EDGE ON e2 "player100"->"team204"@1 SET name = "";`,
		},
		{
			stmt: func() *Statement {
				return New().UpdateEdge(e2{SrcID: "player100", DstID: "team204"}, map[string]interface{}{"start_year": clause.Expr{Str: "start_year + 1"}}).
					When("end_year > ?", 2010).Yield("start_year, end_year")
			},
			want: `UPDATE EDGE ON e2 "player100"->"team204" SET start_year = start_year + 1 WHEN end_year > 2010 YIELD start_year, end_year;`,
		},
		{
			stmt: func() *Statement {
				return New().UpsertEdge(e2{SrcID: "player668", DstID: "team200"}, map[string]interface{}{"start_year": 2000, "end_year": clause.Expr{Str: "end_year + 1"}}).Yield("start_year, end_year")
			},
			want: `UPSERT EDGE ON e2 "player668"->"team200" SET start_year = 2000, end_year = end_year + 1 YIELD start_year, end_year;`,
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

type playerUpdate map[string]interface{}

func (m playerUpdate) VertexTagName() string {
	return "player"
}
