package clause

import (
	"fmt"
	"strings"
)

type Over struct {
	EdgeTypeList []string
	Direction    string
}

const OverName = "OVER"

const (
	OverDirectReversely = "REVERSELY"
	OverDirectBidirect  = "BIDIRECT"
)

func (over Over) Name() string {
	return OverName
}

func (over Over) MergeIn(clause *Clause) {
	exist, ok := clause.Expression.(Over)
	if !ok {
		clause.Expression = over
		return
	}
	// edge type merge, direction override
	exist.EdgeTypeList = append(exist.EdgeTypeList, over.EdgeTypeList...)
	if over.Direction != "" {
		exist.Direction = over.Direction
	}
	clause.Expression = exist
}

func (over Over) Build(nGQL Builder) error {
	edgeTypeList := make([]string, 0, len(over.EdgeTypeList))
	for _, edgeType := range over.EdgeTypeList {
		if edgeType != "" {
			edgeTypeList = append(edgeTypeList, edgeType)
		}
	}
	if len(edgeTypeList) == 0 {
		return fmt.Errorf("nebulaorm: %w, edge type list is empty in over clause", ErrInvalidClauseParams)
	}
	nGQL.WriteString("OVER ")
	nGQL.WriteString(strings.Join(edgeTypeList, ", "))
	if over.Direction != "" {
		if over.Direction != OverDirectReversely && over.Direction != OverDirectBidirect {
			return fmt.Errorf("nebulaorm: %w, over direction must be %s or %s", ErrInvalidClauseParams, OverDirectReversely, OverDirectBidirect)
		}
		nGQL.WriteByte(' ')
		nGQL.WriteString(over.Direction)
	}
	return nil
}
