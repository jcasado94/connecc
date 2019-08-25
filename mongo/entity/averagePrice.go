package entity

import "strconv"

type AveragePrice struct {
	ID       string             `json:"id"`
	NodeId   int                `json:"nodeId"`
	Averages map[string]Average `json:"averages"`
}

type Average struct {
	Avg float64 `json:"avg"`
	N   int     `json:"n"`
}

func NewAveragePrice(s, t int) (avgPrice AveragePrice, price float64) {
	averages := map[string]Average{
		strconv.Itoa(t): Average{
			Avg: 0.0,
			N:   0,
		},
	}
	return AveragePrice{
		NodeId:   s,
		Averages: averages,
	}, 0.0
}

type AveragePriceService interface {
	GetAverage(s, t int) (float64, error)
}
