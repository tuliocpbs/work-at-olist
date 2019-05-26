package controllers

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestVerifyType(t *testing.T) {
	var err error
	var listTestAcc []string
	var listTestNot []string

	listTestAcc = append(listTestAcc, "start", "end")
	for _, s := range listTestAcc {
		if err = verifyType(&s); err != nil {
			t.Error(err)
		}
	}

	listTestNot = append(listTestNot, "startx", "End", "fstart", "end-l")
	for _, l := range listTestNot {
		if err = verifyType(&l); err == nil {
			t.Error("Invalid type accept")
		}
	}
}

func TestVerifyTelephone(t *testing.T) {
	var err error
	var listTestAcc []int64
	var listTestNot []int64

	listTestAcc = append(listTestAcc, 89999234568, 8999923456)
	for _, s := range listTestAcc {
		if err = verifyTelephone(s); err != nil {
			t.Error(err)
		}
	}

	listTestNot = append(listTestNot, 99923456, 999234568, 899992345682)
	for _, l := range listTestNot {
		if err = verifyTelephone(l); err == nil {
			t.Error("Invalid telephone accept")
		}
	}
}

func TestGetDDD(t *testing.T) {
	var listTestAcc []int64

	listTestAcc = append(listTestAcc, 89999234568, 8999923456, 123)
	for _, s := range listTestAcc {
		var ddd int
		ddd, _ = getDDD(s)

		if ddd <= 9 || ddd >= 100 {
			t.Error("Invalid DDD generate")
		}
	}
}

func TestCountCallMinutes(t *testing.T) {
	start, _ := time.Parse("2006-01-02T15:04:05Z", "2019-05-25T15:00:37Z")
	end, _ := time.Parse("2006-01-02T15:04:05Z", "2019-05-25T16:00:00Z")
	minutes := countCallMinutes(start, end)
	assert.Equal(t, minutes, 59)

	end, _ = time.Parse("2006-01-02T15:04:05Z", "2019-05-25T16:00:40Z")
	minutes = countCallMinutes(start, end)
	assert.Equal(t, minutes, 60)

	end, _ = time.Parse("2006-01-02T15:04:05Z", "2019-05-25T22:00:00Z")
	minutes = countCallMinutes(start, end)
	assert.Equal(t, minutes, 419)

	end, _ = time.Parse("2006-01-02T15:04:05Z", "2019-05-25T23:00:00Z")
	minutes = countCallMinutes(start, end)
	assert.Equal(t, minutes, 419)

	end, _ = time.Parse("2006-01-02T15:04:05Z", "2019-05-26T06:00:00Z")
	minutes = countCallMinutes(start, end)
	assert.Equal(t, minutes, 419)

	end, _ = time.Parse("2006-01-02T15:04:05Z", "2019-05-26T15:00:00Z")
	minutes = countCallMinutes(start, end)
	assert.Equal(t, minutes, 959)
}
