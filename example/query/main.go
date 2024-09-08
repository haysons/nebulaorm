package main

import (
	"github.com/haysons/nebulaorm"
	"log"
	"time"
)

// After installing nebula graph locally, a default demo_basketballplayer graph space will be created.
// The examples below are based on the default data that exists in this graph space.

type Player struct {
	// Through the use of the "vertex_id" tag, it is indicated that this field represents the VID of a vertex,
	// so when assigning values, the VID returned by nebula graph will be assigned to this field.
	VID string `norm:"vertex_id"`
	// Other exported fields in the struct are treated as an attribute of the tag.
	// The attribute name can be defined using the `prop` field, with the default
	// name being the snake_case version of the field's camelCase name.
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

// VertexID When the method is implemented in a struct, the struct is considered a vertex.
// This method does not necessarily need to return a predefined VID field.
// It can also generate a VID based on properties, such as using the player's name as the unique identifier for the vertex.
// For example, instead of defining a VID field, it can directly return md5(v.Name).
func (p Player) VertexID() string {
	return p.VID
}

// VertexTagName When the method is implemented in a struct, the struct is considered a tag of the vertex.
// If the vertex has only one tag, a struct can implement both the VertexID and VertexTagName methods.
// VertexTagName should return the name of the tag on the vertex.
func (p Player) VertexTagName() string {
	return "player"
}

type Team struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
}

func (t Team) VertexID() string {
	return t.VID
}

func (t Team) VertexTagName() string {
	return "team"
}

type Serve struct {
	// Using the edge_src_id tag indicates that this field represents the src_id of an edge.
	// When assigning values, the src_id returned by nebula graph will be assigned to this field.
	SrcID string `norm:"edge_src_id"`
	// Using the edge_dst_id tag indicates that this field represents the dst_id of an edge.
	// When assigning values, the dst_id returned by nebula graph will be assigned to this field.
	DstID string `norm:"edge_dst_id"`
	// Using the edge_rank tag indicates that this field represents the rank of an edge.
	// When assigning values, the rank returned by nebula graph will be assigned to this field.
	// If the rank of the edge is the default value, this field does not need to be defined.
	Rank int `norm:"edge_rank"`
	// Other exported fields in the struct are treated as attributes of the edge.
	// The attribute name can be specified using the `prop`, with the default name being
	// the snake_case version of the field's camelCase name.
	StartYear int64 `norm:"prop:start_year"`
	EndYear   int64 `norm:"prop:end_year"`
}

// EdgeTypeName When this method is implemented in a struct, it indicates that the struct represents an edge.
// The method should return the name of the edge.
func (s Serve) EdgeTypeName() string {
	return "serve"
}

type Follow struct {
	SrcID  string `norm:"edge_src_id"`
	DstID  string `norm:"edge_dst_id"`
	Rank   int    `norm:"edge_rank"`
	Degree int64  `norm:"prop:degree"`
}

func (f Follow) EdgeTypeName() string {
	return "follow"
}

// Nebulaorm.DB is concurrency-safe and internally uses the connection pool provided by the nebula graph official SDK.
// Therefore, in general, only a single instance needs to be defined.
var db *nebulaorm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := &nebulaorm.Config{
		Username:    "root",
		Password:    "nebula",
		SpaceName:   "demo_basketballplayer",
		Addresses:   []string{"127.0.0.1:9669"},
		ConnTimeout: 10 * time.Second,
	}
	var err error
	db, err = nebulaorm.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	lookup()
	queryGo()
	fetch()
}
