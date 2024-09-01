package clause

import "errors"

var (
	// ErrInvalidClauseParams indicates that the argument to the clause is invalid
	ErrInvalidClauseParams = errors.New("invalid clause params")
)

// Interface clause interface
type Interface interface {
	// Name 子句的名称
	Name() string

	// MergeIn 合并相同子句至Clause对象之中
	MergeIn(clause *Clause)

	// Build 构造nGQL语句
	Build(nGQL Builder) error
}

// Clause 子句的通用结构，子句均包含名称以及表达式
type Clause struct {
	Name       string
	Expression Expression
}

// Build  clause
func (c Clause) Build(nGQL Builder) error {
	return c.Expression.Build(nGQL)
}

// Options 子句配置项
type Options struct {
	propNames []string // 指定属性列表
	tagName   string   // 指定tag名称
}

type Option func(*Options)

// WithPropNames 指定属性字段名
func WithPropNames(propNames []string) Option {
	return func(o *Options) {
		o.propNames = propNames
	}
}

func WithTagName(tagName string) Option {
	return func(o *Options) {
		o.tagName = tagName
	}
}
