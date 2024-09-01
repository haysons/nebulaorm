package statement

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().Go().From("player102").Over("serve").Yield("dst(edge)")
			},
			want: `GO FROM "player102" OVER serve YIELD dst(edge);`,
		},
		{
			stmt: func() *Statement {
				return New().Go(2).From("player102").Over("follow").Yield("dst(edge)")
			},
			want: `GO 2 STEPS FROM "player102" OVER follow YIELD dst(edge);`,
		},
		{
			stmt: func() *Statement {
				return New().Go(1, 2).From("player100").Over("follow").Yield("dst(edge) AS destination")
			},
			want: `GO 1 TO 2 STEPS FROM "player100" OVER follow YIELD dst(edge) AS destination;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From([]string{"player100", "player102"}).Over("serve").Where("properties(edge).start_year > ?", 1995).Yield("properties($$).name AS team_name, properties(edge).start_year AS start_year, properties($^).name AS player_name", true)
			},
			want: `GO FROM "player100", "player102" OVER serve WHERE properties(edge).start_year > 1995 YIELD DISTINCT properties($$).name AS team_name, properties(edge).start_year AS start_year, properties($^).name AS player_name;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player102").Over("*", clause.OverDirectBidirect).Yield("edge AS e")
			},
			want: `GO FROM "player102" OVER * BIDIRECT YIELD edge AS e;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player100").Over("follow", "serve").Yield("properties(edge).degree, properties(edge).start_year")
			},
			want: `GO FROM "player100" OVER follow, serve YIELD properties(edge).degree, properties(edge).start_year;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player100").Over("follow", clause.OverDirectReversely).Yield("src(edge) AS destination")
			},
			want: `GO FROM "player100" OVER follow REVERSELY YIELD src(edge) AS destination;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player100").Over("follow").Yield("src(edge) AS id").Pipe().
					Go().From(clause.Expr{Str: "$-.id"}).Over("serve").Where("properties($^).age > ?", 20).Yield("properties($^).name AS Friend, properties($$).name AS Team")
			},
			want: `GO FROM "player100" OVER follow YIELD src(edge) AS id | GO FROM $-.id OVER serve WHERE properties($^).age > 20 YIELD properties($^).name AS Friend, properties($$).name AS Team;`,
		},
		{
			stmt: func() *Statement {
				return New().Go(2).From("player100").Over("follow").Yield("src(edge) AS src, dst(edge) AS dst, properties($$).age AS age").GroupBy("$-.dst").Yield("$-.dst AS dst, collect_set($-.src) AS src, collect($-.age) AS age")
			},
			want: `GO 2 STEPS FROM "player100" OVER follow YIELD src(edge) AS src, dst(edge) AS dst, properties($$).age AS age | GROUP BY $-.dst YIELD $-.dst AS dst, collect_set($-.src) AS src, collect($-.age) AS age;`,
		},
		{
			stmt: func() *Statement {
				return New().Go(2).From(clause.Expr{Str: "$a.dst"}).Over("follow").Yield("$a.src AS src, $a.dst, src(edge), dst(edge)").OrderBy("$-.src").Limit(2, 1)
			},
			want: `GO 2 STEPS FROM $a.dst OVER follow YIELD $a.src AS src, $a.dst, src(edge), dst(edge) | ORDER BY $-.src | LIMIT 1, 2;`,
		},
		{
			stmt: func() *Statement {
				return New().Lookup("player").Where("player.name == ?", "Tony Parker").Yield("id(vertex)")
			},
			want: `LOOKUP ON player WHERE player.name == "Tony Parker" YIELD id(vertex);`,
		},
		{
			stmt: func() *Statement {
				return New().Lookup("player").Where("player.name STARTS WITH ?", "B").Where("player.age IN ?", []int{22, 30}).Yield("properties(vertex).name, properties(vertex).age")
			},
			want: `LOOKUP ON player WHERE player.name STARTS WITH "B" AND player.age IN [22, 30] YIELD properties(vertex).name, properties(vertex).age;`,
		},
		{
			stmt: func() *Statement {
				return New().Lookup("player").Where("player.name == ?", "Kobe Bryant").Yield("id(vertex) AS VertexID, properties(vertex).name AS name").Pipe().
					Go().From(clause.Expr{Str: "$-.VertexID"}).Over("serve").Yield("$-.name, properties(edge).start_year, properties(edge).end_year, properties($$).name")
			},
			want: `LOOKUP ON player WHERE player.name == "Kobe Bryant" YIELD id(vertex) AS VertexID, properties(vertex).name AS name | GO FROM $-.VertexID OVER serve YIELD $-.name, properties(edge).start_year, properties(edge).end_year, properties($$).name;`,
		},
		{
			stmt: func() *Statement {
				return New().Lookup("follow").Yield("properties(edge).degree as degree").OrderBy("$-.degree").Limit(10)
			},
			want: `LOOKUP ON follow YIELD properties(edge).degree as degree | ORDER BY $-.degree | LIMIT 10;`,
		},
		{
			stmt: func() *Statement {
				return New().Fetch("player", "player100").Yield("properties(vertex)")
			},
			want: `FETCH PROP ON player "player100" YIELD properties(vertex);`,
		},
		{
			stmt: func() *Statement {
				return New().Fetch("player", []string{"player101", "player102", "player103"}).Yield("properties(vertex)")
			},
			want: `FETCH PROP ON player "player101", "player102", "player103" YIELD properties(vertex);`,
		},
		{
			stmt: func() *Statement {
				return New().FetchMulti([]string{"player", "t1"}, "player100").Yield("vertex AS v")
			},
			want: `FETCH PROP ON player, t1 "player100" YIELD vertex AS v;`,
		},
		{
			stmt: func() *Statement {
				return New().FetchMulti([]string{"player", "t1"}, []string{"player100", "player103"}).Yield("vertex AS v")
			},
			want: `FETCH PROP ON player, t1 "player100", "player103" YIELD vertex AS v;`,
		},
		{
			stmt: func() *Statement {
				return New().Fetch("*", []string{"player100", "player106", "team200"}).Yield("vertex AS v")
			},
			want: `FETCH PROP ON * "player100", "player106", "team200" YIELD vertex AS v;`,
		},
		{
			stmt: func() *Statement {
				return New().Fetch("serve", clause.Expr{Str: `"player100" -> "team204"`}).Yield("properties(edge)")
			},
			want: `FETCH PROP ON serve "player100" -> "team204" YIELD properties(edge);`,
		},
		{
			stmt: func() *Statement {
				return New().Fetch("serve", []*clause.Expr{{Str: `"player100" -> "team204"`}, {Str: `"player133" -> "team202"`}}).Yield("edge AS e")
			},
			want: `FETCH PROP ON serve "player100" -> "team204", "player133" -> "team202" YIELD edge AS e;`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player101").Over("follow").Yield("src(edge) AS s, dst(edge) AS d").Pipe().
					Fetch("follow", clause.Expr{Str: "$-.s -> $-.d"}).Yield("properties(edge).degree")
			},
			want: `GO FROM "player101" OVER follow YIELD src(edge) AS s, dst(edge) AS d | FETCH PROP ON follow $-.s -> $-.d YIELD properties(edge).degree;`,
		},
		{
			stmt: func() *Statement {
				return New().Lookup("player").Where("player.age > 34").Yield("id(vertex) AS v").Pipe().
					Go().From(clause.Expr{Str: "$-.v"}).Over("serve").Yield("serve.start_year AS start_year, serve.end_year AS end_year").Pipe().
					Yield("$-.start_year, $-.end_year, count(*) AS count").Pipe().OrderBy("$-.count DESC").Limit(5)
			},
			want: `LOOKUP ON player WHERE player.age > 34 YIELD id(vertex) AS v | GO FROM $-.v OVER serve YIELD serve.start_year AS start_year, serve.end_year AS end_year | YIELD $-.start_year, $-.end_year, count(*) AS count | ORDER BY $-.count DESC | LIMIT 5;`,
		},
		{
			stmt: func() *Statement {
				return New().Go(3).From("player100").Over("*").Yield("properties($$).name AS NAME, properties($$).age AS Age").Sample(1, 2, 3)
			},
			want: `GO 3 STEPS FROM "player100" OVER * YIELD properties($$).name AS NAME, properties($$).age AS Age SAMPLE [1,2,3];`,
		},
		{
			stmt: func() *Statement {
				return New().Go(1, 3).From("player100").Over("*").Yield("properties($$).name AS NAME, properties($$).age AS Age").Sample(2, 2, 2)
			},
			want: `GO 1 TO 3 STEPS FROM "player100" OVER * YIELD properties($$).name AS NAME, properties($$).age AS Age SAMPLE [2,2,2];`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player100").Over("follow").Where("properties(edge).degree > ?", 90).Or("properties($$).age != ?", 33).Where("properties($$).name != ?", "Tony Parker").Yield("properties($$)")
			},
			want: `GO FROM "player100" OVER follow WHERE properties(edge).degree > 90 OR properties($$).age != 33 AND properties($$).name != "Tony Parker" YIELD properties($$);`,
		},
		{
			stmt: func() *Statement {
				return New().Go().From("player100").Over("follow").Where("properties(edge).degree > ?", 90).Xor("properties($$).age != ?", 33).Not("properties($$).name != ?", "Tony Parker").Yield("properties($$)")
			},
			want: `GO FROM "player100" OVER follow WHERE properties(edge).degree > 90 XOR properties($$).age != 33 NOT properties($$).name != "Tony Parker" YIELD properties($$);`,
		},
		{
			stmt: func() *Statement {
				return New().Yield("rand32(1, 6)")
			},
			want: `YIELD rand32(1, 6);`,
		},
		{
			stmt: func() *Statement {
				return New().Raw(`$var = GO FROM "player100" OVER follow YIELD dst(edge) AS id; GO FROM $var.id OVER serve YIELD properties($$).name AS Team, properties($^).name AS Player`).Yield("test")
			},
			want: `$var = GO FROM "player100" OVER follow YIELD dst(edge) AS id; GO FROM $var.id OVER serve YIELD properties($$).name AS Team, properties($^).name AS Player`,
		},
		{
			stmt: func() *Statement {
				return New().Raw(`GO FROM "player100" OVER follow YIELD dst(edge) AS id | GO FROM $-.id OVER serve YIELD properties($$).name AS Team, properties($^).name AS Player`)
			},
			want: `GO FROM "player100" OVER follow YIELD dst(edge) AS id | GO FROM $-.id OVER serve YIELD properties($$).name AS Team, properties($^).name AS Player`,
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
