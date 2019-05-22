package main

import (
	"work-at-olist/controllers"
	repo "work-at-olist/repository"
)

func main() {
	repo.CreateClient()

	router := controllers.ConfigRouter()

	router.Run(":8080")
}
