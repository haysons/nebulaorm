package clause

import (
	"fmt"
	"testing"
)

func Test_vertexIDExpr(t *testing.T) {
	tests := []struct {
		vid     interface{}
		want    string
		wantErr bool
	}{
		{vid: 101, want: "101"},
		{vid: int64(101), want: "101"},
		{vid: "hayson", want: `"hayson"`},
		{vid: Expr{Str: "hayson"}, want: `hayson`},
		{vid: &Expr{Str: "hayson"}, want: `hayson`},
		{vid: []int{1, 2, 3}, want: `1, 2, 3`},
		{vid: []int64{1, 2, 3}, want: `1, 2, 3`},
		{vid: []string{"n1", "n2"}, want: `"n1", "n2"`},
		{vid: []Expr{{Str: "expr1"}, {Str: "expr2"}}, want: `expr1, expr2`},
		{vid: []*Expr{{Str: "expr1"}, {Str: "expr2"}}, want: `expr1, expr2`},
		{vid: nil, wantErr: true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			got, err := vertexIDExpr(tt.vid)
			if (err != nil) != tt.wantErr {
				t.Errorf("vertexIDExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("vertexIDExpr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
