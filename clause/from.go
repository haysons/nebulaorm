package clause

import (
	"fmt"
)

type From struct {
	VID interface{}
}

const FromName = "FROM"

func (from From) Name() string {
	return FromName
}

func (from From) MergeIn(clause *Clause) {
	clause.Expression = from
}

func (from From) Build(nGQL Builder) error {
	nGQL.WriteString("FROM ")
	vidExpr, err := vertexIDExpr(from.VID)
	if err != nil {
		return fmt.Errorf("nebulaorm: %w, build from clause failed, %v", ErrInvalidClauseParams, err)
	}
	nGQL.WriteString(vidExpr)
	return nil
}
