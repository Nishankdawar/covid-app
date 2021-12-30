package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"covid-app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TOKEN string
var REVERSE_GEOCODING_ENDPOINT string
var MONGO_ENDPOINT string
var COVID_DATA_ENDPOINT string
var DB_NAME string
var COLLECTION_NAME string
var TIME_FORMAT string
var MY_APP_PORT string

var (
	db              *mongo.Database
	covidcollection *mongo.Collection
	ctx             context.Context
)

func init() {
	TOKEN = "ae8e0bf69ec3ee45d71cf075ef867112"
	REVERSE_GEOCODING_ENDPOINT = "http://api.positionstack.com/v1/reverse?access_key=%[1]v&query=%[2]v,%[3]v"
	MONGO_ENDPOINT = "mongodb+srv://nishank:beta!1234@cluster0.1mfkt.mongodb.net/covid-data?retryWrites=true&w=majority"
	COVID_DATA_ENDPOINT = "https://data.covid19india.org/data.json"
	DB_NAME = "covid-data"
	COLLECTION_NAME = "CovidData"
	TIME_FORMAT = "02/01/2006 15:04:05"
	MY_APP_PORT = "8080"

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_ENDPOINT))

	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(DB_NAME)
	covidcollection = db.Collection(COLLECTION_NAME)
}

func convert_string_to_integer(str string) int64 {

	cases, cases_err := strconv.Atoi(str)
	if cases_err != nil {
		log.Fatal(cases_err)
	}
	return int64(cases)
}

func make_data(state_code string, state string, active_cases int64, confirmed_cases int64, last_updated_time int64) models.CovidData {
	var data = models.CovidData{
		State:           strings.ToLower(state),
		StateCode:       strings.ToLower(state_code),
		ActiveCases:     active_cases,
		ConfirmedCases:  confirmed_cases,
		LastUpdatedTime: last_updated_time,
	}
	return data
}

func get_user_region_code(latitude string, longitude string) string {

	reverse_geocoding := fmt.Sprintf(REVERSE_GEOCODING_ENDPOINT, TOKEN, latitude, longitude)
	resp, err := http.Get(reverse_geocoding)

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Some problem is there", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var jsonResponseReverseGeoCoding models.JSONResponseReverseGeoCoding
	err1 := json.Unmarshal([]byte(body), &jsonResponseReverseGeoCoding)

	if err1 != nil {
		log.Fatal(err1)
	}

	return strings.ToLower(jsonResponseReverseGeoCoding.Data[0].RegionCode)
}

func get_covid_data() models.JSONResponse {
	resp, err := http.Get(COVID_DATA_ENDPOINT)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var jsonResponse models.JSONResponse
	err1 := json.Unmarshal([]byte(body), &jsonResponse)

	if err1 != nil {
		log.Fatal(err1)
	}
	return jsonResponse
}

func PopulateData(c echo.Context) error {
	var jsonResponse = get_covid_data()

	result, err := covidcollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)

	for _, element := range jsonResponse.StateWise {

		active_cases := convert_string_to_integer(element.ActiveCases)
		confirmed_cases := convert_string_to_integer(element.ConfirmedCases)
		last_updated_time, err := time.Parse(TIME_FORMAT, element.LastUpdatedTime)
		if err != nil {
			log.Fatal(err)
		}

		var data = make_data(element.StateCode, element.State, active_cases, confirmed_cases, last_updated_time.Unix())

		res, err := covidcollection.InsertOne(context.Background(), data)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(res.InsertedID.(primitive.ObjectID).Timestamp())
	}

	return c.JSONPretty(http.StatusCreated, bson.M{"status": "created"}, " ")
}

func CovidStats(c echo.Context) error {
	latitude, longitude := c.FormValue("lat"), c.FormValue("long")
	region_code := get_user_region_code(latitude, longitude)
	state_codes := []string{region_code, "tt"}
	filter := bson.M{"statecode": bson.M{"$in": state_codes}}
	filterCursor, error_1 := covidcollection.Find(context.Background(), filter)
	if error_1 != nil {
		log.Fatal(error_1)
	}
	var resultsFiltered []bson.M
	for filterCursor.Next(ctx) {
		var document bson.M
		err := filterCursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		resultsFiltered = append(resultsFiltered, document)
	}

	var country_cases int64
	var state_cases int64
	var state_name string
	var lastupdatedtime int64

	for _, ele := range resultsFiltered {
		if ele["statecode"] == "tt" {
			country_cases = ele["confirmed"].(int64)
		}
		if ele["statecode"] == "dl" {
			state_cases = ele["confirmed"].(int64)
			state_name = ele["state"].(string)
		}
		lastupdatedtime = ele["lastupdatedtime"].(int64)
	}

	stats := models.UserStats{
		CountryCases:    country_cases,
		StateCases:      state_cases,
		StateName:       state_name,
		LastUpdatedTime: time.Unix(lastupdatedtime, 0).String(),
	}
	return c.JSONPretty(http.StatusOK, stats, " ")
}

func main() {

	port := os.Getenv("MY_APP_PORT")
	if port == "" {
		port = "8080"
	}
	e := echo.New()

	e.POST("/populate_data", PopulateData)
	e.GET("/covid_stats", CovidStats)

	e.Logger.Print(fmt.Sprintf("Listening to port %s", port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
