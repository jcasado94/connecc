package scraping

import (
	"fmt"
	"time"
)

type Trip struct {
	Price float64
	Legs  []*Leg
}

func (t *Trip) String() string {
	return fmt.Sprintf("Trip{Price:%f, Legs:%v}", t.Price, t.Legs)
}

type Leg struct {
	Dep, Arr string
	DepTime  time.Time
	ArrTime  time.Time
	Id       string
}

func (l *Leg) String() string {
	return fmt.Sprintf("Leg{Dep:%s, Arr: %s, DepTime:%v, ArrTime:%v, Id:%s}", l.Dep, l.Arr, l.DepTime, l.ArrTime, l.Id)
}

func newLeg(dep, arr, id string, depTime, arrTime time.Time) *Leg {
	return &Leg{
		Dep:     dep,
		Arr:     arr,
		DepTime: depTime,
		ArrTime: arrTime,
		Id:      id,
	}
}
