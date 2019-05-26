package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"work-at-olist/models"
	repo "work-at-olist/repository"
	"work-at-olist/utils"
)

func configReceiveCallsRouter(router *gin.RouterGroup) {
	router.POST("/calldetails", addCallDetails)
}

func addCallDetails(c *gin.Context) {
	recordEntry := &models.RecordEntry{}
	if err := utils.ReadJSON(c, recordEntry); err != nil {
		utils.WriteJSON(c, http.StatusForbidden, err.Error())
		return
	}
	status, msg := saveRecordDetails(recordEntry)
	utils.WriteJSON(c, status, msg)
}

func saveRecordDetails(r *models.RecordEntry) (int, error) {
	var err error
	err = saveRecord(r)
	if err != nil {
		return http.StatusForbidden, err
	}

	if strings.Compare(r.Type, "start") == 0 {
		err = saveStartCall(r)
	} else if strings.Compare(r.Type, "end") == 0 {
		err = saveEndCall(r)
	}

	if err != nil {
		return http.StatusForbidden, err
	}

	return http.StatusOK, err
}

func saveRecord(r *models.RecordEntry) (err error) {
	record := &models.Record{}
	record.ID = r.ID
	record.CallID = r.CallID

	if err = verifyType(&r.Type); err != nil {
		return
	}
	record.Type = r.Type

	repo.Save(record)
	return
}

func saveStartCall(r *models.RecordEntry) (err error) {
	callStart := &models.CallDetailsStart{}
	callStart.ID = r.CallID
	callStart.TimeStampStart, err = time.Parse("2006-01-02T15:04:05Z", r.TimeStamp)
	if err != nil {
		return
	}
	if err = verifyTelephone(r.Source); err != nil {
		return
	}
	callStart.Source = r.Source
	callStart.SourceAreaCode, err = getDDD(r.Source)
	if err != nil {
		return
	}
	if err = verifyTelephone(r.Destination); err != nil {
		return
	}
	callStart.Destination = r.Destination
	callStart.DestinationAreaCode, err = getDDD(r.Destination)
	if err != nil {
		return
	}

	repo.Save(callStart)
	return
}

func saveEndCall(r *models.RecordEntry) (err error) {
	callEnd := &models.CallDetailsEnd{}
	callEnd.ID = r.CallID
	callEnd.TimeStampEnd, err = time.Parse("2006-01-02T15:04:05Z", r.TimeStamp)
	if err != nil {
		return
	}
	callEnd.Cost, err = callCost(callEnd.ID, callEnd.TimeStampEnd)

	repo.Save(callEnd)
	return
}

func callCost(callID int, endCallTime time.Time) (float32, error) {

	var cs models.CallDetailsStart
	err := repo.Get(callID, &cs)
	if err != nil {
		return 0, err
	}

	query := repo.TermQueryLE(endCallTime)

	var cost models.Cost
	err = repo.Search(query, &cost)
	if err != nil {
		return 0, err
	}

	callCharge, err := calculateCallCost(endCallTime, cs.TimeStampStart, cost.StandingCharge, cost.MinuteCharge)
	if err != nil {
		return 0, err
	}

	return callCharge, nil
}

func calculateCallCost(ect time.Time, cst time.Time, sc float32, mc float32) (float32, error) {

	if ect.Sub(cst).Minutes() >= 0 {
		minutesCall := countCallMinutes(cst, ect)
		return sc + mc*float32(minutesCall), nil
	}

	return 0, errors.New("Error: EndCallDate smaller than StartCallDate")
}

func countCallMinutes(cst time.Time, ect time.Time) int {
	var durationCall float64
	callTime := cst.Add(time.Minute * 1)
	earlyLimit := time.Date(callTime.Year(), callTime.Month(), callTime.Day(), 6, 0, 0, 0, time.UTC)
	lateLimit := time.Date(callTime.Year(), callTime.Month(), callTime.Day(), 22, 0, 0, 0, time.UTC)

	for ect.Sub(callTime) > 0 {
		if callTime.Sub(earlyLimit).Minutes() > 0.0 && callTime.Sub(lateLimit).Minutes() < 0.0 { // Interval to considerer charge
			durationCall += 1.0
		}

		callTime = callTime.Add(time.Minute * 1)
		if callTime.Day() != earlyLimit.Day() { // Verify if the day change after add 1 minute
			earlyLimit = earlyLimit.AddDate(0, 0, 1)
			lateLimit = lateLimit.AddDate(0, 0, 1)
		}
	}
	callTime = callTime.Add(time.Minute * -1)
	durationCall += ect.Sub(callTime).Minutes()

	return int(durationCall)
}

func verifyType(t *string) error {
	if strings.Compare(*t, "start") == 0 || strings.Compare(*t, "end") == 0 {
		return nil
	}
	return errors.New("Type not valid")
}

func verifyTelephone(t int64) error {
	telephone := strconv.FormatInt(t, 10)
	if len(telephone) >= 10 && len(telephone) <= 11 {
		return nil
	}
	return errors.New("Telephone not valid")
}

func getDDD(t int64) (int, error) {
	telephone := strconv.FormatInt(t, 10)
	ddd, err := strconv.Atoi(telephone[0:2])

	return ddd, err
}
