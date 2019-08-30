package scraping

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

type SpiritScraper struct {
	browser *browser.Browser
}

func NewSpiritScraper() *SpiritScraper {
	browser := surf.NewBrowser()
	browser.SetUserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36")
	browser.Open("https://www.spirit.com/Default.aspx")
	return &SpiritScraper{
		browser: browser,
	}
}

func (sc *SpiritScraper) GetTrips(departure, arrival string, day, month, year, adults, children, infants int) ([]*Trip, error) {
	trips := make([]*Trip, 0)
	var err error

	sc.browser.Post("https://www.spirit.com/Default.aspx?action=search", "application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("bypassHC=False&birthdates=&lapoption=&awardFSNumber=&bookingType=F&hotelOnlyInput=&autoCompleteValueHidden=&carPickUpTime=16&carDropOffTime=16&tripType=oneWay&vacationPackageType=on&from=%s&to=%s&departDate=%d%%2F%d%%2F%d&departDateDisplay=08%%2F31%%2F2019&returnDate=09%%2F03%%2F2019&returnDateDisplay=09%%2F03%%2F2019&ADT=%d&CHD=%d&INF=%d&promoCode=&fromMultiCity1=&toMultiCity1=&dateMultiCity1=&dateMultiCityDisplay1=&fromMultiCity2=&toMultiCity2=&dateMultiCity2=&dateMultiCityDisplay2=&fromMultiCity3=&toMultiCity3=&dateMultiCity3=&dateMultiCityDisplay3=&fromMultiCity4=&toMultiCity4=&dateMultiCity4=&dateMultiCityDisplay4=&redeemMiles=false",
			departure, arrival,
			month, day, year,
			adults, children, infants)))

	sc.browser.Open("https://www.spirit.com/DPPCalendarMarket.aspx")

	sc.browser.Dom().Find(".rowsMarket1").Each(func(_ int, s *goquery.Selection) {
		trip := Trip{}
		trip.Price, err = strconv.ParseFloat(s.Find(".valueToSortBy").First().Text(), 64)
		if err != nil {
			return
		}

		sFlightNumbers := s.Find(".popUpContent .fi-header-text.text-uppercase.text-right")
		var arrTime, depTime time.Time
		s.Find(".flight-info-body").Each(func(i int, s *goquery.Selection) {
			fieldsDates := s.Find(".fi-text-bold")
			depHour, depMin, err := processTime(fieldsDates.Get(1).FirstChild.Data)
			if err != nil {
				return
			}
			if arrTime.IsZero() {
				depTime = time.Date(year, time.Month(month), day, depHour, depMin, 0, 0, location)
			} else {
				depTime = processDayDifference(&arrTime, depHour, depMin)
			}
			arrHour, arrMin, err := processTime(fieldsDates.Get(3).FirstChild.Data)
			if err != nil {
				return
			}
			arrTime = processDayDifference(&depTime, arrHour, arrMin)

			fieldsLocations := s.Find(".fi-text")
			dep := fieldsLocations.Get(0).FirstChild.Data
			arr := fieldsLocations.Get(1).FirstChild.Data

			flightNumberSlice := strings.Split(sFlightNumbers.Get(i).FirstChild.Data, " ")
			flightNumber := fmt.Sprintf("NK%s", flightNumberSlice[len(flightNumberSlice)-1])

			trip.Legs = append(trip.Legs, newLeg(dep, arr, flightNumber, depTime, arrTime))
		})

		trips = append(trips, &trip)
	})

	if err != nil {
		return []*Trip{}, err
	}

	return trips, nil
}

func processTime(time string) (depHour, depMin int, err error) {
	depTimeSlice := strings.Split(time, " ")
	depTimeSlice = strings.Split(depTimeSlice[0], ":")
	depHour, err = strconv.Atoi(depTimeSlice[0])
	depMin, err = strconv.Atoi(depTimeSlice[1])
	if err != nil {
		return 0, 0, err
	}
	if depTimeSlice[1] == "PM" {
		depHour += 12
	}
	return depHour, depMin, nil
}
