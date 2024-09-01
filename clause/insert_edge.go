package clause

import (
	"fmt"
	"github.com/haysons/nebulaorm/resolver"
	"reflect"
)

type InsertEdge struct {
	IfNotExist bool
	Edges      reflect.Value
	edgeSchema *resolver.EdgeSchema
}

const InsertEdgeName = "INSERT_EDGE"

func (ie InsertEdge) Name() string {
	return InsertEdgeName
}

func (ie InsertEdge) MergeIn(clause *Clause) {
	clause.Expression = ie
}

func (ie InsertEdge) Build(nGQL Builder) error {
	nGQL.WriteString("INSERT EDGE ")
	if ie.IfNotExist {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	ie.Edges = reflect.Indirect(ie.Edges)
	switch ie.Edges.Kind() {
	case reflect.Struct:
		var err error
		edgeType := ie.Edges.Type()
		ie.edgeSchema, err = resolver.ParseEdge(edgeType)
		if err != nil {
			return err
		}
		ie.buildPropNames(nGQL)
		nGQL.WriteString(" VALUES ")
		return ie.buildPropValues(ie.Edges, nGQL)
	case reflect.Slice, reflect.Array:
		var err error
		edgeType := ie.Edges.Type().Elem()
		if edgeType.Kind() == reflect.Pointer {
			edgeType = edgeType.Elem()
		}
		ie.edgeSchema, err = resolver.ParseEdge(edgeType)
		if err != nil {
			return err
		}
		ie.buildPropNames(nGQL)
		nGQL.WriteString(" VALUES ")
		edgesLen := ie.Edges.Len()
		for i := 0; i < edgesLen; i++ {
			curValue := reflect.Indirect(ie.Edges.Index(i))
			if err = ie.buildPropValues(curValue, nGQL); err != nil {
				return err
			}
			if i != edgesLen-1 {
				nGQL.WriteString(", ")
			}
		}
	default:
		return fmt.Errorf("nebulaorm: %w, build insert edge clause failed, dest must be struct, slice, array or pointer", ErrInvalidClauseParams)
	}
	return nil
}

func (ie InsertEdge) buildPropNames(nGQL Builder) {
	nGQL.WriteString(ie.edgeSchema.GetTypeName())
	nGQL.WriteString("(")
	for i, prop := range ie.edgeSchema.GetProps() {
		nGQL.WriteString(prop.Name)
		if i != len(ie.edgeSchema.GetProps())-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteByte(')')
}

func (ie InsertEdge) buildPropValues(curValue reflect.Value, nGQL Builder) error {
	nGQL.WriteString(edgeIDExpr(ie.edgeSchema, curValue))
	nGQL.WriteString(":(")
	props := ie.edgeSchema.GetProps()
	for i, prop := range props {
		valueFmt, err := resolver.FormatSimpleValue(prop.NebulaType, curValue.FieldByIndex(prop.StructField.Index))
		if err != nil {
			return err
		}
		nGQL.WriteString(valueFmt)
		if i != len(props)-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteString(")")
	return nil
}
