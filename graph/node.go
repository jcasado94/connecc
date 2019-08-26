package graph

const (
	airportLabel = "Airport"
	cityLabel    = "City"
)

type node interface {
	Id() int
	Equals(n node) bool
}

type airport struct {
	id   int
	code string
}

func newAirport(id int, code string) airport {
	return airport{
		id:   id,
		code: code,
	}
}

func (a airport) Id() int {
	return a.id
}

func (a1 airport) Equals(n node) bool {
	a2 := n.(airport)
	return a1.id == a2.id && a1.code == a2.code
}

type city struct {
	id   int
	name string
}

func newCity(id int, name string) city {
	return city{
		id:   id,
		name: name,
	}
}

func (c city) Id() int {
	return c.id
}

func (c1 city) Equals(n node) bool {
	c2 := n.(city)
	return c1.id == c2.id && c1.name == c2.name
}

func newNode(label string, id int, params map[string]interface{}) node {
	if label == airportLabel {
		return newAirport(id, params["code"].(string))
	}
	return newCity(id, params["name"].(string))
}
