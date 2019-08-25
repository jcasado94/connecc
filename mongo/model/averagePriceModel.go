package model

import (
	"github.com/jcasado94/connecc/mongo/entity"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AveragePriceModel struct {
	ID       bson.ObjectId      `bson:"_id,omitempty"`
	NodeId   int                `bson:"nodeId"`
	Averages map[string]Average `bson:"averages"`
}

type Average struct {
	Avg float64 `bson:"avg"`
	N   int     `bson:"n"`
}

func NewAveragePriceModel(ap *entity.AveragePrice) *AveragePriceModel {
	averages := make(map[string]Average)
	for key, value := range ap.Averages {
		averages[key] = Average{
			Avg: value.Avg,
			N:   value.N,
		}
	}
	return &AveragePriceModel{
		NodeId:   ap.NodeId,
		Averages: averages,
	}
}

func (apm *AveragePriceModel) AddAverage(t string, price float64, n int) {
	avg := Average{
		Avg: price,
		N:   n,
	}
	apm.Averages[t] = avg
}

func AveragePriceModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"ID"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}
