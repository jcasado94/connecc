package db

import "github.com/neo4j/neo4j-go-driver/neo4j"

type genConnection struct {
	Id       int
	Price    float64
	Provider int
}

type neighboursGenResponse struct {
	Neighbours []genConnection
}

func newNeighboursGenResponse(result neo4j.Result) (neighboursGenResponse, error) {
	resp := neighboursGenResponse{
		Neighbours: make([]genConnection, 0),
	}

	var next bool
	for next = result.Next(); next; next = result.Next() {
		rec := result.Record()
		resp.Neighbours = append(resp.Neighbours, genConnection{
			int(rec.GetByIndex(0).(int64)),
			rec.GetByIndex(1).(float64),
			int(rec.GetByIndex(2).(int64)),
		})
	}

	return resp, result.Err()
}
