package controllers

import (
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
)

func ConfigRouter() *gin.Engine {
	router := gin.New()

	v1 := router.Group("/v1")
	{
		v1.Use(apmgin.Middleware(router))
		// Health-check to test application and APM Server
		v1.GET("/health-check/", HealthCheck)
	}

	return router
}

func HealthCheck(c *gin.Context) {
	c.SecureJSON(200, "OK")
}
