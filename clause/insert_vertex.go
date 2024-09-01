package clause

import (
	"fmt"
	"github.com/haysons/nebulaorm/resolver"
	"reflect"
)

type InsertVertex struct {
	IfNotExist   bool
	Vertexes     reflect.Value
	vertexSchema *resolver.VertexSchema
}

const InsertVertexName = "INSERT_VERTEX"

func (iv InsertVertex) Name() string {
	return InsertVertexName
}

func (iv InsertVertex) MergeIn(clause *Clause) {
	clause.Expression = iv
}

func (iv InsertVertex) Build(nGQL Builder) error {
	nGQL.WriteString("INSERT VERTEX ")
	if iv.IfNotExist {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	iv.Vertexes = reflect.Indirect(iv.Vertexes)
	switch iv.Vertexes.Kind() {
	case reflect.Struct:
		var err error
		vertexType := iv.Vertexes.Type()
		iv.vertexSchema, err = resolver.ParseVertex(vertexType)
		if err != nil {
			return err
		}
		iv.buildTagProps(nGQL)
		nGQL.WriteString(" VALUES ")
		return iv.buildPropValue(iv.Vertexes, nGQL)
	case reflect.Slice, reflect.Array:
		var err error
		vertexType := iv.Vertexes.Type().Elem()
		if vertexType.Kind() == reflect.Ptr {
			vertexType = vertexType.Elem()
		}
		iv.vertexSchema, err = resolver.ParseVertex(vertexType)
		if err != nil {
			return err
		}
		iv.buildTagProps(nGQL)
		nGQL.WriteString(" VALUES ")
		vertexesLen := iv.Vertexes.Len()
		for i := 0; i < vertexesLen; i++ {
			curValue := reflect.Indirect(iv.Vertexes.Index(i))
			if err = iv.buildPropValue(curValue, nGQL); err != nil {
				return err
			}
			if i != vertexesLen-1 {
				nGQL.WriteString(", ")
			}
		}
		return nil
	default:
		return fmt.Errorf("nebulaorm: %w, build insert vertex clause failed, dest must be struct, slice, array or pointer", ErrInvalidClauseParams)
	}
}

func (iv InsertVertex) buildTagProps(nGQL Builder) {
	tags := iv.vertexSchema.GetTags()
	for i, t := range tags {
		nGQL.WriteString(t.TagName)
		nGQL.WriteString("(")
		props := t.GetProps()
		for j, p := range props {
			nGQL.WriteString(p.Name)
			if j != len(props)-1 {
				nGQL.WriteString(", ")
			}
		}
		nGQL.WriteString(")")
		if i != len(tags)-1 {
			nGQL.WriteString(", ")
		}
	}
}

func (iv InsertVertex) buildPropValue(curValue reflect.Value, nGQL Builder) error {
	tags := iv.vertexSchema.GetTags()
	vid := iv.vertexSchema.GetVIDExpr(curValue)
	nGQL.WriteString(vid)
	nGQL.WriteString(":(")
	for j, t := range tags {
		props := t.GetProps()
		for k, p := range props {
			valueFmt, err := resolver.FormatSimpleValue(p.NebulaType, curValue.FieldByIndex(p.StructField.Index))
			if err != nil {
				return err
			}
			nGQL.WriteString(valueFmt)
			if k != len(props)-1 {
				nGQL.WriteString(", ")
			}
		}
		if j != len(tags)-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteString(")")
	return nil
}
