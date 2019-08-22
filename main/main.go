package main

import (
	"fmt"

	"github.com/jcasado94/connecc"
)

func main() {
	g, err := connecc.NewGenGraph("bolt://localhost:7687", "neo4j", "prod")
	if err != nil {
		panic(err)
	}
	fmt.Println(g.Connections(23))
}
