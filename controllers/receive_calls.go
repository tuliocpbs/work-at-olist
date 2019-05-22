package controllers

import (
	"github.com/gin-gonic/gin"

	"work-at-olist/utils"
)

func configReceiveCallsRouter(router *gin.RouterGroup) {
	router.POST("/calldetails", addCallDetails)
}

func addCallDetails(c *gin.Context) {
	utils.WriteJSON(c, 200, "Ok")
}
