package clause_test

import (
	"errors"
	"github.com/haysons/nebulaorm/clause"
	"github.com/haysons/nebulaorm/statement"
	"strings"
	"testing"
)

func testBuildClauses(t *testing.T, clauses []clause.Interface, gqlWant string, errWant error) {
	buildNames := make([]string, len(clauses))
	buildNamesMap := make(map[string]bool, len(clauses))
	stmtPart := statement.NewPart()
	for _, c := range clauses {
		if _, ok := buildNamesMap[c.Name()]; !ok {
			buildNames = append(buildNames, c.Name())
			buildNamesMap[c.Name()] = true
		}
		stmtPart.AddClause(c)
	}
	stmtPart.SetClausesBuild(buildNames)
	gqlBuilder := new(strings.Builder)
	err := stmtPart.Build(gqlBuilder)
	gql := gqlBuilder.String()
	if !errors.Is(err, errWant) {
		t.Errorf("clause build err exception, want: %v  got: %v", errWant, err)
	}
	if err == nil && gql != gqlWant {
		t.Errorf("clause build exception, want: %s  got: %s", gqlWant, gql)
	}
}
