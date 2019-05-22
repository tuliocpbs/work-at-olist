package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"work-at-olist/models"
)

func ReadJSON(c *gin.Context, v interface{}) error {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, v)
}

func WriteResponse(c *gin.Context, status int, data interface{}, message interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.JSON(status, models.Response{Message: message, Data: data})
}

func WriteJSON(c *gin.Context, status int, data interface{}) {
	switch data.(type) {
	case error:
		WriteResponse(c, status, nil, data.(error).Error())
	default:
		WriteResponse(c, status, data, http.StatusText(status))
	}
}
