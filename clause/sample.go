package clause

import (
	"fmt"
	"strconv"
)

type Sample struct {
	SampleList []int
}

const SampleName = "SAMPLE"

func (sample Sample) Name() string {
	return SampleName
}

func (sample Sample) MergeIn(clause *Clause) {
	clause.Expression = sample
}

func (sample Sample) Build(nGQL Builder) error {
	if len(sample.SampleList) == 0 {
		return fmt.Errorf("nebulaorm: %w, sample list must have at least one item", ErrInvalidClauseParams)
	}
	nGQL.WriteString("SAMPLE [")
	for i, s := range sample.SampleList {
		nGQL.WriteString(strconv.Itoa(s))
		if i != len(sample.SampleList)-1 {
			nGQL.WriteByte(',')
		}
	}
	nGQL.WriteByte(']')
	return nil
}
