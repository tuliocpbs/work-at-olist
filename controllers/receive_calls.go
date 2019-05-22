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
	callEnd.Cost = calculateCallCost()

	repo.Save(callEnd)
	return
}

func calculateCallCost() float32 {
	return 0.0
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
