package clause

import (
	"fmt"
	"strconv"
)

type Limit struct {
	Limit  int
	Offset int
}

const LimitName = "LIMIT"

func (limit Limit) Name() string {
	return LimitName
}

func (limit Limit) MergeIn(clause *Clause) {
	clause.Expression = limit
}

func (limit Limit) Build(nGQL Builder) error {
	if limit.Limit < 0 {
		return fmt.Errorf("nebulaorm: %w, limit can't be negative", ErrInvalidClauseParams)
	}
	nGQL.WriteString("LIMIT ")
	if limit.Offset > 0 {
		nGQL.WriteString(strconv.Itoa(limit.Offset))
		nGQL.WriteString(", ")
	}
	nGQL.WriteString(strconv.Itoa(limit.Limit))
	return nil
}
