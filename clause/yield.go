package clause

import "fmt"

type Yield struct {
	Distinct bool
	ExprList []string
}

const YieldName = "YIELD"

func (y Yield) Name() string {
	return YieldName
}

func (y Yield) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(Yield)
	if !ok {
		clause.Expression = y
		return
	}
	exist.ExprList = append(exist.ExprList, y.ExprList...)
	exist.Distinct = y.Distinct
	clause.Expression = exist
}

func (y Yield) Build(nGQL Builder) error {
	exprList := make([]string, 0, len(y.ExprList))
	for _, expr := range y.ExprList {
		if expr != "" {
			exprList = append(exprList, expr)
		}
	}
	if len(exprList) == 0 {
		return fmt.Errorf("nebulaorm: %w, yield expr is empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString("YIELD ")
	if y.Distinct {
		nGQL.WriteString("DISTINCT ")
	}
	for i, expr := range exprList {
		nGQL.WriteString(expr)
		if i != len(exprList)-1 {
			nGQL.WriteString(", ")
		}
	}
	return nil
}
