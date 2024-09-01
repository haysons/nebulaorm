package clause

import (
	"fmt"
)

type Fetch struct {
	Names []string
	VID   interface{}
}

const FetchName = "FETCH"

func (fetch Fetch) Name() string {
	return FetchName
}

func (fetch Fetch) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(Fetch)
	if !ok {
		clause.Expression = fetch
		return
	}
	existNameSet := make(map[string]struct{}, len(exist.Names))
	for _, name := range exist.Names {
		existNameSet[name] = struct{}{}
	}
	for _, name := range fetch.Names {
		_, ok := existNameSet[name]
		if !ok {
			exist.Names = append(exist.Names, name)
			existNameSet[name] = struct{}{}
		}
	}
	exist.VID = fetch.VID
	clause.Expression = exist
}

func (fetch Fetch) Build(nGQL Builder) error {
	if len(fetch.Names) == 0 {
		return fmt.Errorf("nebulaorm: %w, the names in fetch clause is empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString("FETCH PROP ON ")
	for i, name := range fetch.Names {
		nGQL.WriteString(name)
		if i != len(fetch.Names)-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteByte(' ')
	vidExpr, err := vertexIDExpr(fetch.VID)
	if err != nil {
		return fmt.Errorf("nebulaorm: %w, build fetch clause failed, %v", ErrInvalidClauseParams, err)
	}
	nGQL.WriteString(vidExpr)
	return nil
}
