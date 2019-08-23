package mongo_test

import (
	"log"
	"strconv"
	"testing"

	"github.com/jcasado94/connecc/mongo"
	"github.com/jcasado94/connecc/mongo/entity"
	"github.com/jcasado94/connecc/mongo/service"
)

const (
	mongoUrl       = "localhost:27017"
	dbName         = "testDb"
	collectionName = "testCol"
)

func TestServices(t *testing.T) {
	t.Run("AveragePriceService", AveragePriceService)
}

func AveragePriceService(t *testing.T) {
	t.Run("GetAverage", getAverage_should_get_avg_from_mongo)
}

func getAverage_should_get_avg_from_mongo(t *testing.T) {
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer finishTest(session)
	apService := service.NewAveragePriceService(session.Copy(), dbName, collectionName)

	//populate
	testId, testNodeIdS, testNodeIdT, testAvg, testN := "1111", 3, 2, 4.5, 3
	averagePrice := entity.AveragePrice{
		ID:     testId,
		NodeId: testNodeIdS,
		Averages: map[string]entity.Average{
			strconv.Itoa(testNodeIdT): entity.Average{
				Avg: testAvg,
				N:   testN,
			},
		},
	}

	err = apService.CreateAveragePriceService(&averagePrice)
	if err != nil {
		t.Errorf("Unable to create averagePrice: %s", err)
	}

	//test
	avg, err := apService.GetAverage(testNodeIdS, testNodeIdT)
	if err != nil {
		t.Error(err)
	}
	if avg != testAvg {
		t.Errorf("Unable to get avg from %v to %v", testNodeIdS, testNodeIdT)
	}
}

func connect() *mongo.Session {
	session, err := mongo.NewSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo %s", err)
	}
	return session
}

func finishTest(s *mongo.Session) {
	s.DropDatabase(dbName)
	s.Close()
}
