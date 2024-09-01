package clause_test

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"testing"
)

func TestGo(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.Go{StepStart: -1, StepEnd: -1}},
			gqlWant: `GO`,
		},
		{
			clauses: []clause.Interface{clause.Go{StepStart: 2, StepEnd: -1}},
			gqlWant: `GO 2 STEPS`,
		},
		{
			clauses: []clause.Interface{clause.Go{StepStart: 1, StepEnd: 2}},
			gqlWant: `GO 1 TO 2 STEPS`,
		},
		{
			clauses: []clause.Interface{clause.Go{StepEnd: 2}, clause.Go{StepEnd: -1}},
			gqlWant: `GO 0 STEPS`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
