package connecc

import (
	"github.com/jcasado94/connecc/db"
)

type genGraphConnection struct {
	provider int
}

type genGraph struct {
	driver           db.Driver
	connectionsCache map[int]map[int][]float64
	connectionsInfo  map[int]map[int][]genGraphConnection
}

func NewGenGraph() (genGraph, error) {
	driver, err := db.NewDriver(false)
	if err != nil {
		return genGraph{}, err
	}
	g := genGraph{
		driver:           driver,
		connectionsCache: make(map[int]map[int][]float64),
		connectionsInfo:  make(map[int]map[int][]genGraphConnection),
	}

	return g, nil
}

func (g *genGraph) Connections(n int) map[int][]float64 {

	if _, exists := g.connectionsCache[n]; exists {
		return g.connectionsCache[n]
	}

	neighboursGenResponse, err := g.driver.NeighboursGen(n)
	if err != nil {
		panic(err)
	}

	if _, exists := g.connectionsCache[n]; !exists {
		g.connectionsCache[n] = make(map[int][]float64)
		g.connectionsInfo[n] = make(map[int][]genGraphConnection)
	}

	for _, genConnection := range neighboursGenResponse.Neighbours {
		if _, exists := g.connectionsCache[genConnection.Id]; !exists {
			g.connectionsCache[n][genConnection.Id] = make([]float64, 0)
			g.connectionsInfo[n][genConnection.Id] = make([]genGraphConnection, 0)
		}
		g.connectionsCache[n][genConnection.Id] = append(g.connectionsCache[n][genConnection.Id], genConnection.Price)
		g.connectionsInfo[n][genConnection.Id] = append(g.connectionsInfo[n][genConnection.Id], genGraphConnection{provider: genConnection.Provider})
	}

	return g.connectionsCache[n]

}

