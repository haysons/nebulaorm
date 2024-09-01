package clause

import "fmt"

type Order struct {
	Expr string
}

const OrderName = "ORDER"

func (order Order) Name() string {
	return OrderName
}

func (order Order) MergeIn(clause *Clause) {
	clause.Expression = order
}

func (order Order) Build(nGQL Builder) error {
	if order.Expr == "" {
		return fmt.Errorf("nebulaorm: %w, order by expr is empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString("ORDER BY ")
	nGQL.WriteString(order.Expr)
	return nil
}
