package resolver

import (
	"errors"
	"fmt"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"reflect"
	"strconv"
)

// EdgeTypeNamer specifies the name of the edge type. a structure that implements this interface will be treated as an edge.
type EdgeTypeNamer interface {
	EdgeTypeName() string
}

type EdgeSchema struct {
	srcVIDType       VIDType
	srcVIDFieldIndex int
	dstVIDType       VIDType
	dstVIDFieldIndex int
	edgeTypeName     string
	rankFieldIndex   int
	props            []*Prop
	propByName       map[string]*Prop
}

// ParseEdge parse edge struct
func ParseEdge(destType reflect.Type) (*EdgeSchema, error) {
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	if destType.Kind() != reflect.Struct {
		return nil, errors.New("nebulaorm: parse edge failed, dest should be a struct or a struct pointer")
	}
	edge := &EdgeSchema{
		srcVIDFieldIndex: -1,
		dstVIDFieldIndex: -1,
		rankFieldIndex:   -1,
		props:            make([]*Prop, 0),
		propByName:       make(map[string]*Prop),
	}
	// whether it implements the EdgeTypeNamer interface
	destValue := reflect.New(destType).Interface()
	edgeTypeNamer, ok := destValue.(EdgeTypeNamer)
	if !ok {
		return nil, errors.New("nebulaorm: parse edge failed, need to implement interface resolver.EdgeTypeNamer")
	}
	edge.edgeTypeName = edgeTypeNamer.EdgeTypeName()
	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		if field.Anonymous || !field.IsExported() {
			continue
		}
		setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
		if _, ok := setting[TagSettingIgnore]; ok {
			continue
		}
		if _, isSrcID := setting[TagSettingEdgeSrcID]; isSrcID {
			switch field.Type.Kind() {
			case reflect.String:
				edge.srcVIDType = VIDTypeString
			case reflect.Int64:
				edge.srcVIDType = VIDTypeInt64
			default:
				return nil, errors.New("nebulaorm: parse edge failed, src_id field should be a string or int64")
			}
			edge.srcVIDFieldIndex = i
			continue
		}
		if _, isDstID := setting[TagSettingEdgeDstID]; isDstID {
			switch field.Type.Kind() {
			case reflect.String:
				edge.dstVIDType = VIDTypeString
			case reflect.Int64:
				edge.dstVIDType = VIDTypeInt64
			default:
				return nil, errors.New("nebulaorm: parse edge failed, dst_id field should be a string or int64")
			}
			edge.dstVIDFieldIndex = i
			continue
		}
		if _, isRank := setting[TagSettingEdgeRank]; isRank {
			if !(field.Type.Kind() == reflect.Int64 || field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.Int32 || field.Type.Kind() == reflect.Int8 || field.Type.Kind() == reflect.Int16) {
				return nil, errors.New("nebulaorm: parse edge failed, rank field should be int")
			}
			edge.rankFieldIndex = i
			continue
		}
		// parsing Edge Properties
		propName := GetPropName(field)
		nebulaType := GetValueNebulaType(field)
		prop := &Prop{
			Name:        propName,
			StructField: field,
			Type:        field.Type,
			NebulaType:  nebulaType,
		}
		if _, ok = edge.propByName[propName]; ok {
			continue
		}
		edge.props = append(edge.props, prop)
		edge.propByName[propName] = prop
	}
	if edge.srcVIDFieldIndex < 0 || edge.dstVIDFieldIndex < 0 {
		return nil, errors.New("nebulaorm: parse edge failed, edge must contains src_id field and dst_id field")
	}
	return edge, nil
}

// GetTypeName get edge type name
func (e *EdgeSchema) GetTypeName() string {
	return e.edgeTypeName
}

// GetSrcVID get the src_id of the edge
func (e *EdgeSchema) GetSrcVID(edgeValue reflect.Value) interface{} {
	if e.srcVIDFieldIndex >= 0 {
		edgeValue = reflect.Indirect(edgeValue)
		return edgeValue.Field(e.srcVIDFieldIndex).Interface()
	}
	return nil
}

// GetSrcVIDExpr get the src_id expr of the edge
func (e *EdgeSchema) GetSrcVIDExpr(edgeValue reflect.Value) string {
	srcID := e.GetSrcVID(edgeValue)
	if srcID == nil {
		return ""
	}
	switch e.srcVIDType {
	case VIDTypeString:
		return strconv.Quote(srcID.(string))
	case VIDTypeInt64:
		return strconv.FormatInt(srcID.(int64), 10)
	}
	return ""
}

// GetDstVID get the dst_id of the edge
func (e *EdgeSchema) GetDstVID(edgeValue reflect.Value) interface{} {
	if e.dstVIDFieldIndex >= 0 {
		edgeValue = reflect.Indirect(edgeValue)
		return edgeValue.Field(e.dstVIDFieldIndex).Interface()
	}
	return nil
}

// GetDstVIDExpr get the dst_id expr of the edge
func (e *EdgeSchema) GetDstVIDExpr(edgeValue reflect.Value) string {
	dstID := e.GetDstVID(edgeValue)
	if dstID == nil {
		return ""
	}
	switch e.dstVIDType {
	case VIDTypeString:
		return strconv.Quote(dstID.(string))
	case VIDTypeInt64:
		return strconv.FormatInt(dstID.(int64), 10)
	}
	return ""
}

// GetRank get the rank value of the edge
func (e *EdgeSchema) GetRank(edgeValue reflect.Value) int64 {
	if e.rankFieldIndex >= 0 {
		edgeValue = reflect.Indirect(edgeValue)
		return edgeValue.Field(e.rankFieldIndex).Int()
	}
	return 0
}

// GetProps get a list of attributes for the current edge
func (e *EdgeSchema) GetProps() []*Prop {
	return e.props
}

// Scan assign a value to a target struct
func (e *EdgeSchema) Scan(rl *nebula.Relationship, destValue reflect.Value) error {
	destValue = reflect.Indirect(destValue)
	if !destValue.CanSet() {
		return fmt.Errorf("nebulaorm: edge schema scan dest value failed, %w", ErrValueCannotSet)
	}
	if e.srcVIDFieldIndex >= 0 {
		srcID := rl.GetSrcVertexID()
		if err := ScanSimpleValue(&srcID, destValue.Field(e.srcVIDFieldIndex)); err != nil {
			return err
		}
	}
	if e.dstVIDFieldIndex >= 0 {
		dstID := rl.GetDstVertexID()
		if err := ScanSimpleValue(&dstID, destValue.Field(e.dstVIDFieldIndex)); err != nil {
			return err
		}
	}
	if e.rankFieldIndex >= 0 {
		rank := rl.GetRanking()
		destValue.Field(e.rankFieldIndex).SetInt(rank)
	}
	for propName, propValue := range rl.Properties() {
		eProp, ok := e.propByName[propName]
		if !ok {
			continue
		}
		if err := ScanSimpleValue(propValue, destValue.FieldByIndex(eProp.StructField.Index)); err != nil {
			return err
		}
	}
	return nil
}
