package repository

import (
	"context"
	"log"

	"github.com/olivere/elastic"
)

type ElasticSearch struct {
	client *elastic.Client
}

func (es ElasticSearch) Insert(d interface{}, id string, index string) {
	put, err := es.client.Index().
		Id(id).
		Index(index).
		Type("item").
		BodyJson(d).
		Do(context.Background())

	if err != nil {
		log.Printf("Data not inserted, Error: %s\n", err)
	} else {
		log.Printf("Indexed data %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	}
}
