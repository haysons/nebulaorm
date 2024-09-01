# nebulaorm

[English](README.md)

[![go report card](https://goreportcard.com/badge/haysons/nebulaorm)](https://goreportcard.com/report/github.com/haysons/nebulaorm)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## 简介

nebulaorm 是一个专为 nebula graph 设计的 orm 框架，通过链式调用以更优雅、快速的方式拼接 nGQL 语句，并解析返回的结果集，
将其赋值给开发者提供的变量，旨在提高golang对于nebula graph的使用体验。

## 安装

```
go get github.com/haysons/nebulaorm
```

## 快速开始

``` go
// Player 节点
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

// Team 节点
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

// Serve 边
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

// 注: 目前在声明结构体时还不支持内嵌其他结构体，重复的字段请保留多份
func main() {
    // 初始化db对象
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
    
    // 写入player节点
    player := &Player{
        VID:  "player1001",
        Name: "Kobe Bryant",
        Age:  33,
    }
    if err := db.InsertVertex(player).Exec(); err != nil {
        log.Fatalf("insert player failed: %v", err)
    }
    // 写入team节点
    team := &Team{
        VID:  "team1001",
        Name: "Lakers",
    }
    if err := db.InsertVertex(team).Exec(); err != nil {
        log.Fatalf("insert team failed: %v", err)
    }
    // 写入serve边
    serve := &Serve{
        SrcID:     "player1001",
        DstID:     "team1001",
        StartYear: time.Date(1996, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
        EndYear:   time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
    }
    if err := db.InsertEdge(serve).Exec(); err != nil {
        log.Fatalf("insert serve failed: %v", err)
    }

    // 查询player节点
    player = new(Player)
    err = db.
        Fetch("player", "player1001").
        Yield("vertex as v").
        FindCol("v", player)
    if err != nil {
        log.Fatalf("fetch player failed: %v", err)
    }
    log.Printf("player: %+v", player)
    
    // 统计player节点通过不同边关联到的节点的数量
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

## 特性

* 通过链式调用快速拼接nGQL语句
* 对于复合类型的解析和赋值提供友好地支持，例如：vertex, edge, list, map, set
* 完善的单元测试
* 开发者友好

## 贡献

欢迎您做贡献! 请提交 pull request.

## 致谢

本项目在开发过程中得到了以下开源项目的启发和帮助：

* **gorm**: 适用于 Golang 的梦幻般的 ORM 库，旨在为开发人员提供方便。

感谢这些项目的作者为开源社区做出的贡献！

## 许可证

2024-NOW hayson

使用 [MIT License](./LICENSE)