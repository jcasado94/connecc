package graph

import (
	"strconv"
	"time"

	"github.com/jcasado94/connecc/drivers"
	cmap "github.com/orcaman/concurrent-map"
)

var invalidateAgeGenRel = time.Hour * 24

type genGraph struct {
	mDriver  drivers.MongoDriver
	dbDriver drivers.DbDriver
	cache    genGeaphCache
	s, t     int
}

func NewGenGraph(s, t int, dbEndpoint, dbUsername, dbPw string) (*genGraph, error) {
	driver, err := drivers.NewDbDriver(dbEndpoint, dbUsername, dbPw, false)
	if err != nil {
		return &genGraph{}, err
	}
	mDriver, err := drivers.NewMongoDriver()
	if err != nil {
		return &genGraph{}, err
	}
	g := genGraph{
		mDriver:  mDriver,
		dbDriver: driver,
		s:        s,
		t:        t,
	}

	g.cache = newGenGraphCache(&g)

	err = g.cacheNodeInfo(s)
	if err != nil {
		return &g, err
	}

	return &g, nil
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
		g.cache.setNode(id, &node)
	}
	return result.Err()
}

func (g *genGraph) Connections(n int) map[int][]float64 {

	connections, err := g.cache.getOrInvalidate(n)
	if err != nil {
		panic(err)
	}
	return connections

}

func (g *genGraph) retrieveGenConnections(n int) error {

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
		g.cache.setNode(id, &gcon.n)
		g.cache.setGeneralRelationship(n, id, gcon.Provider, gcon.Price)
	}

	return nil

}

// Get the neighbours through the BelongsTo City node, plus the City node itself, excluding S. City nodes shall return no neighbours, except for S.
func (g *genGraph) retrieveBelongsToConnections(n int) error {

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
		g.cache.setNode(id, &btcon.n)
		g.cache.setBelongsToRelationship(n, id, btcon.Cost)
	}

	return nil
}

func (g *genGraph) S() int {
	return g.s
}

func (g *genGraph) T() int {
	return g.t
}

func (g *genGraph) FValue(n int) float64 {
	avgPrice, err := g.mDriver.GetAvgPrice(n, g.T())
	if err != nil {
		panic(err)
	}
	return avgPrice
}

type genGeaphCache struct {
	infoCache            intCMap // map[int]map[int][]genConnectionInfo
	cache                intCMap // map[int]map[int][]float64
	connectionsTimeStamp intCMap // map[int]time.Time
	nodesCache           intCMap // map[int]node
	g                    *genGraph
}

func newGenGraphCache(g *genGraph) genGeaphCache {
	return genGeaphCache{
		infoCache:            newIntCMap(),
		cache:                newIntCMap(),
		connectionsTimeStamp: newIntCMap(),
		nodesCache:           newIntCMap(),
		g:                    g,
	}
}

type genConnectionInfo struct {
	provider int
}

type intCMap struct {
	cm cmap.ConcurrentMap
}

func newIntCMap() intCMap {
	return intCMap{
		cm: cmap.New(),
	}
}

func (m *intCMap) checkGet(key int) (interface{}, bool) {
	return m.cm.Get(strconv.Itoa(key))
}

func (m *intCMap) get(key int) interface{} {
	val, _ := m.cm.Get(strconv.Itoa(key))
	return val
}

func (m *intCMap) set(key int, val interface{}) {
	m.cm.Set(strconv.Itoa(key), val)
}

func (m *intCMap) checkSet(key int, val interface{}) {
	if _, ok := m.checkGet(key); !ok {
		m.set(key, val)
	}
}

func (c *genGeaphCache) getOrInvalidate(n int) (map[int][]float64, error) {
	var err error
	tInt, ok := c.connectionsTimeStamp.checkGet(n)
	if !ok {
		err = c.initializeCache(n)
	} else if time.Now().Sub(tInt.(time.Time)) > invalidateAgeGenRel {
		err = c.invalidateCache(n)
	}
	return c.cache.get(n).(map[int][]float64), err
}

func (c *genGeaphCache) initializeCache(n int) error {
	c.connectionsTimeStamp.set(n, time.Now())
	c.cache.set(n, make(map[int][]float64))
	c.infoCache.set(n, make(map[int][]genConnectionInfo))
	err := c.g.retrieveGenConnections(n)
	if err != nil {
		return err
	}
	err = c.g.retrieveBelongsToConnections(n)
	if err != nil {
		return err
	}
	return nil
}

func (c *genGeaphCache) invalidateCache(n int) error {
	c.connectionsTimeStamp.set(n, time.Now())
	err := c.g.retrieveGenConnections(n)
	if err != nil {
		return err
	}
	return nil
}

func (c *genGeaphCache) setGeneralRelationship(n, id, provider int, price float64) {
	mCon := c.cache.get(n).(map[int][]float64)
	mConInfo := c.infoCache.get(n).(map[int][]genConnectionInfo)
	if _, exists := mCon[id]; !exists {
		mCon[id] = make([]float64, 0)
		mConInfo[id] = make([]genConnectionInfo, 0)
	}
	mCon[id] = append(mCon[id], price)
	mConInfo[id] = append(mConInfo[id], genConnectionInfo{provider: provider})
}

func (c *genGeaphCache) setBelongsToRelationship(n, id int, cost float64) {
	mCon := c.cache.get(n).(map[int][]float64)
	if _, exists := mCon[id]; !exists {
		mCon[id] = []float64{cost}
	}
}

func (c *genGeaphCache) setNode(id int, n *node) {
	c.nodesCache.checkSet(id, *n)
}
