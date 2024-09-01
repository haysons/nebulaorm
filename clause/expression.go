package clause

import (
	"errors"
	"github.com/haysons/nebulaorm/resolver"
	"reflect"
	"strconv"
	"strings"
)

// Expression expression interface
type Expression interface {
	Build(nGQL Builder) error
}

type Builder interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
}

// Expr raw expression
type Expr struct {
	Str  string
	Vars []interface{}
}

// Build raw expression
func (expr Expr) Build(builder Builder) error {
	var idx int
	for _, v := range []byte(expr.Str) {
		if v == '?' && len(expr.Vars) > idx {
			valFmt, err := expr.formatValue(expr.Vars[idx])
			if err != nil {
				return err
			}
			builder.WriteString(valFmt)
			idx++
		} else {
			builder.WriteByte(v)
		}
	}
	if idx < len(expr.Vars) {
		for _, v := range expr.Vars[idx:] {
			valFmt, err := expr.formatValue(v)
			if err != nil {
				return err
			}
			builder.WriteString(valFmt)
		}
	}
	return nil
}

func (expr Expr) formatValue(value interface{}) (string, error) {
	switch v := value.(type) {
	case Expr:
		exprBuilder := new(strings.Builder)
		err := v.Build(exprBuilder)
		if err != nil {
			return "", err
		}
		return exprBuilder.String(), nil
	case *Expr:
		exprBuilder := new(strings.Builder)
		err := v.Build(exprBuilder)
		if err != nil {
			return "", err
		}
		return exprBuilder.String(), nil
	default:
		return resolver.FormatSimpleValue("", reflect.ValueOf(value))
	}
}

func vertexIDExpr(vid interface{}) (string, error) {
	vidList := make([]string, 0)
	switch id := vid.(type) {
	case int:
		vidList = append(vidList, strconv.Itoa(id))
	case int64:
		vidList = append(vidList, strconv.FormatInt(id, 10))
	case string:
		vidList = append(vidList, strconv.Quote(id))
	case Expr:
		exprBuilder := new(strings.Builder)
		err := id.Build(exprBuilder)
		if err != nil {
			return "", err
		}
		vidList = append(vidList, exprBuilder.String())
	case *Expr:
		exprBuilder := new(strings.Builder)
		err := id.Build(exprBuilder)
		if err != nil {
			return "", err
		}
		vidList = append(vidList, exprBuilder.String())
	case []int:
		for _, v := range id {
			vidList = append(vidList, strconv.Itoa(v))
		}
	case []int64:
		for _, v := range id {
			vidList = append(vidList, strconv.FormatInt(v, 10))
		}
	case []string:
		for _, v := range id {
			vidList = append(vidList, strconv.Quote(v))
		}
	case []Expr:
		for _, v := range id {
			exprBuilder := new(strings.Builder)
			err := v.Build(exprBuilder)
			if err != nil {
				return "", err
			}
			vidList = append(vidList, exprBuilder.String())
		}
	case []*Expr:
		for _, v := range id {
			exprBuilder := new(strings.Builder)
			err := v.Build(exprBuilder)
			if err != nil {
				return "", err
			}
			vidList = append(vidList, exprBuilder.String())
		}
	default:
		return "", errors.New("vertex id must be a int, int64, string, clause.Expr, *clause.Expr or slice made of the above elements")
	}
	var vidExpr strings.Builder
	for i, v := range vidList {
		vidExpr.WriteString(v)
		if i != len(vidList)-1 {
			vidExpr.WriteString(", ")
		}
	}
	return vidExpr.String(), nil
}
