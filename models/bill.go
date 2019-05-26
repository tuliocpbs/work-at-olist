package models

type BillRecord struct {
	Destination   int64  `json:"destination"`
	CallStartDate string `json:"call_start_date"`
	CallStartTime string `json:"call_start_time"`
	CallDuration  string `json:"call_duration"`
	CallPrice     string `json:"call_price"`
}
