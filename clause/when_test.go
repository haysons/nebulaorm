package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
	"time"
)

func TestWhen(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Operator: "AND", Expr: clause.Expr{Str: "v.player.name == ?", Vars: []interface{}{"Tim Duncan"}}}}},
				clause.When{Conditions: []clause.Condition{{Operator: "XOR", Expr: clause.Expr{Str: "v.player.age < ? AND v.player.name == ?", Vars: []interface{}{30, "Yao Ming"}}}}},
				clause.When{Conditions: []clause.Condition{{Operator: "OR NOT", Expr: clause.Expr{Str: "v.player.name == ? OR v.player.name == ?", Vars: []interface{}{"Yao Ming", "Tim Duncan"}}}}},
			},
			gqlWant: `WHEN v.player.name == "Tim Duncan" XOR (v.player.age < 30 AND v.player.name == "Yao Ming") OR NOT (v.player.name == "Yao Ming" OR v.player.name == "Tim Duncan")`,
		},
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Operator: "AND", Expr: clause.Expr{Str: "properties(edge).degree > ?", Vars: []interface{}{90}}}}},
				clause.When{Conditions: []clause.Condition{{Operator: "OR", Expr: clause.Expr{Str: "properties($$).age != ?", Vars: []interface{}{33}}}}},
				clause.When{Conditions: []clause.Condition{{Operator: "AND", Expr: clause.Expr{Str: "properties($$).name != ?", Vars: []interface{}{"Tony Parker"}}}}},
			},
			gqlWant: `WHEN properties(edge).degree > 90 OR properties($$).age != 33 AND properties($$).name != "Tony Parker"`,
		},
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Operator: "AND", Expr: clause.Expr{Str: "exists(v.player.age)"}}}},
			},
			gqlWant: `WHEN exists(v.player.age)`,
		},
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "NOT (v)-[e]->(t:team)"}}}},
			},
			gqlWant: `WHEN NOT (v)-[e]->(t:team)`,
		},
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "v.player.name STARTS WITH ?", Vars: []interface{}{"t"}}}}},
			},
			gqlWant: `WHEN v.player.name STARTS WITH "t"`,
		},
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "player.age IN ?", Vars: []interface{}{[]int{25, 28}}}}}},
			},
			gqlWant: `WHEN player.age IN [25, 28]`,
		},
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "player.name IN ?", Vars: []interface{}{[]string{"Anne", "John"}}}}}},
			},
			gqlWant: `WHEN player.name IN ["Anne", "John"]`,
		},
		{
			clauses: []clause.Interface{
				clause.When{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "v.date1.p3 < ?", Vars: []interface{}{time.Date(1988, 3, 18, 0, 0, 0, 0, time.Local)}}}}},
			},
			gqlWant: `WHEN v.date1.p3 < datetime("1988-03-18T00:00:00")`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
