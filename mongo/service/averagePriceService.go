package service

import (
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

func (aps *AveragePriceService) CreateAveragePrice(ap *entity.AveragePrice) error {
	apm := model.NewAveragePriceModel(ap)
	return aps.collection.Insert(&apm)
}

func (aps *AveragePriceService) GetAverage(s, tInt int) (float64, error) {
	t := strconv.Itoa(tInt)
	query := map[string]int{"nodeId": s}
	var ap model.AveragePriceModel
	err := aps.collection.Find(query).One(&ap)
	if err != nil {
		return 0.0, newAvgDocumentNotFoundError(s)
	}
	if _, exists := ap.Averages[t]; !exists {
		return 0.0, newAvgNotFoundError(s, tInt)
	}
	return ap.Averages[t].Avg, nil
}

func (aps *AveragePriceService) AddAverage(s, tInt int) (price float64, err error) {
	t := strconv.Itoa(tInt)
	query := map[string]int{"nodeId": s}
	var ap model.AveragePriceModel
	err = aps.collection.Find(query).One(&ap)
	if err != nil {
		return 0.0, err
	}
	ap.AddAverage(t, 0.0, 0)
	err = aps.collection.Update(query, ap)
	return 0.0, err
}
