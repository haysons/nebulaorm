package clause

import (
	"strconv"
)

// Go clause, a step of 0 is a legal value; if you want no step, you can set step to a negative number.
type Go struct {
	StepStart int
	StepEnd   int
}

const GoName = "GO"

func (g Go) Name() string {
	return GoName
}

func (g Go) MergeIn(clause *Clause) {
	clause.Expression = g
}

func (g Go) Build(nGQL Builder) error {
	nGQL.WriteString("GO")
	if g.StepStart >= 0 && g.StepEnd >= 0 {
		nGQL.WriteByte(' ')
		nGQL.WriteString(strconv.Itoa(g.StepStart))
		nGQL.WriteString(" TO ")
		nGQL.WriteString(strconv.Itoa(g.StepEnd))
		nGQL.WriteString(" STEPS")
		return nil
	}
	if g.StepStart >= 0 {
		nGQL.WriteByte(' ')
		nGQL.WriteString(strconv.Itoa(g.StepStart))
		nGQL.WriteString(" STEPS")
	}
	return nil
}
