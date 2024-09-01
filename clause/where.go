package clause

import "strings"

type Where struct {
	Conditions []Condition
}

type Condition struct {
	Operator string
	Expr     Expr
}

const (
	WhereName   = "WHERE"
	OperatorAnd = "AND"
	OperatorOr  = "OR"
	OperatorNot = "NOT"
	OperatorXor = "XOR"
)

func (where Where) Name() string {
	return WhereName
}

func (where Where) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(Where)
	if !ok {
		clause.Expression = where
		return
	}
	exist.Conditions = append(exist.Conditions, where.Conditions...)
	clause.Expression = exist
}

func (where Where) Build(nGQL Builder) error {
	nGQL.WriteString("WHERE ")
	return buildConditions(where.Conditions, nGQL)
}

func buildConditions(conditions []Condition, nGQL Builder) error {
	for i, expr := range conditions {
		if i > 0 {
			nGQL.WriteString(expr.Operator)
			nGQL.WriteByte(' ')
		}
		gql := strings.ToUpper(expr.Expr.Str)
		if strings.Contains(gql, " AND ") || strings.Contains(gql, " OR ") || strings.Contains(gql, " NOT ") || strings.Contains(gql, " XOR ") {
			nGQL.WriteByte('(')
			if err := expr.Expr.Build(nGQL); err != nil {
				return err
			}
			nGQL.WriteByte(')')
		} else {
			if err := expr.Expr.Build(nGQL); err != nil {
				return err
			}
		}
		if i < len(conditions)-1 {
			nGQL.WriteByte(' ')
		}
	}
	return nil
}
