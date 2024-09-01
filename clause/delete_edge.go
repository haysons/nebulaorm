package clause

import (
	"fmt"
	"github.com/haysons/nebulaorm/resolver"
	"reflect"
	"strconv"
)

type DeleteEdge struct {
	EdgeTypeName string
	Edges        interface{}
}

const DeleteEdgeName = "DELETE_EDGE"

func (de DeleteEdge) Name() string {
	return DeleteEdgeName
}

func (de DeleteEdge) MergeIn(clause *Clause) {
	clause.Expression = de
}

func (de DeleteEdge) Build(nGQL Builder) error {
	nGQL.WriteString("DELETE EDGE ")
	if de.EdgeTypeName == "" {
		return fmt.Errorf("nebulaorm: %w, build delete_edge clause failed, edge type name is empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString(de.EdgeTypeName)
	nGQL.WriteByte(' ')
	edgeList := make([]string, 0)
	switch edge := de.Edges.(type) {
	case string:
		edgeList = append(edgeList, edge)
	case []string:
		edgeList = edge
	default:
		edgeValue := reflect.Indirect(reflect.ValueOf(edge))
		if !edgeValue.IsValid() {
			return fmt.Errorf("nebulaorm: %w, build delete_edge clause failed, edge must be a string, []string, edge, edge slice or edge array", ErrInvalidClauseParams)
		}
		edgeType := edgeValue.Type()
		switch edgeType.Kind() {
		case reflect.Struct:
			edgeSchema, err := resolver.ParseEdge(edgeType)
			if err != nil {
				return err
			}
			edgeList = append(edgeList, edgeIDExpr(edgeSchema, edgeValue))
		case reflect.Slice, reflect.Array:
			edgeType = edgeType.Elem()
			if edgeType.Kind() == reflect.Ptr {
				edgeType = edgeType.Elem()
			}
			if edgeType.Kind() == reflect.Struct {
				edgeSchema, err := resolver.ParseEdge(edgeType)
				if err != nil {
					return err
				}
				for i := 0; i < edgeValue.Len(); i++ {
					curValue := reflect.Indirect(edgeValue.Index(i))
					edgeList = append(edgeList, edgeIDExpr(edgeSchema, curValue))
				}
			} else {
				return fmt.Errorf("nebulaorm: %w, build delete_edge clause failed, slice element must be a struct or a struct pointer", ErrInvalidClauseParams)
			}
		default:
			return fmt.Errorf("nebulaorm: %w, build delete_edge clause failed, edge must be a string, []string, edge, edge slice or edge array", ErrInvalidClauseParams)
		}
	}
	if len(edgeList) == 0 {
		return fmt.Errorf("nebulaorm: %w, build delete_edge clause failed, edge list is empty", ErrInvalidClauseParams)
	}
	for i, e := range edgeList {
		nGQL.WriteString(e)
		if i < len(edgeList)-1 {
			nGQL.WriteString(", ")
		}
	}
	return nil
}

func edgeIDExpr(edgeSchema *resolver.EdgeSchema, edgeValue reflect.Value) string {
	srcID := edgeSchema.GetSrcVIDExpr(edgeValue)
	dstID := edgeSchema.GetDstVIDExpr(edgeValue)
	rank := edgeSchema.GetRank(edgeValue)
	edgeStr := srcID + "->" + dstID
	if rank > 0 {
		edgeStr += "@" + strconv.Itoa(int(rank))
	}
	return edgeStr
}
