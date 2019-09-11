package graph

import (
	"fmt"
	"reflect"
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
				if len(tc.expectedConnections) != len(connections) {
					t.Errorf("connection maps differ. Want %v, got %v", tc.expectedConnections, connections)
					return
				}
				for key, expectedSlice := range tc.expectedConnections {
					if _, exists := connections[key]; !exists {
						t.Errorf("connection maps differ. Want %v, got %v", tc.expectedConnections, connections)
						continue
					}
					slice := connections[key]
					if len(expectedSlice) != len(slice) {
						t.Errorf("connection maps differ. Want %v, got %v", tc.expectedConnections, connections)
						continue
					}
					for i, expectedValue := range expectedSlice {
						if expectedValue != slice[i] {
							t.Errorf("connection maps differ. Want %v, got %v", tc.expectedConnections, connections)
							break
						}
					}
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
		if len(expectedNodesCache) != len(g.nodesCache) {
			cleanDb(session)
			t.Errorf("nodesCache not properly stored. Want %v, got %v", expectedNodesCache, g.nodesCache)
		}
		for id, expectedNode := range expectedNodesCache {
			if n, exists := g.nodesCache[id]; !exists || !expectedNode.Equals(n) {
				cleanDb(session)
				t.Errorf("nodesCache not properly stored. Want %v, got %v", expectedNodesCache, g.nodesCache)
			}
		}
	})

	cleanDb(session)

}

func TestNewGenGraph(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()

	g, ids := newMockGenGraph(session, t)
	idNewYork := ids[0]

	if _, exists := g.nodesCache[idNewYork]; !exists {
		t.Errorf("No cached node for node s: %d", idNewYork)
	}
	expectedNode := newCity(idNewYork, "New York")
	if !expectedNode.Equals(g.nodesCache[idNewYork]) {
		t.Errorf("Cached s node differs. Expected %v, got %v.", expectedNode, g.nodesCache[idNewYork])
	}

	cleanDb(session)
}

func TestGenConnectionCache(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()
	g, ids := newMockGenGraph(session, t)
	c := &g.connectionsCache
	idYYZ, idJFK, _, idToronto, idNewYork := ids[0], ids[1], ids[2], ids[3], ids[4]
	t.Run("SetBelongsToRelationship", func(t *testing.T) {
		expectedMap := map[int][]float64{
			1: []float64{0.0},
		}
		c.cache.set(idNewYork, make(map[int][]float64))
		c.setBelongsToRelationship(idNewYork, 1, 0.0)
		cons := c.cache.get(idNewYork).(map[int][]float64)
		if !reflect.DeepEqual(expectedMap, cons) {
			t.Errorf("Expected %v,\ngot %v", expectedMap, cons)
		}
		c.cache = newIntCMap()
	})

	t.Run("SetGeneralRelationship", func(t *testing.T) {
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
		c.cache = newIntCMap()
		c.infoCache = newIntCMap()
	})

	t.Run("InvalidateCache", func(t *testing.T) {
		c.cache.set(idYYZ, make(map[int][]float64))
		c.infoCache.set(idYYZ, make(map[int][]genConnectionInfo))
		now := time.Now()
		c.timeStamp.set(idYYZ, now)
		c.invalidateCache(idYYZ)
		expectedCacheMap := map[int][]float64{idJFK: []float64{200.0}}
		expectedInfoCacheMap := map[int][]genConnectionInfo{idJFK: []genConnectionInfo{genConnectionInfo{provider: 0}}}
		if !reflect.DeepEqual(c.cache.get(idYYZ), expectedCacheMap) {
			t.Errorf("Expected %v,\ngot\n %v", expectedCacheMap, c.cache.get(idYYZ))
		}
		if !reflect.DeepEqual(c.infoCache.get(idYYZ), expectedInfoCacheMap) {
			t.Errorf("Expected %v,\ngot\n %v", expectedInfoCacheMap, c.infoCache.get(idYYZ))
		}
		if c.timeStamp.get(idYYZ).(time.Time).Equal(now) {
			t.Error("Timestamp hasn't changed.")
		}
		c.cache = newIntCMap()
		c.infoCache = newIntCMap()
		c.timeStamp = newIntCMap()
	})

	t.Run("InitializeCache", func(t *testing.T) {
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
		if c.timeStamp.get(idYYZ).(time.Time).Sub(now) <= 0 {
			t.Error("Timestamp was not set")
		}
		c.cache = newIntCMap()
		c.infoCache = newIntCMap()
		c.timeStamp = newIntCMap()
	})

	t.Run("GetOrInvalidate", func(t *testing.T) {
		expectedInfoCacheMap := map[int][]genConnectionInfo{idJFK: []genConnectionInfo{genConnectionInfo{provider: 0}}}

		// No timestamp
		expectedCacheMap := map[int][]float64{idJFK: []float64{200.0}, idToronto: []float64{defaultCostBelongsToCity}}
		c.getOrInvalidate(idYYZ)
		if _, ok := c.timeStamp.checkGet(idYYZ); !ok {
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
		c.timeStamp = newIntCMap()

		// Old timestamp
		expectedCacheMap = map[int][]float64{idJFK: []float64{200.0}}
		ts, _ := time.Parse(time.RFC822, time.RFC822)
		c.timeStamp.set(idYYZ, ts)
		c.cache.set(idYYZ, make(map[int][]float64))
		c.infoCache.set(idYYZ, make(map[int][]genConnectionInfo))
		c.getOrInvalidate(idYYZ)
		if c.timeStamp.get(idYYZ).(time.Time).Equal(ts) {
			t.Error("Timestamp did not change")
		}
		if !reflect.DeepEqual(c.cache.get(idYYZ), expectedCacheMap) {
			t.Errorf("Expected %v,\ngot\n %v", expectedCacheMap, c.cache.get(idYYZ))
		}
		if !reflect.DeepEqual(c.infoCache.get(idYYZ), expectedInfoCacheMap) {
			t.Errorf("Expected %v,\ngot\n %v", expectedInfoCacheMap, c.infoCache.get(idYYZ))
		}
		c.cache = newIntCMap()
		c.infoCache = newIntCMap()
		c.timeStamp = newIntCMap()
	})

	cleanDb(session)
}
