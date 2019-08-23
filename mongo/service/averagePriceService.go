package service

import (
	"fmt"
	"strconv"

	"github.com/jcasado94/connecc/mongo"
	"github.com/jcasado94/connecc/mongo/entity"
	"github.com/jcasado94/connecc/mongo/model"
	mgo "gopkg.in/mgo.v2"
)

type AveragePriceService struct {
	collection *mgo.Collection
}

func NewAveragePriceService(session *mongo.Session, dbName, colName string) *AveragePriceService {
	collection := session.GetCollection(dbName, colName)
	collection.EnsureIndex(model.AveragePriceModelIndex())
	return &AveragePriceService{collection}
}

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

func (aps *AveragePriceService) CreateAveragePriceService(ap *entity.AveragePrice) error {
	apm := model.NewAveragePriceModel(ap)
	return aps.collection.Insert(&apm)
}

func (aps *AveragePriceService) GetAverage(s, tInt int) (float64, error) {
	t := strconv.Itoa(tInt)
	query := map[string]int{"nodeId": s}
	var ap model.AveragePriceModel
	err := aps.collection.Find(query).One(&ap)
	if err != nil {
		return 0.0, err
	}
	if _, exists := ap.Averages[t]; !exists {
		return 0.0, newAvgNotFoundError(s, tInt)
	}
	return ap.Averages[t].Avg, nil
}
