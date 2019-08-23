package entity

type AveragePrice struct {
	ID       string             `json:"id"`
	NodeId   int                `json:"nodeId"`
	Averages map[string]Average `json:"averages"`
}

type Average struct {
	Avg float64 `json:"avg"`
	N   int     `json:"n"`
}

type AveragePriceService interface {
	GetAverage(s, t int) (float64, error)
}
