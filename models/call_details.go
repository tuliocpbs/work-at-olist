package models

import "time"

type Record struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	CallID      int    `json:"call_id"`
	CallDetails interface{}
}

type CallDetailsStart struct {
	ID                  int       `json:"id"`
	TimeStampStart      time.Time `json:"timestamp_start"`
	TimeStampEnd        time.Time `json:"timestamp_end"`
	SourceAreaCode      int       `json:"source_area_code"`
	Source              int64     `json:"source"`
	DestinationAreaCode int       `json:"destination_area_code"`
	Destination         int64     `json:"destination"`
	Cost                float32   `json:"cost"`
}

type CallDetailsEnd struct {
	ID           int       `json:"id"`
	TimeStampEnd time.Time `json:"timestamp_end"`
	Cost         float32   `json:"cost"`
}

type RecordEntry struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	TimeStamp   string `json:"timestamp"`
	CallID      int    `json:"call_id"`
	Source      int64  `json:"source"`
	Destination int64  `json:"destination"`
}
