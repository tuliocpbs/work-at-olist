package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"work-at-olist/models"
	repo "work-at-olist/repository"
	"work-at-olist/utils"

	"github.com/gin-gonic/gin"
)

func configBillRouter(router *gin.RouterGroup) {
	router.GET("/bill", addBill)
}

func addBill(c *gin.Context) {
	if _, has := c.GetQuery("subscriber"); !has {
		utils.WriteJSON(c, http.StatusForbidden, errors.New("Subscriber not informed"))
		return
	}
	sub, _ := strconv.Atoi(c.Query("subscriber"))
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("month"))

	var billsRecords []models.BillRecord
	status, err := getBill(&billsRecords, sub, month, year)
	if err != nil {
		utils.WriteJSON(c, http.StatusForbidden, err)
		return
	}

	utils.WriteJSON(c, status, billsRecords)
}

func getBill(b *[]models.BillRecord, sub int, month int, year int) (int, error) {
	start, end, err := getPeriodBoundaries(month, year)
	if err != nil {
		return http.StatusForbidden, err
	}

	query := repo.TermQueryBill(sub, start, end)

	var cdsList interface{}
	cdsList, err = repo.SearchList(query, &models.CallDetailsStart{})
	if err != nil {
		return http.StatusForbidden, err
	}

	for _, item := range cdsList.([]models.CallDetailsStart) {
		*b = append(*b, setNewBillRecord(item))
	}

	return 200, nil
}

func getPeriodBoundaries(m int, y int) (time.Time, time.Time, error) {
	if m > 0 {
		return getMonthBoundaries(m)
	} else if y > 0 {
		return getYearBoundaries(y)
	}

	return time.Time{}, time.Time{}, errors.New("Year or Month not informed")
}

func getMonthBoundaries(m int) (time.Time, time.Time, error) {
	var err error
	var start time.Time
	var end time.Time
	now := time.Now()

	if time.Month(m) == now.Month() {
		err = errors.New("This telephone bill period has not ended")
	} else if time.Month(m) > now.Month() {
		start = time.Date(now.Year()-1, time.Month(m), 1, 0, 0, 0, 0, time.UTC) // Month of previously year
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	} else {
		start = time.Date(now.Year(), time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	}

	return start, end, err
}

func getYearBoundaries(y int) (time.Time, time.Time, error) {
	var err error
	var start time.Time
	var end time.Time
	now := time.Now()

	if y >= now.Year() {
		err = errors.New("This telephone bill period has not ended")
	} else {
		start = time.Date(now.Year()-1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	}

	return start, end, err
}

func setNewBillRecord(c models.CallDetailsStart) (b models.BillRecord) {
	b.Destination = c.Destination
	b.CallStartDate = c.TimeStampStart.Format("02/01/2006")
	b.CallStartTime = c.TimeStampStart.Format("3:04 PM")
	d := c.TimeStampEnd.Sub(c.TimeStampStart)
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	b.CallDuration = fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	b.CallPrice = fmt.Sprintf("R$ %.2f", c.Cost)

	return
}
