package resolver

import (
	"errors"
	"reflect"
)

// RecordSchema parses the record structure provided by the business layer for subsequent assignment of the Record \
// object returned by the nebula graph to the business layer
type RecordSchema struct {
	Name          string
	colFieldIndex map[string][]int
}

func ParseRecord(destType reflect.Type) (*RecordSchema, error) {
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	if destType.Kind() != reflect.Struct {
		return nil, errors.New("nebulaorm: parse record schema failed, dest should be a struct or struct pointer")
	}
	record := &RecordSchema{
		Name:          destType.Name(),
		colFieldIndex: make(map[string][]int),
	}
	for i := 0; i < destType.NumField(); i++ {
		structField := destType.Field(i)
		if structField.Anonymous || !structField.IsExported() {
			continue
		}
		if FieldIgnore(structField) {
			continue
		}
		colName := getColName(structField)
		record.colFieldIndex[colName] = []int{i}
	}
	return record, nil
}

// GetFieldIndexByColName get the index position of a field
func (r *RecordSchema) GetFieldIndexByColName(colName string) []int {
	return r.colFieldIndex[colName]
}

func getColName(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	colName := setting[TagSettingColName]
	if colName == "" {
		colName = camelCaseToUnderscore(field.Name)
	}
	return colName
}
