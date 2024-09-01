package statement

import (
	"github.com/haysons/nebulaorm/clause"
	"strings"
)

// Raw execute any statement
func (stmt *Statement) Raw(raw string) *Statement {
	stmt.nGQL.WriteString(raw)
	stmt.built = true
	return stmt
}

// Go generate go clause
//
// GO
// stmt.Go()
//
// GO 2 STEPS
// stmt.Go(2)
//
// GO 1 TO 2 STEPS
// stmt.Go(1, 2)
func (stmt *Statement) Go(step ...int) *Statement {
	stepStart, stepEnd := -1, -1
	if len(step) >= 1 {
		stepStart = step[0]
	}
	if len(step) >= 2 {
		stepEnd = step[1]
	}
	stmt.AddClause(&clause.Go{
		StepStart: stepStart,
		StepEnd:   stepEnd,
	})
	stmt.SetPartType(PartTypeGo)
	return stmt
}

// From generate from clause
//
// FROM 1
// stmt.From(int(1)) or stmt.From(int64(1))
//
// FROM "player102"
// stmt.From("player102")
//
// FROM 1, 2, 3
// stmt.From([]int{1,2,3}) or stmt.From([]int64{1,2,3})
//
// FROM $-.id
// stmt.From(clause.Expr{Str:"$-.id"}) or stmt.From(&clause.Expr{Str:"$-.id"})
func (stmt *Statement) From(vid interface{}) *Statement {
	stmt.AddClause(&clause.From{
		VID: vid,
	})
	return stmt
}

// Over generate over clause
//
// OVER serve
// stmt.Over("server")
//
// OVER * BIDIRECT
// stmt.Over("*", clause.OverDirectBidirect)
//
// OVER * REVERSELY
// stmt.Over("*", clause.OverDirectReversely)
//
// OVER serve, follow BIDIRECT
// stmt.Over("server", "follow", clause.OverDirectBidirect)
func (stmt *Statement) Over(edgeType ...string) *Statement {
	if len(edgeType) == 0 {
		return stmt
	}
	var overDirect string
	if len(edgeType) > 1 {
		direct := strings.ToUpper(edgeType[len(edgeType)-1])
		if direct == clause.OverDirectBidirect || direct == clause.OverDirectReversely {
			overDirect = direct
			edgeType = edgeType[:len(edgeType)-1]
		}
	}
	stmt.AddClause(&clause.Over{
		EdgeTypeList: edgeType,
		Direction:    overDirect,
	})
	return stmt
}

// Where generate over clause，multiple where clauses are joined using AND
//
// WHERE v.player.name == "Tim Duncan"
// stmt.Where("v.player.name == ?", "Tim Duncan")
//
// WHERE v.player.age > 30
// stmt.Where("v.player.age > ?", 30)
//
// WHERE v.date1.p3 < datetime("1988-03-18T00:00:00")
// stmt.Where("v.date1.p3 < ?", time.Date(1988, 3, 18, 0, 0, 0, 0, time.Local))
//
// WHERE player.name IN ? ["Anne", "John"]
// stmt.Where("v.player.age IN ?", []string{"Anne", "John"})
//
// WHERE v.player.name == "Tim Duncan" AND v.player.age > 30
// stmt.Where("v.player.name == ?", "Tim Duncan").Where("v.player.age > ?", 30)
func (stmt *Statement) Where(query string, args ...interface{}) *Statement {
	stmt.AddClause(&clause.Where{
		Conditions: []clause.Condition{stmt.buildCondition(clause.OperatorAnd, query, args...)},
	})
	return stmt
}

// Or generate or condition in a where clause
//
// WHERE properties(edge).degree > 90 OR properties($$).age != 33
// stmt.Where("properties(edge).degree > ?", 90).Or("properties($$).age != ?", 33)
func (stmt *Statement) Or(query string, args ...interface{}) *Statement {
	stmt.AddClause(&clause.Where{
		Conditions: []clause.Condition{stmt.buildCondition(clause.OperatorOr, query, args...)},
	})
	return stmt
}

// Not generate not condition in a where clause
//
// WHERE NOT (v)-[e]->(t:team)
// stmt.Where("NOT (v)-[e]->(t:team)")
func (stmt *Statement) Not(query string, args ...interface{}) *Statement {
	stmt.AddClause(&clause.Where{
		Conditions: []clause.Condition{stmt.buildCondition(clause.OperatorNot, query, args...)},
	})
	return stmt
}

// Xor generate xor condition in a where clause
//
// WHERE v.player.name == "Tim Duncan" XOR (v.player.age < 30 AND v.player.name == "Yao Ming")
// stmt.Where("v.player.name == ?", "Tim Duncan").Xor("v.player.age < ? AND v.player.name == ?", 30, "Yao Ming")
func (stmt *Statement) Xor(query string, args ...interface{}) *Statement {
	stmt.AddClause(&clause.Where{
		Conditions: []clause.Condition{stmt.buildCondition(clause.OperatorXor, query, args...)},
	})
	return stmt
}

func (stmt *Statement) buildCondition(op string, query string, args ...interface{}) clause.Condition {
	return clause.Condition{
		Operator: op,
		Expr: clause.Expr{
			Str:  query,
			Vars: args,
		},
	}
}

// Sample generate over clause
//
// SAMPLE [1,2,3]
// stmt.Sample(1,2,3)
func (stmt *Statement) Sample(sampleList ...int) *Statement {
	if len(sampleList) == 0 {
		return stmt
	}
	stmt.AddClause(&clause.Sample{
		SampleList: sampleList,
	})
	return stmt
}

// Fetch generate fetch clause
//
// FETCH PROP ON player "player101"
// stmt.Fetch("player", "player101")
//
// FETCH PROP ON player "player101", "player102", "player103"
// stmt.Fetch("player", []string{"player101", "player102", "player103"})
//
// FETCH PROP ON serve "player100" -> "team204"
// stmt.Fetch("serve", clause.Expr{Str: `"player100" -> "team204"`})
//
// FETCH PROP ON serve "player100" -> "team204", "player133" -> "team202"
// stmt.Fetch("serve", []*clause.Expr{{Str: `"player100" -> "team204"`}, {Str: `"player133" -> "team202"`}})
func (stmt *Statement) Fetch(name string, vid interface{}) *Statement {
	if name == "" || vid == nil {
		return stmt
	}
	stmt.AddClause(&clause.Fetch{
		Names: []string{name},
		VID:   vid,
	})
	stmt.SetPartType(PartTypeFetch)
	return stmt
}

// FetchMulti generate fetch clause, multi vertex tag name or edge type name
//
// FETCH PROP ON player, serve "player101", "player102", "player103"
// stmt.FetchMulti([]string{"player", "serve"}, []string{"player101", "player102", "player103"})
func (stmt *Statement) FetchMulti(names []string, vid interface{}) *Statement {
	if len(names) == 0 || vid == nil {
		return stmt
	}
	stmt.AddClause(&clause.Fetch{
		Names: names,
		VID:   vid,
	})
	stmt.SetPartType(PartTypeFetch)
	return stmt
}

// Lookup generate lookup clause
//
// LOOKUP ON player
// stmt.Lookup("player")
func (stmt *Statement) Lookup(name string) *Statement {
	if name == "" {
		return stmt
	}
	stmt.AddClause(&clause.Lookup{
		TypeName: name,
	})
	stmt.SetPartType(PartTypeLookup)
	return stmt
}

// GroupBy generate group by clause
//
// GROUP BY $-.Name
// stmt.GroupBy("$-.Name")
func (stmt *Statement) GroupBy(expr string) *Statement {
	// automatically add a pipe character before using the group by statement
	stmt.Pipe()
	stmt.AddClause(&clause.Group{
		Expr: expr,
	})
	stmt.SetPartType(PartTypeGroup)
	return stmt
}

// Yield generate yield clause
//
// YIELD properties(vertex).name, properties(vertex).age
// stmt.Yield("properties(vertex).name, properties(vertex).age")
//
// YIELD DISTINCT properties(vertex).age as v
// stmt.Yield("properties(vertex).age as v", true)
func (stmt *Statement) Yield(expr string, distinct ...bool) *Statement {
	var distinctOpt bool
	if len(distinct) > 0 {
		distinctOpt = distinct[0]
	}
	stmt.AddClause(&clause.Yield{
		Distinct: distinctOpt,
		ExprList: []string{expr},
	})
	return stmt
}

// OrderBy generate order by clause
//
// ORDER BY $-.age ASC, $-.name DESC
// stmt.OrderBy("$-.age ASC, $-.name DESC")
func (stmt *Statement) OrderBy(expr string) *Statement {
	stmt.Pipe()
	stmt.AddClause(&clause.Order{
		Expr: expr,
	})
	stmt.SetPartType(PartTypeOrder)
	return stmt
}

// Limit generate limit clause
//
// LIMIT 1
// stmt.Limit(1)
//
// LIMIT 3, 5
// stmt.Limit(5, 3)
func (stmt *Statement) Limit(limit int, offset ...int) *Statement {
	stmt.Pipe()
	var offsetOpt int
	if len(offset) > 0 {
		offsetOpt = offset[0]
	}
	stmt.AddClause(&clause.Limit{
		Limit:  limit,
		Offset: offsetOpt,
	})
	stmt.SetPartType(PartTypeLimit)
	return stmt
}
