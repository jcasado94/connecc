package scraping

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

var location, _ = time.LoadLocation("America/New_York")

func processDayDifference(baseTime *time.Time, hour, min int) time.Time {
	fixedTime := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), hour, min, 0, 0, location)
	if hour < baseTime.Hour() {
		fixedTime = time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()+1, hour, min, 0, 0, location)
	}
	return fixedTime
}

type mockRoundTripper struct {
	mockUrl string
}

func newMockRoundTripper(filePath string) *mockRoundTripper {
	return &mockRoundTripper{
		mockUrl: filePath,
	}
}

func (rt *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    r,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-type", "text/html; charset=utf-8")

	dat, err := ioutil.ReadFile(rt.mockUrl)
	if err != nil {
		return nil, err
	}
	response.Body = ioutil.NopCloser(bytes.NewReader(dat))

	return response, nil
}
