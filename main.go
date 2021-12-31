package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Nishankdawar/covid-app/handlers"
	"github.com/Nishankdawar/covid-app/utils"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/swaggo/echo-swagger/example/docs"

	"github.com/labstack/echo/v4"
)

// @title Swagger Covid Stats API
// @version 1.0
// @description This is a sample Covid App server.

// @BasePath /api/v1
func main() {

	port := os.Getenv("MY_APP_PORT")
	build_mode := os.Getenv("BUILD_MODE")
	if port == "" {
		port = "8080"
	}
	if build_mode == "" {
		build_mode = "PROD"
	}

	message_string := fmt.Sprintf("PORT: %s and BUILD_MODE: %s in main function", port, build_mode)
	utils.Logger("INFO", message_string, "main.go", time.Now().UTC())

	e := echo.New()

	utils.Logger("INFO", "Listening to port: "+port, "main.go", time.Now())

	var server_url string
	if build_mode == "PROD" {
		server_url = fmt.Sprintf(":" + port)
	} else {
		server_url = fmt.Sprintf("localhost:" + port)
	}

	// Routes
	base_api := e.Group("/api/v1")
	{
		base_api.POST("/populate_data", handlers.PopulateData)
		base_api.GET("/covid_stats", handlers.CovidStats)
	}
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	utils.Logger("INFO", "Starting server url: "+server_url, "main.go", time.Now())
	e.Logger.Fatal(e.Start(server_url))
}
