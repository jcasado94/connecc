package pathprocessor

import (
	"github.com/jcasado94/connecc/graph"
	"github.com/jcasado94/connecc/scraping"
	"github.com/jcasado94/kstar"
)

type GenPathProcessor struct {
	g *graph.GenGraph
}

func newGenPathProcessor(g *graph.GenGraph) *GenPathProcessor {
	return &GenPathProcessor{
		g: g,
	}
}

func (c *GenPathProcessor) Process() (trips []scraping.Trip) {
	paths := kstar.Run(c.g, kpaths)
	s, t := c.g.S(), c.g.T()
	//concurrent
	for _, path := range paths {
		// convert path into trip
	}
	return trips
}
