package repository

import (
	"encoding/json"
	"errors"
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

func Get(ID int, typ interface{}) (err error) {
	var index string

	switch reflect.TypeOf(typ) {
	case reflect.TypeOf(&models.CallDetailsStart{}):
		index = "call-details"
	default:
		log.Println("Data Type not found")
		return errors.New("GET: Data Type not found")
	}

	id := strconv.Itoa(ID)
	test, err := DB.getByID(id, index)
	if err != nil {
		return
	}

	bytes, err := test.Source.MarshalJSON()
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &typ)
	if err != nil {
		return
	}

	return nil
}

func Search(query elastic.Query, typ interface{}) (err error) {
	switch reflect.TypeOf(typ) {
	case reflect.TypeOf(&models.Cost{}):
		sr, err := DB.search(query, 1, "hist-cost")
		if err != nil {
			return err
		}

		for _, item := range sr.Each(reflect.TypeOf(typ).Elem()) {
			*typ.(*models.Cost) = item.(models.Cost)
		}

		return nil
	default:
		log.Println("Data not found, Error: Model not found")
		return errors.New("Data not found, Error: Model not found")
	}
}

func SearchList(query elastic.Query, typ interface{}) (interface{}, error) {
	switch reflect.TypeOf(typ) {
	case reflect.TypeOf(&models.CallDetailsStart{}):
		sr, err := DB.search(query, 10000, "call-details") // Set max of 10000 items to return
		if err != nil {
			return nil, err
		}

		var cdsList []models.CallDetailsStart
		var cds models.CallDetailsStart
		for _, item := range sr.Each(reflect.TypeOf(cds)) {
			if t, ok := item.(models.CallDetailsStart); ok {
				cdsList = append(cdsList, t)
			}
		}

		return cdsList, nil
	default:
		log.Println("Data not found, Error: Model not found")
		return nil, errors.New("Data not found, Error: Model not found")
	}
}

func Save(data interface{}) {
	switch reflect.TypeOf(data) {
	case reflect.TypeOf(&models.Record{}):
		id := strconv.Itoa(data.(*models.Record).ID)
		DB.insert(data, id, "record")
	case reflect.TypeOf(&models.CallDetailsStart{}):
		id := strconv.Itoa(data.(*models.CallDetailsStart).ID)
		DB.insert(data, id, "call-details")
	case reflect.TypeOf(&models.CallDetailsEnd{}):
		id := strconv.Itoa(data.(*models.CallDetailsEnd).ID)
		DB.update(data, id, "call-details")
	default:
		log.Println("Data not inserted, Error: Model not found")
	}
}
