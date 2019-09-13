package pathprocessor

import "github.com/jcasado94/connecc/scraping"

// this should be a way more sophisticated logic e,e
const kpaths = 100

type PathProcessor interface {
	Process() (trips []scraping.Trip)
}
