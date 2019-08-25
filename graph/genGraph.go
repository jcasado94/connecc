package graph

import (
	"github.com/jcasado94/connecc/db"
)

type genConnectionInfo struct {
	provider int
}

type genGraph struct {
	mDriver            mongoDriver
	dbDriver           db.Driver
	connectionsCache   map[int]map[int][]float64
	genConnectionsInfo map[int]map[int][]genConnectionInfo
	s, t               int
}

func NewGenGraph(s, t int, dbEndpoint, dbUsername, dbPw string) (genGraph, error) {
	driver, err := db.NewDriver(dbEndpoint, dbUsername, dbPw, false)
	if err != nil {
		return genGraph{}, err
	}
	mDriver, err := newMongoDriver()
	if err != nil {
		return genGraph{}, err
	}
	g := genGraph{
		mDriver:            mDriver,
		dbDriver:           driver,
		connectionsCache:   make(map[int]map[int][]float64),
		genConnectionsInfo: make(map[int]map[int][]genConnectionInfo),
		s:                  s,
		t:                  t,
	}

	return g, nil
}

func (g genGraph) Connections(n int) map[int][]float64 {

	if _, exists := g.connectionsCache[n]; exists {
		return g.connectionsCache[n]
	}

	// concurrent?
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

func (g genGraph) retrieveGenConnections(n int) error {
	g.genConnectionsInfo[n] = make(map[int][]genConnectionInfo)

	neighboursGenResult, err := g.dbDriver.NeighboursGen(n)
	if err != nil {
		return err
	}

	gn, err := buildGenNeighbours(neighboursGenResult)
	if err != nil {
		return err
	}
	for _, gcon := range gn {
		if _, exists := g.connectionsCache[gcon.Id]; !exists {
			g.connectionsCache[n][gcon.Id] = make([]float64, 0)
			g.genConnectionsInfo[n][gcon.Id] = make([]genConnectionInfo, 0)
		}
		g.connectionsCache[n][gcon.Id] = append(g.connectionsCache[n][gcon.Id], gcon.Price)
		g.genConnectionsInfo[n][gcon.Id] = append(g.genConnectionsInfo[n][gcon.Id], genConnectionInfo{provider: gcon.Provider})
	}

	return nil

}

// Get the neighbours through the BelongsTo City node, plus the City node itself, excluding S. City nodes shall return no neighbours, except for S.
func (g genGraph) retrieveBelongsToConnections(n int) error {

	//concurrent?
	neighboursBelongsToCityResult, err := g.dbDriver.NeighboursBelongsToCity(n, g.S())
	if err != nil {
		return err
	}

	neighboursBelongsToThroughCityResult, err := g.dbDriver.NeighboursBelongsToThroughCity(n, g.S())
	if err != nil {
		return err
	}

	btn, err := buildBelongsToNeighbours(neighboursBelongsToCityResult, neighboursBelongsToThroughCityResult)
	if err != nil {
		return err
	}
	for _, btcon := range btn {
		if _, exists := g.connectionsCache[n][btcon.Id]; !exists {
			g.connectionsCache[n][btcon.Id] = []float64{btcon.Cost}
		}
	}

	return nil
}

func (g genGraph) S() int {
	return g.s
}

func (g genGraph) T() int {
	return g.t
}

func (g genGraph) FValue(n int) float64 {
	avgPrice, err := g.mDriver.getAvgPrice(n, g.T())
	if err != nil {
		panic(err)
	}
	return avgPrice
}
