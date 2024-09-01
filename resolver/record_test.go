package resolver

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseRecord(t *testing.T) {
	tests := []struct {
		record  interface{}
		want    *RecordSchema
		wangErr bool
	}{
		{
			record: record1{},
			want:   &RecordSchema{Name: "record1", colFieldIndex: map[string][]int{"name": {0}, "age": {1}, "c": {3}}},
		},
		{
			record: record2{},
			want:   &RecordSchema{Name: "record2", colFieldIndex: map[string][]int{"col1": {1}, "names": {2}}},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			recordType := reflect.TypeOf(tt.record)
			schemaGot, err := ParseRecord(recordType)
			if (err != nil) != tt.wangErr {
				t.Errorf("ParseRecord() error = %v, wantErr %v", err, tt.wangErr)
				return
			}
			if !reflect.DeepEqual(schemaGot, tt.want) {
				t.Errorf("ParseRecord() got = %#v, want %#v", schemaGot, tt.want)
			}
		})
	}
}

type record1 struct {
	Name     string
	Age      int
	gender   int
	Class    string `norm:"col:c"`
	Pleasure string `norm:"-"`
}

type record2 struct {
	record1
	Col1  *record1 `norm:"col:col1"`
	Names []string `norm:"col:names"`
}
