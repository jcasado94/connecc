package drivers

import (
	"github.com/jcasado94/connecc/mongo"
	mongoEntity "github.com/jcasado94/connecc/mongo/entity"
	mongoService "github.com/jcasado94/connecc/mongo/service"
)

const (
	mongoEndpoint    = "localhost:27017"
	mongoDb          = "tripz"
	mongoAvgPriceCol = "averagePrice"
)

type MongoDriver struct {
	session   *mongo.Session
	apService *mongoService.AveragePriceService
}

func NewMongoDriver() (MongoDriver, error) {
	session, err := mongo.NewSession(mongoEndpoint)
	if err != nil {
		return MongoDriver{}, err
	}
	return MongoDriver{
		session:   session,
		apService: mongoService.NewAveragePriceService(session, mongoDb, mongoAvgPriceCol),
	}, nil
}

func (md *MongoDriver) GetAvgPrice(s, t int) (float64, error) {
	price, err := md.apService.GetAverage(s, t)
	_, missingDocument := err.(mongoService.AvgDocumentNotFoundError)
	_, missingEntry := err.(mongoService.AvgNotFoundError)
	if missingDocument {
		return md.createAvgPriceDocument(s, t)
	} else if missingEntry {
		return md.createAvgPriceEntry(s, t)
	} else if err != nil {
		return 0.0, err
	}
	return price, nil
}

func (md *MongoDriver) createAvgPriceDocument(s, t int) (price float64, err error) {
	item, price := mongoEntity.NewAveragePrice(s, t)
	return price, md.apService.CreateAveragePrice(&item)
}

func (md *MongoDriver) createAvgPriceEntry(s, t int) (price float64, err error) {
	return md.apService.AddAverage(s, t)
}
