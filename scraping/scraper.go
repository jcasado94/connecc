package scraping

type Scraper interface {
	GetTrips(year, month, day, passengers int) ([]*Trip, error)
}
