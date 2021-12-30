package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Nishankdawar/covid-app/literals"
	"github.com/Nishankdawar/covid-app/models"
	"github.com/Nishankdawar/covid-app/services"
	"github.com/Nishankdawar/covid-app/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func logRequest(request *http.Request) {
	uri := request.RequestURI
	method := request.Method
	log.Println("Received request from : ", method, uri)
}

func PopulateData(c echo.Context) error {
	utils.Logger("INFO", "Inside PopulateData function", "handlers.go", time.Now())
	logRequest(c.Request())

	mongo_client := utils.GetMongoClient()
	mongo_collection := utils.GetMongoCollection(mongo_client)

	var jsonResponse = services.GetCovidData()

	result, err := mongo_collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		utils.ErrorLogger(err)
	}
	deleted_string := fmt.Sprintf("DeleteMany removed %v document(s)\n", result.DeletedCount)
	utils.Logger("INFO", deleted_string, "handlers.go", time.Now())

	utils.Logger("INFO", "Inserting elements in db", "handlers.go", time.Now())
	for _, element := range jsonResponse.StateWise {

		active_cases := utils.ConvertStringToInteger(element.ActiveCases)
		confirmed_cases := utils.ConvertStringToInteger(element.ConfirmedCases)
		last_updated_time, err := time.Parse(literals.TIME_FORMAT, element.LastUpdatedTime)
		if err != nil {
			utils.ErrorLogger(err)
		}

		var data = utils.MakeData(element.StateCode, element.State, active_cases, confirmed_cases, last_updated_time.Unix())

		res, err := mongo_collection.InsertOne(context.Background(), data)
		if err != nil {
			utils.ErrorLogger(err)
		}
		log.Println(res.InsertedID.(primitive.ObjectID).Timestamp())
		utils.Logger("INFO", "Inserting elements in db completed", "handlers.go", time.Now())
	}

	return c.JSONPretty(http.StatusCreated, bson.M{literals.STATUS: literals.CREATED}, " ")
}

func CovidStats(c echo.Context) error {
	logRequest(c.Request())
	utils.Logger("INFO", "Inside CovidStats function", "handlers.go", time.Now())
	latitude, longitude := c.FormValue("lat"), c.FormValue("long")
	region_code := services.GetUserRegionCode(latitude, longitude)
	state_codes := []string{region_code, literals.INDIA_CODE}
	filter := bson.M{"statecode": bson.M{"$in": state_codes}}

	mongo_client := utils.GetMongoClient()
	mongo_collection := utils.GetMongoCollection(mongo_client)

	utils.Logger("INFO", "Finding inside collection using filter", "handlers.go", time.Now())
	filterCursor, error_1 := mongo_collection.Find(context.Background(), filter)
	if error_1 != nil {
		utils.ErrorLogger(error_1)
	}
	var resultsFiltered []bson.M
	for filterCursor.Next(context.Background()) {
		var document bson.M
		err := filterCursor.Decode(&document)
		if err != nil {
			utils.ErrorLogger(err)
		}
		resultsFiltered = append(resultsFiltered, document)
	}

	var country_cases int64
	var state_cases int64
	var state_name string
	var lastupdatedtime int64

	for _, ele := range resultsFiltered {
		if ele["statecode"] == literals.INDIA_CODE {
			country_cases = ele["confirmed"].(int64)
		}
		if ele["statecode"] == region_code {
			state_cases = ele["confirmed"].(int64)
			state_name = ele["state"].(string)
		}
		lastupdatedtime = ele["lastupdatedtime"].(int64)
	}

	stats := models.UserStatsResponse{
		CountryCases:    country_cases,
		StateCases:      state_cases,
		StateName:       state_name,
		LastUpdatedTime: time.Unix(lastupdatedtime, 0).String(),
	}
	return c.JSONPretty(http.StatusOK, stats, " ")
}
