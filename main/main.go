package main

import (
	"fmt"

	"github.com/jcasado94/connecc"
)

func main() {
	g, err := connecc.NewGenGraph()
	if err != nil {
		panic(err)
	}
	fmt.Println(g.Connections(23))
}
