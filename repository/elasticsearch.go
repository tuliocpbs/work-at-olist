package repository

import (
	"context"
	"log"

	"github.com/olivere/elastic"
)

type ElasticSearch struct {
	client *elastic.Client
}

func (es ElasticSearch) insert(d interface{}, id string, index string) {
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

func (es ElasticSearch) update(d interface{}, id string, index string) {
	update, err := es.client.Update().
		Id(id).
		Index(index).
		Type("item").
		Doc(d).
		Do(context.Background())

	if err != nil {
		log.Printf("Data not updated, Error: %s\n", err)
	} else {
		log.Printf("Updated data %s to index %s, type %s\n", update.Id, update.Index, update.Type)
	}

}

func (es ElasticSearch) search(query elastic.Query, maxSize int, index string) (*elastic.SearchResult, error) {
	var err error

	searchResult, err := es.client.Search().
		Index(index).
		Type("item").
		Query(query).
		Size(maxSize). // Max quantity value return
		Do(context.Background())

	if err != nil {
		log.Printf("Search | Error: %s\n", err)
		return nil, err
	}

	if n := searchResult.TotalHits(); n >= 0 {
		log.Printf("Found a total of %d item\n", n)
	}

	return searchResult, err
}

func (es ElasticSearch) getByID(id string, index string) (*elastic.GetResult, error) {
	get1, err := es.client.Get().
		Index(index).
		Type("item").
		Id(id).
		Do(context.Background())

	if err != nil {
		log.Printf("Data not found, Error: %s\n", err)
		return nil, err
	}

	if get1.Found {
		log.Println("Item found")
	} else {
		log.Println("Item not found")
	}

	return get1, nil
}
