package graph

import "github.com/neo4j/neo4j-go-driver/neo4j"

const defaultCostBelongsToCity = 0.0
const defaultCostBelongsTo = 100.0

type genConnection struct {
	Id       int
	Price    float64
	Provider int
}

type genNeighbours []genConnection

func buildGenNeighbours(result neo4j.Result) (genNeighbours, error) {
	resp := make(genNeighbours, 0)

	var next bool
	for next = result.Next(); next; next = result.Next() {
		rec := result.Record()
		resp = append(resp, genConnection{
			Id:       int(rec.GetByIndex(0).(int64)),
			Price:    rec.GetByIndex(1).(float64),
			Provider: int(rec.GetByIndex(2).(int64)),
		})
	}

	return resp, result.Err()
}

type belongsToConnection struct {
	Id   int
	Cost float64
}

type belongsToNeighbours []belongsToConnection

func buildBelongsToNeighbours(resultCity, resultThroughCity neo4j.Result) (belongsToNeighbours, error) {
	resp := make(belongsToNeighbours, 0)

	var next bool
	for next = resultCity.Next(); next; next = resultCity.Next() {
		resp = append(resp, belongsToConnection{
			Id:   int(resultCity.Record().GetByIndex(0).(int64)),
			Cost: defaultCostBelongsToCity,
		})
	}
	if resultCity.Next() {

	}

	if resultCity.Err() != nil {
		return nil, resultCity.Err()
	}

	for next = resultThroughCity.Next(); next; next = resultThroughCity.Next() {
		resp = append(resp, belongsToConnection{
			int(resultThroughCity.Record().GetByIndex(0).(int64)),
			defaultCostBelongsTo, // find cost function (google maps?)
		})
	}

	if resultThroughCity.Err() != nil {
		return nil, resultThroughCity.Err()
	}

	return resp, nil
}
