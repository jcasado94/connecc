package graph

import "github.com/neo4j/neo4j-go-driver/neo4j"

const defaultCostBelongsToCity = 0.0
const defaultCostBelongsTo = 100.0

type genConnection struct {
	Price    float64
	Provider int
	n        node
}

type genNeighbours []genConnection

func buildGenNeighbours(result neo4j.Result) (genNeighbours, error) {
	resp := make(genNeighbours, 0)

	var next bool
	for next = result.Next(); next; next = result.Next() {
		rec := result.Record()
		resp = append(resp, genConnection{
			Price:    rec.GetByIndex(0).(float64),
			Provider: int(rec.GetByIndex(1).(int64)),
			n:        newNode(rec.GetByIndex(2).(string), int(rec.GetByIndex(3).(int64)), rec.GetByIndex(4).(map[string]interface{})),
		})
	}

	return resp, result.Err()
}

type belongsToConnection struct {
	Cost float64
	n    node
}

type belongsToNeighbours []belongsToConnection

func buildBelongsToNeighbours(resultCity, resultThroughCity neo4j.Result) (belongsToNeighbours, error) {
	resp := make(belongsToNeighbours, 0)

	var next bool
	for next = resultCity.Next(); next; next = resultCity.Next() {
		rec := resultCity.Record()
		resp = append(resp, belongsToConnection{
			Cost: defaultCostBelongsToCity,
			n:    newNode(rec.GetByIndex(0).(string), int(rec.GetByIndex(1).(int64)), rec.GetByIndex(2).(map[string]interface{})),
		})
	}
	if resultCity.Next() {

	}

	if resultCity.Err() != nil {
		return nil, resultCity.Err()
	}

	for next = resultThroughCity.Next(); next; next = resultThroughCity.Next() {
		rec := resultThroughCity.Record()
		resp = append(resp, belongsToConnection{
			Cost: defaultCostBelongsTo, // find cost function (google maps?)
			n:    newNode(rec.GetByIndex(0).(string), int(rec.GetByIndex(1).(int64)), rec.GetByIndex(2).(map[string]interface{})),
		})
	}

	if resultThroughCity.Err() != nil {
		return nil, resultThroughCity.Err()
	}

	return resp, nil
}
