package controllers

import (
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
)

func ConfigRouter() *gin.Engine {
	router := gin.New()

	router.GET("/health-check/", healthCheck)

	v1 := router.Group("/v1")
	{
		v1.Use(apmgin.Middleware(router))
		configReceiveCallsRouter(v1)
	}

	return router
}

func healthCheck(c *gin.Context) {
	c.SecureJSON(200, "OK")
}
