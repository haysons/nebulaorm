package main

import "log"

// T1 does not contain any attributes
type T1 struct {
	VID string `norm:"vertex_id"`
}

func (t *T1) VertexID() string {
	return t.VID
}

func (t *T1) VertexTagName() string {
	return "t1"
}

// T2 contains two attributes
type T2 struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

func (t *T2) VertexID() string {
	return t.VID
}

func (t *T2) VertexTagName() string {
	return "t2"
}

// V1 is a vertex with two tags, so V1 itself needs to implement the VertexID method
// to indicate that it is a vertex and specify its vid.
// Additionally, V1 has two exported attributes, Tag1 and Tag2, which both implement
// the VertexTagName method, indicating that vertex V1 has two tags.
type V1 struct {
	VID  string `norm:"vertex_id"`
	Tag1 T3
	Tag2 T4
}

func (v *V1) VertexID() string {
	return v.VID
}

type T3 struct {
	P1 int `norm:"prop:p1"`
}

func (t *T3) VertexTagName() string {
	return "t3"
}

type T4 struct {
	P2 string `norm:"prop:p2"`
}

func (t *T4) VertexTagName() string {
	return "t4"
}

func insertVertex() {
	// Insert a vertex without attributes
	// CREATE TAG IF NOT EXISTS t_empty();
	// INSERT VERTEX t1() VALUES "10":();
	if err := db.Raw("CREATE TAG IF NOT EXISTS t_empty();").Exec(); err != nil {
		log.Fatal(err)
	}
	// When inserting only one vertex, simply insert the corresponding struct
	t1 := &T1{VID: "10"}
	if err := db.InsertVertex(t1).Exec(); err != nil {
		log.Fatal(err)
	}

	// Insert two vertexes with attributes in a single operation
	// CREATE TAG IF NOT EXISTS t2 (name string, age int);
	// INSERT VERTEX t2 (name, age) VALUES "13":("n3", 12), "14":("n4", 8);
	if err := db.Raw("CREATE TAG IF NOT EXISTS t2 (name string, age int);").Exec(); err != nil {
		log.Fatal(err)
	}
	t2 := []*T2{{VID: "13", Name: "n3", Age: 12}, {VID: "14", Name: "n4", Age: 8}}
	if err := db.InsertVertex(t2).Exec(); err != nil {
		log.Fatal(err)
	}

	// Insert two Tag attributes into the same vertex in a single operation
	// CREATE TAG IF NOT EXISTS t3(p1 int);
	// CREATE TAG IF NOT EXISTS t4(p2 string);
	// INSERT VERTEX  t3 (p1), t4(p2) VALUES "21": (321, "hello");
	if err := db.Raw("CREATE TAG IF NOT EXISTS t3(p1 int);CREATE TAG IF NOT EXISTS t4(p2 string);").Exec(); err != nil {
		log.Fatal(err)
	}
	v1 := &V1{VID: "21", Tag1: T3{P1: 321}, Tag2: T4{P2: "hello"}}
	if err := db.InsertVertex(v1).Exec(); err != nil {
		log.Fatal(err)
	}

	// When inserting an existing vertex with IF NOT EXISTS, no modification will be made
	// INSERT VERTEX IF NOT EXISTS t2 (name, age) VALUES "1":("n3", 14);
	t22 := &T2{VID: "1", Name: "n3", Age: 14}
	if err := db.InsertVertex(t22, true).Exec(); err != nil {
		log.Fatal(err)
	}
}
