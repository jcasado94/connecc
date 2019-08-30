package main

import (
	"fmt"

	"github.com/jcasado94/connecc/scraping"
)

func main() {
	scraper := scraping.NewSpiritScraper()
	trips, _ := scraper.GetTrips("BOS", "DEN", 31, 8, 2019, 1, 0, 0)
	for _, trip := range trips {
		fmt.Println(trip)
	}
	// scraper.GetTrips("BOS", "DEN", 25, 9, 2019, 1, 0, 0)
}
