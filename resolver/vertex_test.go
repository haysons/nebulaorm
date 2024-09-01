package resolver

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type prop struct {
	name       string
	index      []int
	nebulaType string
}

func TestParseVertex(t *testing.T) {
	tests := []struct {
		dest                 interface{}
		wantVIDType          VIDType
		wantVIDIndex         int
		wantVIDMethodIndex   int
		wantVIDReceiverIsPtr bool
		wantTag              map[string][]prop
		wantErr              bool
	}{
		{dest: vertex1{}, wantVIDType: VIDTypeString, wantVIDIndex: 2, wantVIDMethodIndex: 0, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag1": {
				{"name", []int{0}, ""},
				{"age", []int{1}, ""},
			},
		}},
		{dest: vertex2{}, wantVIDType: VIDTypeString, wantVIDIndex: -1, wantVIDMethodIndex: 1, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag2": {
				{"name", []int{0}, ""},
				{"age", []int{1}, ""},
				{"gender", []int{2}, "string"},
			},
		}},
		{dest: &vertex3{}, wantVIDType: VIDTypeInt64, wantVIDIndex: -1, wantVIDMethodIndex: 1, wantVIDReceiverIsPtr: false, wantTag: map[string][]prop{
			"vertex_tag3": {
				{"name", []int{0}, ""},
				{"age", []int{1}, ""},
			},
		}},
		{dest: &vertex4{}, wantVIDType: VIDTypeString, wantVIDIndex: 4, wantVIDMethodIndex: 0, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag1": {
				{"name", []int{2, 0}, ""},
				{"age", []int{2, 1}, ""},
			},
			"vertex_tag2": {
				{"name", []int{3, 0}, ""},
				{"age", []int{3, 1}, ""},
				{"gender", []int{3, 2}, "string"},
			},
		}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			got, err := ParseVertex(reflect.TypeOf(tt.dest))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseVertex() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.vidType != tt.wantVIDType {
				t.Errorf("ParseVertex().vidType = %v, want %v", got.vidType, tt.wantVIDType)
				return
			}
			if got.vidFieldIndex != tt.wantVIDIndex {
				t.Errorf("ParseVertex().vidFieldIndex = %v, want %v", got.vidFieldIndex, tt.wantVIDIndex)
				return
			}
			if got.vidMethodIndex != tt.wantVIDMethodIndex {
				t.Errorf("ParseVertex().vidMethodIndex = %v, want %v", got.vidMethodIndex, tt.wantVIDMethodIndex)
				return
			}
			if got.vidReceiverIsPtr != tt.wantVIDReceiverIsPtr {
				t.Errorf("ParseVertex().vidReceiverIsPtr = %v, want %v", got.vidReceiverIsPtr, tt.wantVIDReceiverIsPtr)
				return
			}
			gotTag := make(map[string][]prop)
			for _, tag := range got.tags {
				props := make([]prop, 0)
				for _, p := range tag.props {
					props = append(props, prop{
						name:       p.Name,
						index:      p.StructField.Index,
						nebulaType: p.NebulaType,
					})
				}
				gotTag[tag.TagName] = props
			}
			if !reflect.DeepEqual(gotTag, tt.wantTag) {
				t.Errorf("ParseVertex().tags = %+v, want %+v", gotTag, tt.wantTag)
			}
		})
	}
}

func TestGetVertexInfo(t *testing.T) {
	v1 := vertex4{
		Tag3: vertex1{
			Name: "name11",
		},
		Tag4: &vertex2{
			Name: "name21",
		},
		VID: "v1",
	}
	v2 := &vertex4{
		Tag3: vertex1{
			Name: "name12",
		},
		Tag4: &vertex2{
			Name: "name22",
		},
		VID: "v2",
	}
	v3 := &vertex4{
		Tag3: vertex1{
			Name: "name13",
		},
		Tag4: &vertex2{
			Name: "name23",
		},
		VID: "v3",
	}
	vertexSchema, err := ParseVertex(reflect.TypeOf(v1))
	if err != nil {
		t.Errorf("ParseVertex() error = %v", err)
		return
	}
	tests := []struct {
		v            interface{}
		wantVIDExpr  string
		wantPropExpr string
	}{
		{v: v1, wantVIDExpr: `"v1"`, wantPropExpr: `"name11" "name21"`},
		{v: v2, wantVIDExpr: `"v2"`, wantPropExpr: `"name12" "name22"`},
		{v: v3, wantVIDExpr: `"v3"`, wantPropExpr: `"name13" "name23"`},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			vertexValue := reflect.ValueOf(tt.v)
			vidExpr := vertexSchema.GetVIDExpr(vertexValue)
			if vidExpr != tt.wantVIDExpr {
				t.Errorf("GetVIDExpr = %v, want %v", vidExpr, tt.wantVIDExpr)
				return
			}
			propStr := ""
			vertexValue = reflect.Indirect(vertexValue)
			for _, tag := range vertexSchema.GetTags() {
				for _, p := range tag.GetProps() {
					f := vertexValue.FieldByIndex(p.StructField.Index)
					if !f.IsZero() {
						s, _ := FormatSimpleValue("", f)
						propStr += " " + s
					}
				}
			}
			if strings.TrimSpace(propStr) != tt.wantPropExpr {
				t.Errorf("GetPropExpr = %v, want %v", propStr, tt.wantPropExpr)
			}
		})
	}
}

type vertex1 struct {
	Name     string `norm:"prop:name"`
	Age      int    `norm:"prop:age"`
	VID      string `norm:"vertex_id"`
	gender   int
	Pleasure string `norm:"-"`
}

func (v *vertex1) VertexID() string {
	return v.VID
}

func (v *vertex1) VertexTagName() string {
	return "vertex_tag1"
}

type vertex2 struct {
	Name   string `norm:"prop:name"`
	Age    int
	Gender int `norm:"datatype:string"`
}

func (v *vertex2) A() string {
	return v.Name
}

func (v *vertex2) VertexTagName() string {
	return "vertex_tag2"
}

func (v *vertex2) VertexID() string {
	return v.Name
}

type vertex3 struct {
	Name string `norm:"prop:name"`
	Age  int64
}

func (v vertex3) A() string {
	return v.Name
}

func (v vertex3) VertexTagName() string {
	return "vertex_tag3"
}

func (v vertex3) VertexID() int64 {
	return v.Age
}

type vertex4 struct {
	tag1 *vertex3
	Tag2 *vertex3 `norm:"-"`
	Tag3 vertex1
	Tag4 *vertex2
	VID  string `norm:"vertex_id"`
}

func (v *vertex4) VertexID() string {
	return v.VID
}
