package main

import (
	"work-at-olist/controllers"
)

func main() {
	router := controllers.ConfigRouter()

	router.Run(":8080")
}
