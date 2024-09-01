# nebulaorm

## Introduction

nebulaorm is an orm framework specially designed for nebula graph, which can be used to splice nGQL statements by
chained calls in a more elegant and faster way. It parses the result set returned by nebula graph and assigns it to the
variables provided by the developer, hoping to ultimately improve the experience of using nebula graph.

## Installation

```
go get github.com/haysons/nebulaorm
```

## Quick Start

``` go
// Player vertex player
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

// Team vertex team
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

// Serve edge Serve
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

// NOTE: Currently struct does not support embedded fields, so do not use embedded fields when declaring struct.
func main() {
    // initialize the db object
    conf := &nebulaorm.Config{
        Username:    "root",
        Password:    "nebula",
        SpaceName:   "demo_basketballplayer",
        Addresses:   []string{"127.0.0.1:9669"},
    }
    db, err := nebulaorm.Open(conf)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // insert the player vertex
    player := &Player{
        VID:  "player1001",
        Name: "Kobe Bryant",
        Age:  33,
    }
    if err := db.InsertVertex(player).Exec(); err != nil {
        log.Fatalf("insert player failed: %v", err)
    }
    // insert the team vertex
    team := &Team{
        VID:  "team1001",
        Name: "Lakers",
    }
    if err := db.InsertVertex(team).Exec(); err != nil {
        log.Fatalf("insert team failed: %v", err)
    }
    // insert the serve edge
    serve := &Serve{
        SrcID:     "player1001",
        DstID:     "team1001",
        StartYear: time.Date(1996, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
        EndYear:   time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
    }
    if err := db.InsertEdge(serve).Exec(); err != nil {
        log.Fatalf("insert serve failed: %v", err)
    }

    // find the player vertex
    player = new(Player)
    err = db.
        Fetch("player", "player1001").
        Yield("vertex as v").
        FindCol("v", player)
    if err != nil {
        log.Fatalf("fetch player failed: %v", err)
    }
    log.Printf("player: %+v", player)
    
    // count the number of vertexes that the player vertex connects through different edges
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
}
```

## Features

* Fast splicing of nGQL by chained calls
* Friendly support for parsing and assigning compound types such as vertex, edge, list, map, set
* Fully unit tested
* Developer Friendly

## Contributing

Contributions are welcome! Please submit a pull request.

## Acknowledgements

This project was inspired and helped by the following open source projects during the development process:

* **gorm**: The fantastic ORM library for Golang, aims to be developer friendly.

Thanks to the authors of these projects for their contributions to the open source community!

## License

2024-NOW hayson

Released under the [MIT License](./LICENSE)