package scraping

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTripsMegabus(t *testing.T) {
	sc := newMegabusScraper()
	sc.client.Transport = newMultipleMockRoundTripper(urlToFilePath(), urlToContentType())
	expectedTrips := []Trip{
		Trip{
			Fares: []Fare{Fare{Price: 99.0, Type: "standard"}},
			Legs: []Leg{Leg{Dep: "123", Arr: "142", DepTime: time.Date(2019, time.Month(9), 8, 2, 0, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 8, 7, 30, 0, 0, time.UTC)},
				Leg{Dep: "142", Arr: "289", DepTime: time.Date(2019, time.Month(9), 8, 10, 5, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 9, 00, 25, 0, 0, time.UTC)},
			},
		},
		Trip{
			Fares: []Fare{Fare{Price: 99.0, Type: "standard"}},
			Legs: []Leg{Leg{Dep: "123", Arr: "142", DepTime: time.Date(2019, time.Month(9), 8, 8, 0, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 8, 12, 15, 0, 0, time.UTC)},
				Leg{Dep: "142", Arr: "289", DepTime: time.Date(2019, time.Month(9), 8, 15, 30, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 9, 7, 25, 0, 0, time.UTC)},
			},
		},
		Trip{
			Fares: []Fare{Fare{Price: 99.0, Type: "standard"}},
			Legs: []Leg{Leg{Dep: "123", Arr: "142", DepTime: time.Date(2019, time.Month(9), 8, 9, 0, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 8, 13, 40, 0, 0, time.UTC)},
				Leg{Dep: "142", Arr: "289", DepTime: time.Date(2019, time.Month(9), 8, 15, 30, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 9, 7, 25, 0, 0, time.UTC)},
			},
		},
		Trip{
			Fares: []Fare{Fare{Price: 99.0, Type: "standard"}},
			Legs: []Leg{Leg{Dep: "123", Arr: "142", DepTime: time.Date(2019, time.Month(9), 8, 16, 0, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 8, 20, 15, 0, 0, time.UTC)},
				Leg{Dep: "142", Arr: "289", DepTime: time.Date(2019, time.Month(9), 8, 23, 20, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 9, 13, 45, 0, 0, time.UTC)},
			},
		},
		Trip{
			Fares: []Fare{Fare{Price: 99.0, Type: "standard"}},
			Legs: []Leg{Leg{Dep: "123", Arr: "142", DepTime: time.Date(2019, time.Month(9), 8, 17, 0, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 8, 21, 15, 0, 0, time.UTC)},
				Leg{Dep: "142", Arr: "289", DepTime: time.Date(2019, time.Month(9), 8, 23, 20, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 9, 13, 45, 0, 0, time.UTC)},
			},
		},
		Trip{
			Fares: []Fare{Fare{Price: 99.0, Type: "standard"}},
			Legs: []Leg{Leg{Dep: "123", Arr: "142", DepTime: time.Date(2019, time.Month(9), 8, 23, 0, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 9, 4, 0, 0, 0, time.UTC)},
				Leg{Dep: "142", Arr: "289", DepTime: time.Date(2019, time.Month(9), 9, 6, 5, 0, 0, time.UTC), ArrTime: time.Date(2019, time.Month(9), 9, 21, 35, 0, 0, time.UTC)},
			},
		},
	}
	trips, err := sc.GetTrips("123", "289", 8, 9, 2019, 1, 0, 0)
	if err != nil {
		t.Errorf("Couldn't retrieve trips.\n%v", err)
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
				if !l.Equals(&el) {
					t.Errorf("Legs slices differ. Want \n%v, \ngot \n%v", want.Legs, have.Legs)
				}
			}
			for j, ef := range want.Fares {
				f := have.Fares[j]
				if f != ef {
					t.Errorf("Fares slices differ. Want \n%v, \ngot \n%v", want.Fares, have.Fares)
				}
			}
		})
	}
}

func urlToFilePath() map[string]string {
	return map[string]string{
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1413647":                                                                                                                                                                                                                    "./testScrapingSites/megabusItinerary0.json",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1410032":                                                                                                                                                                                                                    "./testScrapingSites/megabusItinerary1.json",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1407252":                                                                                                                                                                                                                    "./testScrapingSites/megabusItinerary2.json",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1417628":                                                                                                                                                                                                                    "./testScrapingSites/megabusItinerary3.json",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1419463":                                                                                                                                                                                                                    "./testScrapingSites/megabusItinerary4.json",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1403592":                                                                                                                                                                                                                    "./testScrapingSites/megabusItinerary5.json",
		"https://us.megabus.com/journey-planner/journeys?days=1&concessionCount=0&departureDate=2019-9-8&destinationId=289&inboundOtherDisabilityCount=0&inboundPcaCount=0&inboundWheelchairSeated=0&nusCount=0&originId=123&otherDisabilityCount=0&pcaCount=0&totalPassengers=1&wheelchairSeated=0": "./testScrapingSites/megabusTrips.html",
	}
}

func urlToContentType() map[string]string {
	return map[string]string{
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1413647":                                                                                                                                                                                                                    "application/json; charset=utf-8",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1410032":                                                                                                                                                                                                                    "application/json; charset=utf-8",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1407252":                                                                                                                                                                                                                    "application/json; charset=utf-8",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1403592":                                                                                                                                                                                                                    "application/json; charset=utf-8",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1417628":                                                                                                                                                                                                                    "application/json; charset=utf-8",
		"https://us.megabus.com/journey-planner/api/itinerary?journeyId=*1419463":                                                                                                                                                                                                                    "application/json; charset=utf-8",
		"https://us.megabus.com/journey-planner/journeys?days=1&concessionCount=0&departureDate=2019-9-8&destinationId=289&inboundOtherDisabilityCount=0&inboundPcaCount=0&inboundWheelchairSeated=0&nusCount=0&originId=123&otherDisabilityCount=0&pcaCount=0&totalPassengers=1&wheelchairSeated=0": "text/html; charset=utf-8",
	}
}
