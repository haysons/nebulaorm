package main

import (
	"errors"
	"github.com/haysons/nebulaorm"
	"github.com/haysons/nebulaorm/clause"
	"log"
)

func fetch() {
	// FETCH PROP ON player "player100" YIELD properties(vertex);
	// for the full attributes of a vertex, you can use map[string]interface{} to receive
	prop := make(map[string]interface{})
	err := db.
		Fetch("player", "player100").
		Yield("properties(vertex) as p").
		FindCol("p", &prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("prop: %+v", prop)

	// It is also possible to get all the information about a vertex directly, which requires the struct to implement the VertexID
	// and VertexTagNamer interfaces to indicate that it is a vertex. Query based on the vid, expect to get and only get one row,
	// then use the Take related methods, then return nebulaorm.ErrRecordNotFound if you can't find it.
	player := new(Player)
	err = db.
		Fetch("player", "player100").
		Yield("vertex as v").
		TakeCol("v", player)
	if err != nil && !errors.Is(err, nebulaorm.ErrRecordNotFound) {
		log.Fatal(err)
	}
	log.Printf("player: %+v", player)

	// FETCH PROP ON player "player100" \
	// YIELD properties(vertex).name AS name;
	// Get only the name attribute, which can be assigned directly to a string variable.
	var playerName string
	err = db.
		Fetch("player", "player100").
		Yield("properties(vertex).name AS name").
		TakeCol("name", &playerName)
	if err != nil && !errors.Is(err, nebulaorm.ErrRecordNotFound) {
		log.Fatal(err)
	}
	log.Printf("playerName: %+v", playerName)

	// FETCH PROP ON player "player101", "player102", "player103" YIELD properties(vertex);
	props := make([]map[string]interface{}, 0)
	err = db.
		Fetch("player", []string{"player101", "player102", "player103"}).
		Yield("properties(vertex) as p").
		FindCol("p", &props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("props: %+v", props)

	// When there are multiple tags for a vertex, you can query the values of the tags in the following way. The VertexPlayer
	// implements the VertexID interface, which is treated as a vertex, and also contains two exported fields, each of which
	// implements the VertexTagNamer interface, and is therefore treated as tags, which are assigned to a tag. So the two
	// tags will be assigned separate values.
	// CREATE TAG IF NOT EXISTS t1(a string, b int);
	// INSERT VERTEX t1(a, b) VALUES "player100":("Hello", 100);
	// FETCH PROP ON player, t1 "player100" YIELD vertex AS v;
	if err = db.Raw("CREATE TAG IF NOT EXISTS t1(a string, b int);").Exec(); err != nil {
		log.Fatal(err)
	}
	vPlayer := new(VertexPlayer)
	err = db.
		FetchMulti([]string{"player", "t1"}, "player100").
		Yield("vertex AS v").
		TakeCol("v", vPlayer)
	if err != nil && !errors.Is(err, nebulaorm.ErrRecordNotFound) {
		log.Fatal(err)
	}
	log.Printf("vPlayer: %+v", vPlayer)

	// FETCH PROP ON serve "player100"->"team204" YIELD properties(edge);
	// The current fetch processing logic is mainly designed for vertexes, so querying the properties of edges is a bit tricky.
	prop = make(map[string]interface{})
	err = db.
		Fetch("serve", clause.Expr{Str: `"player100" -> "team204"`}).
		Yield("properties(edge) as p").
		FindCol("p", &prop)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("prop: %+v", prop)

	// Get information about multiple edges
	// FETCH PROP ON serve "player100" → "team204", "player133" → "team202" YIELD edge AS e;
	edges := make([]Serve, 0)
	err = db.
		Fetch("serve", []*clause.Expr{{Str: `"player100" -> "team204"`}, {Str: `"player133" -> "team202"`}}).
		Yield("edge as e").
		FindCol("e", &edges)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("edges: %+v", edges)

	// Using FETCH in Compound Statements
	// GO FROM "player101" OVER follow \
	// YIELD src(edge) AS s, dst(edge) AS d \
	//  | FETCH PROP ON follow $-.s -> $-.d \
	// YIELD properties(edge).degree;
	degrees := make([]int, 0)
	err = db.
		Go().
		From("player101").
		Over("follow").
		Yield("src(edge) AS s, dst(edge) AS d").
		Pipe().
		Fetch("follow", clause.Expr{Str: `$-.s -> $-.d`}).
		Yield("properties(edge).degree as d").
		FindCol("d", &degrees)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("degrees: %+v", degrees)
}

type T1 struct {
	A string `norm:"prop:a"`
	B int    `norm:"prop:b"`
}

func (t T1) VertexTagName() string {
	return "t1"
}

// VertexPlayer is a vertex with two tags, Player and T1 both implement the VertexTagNamer interface, so both are treated as vertex tags,
// and the names of the fields for the two tags are arbitrary, as long as the fields are exported.
type VertexPlayer struct {
	TagPlayer Player
	TagT1     T1
}

func (p VertexPlayer) VertexID() string {
	return p.TagPlayer.VID
}
