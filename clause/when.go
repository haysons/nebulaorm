package clause

type When struct {
	Conditions []Condition
}

const WhenName = "WHEN"

func (when When) Name() string {
	return WhenName
}

func (when When) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(When)
	if !ok {
		clause.Expression = when
		return
	}
	exist.Conditions = append(exist.Conditions, when.Conditions...)
	clause.Expression = exist
}

func (when When) Build(nGQL Builder) error {
	nGQL.WriteString("WHEN ")
	return buildConditions(when.Conditions, nGQL)
}
