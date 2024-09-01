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

// NOTE: No embedded field is supported for struct, so do not use embedded field when declaring struct at the present time.

var db *nebulaorm.DB

func main() {
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

	Insert()

	Query()

	Update()

	Delete()
}

func Insert() {
	player := &Player{
		VID:  "player1001",
		Name: "Kobe Bryant",
		Age:  33,
	}
	if err := db.InsertVertex(player).Exec(); err != nil {
		log.Fatalf("insert player failed: %v", err)
	}
	team := &Team{
		VID:  "team1001",
		Name: "Lakers",
	}
	if err := db.InsertVertex(team).Exec(); err != nil {
		log.Fatalf("insert team failed: %v", err)
	}
	serve := &Serve{
		SrcID:     "player1001",
		DstID:     "team1001",
		StartYear: time.Date(1996, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
		EndYear:   time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
	}
	if err := db.InsertEdge(serve).Exec(); err != nil {
		log.Fatalf("insert serve failed: %v", err)
	}
}

func Query() {
	player := new(Player)
	err := db.Fetch("player", "player1001").
		Yield("vertex as v").
		FindCol("v", player)
	if err != nil {
		log.Fatalf("fetch player failed: %v", err)
	}
	log.Printf("player: %+v", player)

	team := new(Team)
	err = db.Go().
		From("player1001").
		Over("serve").
		Yield("$$ as t").
		FindCol("t", team)
	if err != nil {
		log.Fatalf("fetch team failed: %v", err)
	}
	log.Printf("team: %+v", team)

	serve := make([]*Serve, 0)
	err = db.Go().
		From("player1001").
		Over("serve").
		Yield("edge as e").
		FindCol("e", &serve)
	if err != nil {
		log.Fatalf("fetch serve failed: %v", err)
	}
	for _, s := range serve {
		log.Printf("serve: %+v", s)
	}

	type edgeCnt struct {
		Edge string `norm:"col:e"`
		Cnt  int    `norm:"col:cnt"`
	}
	edgesCnt := make([]*edgeCnt, 0)
	err = db.Go().
		From("player1001").
		Over("*").
		Yield("type(edge) as t").
		GroupBy("$-.t").
		Yield("$-.t as e, count(*) as cnt").
		Find(&edgesCnt)
	if err != nil {
		log.Fatalf("get edge cnt failed: %v", err)
	}
	for _, c := range edgesCnt {
		log.Printf("edge cnt: %+v\n", c)
	}

	type edgeVertex struct {
		ID string `norm:"col:id"`
		T  *Team  `norm:"col:t"`
	}
	edgeVertexes := make([]*edgeVertex, 0)
	err = db.Go().
		From("player1001").
		Over("serve").
		Yield("id($^) as id, $$ as t").
		Find(&edgeVertexes)
	if err != nil {
		log.Fatalf("get edge vertex failed: %v", err)
	}
	for _, v := range edgeVertexes {
		log.Printf("id: %v, t: %+v", v.ID, v.T)
	}
}

func Update() {
	if err := db.UpdateVertex("player1001", &Player{Age: 23}).Exec(); err != nil {
		log.Fatalf("update player failed: %v", err)
	}
	prop := make(map[string]interface{})
	err := db.Fetch("player", "player1001").
		Yield("properties(vertex) as p").
		FindCol("p", &prop)
	if err != nil {
		log.Fatalf("fetch player failed: %v", err)
	}
	log.Printf("vertex prop after update: %+v", prop)

	if err = db.UpdateEdge(Serve{SrcID: "player1001", DstID: "team1001"}, &Serve{StartYear: 160123456}).Exec(); err != nil {
		log.Fatalf("update edge serve failed: %v", err)
	}
	prop = make(map[string]interface{})
	err = db.Fetch("serve", clause.Expr{Str: `"player1001"->"team1001"`}).
		Yield("properties(edge) as p").
		FindCol("p", &prop)
	if err != nil {
		log.Fatalf("fetch serve failed: %v", err)
	}
	log.Printf("edge prop after update: %+v", prop)
}

func Delete() {
	if err := db.DeleteVertex("player1001").Exec(); err != nil {
		log.Fatalf("delete player failed: %v", err)
	}
	player := new(Player)
	err := db.Fetch("player", "player1001").
		Yield("vertex as v").
		TakeCol("v", player)
	if err != nil {
		log.Printf("after delete, fetch player failed: %v", err)
	}

	if err = db.DeleteEdge("serve", Serve{SrcID: "player1001", DstID: "team1001"}).Exec(); err != nil {
		log.Fatalf("delete edge serve failed: %v", err)
	}
	serve := new(Serve)
	err = db.Fetch("serve", clause.Expr{Str: `"player1001"->"team1001"`}).
		Yield("edge as e").
		TakeCol("e", serve)
	if err != nil {
		log.Printf("after delete, fetch server failed: %v", err)
	}
}
