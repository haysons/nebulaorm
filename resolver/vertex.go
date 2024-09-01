package resolver

import (
	"errors"
	"fmt"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"reflect"
	"strconv"
)

// VertexIDStr a structure that implements this interface is treated as a vertex and has a vertex_id of type string
type VertexIDStr interface {
	VertexID() string
}

// VertexIDInt64 a structure that implements this interface is treated as a vertex and has a vertex_id of type int64
type VertexIDInt64 interface {
	VertexID() int64
}

// VertexTagNamer the structure that implements this interface is treated as a tag of the vertex, and if the vertex has only one tag,
// it can be implemented by the same structure at the same time VertexIDStr(VertexIDInt64) interface and the VertexTagNamer interface,
// which defines the attributes of the vertex_id and the various tags in the same structure.
type VertexTagNamer interface {
	VertexTagName() string
}

type VIDType int

const (
	VIDTypeString VIDType = iota + 1
	VIDTypeInt64
)

type VertexSchema struct {
	tags             []*VertexTag
	tagByName        map[string]*VertexTag
	vidType          VIDType
	vidFieldIndex    int
	vidMethodIndex   int
	vidReceiverIsPtr bool
}

func ParseVertex(destType reflect.Type) (*VertexSchema, error) {
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	if destType.Kind() != reflect.Struct {
		return nil, errors.New("nebulaorm: parse vertex failed, dest should be a struct or a struct pointer")
	}
	vertex := &VertexSchema{
		vidFieldIndex:  -1,
		vidMethodIndex: -1,
		tagByName:      make(map[string]*VertexTag),
	}
	if err := vertex.parseVID(destType); err != nil {
		return nil, err
	}
	if err := vertex.parseTag(destType, -1); err != nil {
		return nil, err
	}
	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		if field.Anonymous || !field.IsExported() || FieldIgnore(field) {
			continue
		}
		if err := vertex.parseTag(destType.Field(i).Type, i); err != nil {
			return nil, err
		}
	}
	if len(vertex.tags) == 0 {
		return nil, errors.New("nebulaorm: parse vertex failed, vertex has no tags")
	}
	return vertex, nil
}

func (v *VertexSchema) parseVID(vertexType reflect.Type) error {
	vertexIface := reflect.New(vertexType).Interface()
	if _, ok := vertexIface.(VertexIDStr); ok {
		v.vidType = VIDTypeString
	} else if _, ok := vertexIface.(VertexIDInt64); ok {
		v.vidType = VIDTypeInt64
	} else {
		return fmt.Errorf("nebulaorm: parse vertex failed, need to implement interface resolver.VertexIDStr or resolver.VertexIDInt64")
	}
	if vidMethod, ok := vertexType.MethodByName("VertexID"); ok {
		v.vidMethodIndex = vidMethod.Index
		v.vidReceiverIsPtr = false
	} else if vidMethod, ok := reflect.PointerTo(vertexType).MethodByName("VertexID"); ok {
		v.vidMethodIndex = vidMethod.Index
		v.vidReceiverIsPtr = true
	} else {
		return fmt.Errorf("nebulaorm: parse vertex failed, cannnot get vertex_id method")
	}
	// if the struct contains a vid field, save it and assign it to it when scanning
	for i := 0; i < vertexType.NumField(); i++ {
		field := vertexType.Field(i)
		if field.Anonymous || !field.IsExported() {
			continue
		}
		setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
		if _, ok := setting[TagSettingVertexID]; ok {
			v.vidFieldIndex = i
			break
		}
	}
	return nil
}

func (v *VertexSchema) parseTag(destType reflect.Type, superIndex int) error {
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	if destType.Kind() != reflect.Struct {
		return nil
	}
	// if the structure implements the VertexTagNamer interface, treat the structure as a tag and get the name of its tag
	destValue := reflect.New(destType).Interface()
	tagNamer, ok := destValue.(VertexTagNamer)
	if !ok {
		return nil
	}
	tagName := tagNamer.VertexTagName()
	if _, ok := v.tagByName[tagName]; !ok {
		tag := &VertexTag{
			TagName:    tagName,
			props:      make([]*Prop, 0),
			propByName: make(map[string]*Prop),
		}
		v.tagByName[tagName] = tag
		v.tags = append(v.tags, tag)
	}
	// parse each property of the current tag
	for i := 0; i < destType.NumField(); i++ {
		structField := destType.Field(i)
		if structField.Anonymous || !structField.IsExported() {
			continue
		}
		setting := ParseTagSetting(structField.Tag.Get(TagSettingKey))
		if _, ok := setting[TagSettingIgnore]; ok {
			continue
		}
		if _, ok := setting[TagSettingVertexID]; ok {
			continue
		}
		propName := GetPropName(structField)
		nebulaType := GetValueNebulaType(structField)
		// tag may exist in a multi-level structure, the index value of the field needs to be added to the index value of the parent field
		if superIndex >= 0 {
			structField.Index = append([]int{superIndex}, structField.Index...)
		}
		prop := &Prop{
			Name:        propName,
			StructField: structField,
			Type:        structField.Type,
			NebulaType:  nebulaType,
		}
		if _, ok := v.tagByName[tagName].propByName[propName]; ok {
			continue
		}
		v.tagByName[tagName].props = append(v.tagByName[tagName].props, prop)
		v.tagByName[tagName].propByName[propName] = prop
	}
	return nil
}

// GetVID get the vid value of vertexValue
func (v *VertexSchema) GetVID(vertexValue reflect.Value) interface{} {
	if v.vidReceiverIsPtr && vertexValue.Kind() != reflect.Ptr {
		vertexNew := reflect.New(vertexValue.Type())
		vertexNew.Elem().Set(vertexValue)
		vertexValue = vertexNew
	}
	out := vertexValue.Method(v.vidMethodIndex).Call([]reflect.Value{})
	return out[0].Interface()
}

// GetVIDType get the vid type of the vertex
func (v *VertexSchema) GetVIDType() VIDType {
	return v.vidType
}

// GetVIDExpr get the string expression of vid
func (v *VertexSchema) GetVIDExpr(vertexValue reflect.Value) string {
	vid := v.GetVID(vertexValue)
	if vid == nil {
		return ""
	}
	switch v.vidType {
	case VIDTypeString:
		return strconv.Quote(vid.(string))
	case VIDTypeInt64:
		return fmt.Sprintf("%d", vid)
	}
	return ""
}

// GetTags get a list of the vertex's tags
func (v *VertexSchema) GetTags() []*VertexTag {
	return v.tags
}

// Scan assign the nodes returned by the nebula graph to the vertex data in the business layer
func (v *VertexSchema) Scan(node *nebula.Node, destValue reflect.Value) error {
	// schema parsing and assignment can support structs or struct pointers
	destValue = reflect.Indirect(destValue)
	if !destValue.CanSet() {
		return fmt.Errorf("nebulaorm: vertex schema scan dest value failed, %w", ErrValueCannotSet)
	}
	// if a vid field exists in the structure, it is assigned to it
	if v.vidFieldIndex >= 0 {
		vid := node.GetID()
		if err := ScanSimpleValue(&vid, destValue.Field(v.vidFieldIndex)); err != nil {
			return err
		}
	}
	for _, vTag := range v.GetTags() {
		propValueMap, err := node.Properties(vTag.TagName)
		if err != nil {
			return err
		}
		for _, prop := range vTag.GetProps() {
			propValue, ok := propValueMap[prop.Name]
			if !ok {
				continue
			}
			if err = ScanSimpleValue(propValue, destValue.FieldByIndex(prop.StructField.Index)); err != nil {
				return err
			}
		}
	}
	return nil
}

type VertexTag struct {
	TagName    string
	props      []*Prop
	propByName map[string]*Prop // key: prop name
}

type Prop struct {
	Name        string
	StructField reflect.StructField
	Type        reflect.Type
	NebulaType  string
}

// GetProps get all attributes of the tag
func (t *VertexTag) GetProps() []*Prop {
	return t.props
}
