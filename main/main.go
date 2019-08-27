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
	paths := kstar.Run(g, 1)
	pathsObjects := make([][]kstar.Edge, 0)
	for _, p := range paths {
		path := make([]kstar.Edge, 0)
		for _, e := range p {
			path = append(path, *e)
		}
		pathsObjects = append(pathsObjects, path)
	}
	fmt.Println(pathsObjects)
}
