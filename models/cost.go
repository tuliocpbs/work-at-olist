package models

import "time"

type Cost struct {
	StandingCharge float32   `json:"standing_charge"`
	MinuteCharge   float32   `json:"minute_charge"`
	DateTime       time.Time `json:"timestamp"`
}
