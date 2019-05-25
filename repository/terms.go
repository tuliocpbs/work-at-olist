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
