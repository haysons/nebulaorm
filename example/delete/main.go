package main

import (
	"github.com/haysons/nebulaorm"
	"github.com/haysons/nebulaorm/clause"
	"log"
	"time"
)

type Player struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

func (p Player) VertexID() string {
	return p.VID
}

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
	SrcID     string `norm:"edge_src_id"`
	DstID     string `norm:"edge_dst_id"`
	Rank      int    `norm:"edge_rank"`
	StartYear int64  `norm:"prop:start_year"`
	EndYear   int64  `norm:"prop:end_year"`
}

func (s Serve) EdgeTypeName() string {
	return "serve"
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
	// Write a new vertex and edge
	player := &Player{VID: "player10003", Name: "n", Age: 23}
	if err = db.InsertVertex(player).Exec(); err != nil {
		log.Fatal(err)
	}
	team := &Team{VID: "team10003", Name: "t"}
	if err = db.InsertVertex(team).Exec(); err != nil {
		log.Fatal(err)
	}
	serve := &Serve{SrcID: "player10003", DstID: "team10003", Rank: 1, StartYear: 2001, EndYear: 2019}
	if err = db.InsertEdge(serve).Exec(); err != nil {
		log.Fatal(err)
	}
	// Query the team connected to player10003
	team = new(Team)
	err = db.Go().
		From("player10003").
		Over("serve").
		Yield("$$ as v").
		FindCol("v", team)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", team)

	// Delete the team associated with player10003
	// GO FROM "player10003" OVER serve \
	// YIELD id($$) as vid \
	// | DELETE VERTEX $-.vid;
	err = db.Go().
		From("player10003").
		Over("serve").
		Yield("id($$) as vid").
		Pipe().
		DeleteVertex(clause.Expr{Str: "$-.vid"}).Exec()
	if err != nil {
		log.Fatal(err)
	}
	// Delete the vertex player10003 and its associated edges
	err = db.DeleteVertex("player10003", true).Exec()
	if err != nil {
		log.Fatal(err)
	}
	// Query the remaining nodes and edges after deleting player10003
	player = new(Player)
	err = db.Fetch("player", "player10003").
		Yield("vertex as v").
		TakeCol("v", player)
	if err != nil {
		log.Printf("after deleting player10003, err:%v", err)
	} else {
		log.Printf("%+v", player)
	}
	team = new(Team)
	err = db.Fetch("team", "team10003").
		Yield("vertex as v").
		TakeCol("v", team)
	if err != nil {
		log.Printf("after deleting team10003, err:%v", err)
	} else {
		log.Printf("%+v", team)
	}
	serve = new(Serve)
	err = db.Fetch("serve", clause.Expr{Str: `"player10003"->"team10003"@1`}).
		Yield("edge as e").
		TakeCol("e", serve)
	if err != nil {
		log.Printf("after deleting serve, err:%v", err)
	} else {
		log.Printf("%+v", serve)
	}
}
