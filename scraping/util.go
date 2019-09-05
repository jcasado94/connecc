package scraping

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func processDayDifference(baseTime *time.Time, hour, min int) time.Time {
	fixedTime := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), hour, min, 0, 0, time.UTC)
	if hour < baseTime.Hour() {
		fixedTime = time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()+1, hour, min, 0, 0, time.UTC)
	}
	return fixedTime
}

type singularMockRoundTripper struct {
	mockUrl     string
	contentType string
}

func newSingularMockRoundTripper(filePath, contentType string) *singularMockRoundTripper {
	return &singularMockRoundTripper{
		mockUrl:     filePath,
		contentType: contentType,
	}
}

func (rt *singularMockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    r,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-type", rt.contentType)

	dat, err := ioutil.ReadFile(rt.mockUrl)
	if err != nil {
		return nil, err
	}
	response.Body = ioutil.NopCloser(bytes.NewReader(dat))

	return response, nil
}

type multipleMockRoundTripper struct {
	urlToFilePath    map[string]string
	urlToContentType map[string]string
}

func newMultipleMockRoundTripper(urlToFilePath, urlToContentType map[string]string) *multipleMockRoundTripper {
	return &multipleMockRoundTripper{
		urlToFilePath:    urlToFilePath,
		urlToContentType: urlToContentType,
	}
}

func (rt *multipleMockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    r,
		StatusCode: http.StatusOK,
	}
	url := r.URL.String()
	path := rt.urlToFilePath[url]
	contentType := rt.urlToContentType[url]

	response.Header.Set("Content-type", contentType)

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	response.Body = ioutil.NopCloser(bytes.NewReader(dat))

	return response, nil
}
