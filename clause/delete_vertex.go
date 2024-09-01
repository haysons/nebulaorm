package clause

import (
	"fmt"
)

type DeleteVertex struct {
	VID      interface{}
	WithEdge bool
}

const DeleteVertexName = "DELETE_VERTEX"

func (dv DeleteVertex) Name() string {
	return DeleteVertexName
}

func (dv DeleteVertex) MergeIn(clause *Clause) {
	clause.Expression = dv
}

func (dv DeleteVertex) Build(nGQL Builder) error {
	nGQL.WriteString("DELETE VERTEX ")
	vidExpr, err := vertexIDExpr(dv.VID)
	if err != nil {
		return fmt.Errorf("nebulaorm: %w, build delete_vertex clause failed, %v", ErrInvalidClauseParams, err)
	}
	nGQL.WriteString(vidExpr)
	if dv.WithEdge {
		nGQL.WriteString(" WITH EDGE")
	}
	return nil
}
