package main

import (
	"github.com/haysons/nebulaorm/clause"
	"log"
)

type Player struct {
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

func (p Player) VertexTagName() string {
	return "player"
}

// // You can also express the desired tag names for updating by implementing the VertexTagName method on a map
type playerUpdate map[string]interface{}

func (m playerUpdate) VertexTagName() string {
	return "player"
}

func updateVertex() {
	log.Printf("before update, player101 attributes: %v", getProp("player101"))
	if err := db.UpdateVertex("player101", &Player{Age: 37}).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player101 attributes: %v", getProp("player101"))

	// You can also use map[string]interface{} for updates, but you will need to explicitly inform the framework
	// of the tag names that need to be updated
	if err := db.UpdateVertex("player101", map[string]interface{}{"age": 36}, clause.WithTagName("player")).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player101 attributes: %v", getProp("player101"))

	// Use a map[string]interface{} with tag names for updates
	if err := db.UpdateVertex("player101", playerUpdate{"age": 37}).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player101 attributes: %v", getProp("player101"))

	// Use an expression to update the value of age and retrieve the updated value
	// UPDATE VERTEX ON player "player101" \
	// SET age = age - 1 \
	// WHEN name == "Tony Parker" \
	// YIELD name AS Name, age AS Age;
	prop := make(map[string]interface{})
	err := db.UpdateVertex("player101", playerUpdate{"age": clause.Expr{Str: "age - 1"}}).
		When("name == ?", "Tony Parker").
		Yield("name AS Name, age AS Age").
		Find(prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("after update, player101 attributes: %v", prop)
}

func upsertVertex() {
	// If the vertex does not exist, write a new vertex
	// UPSERT VERTEX ON player "player666" \
	// SET age = 30 \
	// WHEN name == "Joe" \
	// YIELD name AS Name, age AS Age;
	prop := make(map[string]interface{})
	err := db.UpsertVertex("player666", &Player{Age: 30}).
		When("name == ?", "Joe").
		Yield("name AS Name, age AS Age").
		Find(prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("after the first upsert, player666 attributes: %v", prop)

	err = db.UpsertVertex("player666", &Player{Name: "player"}).
		When("name is ?", clause.Expr{Str: "NULL"}).
		Yield("name AS Name, age AS Age").
		Find(prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("after the second upsert, player666 attributes: %v", prop)
}

func getProp(vid interface{}) map[string]interface{} {
	// FETCH PROP ON player "player101" YIELD properties(vertex);
	prop := make(map[string]interface{})
	err := db.Fetch("player", vid).
		Yield("properties(vertex) as p").
		TakeCol("p", prop)
	if err != nil {
		log.Fatal(err)
	}
	return prop
}
