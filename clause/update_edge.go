package clause

import (
	"fmt"
	"github.com/haysons/nebulaorm/resolver"
	"reflect"
)

type UpdateEdge struct {
	IsUpsert    bool
	Edge        interface{}
	PropsUpdate interface{}
	Opts        Options
}

const UpdateEdgeName = "UPDATE_EDGE"

func (ue UpdateEdge) Name() string {
	return UpdateEdgeName
}

func (ue UpdateEdge) MergeIn(clause *Clause) {
	clause.Expression = ue
}

func (ue UpdateEdge) Build(nGQL Builder) error {
	if ue.IsUpsert {
		nGQL.WriteString("UPSERT EDGE ON ")
	} else {
		nGQL.WriteString("UPDATE EDGE ON ")
	}
	var edgeStr string
	switch edge := ue.Edge.(type) {
	case string:
		edgeStr = edge
	default:
		edgeValue := reflect.Indirect(reflect.ValueOf(ue.Edge))
		edgeType := edgeValue.Type()
		switch edgeType.Kind() {
		case reflect.Struct:
			edgeSchema, err := resolver.ParseEdge(edgeType)
			if err != nil {
				return err
			}
			edgeStr = edgeSchema.GetTypeName() + " " + edgeIDExpr(edgeSchema, edgeValue)
		default:
			return fmt.Errorf("nebulaorm: %w, build update edge clause failed, dest edge must be struct or struct pointer", ErrInvalidClauseParams)
		}
	}
	// manually specify the name of the property to be updated
	propsName := make(map[string]bool, len(ue.Opts.propNames))
	for _, propName := range ue.Opts.propNames {
		propsName[propName] = true
	}
	propsUpdate, err := getPropsUpdateSet(ue.PropsUpdate, propsName)
	if err != nil {
		return fmt.Errorf("nebulaorm: %w, build update edge clause failed, %v", ErrInvalidClauseParams, err)
	}
	if edgeStr == "" {
		return fmt.Errorf("nebulaorm: %w, build update edge clause failed, the edge want to update empty", ErrInvalidClauseParams)
	}
	if len(propsUpdate) == 0 {
		return fmt.Errorf("nebulaorm: %w, build update edge clause failed, the values want to update empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString(edgeStr)
	nGQL.WriteString(" SET ")
	for i, update := range propsUpdate {
		nGQL.WriteString(update[0])
		nGQL.WriteString(" = ")
		nGQL.WriteString(update[1])
		if i < len(propsUpdate)-1 {
			nGQL.WriteString(", ")
		}
	}
	return nil
}
