package service

import "fmt"

type AvgNotFoundError struct {
	What string
}

func newAvgNotFoundError(s, t int) AvgNotFoundError {
	return AvgNotFoundError{
		What: fmt.Sprintf("Couldn't find avg for %v in %v", t, s),
	}
}

func (e AvgNotFoundError) Error() string {
	return e.What
}

type AvgDocumentNotFoundError struct {
	What string
}

func newAvgDocumentNotFoundError(s int) AvgDocumentNotFoundError {
	return AvgDocumentNotFoundError{
		What: fmt.Sprintf("No document found in averagePrice collection for %v", s),
	}
}

func (e AvgDocumentNotFoundError) Error() string {
	return e.What
}
