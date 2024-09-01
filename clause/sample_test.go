package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestSample(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Sample{SampleList: []int{1, 2, 4}}},
			gqlWant: "SAMPLE [1,2,4]",
		},
		{
			clauses: []clause.Interface{clause.Sample{SampleList: []int{1, 2, 4}}, clause.Sample{SampleList: []int{1, 2, 3}}},
			gqlWant: "SAMPLE [1,2,3]",
		},
		{
			clauses: []clause.Interface{clause.Sample{}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
