package statement

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().DeleteVertex("team1")
			},
			want: `DELETE VERTEX "team1";`,
		},
		{
			stmt: func() *Statement {
				return New().DeleteVertex([]string{"team1", "team2"}, true)
			},
			want: `DELETE VERTEX "team1", "team2" WITH EDGE;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player100").Over("serve").Where("properties(edge).start_year == ?", "2021").Yield("dst(edge) AS id").Pipe().
					DeleteVertex(clause.Expr{Str: "$-.id"})
			},
			want: `GO FROM "player100" OVER serve WHERE properties(edge).start_year == "2021" YIELD dst(edge) AS id | DELETE VERTEX $-.id;`,
		},
		{
			stmt: func() *Statement {
				return New().DeleteEdge("serve", `"player100" -> "team204"@0`)
			},
			want: `DELETE EDGE serve "player100" -> "team204"@0;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player100").Over("follow").Where("dst(edge) == ?", "player101").Yield("src(edge) AS src, dst(edge) AS dst, rank(edge) AS rank").Pipe().
					DeleteEdge("follow", "$-.src -> $-.dst @ $-.rank")
			},
			want: `GO FROM "player100" OVER follow WHERE dst(edge) == "player101" YIELD src(edge) AS src, dst(edge) AS dst, rank(edge) AS rank | DELETE EDGE follow $-.src -> $-.dst @ $-.rank;`,
		},
		{
			stmt: func() *Statement {
				return New().DeleteEdge("serve", []string{`"player100" -> "team204"@0`, `"player100" -> "team204"@1`})
			},
			want: `DELETE EDGE serve "player100" -> "team204"@0, "player100" -> "team204"@1;`,
		},
		{
			stmt: func() *Statement {
				return New().DeleteEdge("serve", []string{`"player100" -> "team204"@0`, `"player100" -> "team204"@1`})
			},
			want: `DELETE EDGE serve "player100" -> "team204"@0, "player100" -> "team204"@1;`,
		},
		{
			stmt: func() *Statement {
				return New().DeleteEdge("serve", &edgeServe{SrcID: "player100", DstID: "team204", Rank: 1})
			},
			want: `DELETE EDGE serve "player100"->"team204"@1;`,
		},
		{
			stmt: func() *Statement {
				return New().DeleteEdge("serve", []edgeServe{{SrcID: "player100", DstID: "team204"}, {SrcID: "player101", DstID: "team204", Rank: 1}})
			},
			want: `DELETE EDGE serve "player100"->"team204", "player101"->"team204"@1;`,
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

type edgeServe struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
}

func (e edgeServe) EdgeTypeName() string {
	return "serve"
}
