package scraping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MegabusScraper struct {
	client http.Client
}

func newMegabusScraper() *MegabusScraper {
	return &MegabusScraper{
		client: http.Client{},
	}
}

func (sc *MegabusScraper) GetTrips(departure, arrival string, day, month, year, adults, children, infants int) ([]*Trip, error) {
	trips := make([]*Trip, 0)
	var err error

	url := fmt.Sprintf("https://us.megabus.com/journey-planner/journeys?days=1&concessionCount=0&departureDate=%d-%d-%d&destinationId=%s&inboundOtherDisabilityCount=0&inboundPcaCount=0&inboundWheelchairSeated=0&nusCount=0&originId=%s&otherDisabilityCount=0&pcaCount=0&totalPassengers=%d&wheelchairSeated=0",
		year, month, day, arrival, departure, adults+children+infants)

	resp, err := sc.client.Get(url)
	if err != nil {
		return []*Trip{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*Trip{}, err
	}
	document := string(body)

	r, err := regexp.Compile(`window.SEARCH_RESULTS\s?=\s?(?P<Json>{.*})`)
	if err != nil {
		return []*Trip{}, err
	}
	journies := r.FindStringSubmatch(document)
	journiesBytes := []byte(journies[1])
	var js JsonMbJournies
	err = json.Unmarshal(journiesBytes, &js)

	for _, j := range js.Journeys {
		if len(j.Legs) == 1 {
			depTime, err := time.Parse(time.RFC3339, j.Legs[0].DepartureDateTime)
			if err != nil {
				return []*Trip{}, err
			}
			arrTime, err := time.Parse(time.RFC3339, j.Legs[0].ArrivalDateTime)
			if err != nil {
				return []*Trip{}, err
			}
			trips = append(trips, newTrip(
				[]*Fare{newFare("standard", j.Price)},
				[]*Leg{newLeg(j.Legs[0].Origin.CityId, j.Legs[0].Destination.CityId, "", depTime, arrTime)}))
		} else {
			// send query to get mid cities
			url = fmt.Sprintf("https://us.megabus.com/journey-planner/api/itinerary?journeyId=%s", j.JourneyId)
			resp, err = sc.client.Get(url)
			if err != nil {
				log.Printf("Couldn't retrieve itinerary for journeyId:%s. [%s --> %s], %d/%d/%d", j.JourneyId, departure, arrival, day, month, year)
				continue
			}
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Couldn't retrieve itinerary for journeyId:%s. [%s --> %s], %d/%d/%d", j.JourneyId, departure, arrival, day, month, year)
				continue
			}
			var its JsonMbItineraries
			err = json.Unmarshal(body, &its)
			if err != nil {
				log.Printf("Couldn't retrieve itinerary for journeyId:%s. [%s --> %s], %d/%d/%d", j.JourneyId, departure, arrival, day, month, year)
				continue
			}
			var dep, arr string
			var depTime, arrTime time.Time
			var legs []*Leg
			for i, it := range its.ScheduledStops {
				if it.Ordinal == 0 {
					if dep != "" {
						arrStop := its.ScheduledStops[i-1]
						arr = arrStop.CityId
						arrivalTimeSlice := strings.Split(arrStop.ArrivalTime, ":")
						arrivalHourString, arrivalMinString := arrivalTimeSlice[0], arrivalTimeSlice[1]
						arrivalHour, err := strconv.Atoi(arrivalHourString)
						arrivalMin, err := strconv.Atoi(arrivalMinString)
						if err != nil {
							log.Printf("Couldn't retrieve itinerary for journeyId:%s. [%s --> %s], %d/%d/%d", j.JourneyId, departure, arrival, day, month, year)
							continue
						}
						arrTime = processDayDifference(&depTime, arrivalHour, arrivalMin)
						legs = append(legs, newLeg(dep, arr, "", depTime, arrTime))
					}
					depTimeSlice := strings.Split(it.DepartureTime, ":")
					depHourString, depMinString := depTimeSlice[0], depTimeSlice[1]
					depHour, err := strconv.Atoi(depHourString)
					depMin, err := strconv.Atoi(depMinString)
					if err != nil {
						log.Printf("Couldn't retrieve itinerary for journeyId:%s. [%s --> %s], %d/%d/%d", j.JourneyId, departure, arrival, day, month, year)
						continue
					}
					if arrTime.IsZero() {
						depTime = time.Date(year, time.Month(month), day, depHour, depMin, 0, 0, time.UTC)
					} else {
						depTime = processDayDifference(&arrTime, depHour, depMin)
					}
					dep = it.CityId
				} else if i == len(its.ScheduledStops)-1 {
					arrivalTimeSlice := strings.Split(it.ArrivalTime, ":")
					arrivalHourString, arrivalMinString := arrivalTimeSlice[0], arrivalTimeSlice[1]
					arrivalHour, err := strconv.Atoi(arrivalHourString)
					arrivalMin, err := strconv.Atoi(arrivalMinString)
					if err != nil {
						log.Printf("Couldn't retrieve itinerary for journeyId:%s. [%s --> %s], %d/%d/%d", j.JourneyId, departure, arrival, day, month, year)
						continue
					}
					arrTime = processDayDifference(&depTime, arrivalHour, arrivalMin)
					arr = it.CityId
					legs = append(legs, newLeg(dep, arr, "", depTime, arrTime))
				}
			}
			trips = append(trips, newTrip(
				[]*Fare{newFare("standard", j.Price)},
				legs,
			))
		}
	}

	return trips, nil

}

type JsonMbJournies struct {
	Journeys []JsonMbJourney
}

type JsonMbJourney struct {
	JourneyId         string
	DepartureDateTime string
	ArrivalDateTime   string
	Price             float64
	Legs              []JsonMbLeg
}

type JsonMbLeg struct {
	DepartureDateTime string
	ArrivalDateTime   string
	Origin            JsonMbStop
	Destination       JsonMbStop
}

type JsonMbStop struct {
	CityId string
}

type JsonMbItineraries struct {
	ScheduledStops []JsonMbItineraryStop
}

type JsonMbItineraryStop struct {
	CityId        string
	Ordinal       int64
	DepartureTime string
	ArrivalTime   string
}
