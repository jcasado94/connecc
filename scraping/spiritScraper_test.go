package scraping

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTripsSpirit(t *testing.T) {
	sc := NewSpiritScraper()
	sc.browser.SetTransport(newSingularMockRoundTripper("./testScrapingSites/spiritAirlines.html", "text/html; charset=utf-8"))
	trips, err := sc.GetTrips("BOS", "DEN", 13, 9, 2019, 1, 0, 0)
	if err != nil {
		t.Error("Error while getting the trips")
		return
	}

	expectedTrips := []Trip{
		Trip{
			Fares: []*Fare{&Fare{Price: 158.18, Type: "standard"}},
			Legs: []*Leg{
				&Leg{Dep: "Boston, MA", Arr: "Baltimore, MD / Washington, DC AREA", DepTime: time.Date(2019, time.Month(9), 13, 7, 45, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 13, 9, 24, 0, 0, time.UTC), Id: "NK2025"},
				&Leg{Dep: "Baltimore, MD / Washington, DC AREA", Arr: "Minneapolis/St. Paul, MN", DepTime: time.Date(2019, time.Month(9), 13, 11, 55, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 13, 13, 32, 0, 0, time.UTC), Id: "NK381"},
				&Leg{Dep: "Minneapolis/St. Paul, MN", Arr: "Denver, CO", DepTime: time.Date(2019, time.Month(9), 13, 14, 37, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 13, 15, 48, 0, 0, time.UTC), Id: "NK381"},
			},
		},
		Trip{
			Fares: []*Fare{&Fare{Price: 153.98, Type: "standard"}},
			Legs: []*Leg{
				&Leg{Dep: "Boston, MA", Arr: "Baltimore, MD / Washington, DC AREA", DepTime: time.Date(2019, time.Month(9), 13, 7, 45, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 13, 9, 24, 0, 0, time.UTC), Id: "NK2025"},
				&Leg{Dep: "Baltimore, MD / Washington, DC AREA", Arr: "Denver, CO", DepTime: time.Date(2019, time.Month(9), 13, 20, 19, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 13, 22, 17, 0, 0, time.UTC), Id: "NK115"},
			},
		},
		Trip{
			Fares: []*Fare{&Fare{Price: 122.08, Type: "9Dollar"}, &Fare{Price: 171.98, Type: "standard"}},
			Legs: []*Leg{
				&Leg{Dep: "Boston, MA", Arr: "Fort Lauderdale, FL / Miami, FL AREA", DepTime: time.Date(2019, time.Month(9), 13, 15, 35, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 13, 19, 04, 0, 0, time.UTC), Id: "NK1611"},
				&Leg{Dep: "Fort Lauderdale, FL / Miami, FL AREA", Arr: "Denver, CO", DepTime: time.Date(2019, time.Month(9), 13, 21, 45, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 13, 23, 59, 0, 0, time.UTC), Id: "NK355"},
			},
		},
	}
	if len(expectedTrips) != len(trips) {
		t.Errorf("Trip slices lengths differ. Want \n%v, \ngot %v", expectedTrips, trips)
	}
	for i, want := range expectedTrips {
		have := trips[i]
		t.Run(fmt.Sprintf("Trip %d", i), func(t *testing.T) {
			if len(want.Legs) != len(have.Legs) {
				t.Errorf("Legs slices differ. Want \n%v, \ngot \n%v", want.Legs, have.Legs)
				return
			}
			if len(want.Fares) != len(have.Fares) {
				t.Errorf("Fares slices differ. Want \n%v, \ngot \n%v", want.Fares, have.Fares)
				return
			}
			for j, el := range want.Legs {
				l := have.Legs[j]
				if !l.Equals(el) {
					t.Errorf("Legs slices differ. Want \n%v, \ngot \n%v", want.Legs, have.Legs)
				}
			}
			for j, ef := range want.Fares {
				f := have.Fares[j]
				if *f != *ef {
					t.Errorf("Fares slices differ. Want \n%v, \ngot \n%v", want.Fares, have.Fares)
				}
			}
		})
	}
}
