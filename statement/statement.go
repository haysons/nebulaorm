package statement

import (
	"github.com/haysons/nebulaorm/clause"
	"strings"
)

// Statement is an nGQL statement that needs to be constructed.
// A statement consists of multiple parts, which may be separated by '|', and each part consists of multiple clauses
// that independently construct their own part of the statement. The statement object is not concurrency safe.
type Statement struct {
	parts []*Part
	nGQL  *strings.Builder
	built bool
	err   error
}

func New() *Statement {
	return &Statement{
		parts: make([]*Part, 0),
		nGQL:  new(strings.Builder),
	}
}

// LastPart gets the last part of the current statement.
func (stmt *Statement) LastPart() *Part {
	if len(stmt.parts) == 0 {
		stmt.AddPart(NewPart())
	}
	return stmt.parts[len(stmt.parts)-1]
}

// AddPart add a new part at the end of the current statement
func (stmt *Statement) AddPart(part *Part) {
	stmt.parts = append(stmt.parts, part)
}

// Pipe denotes the pipe character in nGQL, which opens a new part and separates it from the previous one using
// the pipe character.
func (stmt *Statement) Pipe() *Statement {
	part := NewPart()
	part.SetCompType(CompositeTypePipe)
	stmt.AddPart(part)
	return stmt
}

// AddClause adds a clause to the last part of the statement.
func (stmt *Statement) AddClause(v clause.Interface) {
	part := stmt.LastPart()
	part.AddClause(v)
}

// SetPartType sets the type of the last part of the current statement. The part type determines which types to use
// and in which order to build the statement.
func (stmt *Statement) SetPartType(typ PartType) {
	part := stmt.LastPart()
	part.SetType(typ)
}

// SetClausesBuild manually specifies the type and order of statements that need to be built for the final part.
func (stmt *Statement) SetClausesBuild(clauses []string) {
	part := stmt.LastPart()
	part.SetClausesBuild(clauses)
}

// Build the current statement; any problems during build will return erring
func (stmt *Statement) Build() error {
	if stmt.err != nil || stmt.built {
		return stmt.err
	}
	stmt.nGQL.Reset()
	stmt.nGQL.Grow(100 * len(stmt.parts))
	var firstPartBuilt bool
	// generate statements for each part in turn
	for _, part := range stmt.parts {
		if len(part.clauses) == 0 {
			continue
		}
		// add a connector between multiple statements based on the type of the compound statement
		if firstPartBuilt {
			switch part.compType {
			case CompositeTypePipe:
				stmt.nGQL.WriteString(" | ")
			}
		}
		firstPartBuilt = true
		if err := part.Build(stmt.nGQL); err != nil {
			stmt.err = err
			break
		}
	}
	stmt.nGQL.WriteByte(';')
	stmt.built = true
	return stmt.err
}

// NGQL build and return the nGQL statement, returning erring if there is a problem with the build
func (stmt *Statement) NGQL() (string, error) {
	if err := stmt.Build(); err != nil {
		return "", err
	}
	return stmt.nGQL.String(), nil
}

// Part is the part of the statement that actually contains the clause to be constructed and completes the construction
// of the statement by calling the clause's Build method. Because the concept of a compound statement exists in nGQL,
// it is necessary to add another layer to the statement concept to generate each part of the compound statement
// independently of each other.
type Part struct {
	typ          PartType
	setType      bool
	compType     CompositeType
	clauses      map[string]clause.Clause
	clausesBuild []string
}

func NewPart() *Part {
	return &Part{
		clauses:      make(map[string]clause.Clause),
		clausesBuild: make([]string, 0),
	}
}

// SetType sets the type of the current part, which is used to specify the list of clauses to be built;
// only the first call takes effect when called multiple times.
func (p *Part) SetType(typ PartType) {
	if p.setType {
		return
	}
	p.typ = typ
	p.setType = true
}

// GetType get the current part type.
func (p *Part) GetType() PartType {
	return p.typ
}

// SetCompType sets the composite type of the current part. The composite type mainly determines how multiple
// parts are separated in a composite statement, e.g., by the use of a pipe character.
func (p *Part) SetCompType(typ CompositeType) {
	p.compType = typ
}

// GetCompType get the composite type of the current part
func (p *Part) GetCompType() CompositeType {
	return p.compType
}

// SetClausesBuild manually specifies the list of clauses that the current part needs to build.
func (p *Part) SetClausesBuild(clauses []string) {
	p.clausesBuild = clauses
}

func (p *Part) AddClause(v clause.Interface) {
	name := v.Name()
	c := p.clauses[name]
	c.Name = name
	v.MergeIn(&c)
	p.clauses[name] = c
}

func (p *Part) Build(nGQL clause.Builder) error {
	var firstClauseWritten bool
	for _, name := range p.getClausesBuild() {
		if c, ok := p.clauses[name]; ok {
			if firstClauseWritten {
				nGQL.WriteByte(' ')
			}
			firstClauseWritten = true
			if err := c.Build(nGQL); err != nil {
				return err
			}
		}
	}
	return nil
}

type CompositeType int

const (
	CompositeTypePipe CompositeType = iota + 1
)

type PartType int

const (
	PartTypeGo PartType = iota + 1
	PartTypeFetch
	PartTypeLookup
	PartTypeGroup
	PartTypeOrder
	PartTypeLimit
	PartTypeInsertVertex
	PartTypeUpdateVertex
	PartTypeDeleteVertex
	PartTypeInsertEdge
	PartTypeUpdateEdge
	PartTypeDeleteEdge
)

func (p *Part) getClausesBuild() []string {
	if len(p.clausesBuild) > 0 {
		return p.clausesBuild
	}
	switch p.typ {
	case PartTypeGo:
		return []string{clause.GoName, clause.FromName, clause.OverName, clause.WhereName, clause.YieldName, clause.SampleName}
	case PartTypeFetch:
		return []string{clause.FetchName, clause.YieldName}
	case PartTypeLookup:
		return []string{clause.LookupName, clause.WhereName, clause.YieldName}
	case PartTypeGroup:
		return []string{clause.GroupName, clause.YieldName}
	case PartTypeOrder:
		return []string{clause.OrderName}
	case PartTypeLimit:
		return []string{clause.LimitName}
	case PartTypeInsertVertex:
		return []string{clause.InsertVertexName}
	case PartTypeUpdateVertex:
		return []string{clause.UpdateVertexName, clause.WhenName, clause.YieldName}
	case PartTypeDeleteVertex:
		return []string{clause.DeleteVertexName}
	case PartTypeInsertEdge:
		return []string{clause.InsertEdgeName}
	case PartTypeUpdateEdge:
		return []string{clause.UpdateEdgeName, clause.WhenName, clause.YieldName}
	case PartTypeDeleteEdge:
		return []string{clause.DeleteEdgeName}
	default:
		// The following clauses may not belong to a specific type of statement and can be used separately
		return []string{clause.GroupName, clause.YieldName, clause.OrderName, clause.LimitName}
	}
}
