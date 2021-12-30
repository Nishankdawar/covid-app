package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Nishankdawar/covid-app/handlers"
	"github.com/Nishankdawar/covid-app/utils"

	"github.com/labstack/echo/v4"
)

func main() {

	utils.Logger("INFO", "Before fetching PORT in main function", "main.go", time.Now().UTC())
	port := os.Getenv("MY_APP_PORT")
	build_mode := os.Getenv("BUILD_MODE")
	if port == "" {
		port = "8080"
	}
	if build_mode == "" {
		build_mode = "PROD"
	}
	utils.Logger("INFO", "After fetching PORT in main function", "main.go", time.Now().UTC())

	e := echo.New()
	utils.Logger("INFO", "Creating instance of echo server in main function", "main.go", time.Now().UTC())

	base_api := e.Group("/api/v1")

	base_api.POST("/populate_data", handlers.PopulateData)
	base_api.GET("/covid_stats", handlers.CovidStats)

	utils.Logger("INFO", "Listening to port: "+port, "main.go", time.Now())

	var server_url string
	if build_mode == "PROD" {
		server_url = fmt.Sprintf(":" + port)
	} else {
		server_url = fmt.Sprintf("localhost:" + port)
	}

	utils.Logger("INFO", "Starting server url: "+server_url, "main.go", time.Now())
	e.Logger.Fatal(e.Start(server_url))
}
