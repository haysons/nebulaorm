package clause

import "fmt"

type Lookup struct {
	TypeName string
}

const LookupName = "LOOKUP"

func (lookup Lookup) Name() string {
	return LookupName
}

func (lookup Lookup) MergeIn(clause *Clause) {
	clause.Expression = lookup
}

func (lookup Lookup) Build(nGQL Builder) error {
	if lookup.TypeName == "" {
		return fmt.Errorf("nebulaorm: %w, the vertex tag or edge type in lookup clause is empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString("LOOKUP ON ")
	nGQL.WriteString(lookup.TypeName)
	return nil
}
