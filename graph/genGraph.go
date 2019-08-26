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
	genConnectionsInfo map[int]map[int][]*genConnectionInfo
	nodesCache         map[int]node
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
		genConnectionsInfo: make(map[int]map[int][]*genConnectionInfo),
		nodesCache:         make(map[int]node),
		s:                  s,
		t:                  t,
	}

	err = g.cacheNodeInfo(s)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (g *genGraph) cacheNodeInfo(id int) error {
	result, err := g.dbDriver.NodeInfo(id)
	if err != nil {
		return err
	}
	if result.Next() {
		rec := result.Record()
		label := rec.GetByIndex(0).(string)
		params := rec.GetByIndex(1).(map[string]interface{})
		node := newNode(label, id, params)
		g.nodesCache[id] = node
	}
	return result.Err()
}

func (g genGraph) Connections(n int) map[int][]float64 {

	if _, exists := g.connectionsCache[n]; exists {
		return g.connectionsCache[n]
	}

	g.connectionsCache[n] = make(map[int][]float64)

	// concurrent?
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
	g.genConnectionsInfo[n] = make(map[int][]*genConnectionInfo)

	neighboursGenResult, err := g.dbDriver.NeighboursGen(n)
	if err != nil {
		return err
	}

	gn, err := buildGenNeighbours(neighboursGenResult)
	if err != nil {
		return err
	}
	for _, gcon := range gn {
		id := gcon.n.Id()
		if _, exists := g.nodesCache[id]; !exists {
			g.nodesCache[id] = gcon.n
		}
		if _, exists := g.connectionsCache[id]; !exists {
			g.connectionsCache[n][id] = make([]float64, 0)
			g.genConnectionsInfo[n][id] = make([]*genConnectionInfo, 0)
		}
		g.connectionsCache[n][id] = append(g.connectionsCache[n][id], gcon.Price)
		g.genConnectionsInfo[n][id] = append(g.genConnectionsInfo[n][id], &genConnectionInfo{provider: gcon.Provider})
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
		id := btcon.n.Id()
		if _, exists := g.nodesCache[id]; !exists {
			g.nodesCache[id] = btcon.n
		}
		if _, exists := g.connectionsCache[n][id]; !exists {
			g.connectionsCache[n][id] = []float64{btcon.Cost}
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
