package graph

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	dbTestEndpoint = "bolt://localhost:7687"
	dbTestUsername = "neo4j"
	dbTestPw       = "test"
)

func newSessionTestGraph() (neo4j.Driver, neo4j.Session) {
	driver, err := neo4j.NewDriver(dbTestEndpoint, neo4j.BasicAuth(dbTestUsername, dbTestPw, ""))
	if err != nil {
		panic(err)
	}
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		panic(err)
	}
	return driver, session
}

func cleanDb(session neo4j.Session) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH ()-[r]-() DELETE r",
			nil,
		)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return err
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (n) DELETE n",
			nil,
		)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	return nil

}

func graphMock(session neo4j.Session) (result []int, err error) {

	response, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"CREATE (a:Airport{code: $code1})-[r:Gen{price:$price, provider:$provider}]->(b:Airport{code: $code2}), (e:Airport{code: $code3}), (c:City{name: $city1}), (d:City{name:$city2}) RETURN id(a), id(b), id(e), id(c), id(d)",
			map[string]interface{}{"code1": "YYZ", "price": 200.0, "provider": 0, "code2": "JFK", "code3": "LGA", "city1": "Toronto", "city2": "New York"},
		)
		if err != nil {
			return nil, err
		}
		result.Next()
		record := result.Record()
		return []int{int(record.GetByIndex(0).(int64)), int(record.GetByIndex(1).(int64)), int(record.GetByIndex(2).(int64)), int(record.GetByIndex(3).(int64)), int(record.GetByIndex(4).(int64))}, nil
	})

	if err != nil {
		return []int{}, err
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Airport), (b:City), (c:Airport), (d:City), (e:Airport) WHERE a.code=$code1 AND b.name=$name1 AND c.code=$code2 AND e.code=$code3 AND d.name=$name2 CREATE (a)-[r:BelongsTo]->(b), (c)-[s:BelongsTo]->(d), (e)-[t:BelongsTo]->(b)",
			map[string]interface{}{"code1": "JFK", "name1": "New York", "code2": "YYZ", "code3": "LGA", "name2": "Toronto"},
		)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return []int{}, err
	}

	return response.([]int), err
}

func newMockGenGraph(session neo4j.Session, t *testing.T) (g *genGraph, ids []int) {
	ids, err := graphMock(session)
	if err != nil {
		t.Fail()
	}

	idYYZ, idJFK, idLGA, idToronto, idNewYork := ids[0], ids[1], ids[2], ids[3], ids[4]
	t.Logf("IdYYZ: %d\nIdJFK: %d\nIdLGA: %d\nIdToronto: %d\nIdNewYork: %d\n", idYYZ, idJFK, idLGA, idToronto, idNewYork)

	g, err = NewGenGraph(idNewYork, idToronto, dbTestEndpoint, dbTestUsername, dbTestPw)
	if err != nil {
		t.Fail()
	}
	return g, ids
}

func TestConnections(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()

	g, ids := newMockGenGraph(session, t)
	idYYZ, idJFK, idLGA, idToronto, idNewYork := ids[0], ids[1], ids[2], ids[3], ids[4]

	t.Run("Test Connections result", func(t *testing.T) {
		testCases := []struct {
			id                  int
			expectedConnections map[int][]float64
		}{
			{idYYZ, map[int][]float64{idJFK: []float64{200.0}, idToronto: []float64{defaultCostBelongsToCity}}},
			{idJFK, map[int][]float64{idLGA: []float64{defaultCostBelongsTo}}},
			{idNewYork, map[int][]float64{idJFK: []float64{defaultCostBelongsToCity}, idLGA: []float64{defaultCostBelongsToCity}}},
			{idToronto, make(map[int][]float64)},
		}
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Connections for %d", tc.id), func(t *testing.T) {
				connections := g.Connections(tc.id)
				if !reflect.DeepEqual(tc.expectedConnections, connections) {
					t.Errorf("Expected %v\ngot\n%v", tc.expectedConnections, connections)
					cleanDb(session)
					return
				}
			})
		}
	})

	t.Run("Test stored cache", func(t *testing.T) {
		expectedNodesCache := map[int]node{
			idJFK:     newAirport(idJFK, "JFK"),
			idLGA:     newAirport(idLGA, "LGA"),
			idToronto: newCity(idToronto, "Toronto"),
			idNewYork: newCity(idNewYork, "New York"),
		}
		nodesCache := make(map[int]node)
		for t := range g.cache.nodesCache.cm.Iter() {
			k, _ := strconv.Atoi(t.Key)
			nodesCache[k] = t.Val.(node)
		}
		if !reflect.DeepEqual(expectedNodesCache, nodesCache) {
			t.Errorf("Expected %v\ngot\n%v", expectedNodesCache, nodesCache)
		}
	})

	cleanDb(session)

}

func TestNewGenGraph(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()

	g, ids := newMockGenGraph(session, t)
	idNewYork := ids[4]

	if _, exists := g.cache.nodesCache.checkGet(idNewYork); !exists {
		t.Errorf("No cached node for node s: %d", idNewYork)
	}
	expectedNode := newCity(idNewYork, "New York")
	n := g.cache.nodesCache.get(idNewYork).(node)
	if !expectedNode.Equals(n) {
		t.Errorf("Cached s node differs. Expected %v, got %v.", expectedNode, n)
	}

	cleanDb(session)
}

func TestSetBelongsToRelationship(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()
	g, ids := newMockGenGraph(session, t)
	c := &g.cache
	idNewYork := ids[4]
	expectedMap := map[int][]float64{
		1: []float64{0.0},
	}
	c.cache.set(idNewYork, make(map[int][]float64))
	c.setBelongsToRelationship(idNewYork, 1, 0.0)
	cons := c.cache.get(idNewYork).(map[int][]float64)
	if !reflect.DeepEqual(expectedMap, cons) {
		t.Errorf("Expected %v,\ngot %v", expectedMap, cons)
	}
	cleanDb(session)
}

func TestSetGeneralRelationship(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()
	g, ids := newMockGenGraph(session, t)
	c := &g.cache
	idNewYork := ids[4]
	c.cache.set(idNewYork, make(map[int][]float64))
	c.infoCache.set(idNewYork, make(map[int][]genConnectionInfo))
	expectedCacheMap := map[int][]float64{
		1: []float64{1.0},
	}
	expectedInfoCacheMap := map[int][]genConnectionInfo{
		1: []genConnectionInfo{genConnectionInfo{provider: 0}},
	}
	c.setGeneralRelationship(idNewYork, 1, 0, 1.0)
	consCache := c.cache.get(idNewYork).(map[int][]float64)
	consInfoCache := c.infoCache.get(idNewYork).(map[int][]genConnectionInfo)
	if !reflect.DeepEqual(expectedCacheMap, consCache) {
		t.Errorf("Expected %v,\ngot\n %v", expectedCacheMap, consCache)
	}
	if !reflect.DeepEqual(expectedInfoCacheMap, consInfoCache) {
		t.Errorf("Expected %v,\ngot\n %v", expectedInfoCacheMap, consInfoCache)
	}
	cleanDb(session)
}

func TestInvlaidateCache(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()
	g, ids := newMockGenGraph(session, t)
	c := &g.cache
	idYYZ, idJFK := ids[0], ids[1]
	c.cache.set(idYYZ, make(map[int][]float64))
	c.infoCache.set(idYYZ, make(map[int][]genConnectionInfo))
	now := time.Now()
	c.connectionsTimeStamp.set(idYYZ, now)
	c.invalidateCache(idYYZ)
	expectedCacheMap := map[int][]float64{idJFK: []float64{200.0}}
	expectedInfoCacheMap := map[int][]genConnectionInfo{idJFK: []genConnectionInfo{genConnectionInfo{provider: 0}}}
	if !reflect.DeepEqual(c.cache.get(idYYZ), expectedCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedCacheMap, c.cache.get(idYYZ))
	}
	if !reflect.DeepEqual(c.infoCache.get(idYYZ), expectedInfoCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedInfoCacheMap, c.infoCache.get(idYYZ))
	}
	if c.connectionsTimeStamp.get(idYYZ).(time.Time).Equal(now) {
		t.Error("Timestamp hasn't changed.")
	}
	cleanDb(session)
}

func TestInitializeCache(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()
	g, ids := newMockGenGraph(session, t)
	c := &g.cache
	idYYZ, idJFK, idToronto := ids[0], ids[1], ids[3]
	now := time.Now()
	expectedCacheMap := map[int][]float64{idJFK: []float64{200.0}, idToronto: []float64{defaultCostBelongsToCity}}
	expectedInfoCacheMap := map[int][]genConnectionInfo{idJFK: []genConnectionInfo{genConnectionInfo{provider: 0}}}
	c.initializeCache(idYYZ)
	if !reflect.DeepEqual(c.cache.get(idYYZ), expectedCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedCacheMap, c.cache.get(idYYZ))
	}
	if !reflect.DeepEqual(c.infoCache.get(idYYZ), expectedInfoCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedInfoCacheMap, c.infoCache.get(idYYZ))
	}
	if c.connectionsTimeStamp.get(idYYZ).(time.Time).Sub(now) <= 0 {
		t.Error("Timestamp was not set")
	}
	cleanDb(session)
}

func TestGetOrInvalidate(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()
	g, ids := newMockGenGraph(session, t)
	c := &g.cache
	idYYZ, idJFK, idToronto := ids[0], ids[1], ids[3]
	expectedInfoCacheMap := map[int][]genConnectionInfo{idJFK: []genConnectionInfo{genConnectionInfo{provider: 0}}}

	// No timestamp
	expectedCacheMap := map[int][]float64{idJFK: []float64{200.0}, idToronto: []float64{defaultCostBelongsToCity}}
	c.getOrInvalidate(idYYZ)
	if _, ok := c.connectionsTimeStamp.checkGet(idYYZ); !ok {
		t.Error("Timestamp was not set")
	}
	if !reflect.DeepEqual(c.cache.get(idYYZ), expectedCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedCacheMap, c.cache.get(idYYZ))
	}
	if !reflect.DeepEqual(c.infoCache.get(idYYZ), expectedInfoCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedInfoCacheMap, c.infoCache.get(idYYZ))
	}
	c.cache = newIntCMap()
	c.infoCache = newIntCMap()
	c.connectionsTimeStamp = newIntCMap()

	// Old timestamp
	expectedCacheMap = map[int][]float64{idJFK: []float64{200.0}}
	ts, _ := time.Parse(time.RFC822, time.RFC822)
	c.connectionsTimeStamp.set(idYYZ, ts)
	c.cache.set(idYYZ, make(map[int][]float64))
	c.infoCache.set(idYYZ, make(map[int][]genConnectionInfo))
	c.getOrInvalidate(idYYZ)
	if c.connectionsTimeStamp.get(idYYZ).(time.Time).Equal(ts) {
		t.Error("Timestamp did not change")
	}
	if !reflect.DeepEqual(c.cache.get(idYYZ), expectedCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedCacheMap, c.cache.get(idYYZ))
	}
	if !reflect.DeepEqual(c.infoCache.get(idYYZ), expectedInfoCacheMap) {
		t.Errorf("Expected %v,\ngot\n %v", expectedInfoCacheMap, c.infoCache.get(idYYZ))
	}
	cleanDb(session)
}

func TestSetNode(t *testing.T) {
	c := newGenGraphCache(nil)
	c.nodesCache = newIntCMap()
	n := newNode("Airport", 0, map[string]interface{}{"code": "NYZ"})
	c.setNode(0, &n)
	if _, ok := c.nodesCache.checkGet(0); !ok {
		t.Error("Didn't store node correctly")
	}
}

