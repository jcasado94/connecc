package db

import (
	neo4j "github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	endpoint = "bolt://localhost:7687"
)

type Driver struct {
	driver  neo4j.Driver
	session neo4j.Session
}

func NewDriver(write bool) (Driver, error) {
	driver, err := neo4j.NewDriver(endpoint, neo4j.BasicAuth("neo4j", "pwd123", ""))
	if err != nil {
		return Driver{}, err
	}

	var session neo4j.Session
	if write {
		session, err = driver.Session(neo4j.AccessModeWrite)
	} else {
		session, err = driver.Session(neo4j.AccessModeRead)
	}
	if err != nil {
		return Driver{}, err
	}

	return Driver{
		driver,
		session,
	}, err
}

func (d *Driver) NeighboursGen(id int) (neighboursGenResponse, error) {
	response, err := d.session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a)-[r:Gen]->(b)	WHERE id(a)=$id RETURN id(b), r.price, r.provider ORDER BY id(r)",
			map[string]interface{}{"id": id})

		if err != nil {
			return nil, err
		}

		resp, err := newNeighboursGenResponse(result)
		if err != nil {
			return nil, err
		}

		return resp, nil

	})

	if err != nil {
		return neighboursGenResponse{}, err
	}

	return response.(neighboursGenResponse), nil
}

func (d *Driver) close() {
	d.driver.Close()
	d.session.Close()
}
