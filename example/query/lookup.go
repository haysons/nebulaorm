package main

import (
	"github.com/haysons/nebulaorm/clause"
	"log"
)

func lookup() {
	// LOOKUP ON player \
	// WHERE player.name == "Steve Nash" \
	// YIELD id(vertex);
	// Query the vertex's ID based on the where condition. In this case, we only care about the "vid" column value.
	// However, nebula graph always returns the result in key-value form.
	// We can use FindCol to extract only the column value without defining an additional struct for the outer layer.
	vidList := make([]string, 0)
	err := db.Lookup("player").
		Where("player.name == ?", "Steve Nash").
		Yield("id(vertex) as vid").
		FindCol("vid", &vidList)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Steve Nash vidList: %v", vidList)

	// If you're certain that this query will return only one result, you can directly define a "vid" variable of type string.
	// This approach makes it more convenient to use.
	var vid string
	err = db.Lookup("player").
		Where("player.name == ?", "Steve Nash").
		Yield("id(vertex) as vid").
		FindCol("vid", &vid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Steve Nash vid: %v", vid)

	// LOOKUP ON player \
	// WHERE player.name == "Tony Parker" \
	// YIELD properties(vertex).name AS name, properties(vertex).age AS age;
	// To retrieve specific attribute values, you can use the following approach.
	// In this case, since we're interested in multiple fields of the result,
	// you can use the Find method to assign values to a struct.
	type record struct {
		Name string
		Age  int
	}
	var r record
	err = db.Lookup("player").
		Where("player.name == ?", "Steve Nash").
		Yield("properties(vertex).name AS name, properties(vertex).age AS age").
		Find(&r)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Steve Nash name and age: %+v", r)

	// If you want to retrieve all attribute values at once, you can use map[string]interface{} to receive the result.
	// In this case, you're only focusing on one field, but since there might be multiple rows,
	// you should use []map[string]interface{} to handle the response.
	curProp := make([]map[string]interface{}, 0)
	err = db.Lookup("player").
		Where("player.name == ?", "Steve Nash").
		Yield("properties(vertex) as properties").
		FindCol("properties", &curProp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Steve Nash props: %+v", curProp)

	// Return both attribute values and other data
	type record2 struct {
		PlayerName string
		TeamName   string
		StartYear  int64
		EndYear    int64
	}
	records := make([]*record2, 0)
	// LOOKUP ON player \
	// WHERE player.name == "Steve Nash"\
	// YIELD id(vertex) AS VertexID, properties(vertex).name AS name | \
	// GO FROM $-.VertexID OVER serve \
	// YIELD $-.name, properties(edge).start_year, properties(edge).end_year, properties($$).name;
	// The Pipe() method is used here to add a pipe symbol to the statement.
	err = db.Lookup("player").
		Where("player.name == ?", "Steve Nash").
		Yield("id(vertex) AS VertexID, properties(vertex).name AS name").Pipe().
		Go().
		From(clause.Expr{Str: "$-.VertexID"}).
		Over("serve").
		Yield("$-.name as player_name, properties(edge).start_year as start_year, properties(edge).end_year as end_year, properties($$).name as team_name").
		Find(&records)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range records {
		log.Printf("player and team: %+v", r)
	}

	// LOOKUP ON follow WHERE follow.degree == 90 YIELD edge AS e;
	// Query and return all edges where the degree equals 90. Since edges are returned directly,
	// the variable used for assignment must implement the EdgeTypeName interface to indicate that it represents an edge.
	edges := make([]Follow, 0)
	err = db.Lookup("follow").
		Where("follow.degree == ?", 90).
		Yield("edge AS e").
		FindCol("e", &edges)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("edges: %+v\n", edges)

	// LOOKUP ON follow YIELD properties(edge).degree as degree | ORDER BY $-.degree | LIMIT 10;
	// Sort by the "degree" attribute in ascending order and return the top 10 entries with the lowest "degree" values.
	degrees := make([]int, 0)
	err = db.Lookup("follow").
		Yield("properties(edge).degree as degree").
		OrderBy("$-.degree").
		Limit(10).
		FindCol("degree", &degrees)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("degrees: %+v\n", degrees)

	// LOOKUP ON player YIELD id(vertex) | YIELD COUNT(*) AS Player_Number;
	// Count the total number of nodes with the tag "player".
	var playerCnt int64
	err = db.Lookup("player").
		Yield("id(vertex)").
		Pipe().
		Yield("COUNT(*) AS Player_Number").
		FindCol("Player_Number", &playerCnt)
	log.Printf("playerCnt: %v\n", playerCnt)
}
