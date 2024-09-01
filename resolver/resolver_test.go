package resolver

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFormatSimpleValue(t *testing.T) {
	a := "hello"
	tests := []struct {
		nebulaType string
		value      []interface{}
		want       string
		wantErr    bool
	}{
		{
			value: []interface{}{1, int8(1), int16(1), int32(1), int64(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1)},
			want:  "1",
		},
		{
			value: []interface{}{-1, int8(-1), int16(-1), int32(-1), int64(-1)},
			want:  "-1",
		},
		{
			value: []interface{}{1000.234, float32(1000.234)},
			want:  "1000.234",
		},
		{
			value: []interface{}{-1000.234, float32(-1000.234)},
			want:  "-1000.234",
		},
		{
			nebulaType: NebulaDataTypeFloat,
			value:      []interface{}{100},
			want:       "100",
		},
		{
			nebulaType: NebulaDataTypeInt,
			value:      []interface{}{100.1234},
			want:       "100",
		},
		{
			value: []interface{}{true},
			want:  "true",
		},
		{
			value: []interface{}{false},
			want:  "false",
		},
		{
			value: []interface{}{"hello world 你好 世界"},
			want:  `"hello world 你好 世界"`,
		},
		{
			value: []interface{}{"Hello \\ world 你好 \" \t t 世界"},
			want:  `"Hello \\ world 你好 \" \t t 世界"`,
		},
		{
			value: []interface{}{`Hello \ world 你好 " t 世界`},
			want:  `"Hello \\ world 你好 \" t 世界"`,
		},
		{
			nebulaType: NebulaDataTypeDatetime,
			value:      []interface{}{`2024-08-01T00:00:00`},
			want:       `datetime("2024-08-01T00:00:00")`,
		},
		{
			nebulaType: NebulaDataTypeDate,
			value:      []interface{}{`2023-12-12`},
			want:       `date("2023-12-12")`,
		},
		{
			nebulaType: NebulaDataTypeTime,
			value:      []interface{}{`11:00:51.457000`},
			want:       `time("11:00:51.457000")`,
		},
		{
			value: []interface{}{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local)},
			want:  `datetime("2024-08-20T11:16:30.000010")`,
		},
		{
			nebulaType: NebulaDataTypeDate,
			value:      []interface{}{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local)},
			want:       `date("2024-08-20")`,
		},
		{
			nebulaType: NebulaDataTypeTime,
			value:      []interface{}{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local)},
			want:       `time("11:16:30.000010")`,
		},
		{
			value: []interface{}{[]int{1, -1, 2, -2, 0}},
			want:  "[1, -1, 2, -2, 0]",
		},
		{
			value: []interface{}{make([]int, 0)},
			want:  "[]",
		},
		{
			value: []interface{}{[]string{"h", "e", "l", "l", "o"}},
			want:  `["h", "e", "l", "l", "o"]`,
		},
		{
			value: []interface{}{[]time.Time{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local), time.Date(2024, 8, 20, 11, 16, 16, 0, time.Local)}},
			want:  `[datetime("2024-08-20T11:16:30.000010"), datetime("2024-08-20T11:16:16")]`,
		},
		{
			value: []interface{}{[5]string{"h", "e", "l", "l", "o"}},
			want:  `["h", "e", "l", "l", "o"]`,
		},
		{
			nebulaType: NebulaDataTypeSet,
			value:      []interface{}{[]int{1, -1, 2, -2, 0}, []int64{1, -1, 2, -2, 0}},
			want:       "set{1, -1, 2, -2, 0}",
		},
		{
			nebulaType: NebulaDataTypeSet,
			value:      []interface{}{[]int(nil)},
			want:       "set{}",
		},
		{
			value: []interface{}{map[string]int{"c": 3}},
			want:  `map{c: 3}`,
		},
		{
			value: []interface{}{map[string]interface{}{"d": map[string]int{"age": 18}}},
			want:  `map{d: map{age: 18}}`,
		},
		{
			nebulaType: NebulaDataTypeSet,
			value:      []interface{}{map[int]struct{}{1: {}}},
			want:       `set{1}`,
		},
		{
			value: []interface{}{(*int)(nil)},
			want:  `NULL`,
		},
		{
			value: []interface{}{&a},
			want:  `"hello"`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			for _, v := range tt.value {
				got, err := FormatSimpleValue(tt.nebulaType, reflect.ValueOf(v))
				if (err != nil) != tt.wantErr {
					t.Errorf("FormatSimpleValue() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("FormatSimpleValue() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
