package clause

import "fmt"

type Group struct {
	Expr string
}

const GroupName = "GROUP"

func (group Group) Name() string {
	return GroupName
}

func (group Group) MergeIn(clause *Clause) {
	clause.Expression = group
}

func (group Group) Build(nGQL Builder) error {
	if group.Expr == "" {
		return fmt.Errorf("nebulaorm: %w, group by expr is empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString("GROUP BY ")
	nGQL.WriteString(group.Expr)
	return nil
}
