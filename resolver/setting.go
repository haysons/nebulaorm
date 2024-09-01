package resolver

import (
	"reflect"
	"strings"
	"time"
	"unicode"
)

const (
	TagSettingKey       = "norm"        // nebulaorm struct tag key
	TagSettingColName   = "col"         // name of the field in the record
	TagSettingVertexID  = "vertex_id"   // annotate that the field is a vertex id
	TagSettingEdgeSrcID = "edge_src_id" // annotate that the field is an edge source id
	TagSettingEdgeDstID = "edge_dst_id" // annotate that the field is an edge dest id
	TagSettingEdgeRank  = "edge_rank"   // annotate that the field is an edge rank
	TagSettingPropName  = "prop"        // property name, vertex or edge
	TagSettingDataType  = "datatype"    // specify the data type (in this case the data type specified in github.com/vesoft-inc/nebula-go/v3)
	TagSettingIgnore    = "-"           // nebulaorm will ignore this field
)

func ParseTagSetting(s string) map[string]string {
	m := make(map[string]string)
	tags := strings.Split(s, ";")
	for _, tag := range tags {
		kv := strings.Split(tag, ":")
		k := strings.TrimSpace(strings.ToLower(kv[0]))
		if k == "" {
			continue
		}
		if len(kv) >= 2 {
			m[k] = strings.Join(kv[1:], ":")
		} else {
			m[k] = k
		}
	}
	return m
}

func GetPropName(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	propName := setting[TagSettingPropName]
	if propName == "" {
		propName = camelCaseToUnderscore(field.Name)
	}
	return propName
}

func GetValueNebulaType(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	return setting[TagSettingDataType]
}

func FieldIgnore(field reflect.StructField) bool {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	return setting[TagSettingIgnore] != ""
}

func camelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}

var timezoneDefault = time.Local

func SetTimezone(loc *time.Location) {
	if loc != nil {
		timezoneDefault = loc
	}
}
