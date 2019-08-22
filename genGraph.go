package connecc

import (
	"github.com/jcasado94/connecc/db"
)

type genConnectionInfo struct {
	provider int
}

type genGraph struct {
	driver             db.Driver
	connectionsCache   map[int]map[int][]float64
	genConnectionsInfo map[int]map[int][]genConnectionInfo
}

func NewGenGraph(dbEndpoint, dbUsername, dbPw string) (genGraph, error) {
	driver, err := db.NewDriver(dbEndpoint, dbUsername, dbPw, false)
	if err != nil {
		return genGraph{}, err
	}
	g := genGraph{
		driver:             driver,
		connectionsCache:   make(map[int]map[int][]float64),
		genConnectionsInfo: make(map[int]map[int][]genConnectionInfo),
	}

	return g, nil
}

func (g *genGraph) Connections(n int) map[int][]float64 {

	if _, exists := g.connectionsCache[n]; exists {
		return g.connectionsCache[n]
	}

	g.connectionsCache[n] = make(map[int][]float64)
	err := g.retrieveGenConnections(n)
	if err != nil {
		panic(err)
	}

	err = g.retrieveBelongsToConnections(n)
	if err != nil {
		panic(err)
	}

	return g.connectionsCache[n]

}

func (g *genGraph) retrieveGenConnections(n int) error {
	g.genConnectionsInfo[n] = make(map[int][]genConnectionInfo)

	neighboursGenResponse, err := g.driver.NeighboursGen(n)
	if err != nil {
		return err
	}

	for _, genConnection := range neighboursGenResponse.Neighbours {
		if _, exists := g.connectionsCache[genConnection.Id]; !exists {
			g.connectionsCache[n][genConnection.Id] = make([]float64, 0)
			g.genConnectionsInfo[n][genConnection.Id] = make([]genConnectionInfo, 0)
		}
		g.connectionsCache[n][genConnection.Id] = append(g.connectionsCache[n][genConnection.Id], genConnection.Price)
		g.genConnectionsInfo[n][genConnection.Id] = append(g.genConnectionsInfo[n][genConnection.Id], genConnectionInfo{provider: genConnection.Provider})
	}

	return nil

}

func (g *genGraph) retrieveBelongsToConnections(n int) error {
	neighboursBelongsToResponse, err := g.driver.NeighboursBelongsTo(n)
	if err != nil {
		return err
	}

	for _, belongsToConnection := range neighboursBelongsToResponse.Neighbours {
		if _, exists := g.connectionsCache[belongsToConnection.Id]; !exists {
			g.connectionsCache[n][belongsToConnection.Id] = []float64{belongsToConnection.Cost}
		}
	}

	return nil
}
