package scraping

import (
	"fmt"
	"time"
)

type Trip struct {
	Fares []*Fare
	Legs  []*Leg
}

func newTrip(fares []*Fare, legs []*Leg) *Trip {
	return &Trip{
		Fares: fares,
		Legs:  legs,
	}
}

func (t *Trip) String() string {
	return fmt.Sprintf("Trip{Price:%v, Legs:%v}", t.Fares, t.Legs)
}

type Leg struct {
	Dep, Arr string
	DepTime  time.Time
	ArrTime  time.Time
	Id       string
}

func (l *Leg) Equals(l2 *Leg) bool {
	return l.Dep == l2.Dep && l.Arr == l2.Arr && l.DepTime.Equal(l2.DepTime) && l.ArrTime.Equal(l2.ArrTime) && l.Id == l2.Id
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

type Fare struct {
	Type  string
	Price float64
}

func newFare(t string, p float64) *Fare {
	return &Fare{
		Type:  t,
		Price: p,
	}
}

func (f *Fare) String() string {
	return fmt.Sprintf("%s: %f", f.Type, f.Price)
}
