package main

import (
	"github.com/haysons/nebulaorm/clause"
	"log"
)

type Serve struct {
	SrcID     string `norm:"edge_src_id"`
	DstID     string `norm:"edge_dst_id"`
	Rank      int    `norm:"edge_rank"`
	StartYear int64  `norm:"prop:start_year"`
	EndYear   int64  `norm:"prop:end_year"`
}

func (s Serve) EdgeTypeName() string {
	return "serve"
}

func updateEdge() {
	log.Printf("before update, player100 serve edge attributes: %v", getServeProp("player100"))
	// Directly use edge expressions to update edge attribute values
	if err := db.UpdateEdge(`serve "player100" -> "team204"@0`, Serve{EndYear: 2017}).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player100 serve edge attributes: %v", getServeProp("player100"))

	// Concatenating edge expressions can be cumbersome. You can pass an edge type variable,
	// and the framework will automatically handle the expression concatenation at a lower level.
	if err := db.UpdateEdge(Serve{SrcID: "player100", DstID: "team204"}, Serve{EndYear: 2016}).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player100 serve edge attributes: %v", getServeProp("player100"))

	// You can also use map[string]interface{} to specify the attributes to be updated
	if err := db.UpdateEdge(Serve{SrcID: "player100", DstID: "team204"}, map[string]interface{}{"end_year": 2017}).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player100 serve edge attributes: %v", getServeProp("player100"))

	// Update attribute values based on expressions
	// UPDATE EDGE ON serve "player100" -> "team204"@0 \
	// SET end_year = end_year - 1 \
	// WHEN end_year > 2010 \
	// YIELD start_year, end_year;
	prop := make(map[string]interface{})
	err := db.UpdateEdge(&Serve{SrcID: "player100", DstID: "team204"}, map[string]interface{}{"end_year": clause.Expr{Str: "end_year - 1"}}).
		When("end_year > ?", 2010).
		Yield("start_year, end_year").
		Find(prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player100 serve edge attributes: %v", prop)
}

func upsertEdge() {
	// If the edge does not exist, a new edge will be written
	prop := make(map[string]interface{})
	err := db.UpsertEdge(Serve{SrcID: "player666", DstID: "team200"}, &Serve{EndYear: 2021}).
		When("end_year == ?", 2010).
		Yield("start_year, end_year").
		Find(prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("after the first upsert, serve edge attributes: %v", prop)

	err = db.UpsertEdge(Serve{SrcID: "player666", DstID: "team200"}, &Serve{StartYear: 1997}).
		When("start_year is null").
		Yield("start_year, end_year").
		Find(prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("after the second upsert, serve edge attributes: %v", prop)
}

func getServeProp(vid interface{}) map[string]interface{} {
	// GO FROM "player100" \
	// OVER serve \
	// YIELD properties(edge).start_year, properties(edge).end_year;
	prop := make(map[string]interface{})
	err := db.Go().
		From(vid).
		Over("serve").
		Yield("properties(edge).start_year as start_year, properties(edge).end_year as end_year").
		Find(prop)
	if err != nil {
		log.Fatal(err)
	}
	return prop
}
