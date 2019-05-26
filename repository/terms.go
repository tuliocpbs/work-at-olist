package repository

import (
	"time"

	"github.com/olivere/elastic"
)

func TermQueryLE(t time.Time) *elastic.BoolQuery {
	rangeQ := elastic.NewRangeQuery("timestamp").
		Lte(t)
	query := elastic.NewBoolQuery().
		Filter(rangeQ)

	return query
}

func TermQueryBill(sub int, start time.Time, end time.Time) *elastic.BoolQuery {
	rangeQ := elastic.NewRangeQuery("timestamp_end").
		Gte(start).
		Lte(end)
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("source", sub)).
		Filter(rangeQ)

	return query
}
