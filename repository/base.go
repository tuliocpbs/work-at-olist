package repository

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic"

	"work-at-olist/models"
)

type Repository interface {
	Insert()
}

var (
	DB ElasticSearch
)

func CreateClient() {
	url := os.Getenv("ELASTIC_URL")
	if len(url) == 0 {
		fmt.Println("Elastic url not informed.")
	}

	timeout := time.Duration(3 * time.Second)
	clientSetTimeout := &http.Client{
		Timeout: timeout,
	}

	startupTimeout := time.Duration(40 * time.Second)
	client, err := elastic.NewClient(
		elastic.SetHttpClient(clientSetTimeout),
		elastic.SetSniff(false),
		elastic.SetURL(url),
		elastic.SetMaxRetries(3),
		elastic.SetHealthcheckTimeoutStartup(startupTimeout), // Wait for elasticsearch be up by docker-compose
	)
	if err != nil {
		panic(err)
	}

	DB.client = client
}

func Save(data interface{}) {
	switch reflect.TypeOf(data) {
	case reflect.TypeOf(&models.Record{}):
		fmt.Println("Record")
		id := strconv.Itoa(data.(*models.Record).ID)
		DB.Insert(data, id, "record")
	case reflect.TypeOf(&models.CallDetailsStart{}):
		id := strconv.Itoa(data.(*models.CallDetailsStart).ID)
		DB.Insert(data, id, "call-details")
	case reflect.TypeOf(&models.CallDetailsEnd{}):
		id := strconv.Itoa(data.(*models.CallDetailsEnd).ID)
		DB.Insert(data, id, "call-details")
	default:
		log.Println("Data not inserted, Error: Model not found")
	}
}
