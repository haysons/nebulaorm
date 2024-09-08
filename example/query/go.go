package main

import (
	"github.com/haysons/nebulaorm/clause"
	"log"
)

func queryGo() {
	// GO FROM "player102" OVER serve YIELD dst(edge);
	teamIDs := make([]string, 0)
	err := db.Go().
		From("player102").
		Over("serve").
		Yield("dst(edge) as t").
		FindCol("t", &teamIDs)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("teamIDs: %+v\n", teamIDs)

	// If you want to retrieve not just the IDs of the team nodes but the list of the vertexes themselves, you can use the following approach.
	// The variable that you want to assign values to need to implement both the VertexID and VertexTagName methods.
	teams := make([]Team, 0)
	err = db.Go().
		From("player102").
		Over("serve").
		Yield("$$ as t").
		FindCol("t", &teams)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("teams: %+v\n", teams)

	// GO FROM "player100", "player102" OVER serve \
	// WHERE properties(edge).start_year > 1995 \
	// YIELD DISTINCT properties($$).name AS team_name, properties(edge).start_year AS start_year, properties($^).name AS player_name;
	// In this case, since we are interested in multiple fields of the return value, use Find to assign values instead of using FindCol.
	type record struct {
		TeamName   string
		StartYear  int64
		PlayerName string
	}
	records := make([]record, 0)
	err = db.Go().
		From([]string{"player100", "player102"}).
		Over("serve").
		Where("properties(edge).start_year > ?", 1995).
		Yield("properties($$).name AS team_name, properties(edge).start_year AS start_year, properties($^).name AS player_name", true).
		Find(&records)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("records: %+v\n", records)

	// GO FROM "player100" OVER follow, serve \
	// YIELD properties(edge).degree, properties(edge).start_year;
	// When querying multiple edge attributes, properties that do not exist on the edge will return null.
	// In Go, these will be assigned the zero value for their respective types.
	type record2 struct {
		Degree    int
		StartYear int64
	}
	records2 := make([]record2, 0)
	err = db.Go().
		From("player100").
		Over("serve", "follow").
		Yield("properties(edge).degree as degree, properties(edge).start_year as start_year").
		Find(&records2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("records: %+v\n", records2)
	// If you need to explicitly know which field values are null, you can make the field types pointers.
	// This will map null values to nil.
	type record3 struct {
		Degree    *int
		StartYear int64
	}
	records3 := make([]record3, 0)
	err = db.Go().
		From("player100").
		Over("serve", "follow").
		Yield("properties(edge).degree as degree, properties(edge).start_year as start_year").
		Find(&records3)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("records: %+v\n", records3)

	// A subquery can be used as the starting point for a graph traversal.
	// GO FROM "player100" OVER follow \
	// YIELD src(edge) AS id | \
	// GO FROM $-.id OVER serve \
	// WHERE properties($^).age > 20 \
	// YIELD properties($^).name AS Friend, properties($$).name AS Team;
	type record4 struct {
		Friend string `norm:"col:Friend"`
		Team   string `norm:"col:Team"`
	}
	record4s := make([]record4, 0)
	err = db.Go().
		From("player100").
		Over("follow").
		Yield("src(edge) AS id").Pipe().
		Go().
		From(clause.Expr{Str: "$-.id"}).
		Over("serve").
		Where("properties($^).age > ?", 20).
		Yield("properties($^).name AS Friend, properties($$).name AS Team").
		Find(&record4s)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("records: %+v\n", record4s)

	// GO 2 STEPS FROM "player100" OVER follow \
	// YIELD src(edge) AS src, dst(edge) AS dst, properties($$).age AS age \
	//  | GROUP BY $-.dst \
	// YIELD $-.dst AS dst, collect_set($-.src) AS src, collect($-.age) AS age;
	// The return values here involve list and set, which can be received using slice.
	type record5 struct {
		Dst string
		Src []string
		Age []int
	}
	record5s := make([]*record5, 0)
	err = db.Go(2).
		From("player100").
		Over("follow").
		Yield("src(edge) AS src, dst(edge) AS dst, properties($$).age AS age").
		GroupBy("$-.dst").
		Yield("$-.dst AS dst, collect_set($-.src) AS src, collect($-.age) AS age").
		Find(&record5s)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range record5s {
		log.Printf("record5: %+v\n", r)
	}
}
