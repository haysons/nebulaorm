package nebulaorm

import "github.com/haysons/nebulaorm/clause"

// Raw exec nGQL statements natively
// see more information on the method of the same name in statement.Statement
func (db *DB) Raw(raw string) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Raw(raw)
	return tx
}

// Go generate go clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Go(step ...int) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Go(step...)
	return
}

// From generate from clause
// see more information on the method of the same name in statement.Statement
func (db *DB) From(vid interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.From(vid)
	return
}

// Over generate over clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Over(edgeType ...string) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Over(edgeType...)
	return
}

// Where generate where clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Where(query string, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Where(query, args...)
	return
}

// Or generate or clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Or(query string, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Or(query, args...)
	return
}

// Not generate not clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Not(query string, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Not(query, args...)
	return
}

// Xor generate xor clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Xor(query string, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Xor(query, args...)
	return
}

// Sample generate sample clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Sample(sampleList ...int) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Sample(sampleList...)
	return
}

// Fetch generate fetch clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Fetch(name string, vid interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Fetch(name, vid)
	return
}

// FetchMulti generate fetch clause，multiple tag or edge types
// see more information on the method of the same name in statement.Statement
func (db *DB) FetchMulti(names []string, vid interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.FetchMulti(names, vid)
	return
}

// Lookup generate lookup clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Lookup(name string) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Lookup(name)
	return
}

// GroupBy generate group by clause
// see more information on the method of the same name in statement.Statement
func (db *DB) GroupBy(expr string) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.GroupBy(expr)
	return
}

// Yield generate yield clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Yield(expr string, distinct ...bool) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Yield(expr, distinct...)
	return
}

// OrderBy generate order by clause
// see more information on the method of the same name in statement.Statement
func (db *DB) OrderBy(expr string) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.OrderBy(expr)
	return
}

// Limit generate limit clause
// see more information on the method of the same name in statement.Statement
func (db *DB) Limit(limit int, offset ...int) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Limit(limit, offset...)
	return
}

// InsertVertex generate insert vertex clause
// see more information on the method of the same name in statement.Statement
func (db *DB) InsertVertex(vertexes interface{}, ifNotExist ...bool) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.InsertVertex(vertexes, ifNotExist...)
	return
}

// UpdateVertex generate update vertex clause
// see more information on the method of the same name in statement.Statement
func (db *DB) UpdateVertex(vid interface{}, propsUpdate interface{}, opts ...clause.Option) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.UpdateVertex(vid, propsUpdate, opts...)
	return
}

// UpsertVertex generate upsert vertex clause
// see more information on the method of the same name in statement.Statement
func (db *DB) UpsertVertex(vid interface{}, propsUpdate interface{}, opts ...clause.Option) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.UpsertVertex(vid, propsUpdate, opts...)
	return
}

// DeleteVertex generate delete vertex clause
// see more information on the method of the same name in statement.Statement
func (db *DB) DeleteVertex(vid interface{}, withEdge ...bool) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.DeleteVertex(vid, withEdge...)
	return
}

// InsertEdge generate insert edge clause
// see more information on the method of the same name in statement.Statement
func (db *DB) InsertEdge(edges interface{}, ifNotExist ...bool) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.InsertEdge(edges, ifNotExist...)
	return
}

// UpdateEdge generate update edge clause
// see more information on the method of the same name in statement.Statement
func (db *DB) UpdateEdge(edge interface{}, propsUpdate interface{}, opts ...clause.Option) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.UpdateEdge(edge, propsUpdate, opts...)
	return
}

// UpsertEdge generate upsert edge clause
// see more information on the method of the same name in statement.Statement
func (db *DB) UpsertEdge(edge interface{}, propsUpdate interface{}, opts ...clause.Option) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.UpsertEdge(edge, propsUpdate, opts...)
	return
}

// DeleteEdge generate delete edge clause
// see more information on the method of the same name in statement.Statement
func (db *DB) DeleteEdge(edgeTypeName string, edge interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.DeleteEdge(edgeTypeName, edge)
	return
}

// When generate when edge clause
// see more information on the method of the same name in statement.Statement
func (db *DB) When(query string, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.When(query, args...)
	return
}

// Pipe add a pipe character in current nGQL
func (db *DB) Pipe() (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Pipe()
	return
}
