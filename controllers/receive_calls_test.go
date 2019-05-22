package controllers

import (
	"testing"
)

func TestverifyType(t *testing.T) {
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

func TestverifyTelephone(t *testing.T) {
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

func TestgetDDD(t *testing.T) {
	var listTestAcc []int64

	listTestAcc = append(listTestAcc, 89999234568, 8999923456, 123)
	for _, s := range listTestAcc {
		var ddd int
		ddd, _ = getDDD(s)

		if ddd > 9 && ddd < 100 {
			t.Error("Invalid DDD generate")
		}
	}
}
