package db

import (
	neo4j "github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	neighboursBelongsToCityQuery = "MATCH (a)-[r:BelongsTo]-(b:City) WHERE id(a)=$id RETURN labels(b)[0], id(b), properties(b) " +
		"UNION MATCH (a:City)-[r:BelongsTo]-(b) WHERE id(a)=$id AND id(a)=$s RETURN labels(b)[0], id(b), properties(b)"
	neighboursBelongsToThroughCityQuery = "MATCH (a)-[r1:BelongsTo]->(b:City)-[r2:BelongsTo]-(c) WHERE id(a)=$id RETURN labels(c)[0], id(c), properties(c) ORDER BY id(r2)"
	neighboursGenQuery                  = "MATCH (a)-[r:Gen]->(b)	WHERE id(a)=$id RETURN r.price, r.provider, labels(b)[0], id(b), properties(b) ORDER BY id(r)"

	nodeInfoQuery = "MATCH (n) WHERE id(n)=$id RETURN labels(n)[0], properties(n)"
)

type Driver struct {
	driver  neo4j.Driver
	session neo4j.Session
}

func NewDriver(dbEndpoint, dbUsername, dbPw string, write bool) (Driver, error) {
	driver, err := neo4j.NewDriver(dbEndpoint, neo4j.BasicAuth(dbUsername, dbPw, ""))
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

func (d *Driver) NodeInfo(id int) (neo4j.Result, error) {
	response, err := d.session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			nodeInfoQuery,
			map[string]interface{}{"id": id})

		if err != nil {
			return nil, err
		}

		return result, nil
	})

	if err != nil {
		return nil, err
	}

	return response.(neo4j.Result), nil
}

func (d *Driver) NeighboursGen(id int) (neo4j.Result, error) {
	response, err := d.session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			neighboursGenQuery,
			map[string]interface{}{"id": id})

		if err != nil {
			return nil, err
		}

		return result, nil

	})

	if err != nil {
		return nil, err
	}

	return response.(neo4j.Result), nil
}

func (d *Driver) neighboursBelongsToCity(id, s int) (neo4j.Result, error) {
	response, err := d.session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			neighboursBelongsToCityQuery,
			map[string]interface{}{"id": id, "s": s})

		if err != nil {
			return -1, err
		}
		return result, nil
	})

	if err != nil {
		return nil, err
	}

	return response.(neo4j.Result), nil
}

func (d *Driver) neighboursBelongsToThroughCity(id int) (neo4j.Result, error) {
	response, err := d.session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			neighboursBelongsToThroughCityQuery,
			map[string]interface{}{"id": id})

		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return nil, err
	}

	return response.(neo4j.Result), nil
}

func (d *Driver) NeighboursBelongsToCity(id, s int) (neo4j.Result, error) {
	resultCity, err := d.neighboursBelongsToCity(id, s)
	if err != nil {
		return nil, err
	}
	return resultCity, nil
}

func (d *Driver) NeighboursBelongsToThroughCity(id, s int) (neo4j.Result, error) {
	resultThroughCity, err := d.neighboursBelongsToThroughCity(id)
	if err != nil {
		return nil, err
	}
	return resultThroughCity, nil
}

func (d *Driver) close() {
	d.driver.Close()
	d.session.Close()
}
