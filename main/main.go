package main

import (
	"fmt"

	"github.com/jcasado94/connecc/graph"
	"github.com/jcasado94/kstar"
)

func main() {
	g, err := graph.NewGenGraph(3, 0, "bolt://localhost:7687", "neo4j", "prod")
	if err != nil {
		panic(err)
	}
	fmt.Println(kstar.Run(g, 2))
}
