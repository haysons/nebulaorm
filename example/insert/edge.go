package main

import "log"

// E1 does not contain attributes
type E1 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
}

func (e *E1) EdgeTypeName() string {
	return "e1"
}

// E2 contains two attributes
type E2 struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
	Name  string `norm:"prop:name"`
	Age   int    `norm:"prop:age"`
}

func (e *E2) EdgeTypeName() string {
	return "e2"
}

func insertEdge() {
	// Insert an edge without attributes
	// CREATE EDGE IF NOT EXISTS e1();
	// INSERT EDGE e1 () VALUES "10"->"11":();
	if err := db.Raw("CREATE EDGE IF NOT EXISTS e1();").Exec(); err != nil {
		log.Fatal(err)
	}
	if err := db.InsertEdge(&E1{SrcID: "10", DstID: "11"}).Exec(); err != nil {
		log.Fatal(err)
	}
	// Insert an edge with rank 1
	if err := db.InsertEdge(&E1{SrcID: "10", DstID: "11", Rank: 1}).Exec(); err != nil {
		log.Fatal(err)
	}

	// Insert two edges with attributes in a single operation
	if err := db.Raw("CREATE EDGE IF NOT EXISTS e2 (name string, age int);").Exec(); err != nil {
		log.Fatal(err)
	}
	e2 := []*E2{{SrcID: "12", DstID: "13", Name: "n1", Age: 1}, {SrcID: "13", DstID: "14", Name: "n2", Age: 2}}
	if err := db.InsertEdge(&e2).Exec(); err != nil {
		log.Fatal(err)
	}

	// When inserting an existing edge with IF NOT EXISTS, no modification will be made
	if err := db.InsertEdge(&E2{SrcID: "14", DstID: "15", Rank: 1, Name: "n2", Age: 13}, true).Exec(); err != nil {
		log.Fatal(err)
	}
}
