package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
	"time"
)

func TestWhere(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Operator: clause.OperatorAnd, Expr: clause.Expr{Str: "v.player.name == ?", Vars: []interface{}{"Tim Duncan"}}}}},
				clause.Where{Conditions: []clause.Condition{{Operator: clause.OperatorXor, Expr: clause.Expr{Str: "v.player.age < ? AND v.player.name == ?", Vars: []interface{}{30, "Yao Ming"}}}}},
				clause.Where{Conditions: []clause.Condition{{Operator: "OR NOT", Expr: clause.Expr{Str: "v.player.name == ? OR v.player.name == ?", Vars: []interface{}{"Yao Ming", "Tim Duncan"}}}}},
			},
			gqlWant: `WHERE v.player.name == "Tim Duncan" XOR (v.player.age < 30 AND v.player.name == "Yao Ming") OR NOT (v.player.name == "Yao Ming" OR v.player.name == "Tim Duncan")`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Operator: "AND", Expr: clause.Expr{Str: "properties(edge).degree > ?", Vars: []interface{}{90}}}}},
				clause.Where{Conditions: []clause.Condition{{Operator: "OR", Expr: clause.Expr{Str: "properties($$).age != ?", Vars: []interface{}{33}}}}},
				clause.Where{Conditions: []clause.Condition{{Operator: "AND", Expr: clause.Expr{Str: "properties($$).name != ?", Vars: []interface{}{"Tony Parker"}}}}},
			},
			gqlWant: `WHERE properties(edge).degree > 90 OR properties($$).age != 33 AND properties($$).name != "Tony Parker"`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Operator: "AND", Expr: clause.Expr{Str: "exists(v.player.age)"}}}},
			},
			gqlWant: `WHERE exists(v.player.age)`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "NOT (v)-[e]->(t:team)"}}}},
			},
			gqlWant: `WHERE NOT (v)-[e]->(t:team)`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "v.player.name STARTS WITH ?", Vars: []interface{}{"t"}}}}},
			},
			gqlWant: `WHERE v.player.name STARTS WITH "t"`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "player.age IN ?", Vars: []interface{}{[]int{25, 28}}}}}},
			},
			gqlWant: `WHERE player.age IN [25, 28]`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "player.name IN ?", Vars: []interface{}{[]string{"Anne", "John"}}}}}},
			},
			gqlWant: `WHERE player.name IN ["Anne", "John"]`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "v.date1.p3 < ?", Vars: []interface{}{time.Date(1988, 3, 18, 0, 0, 0, 0, time.Local)}}}}},
			},
			gqlWant: `WHERE v.date1.p3 < datetime("1988-03-18T00:00:00")`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "v.date1.p3 < ?", Vars: []interface{}{clause.Expr{Str: `datetime("1988-03-18T00:00:00")`}}}}}},
			},
			gqlWant: `WHERE v.date1.p3 < datetime("1988-03-18T00:00:00")`,
		},
		{
			clauses: []clause.Interface{
				clause.Where{Conditions: []clause.Condition{{Expr: clause.Expr{Str: "v.date1.p3 < ?", Vars: []interface{}{&clause.Expr{Str: `datetime("1988-03-18T00:00:00")`}}}}}},
			},
			gqlWant: `WHERE v.date1.p3 < datetime("1988-03-18T00:00:00")`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
