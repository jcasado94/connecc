package connecc_test

import (
	"testing"

	"github.com/jcasado94/connecc"
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

func testConnectionsGraphMock(session neo4j.Session) (result []int, err error) {

	response, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"CREATE (a:Airport{code: $code1})-[r:Gen{price:$price, provider:$provider}]->(b:Airport{code: $code2}), (c:City{name: $city1}), (d:City{name:$city2}) RETURN id(a), id(b), id(c)",
			map[string]interface{}{"code1": "YYZ", "price": 200.0, "provider": 0, "code2": "JFK", "city1": "Toronto", "city2": "New York"},
		)
		if err != nil {
			return nil, err
		}
		result.Next()
		record := result.Record()
		return []int{int(record.GetByIndex(0).(int64)), int(record.GetByIndex(1).(int64)), int(record.GetByIndex(2).(int64))}, nil
	})

	if err != nil {
		return []int{}, err
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Airport), (b:City), (c:Airport), (d:City) WHERE a.code=$code1 AND b.name=$name1 AND c.code=$code2 AND d.name=$name2 CREATE (a)-[r:BelongsTo]->(b), (c)-[s:BelongsTo]->(d)",
			map[string]interface{}{"code1": "JFK", "name1": "New York", "code2": "YYZ", "name2": "Toronto"},
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

func TestConnections(t *testing.T) {
	driver, session := newSessionTestGraph()
	defer driver.Close()
	defer session.Close()

	ids, err := testConnectionsGraphMock(session)
	if err != nil {
		t.Error(err)
	}
	idYYZ, idJFK, idToronto := ids[0], ids[1], ids[2]

	g, err := connecc.NewGenGraph(dbTestEndpoint, dbTestUsername, dbTestPw)
	if err != nil {
		t.Error(err)
	}

	yyzConnections := g.Connections(idYYZ)
	expectedConnections := map[int][]float64{idJFK: []float64{200.0}, idToronto: []float64{0.0}}
	for key, expectedSlice := range expectedConnections {
		if _, exists := yyzConnections[key]; !exists {
			t.Error("genGraph.Connections did not return the expected connections")
		}
		slice := yyzConnections[key]
		for i, expectedValue := range expectedSlice {
			if expectedValue != slice[i] {
				t.Error("genGraph.Connections did not return the expected connections")
			}
		}
	}

	err = cleanDb(session)
	if err != nil {
		t.Error(err)
	}
}

