package literals

import (
	"fmt"
	"os"
)

var USERNAME = os.Getenv("USERNAME")
var PASSWORD = os.Getenv("PASSWORD")
var REVERSE_GEOCODING_TOKEN = "ae8e0bf69ec3ee45d71cf075ef867112"
var REVERSE_GEOCODING_ENDPOINT = "http://api.positionstack.com/v1/reverse?access_key=%[1]v&query=%[2]v,%[3]v"
var MONGO_ENDPOINT_PROD = fmt.Sprintf("mongodb+srv://%s:%s@cluster0.1mfkt.mongodb.net/covid-data?retryWrites=true&w=majority", USERNAME, PASSWORD)
var MONGO_ENDPOINT_LOCAL = "mongodb://localhost:27017"
var COVID_DATA_ENDPOINT = "https://data.covid19india.org/data.json"
var DB_NAME = "covid-data"
var COLLECTION_NAME = "CovidData"
var TIME_FORMAT = "02/01/2006 15:04:05"
var MY_APP_PORT = "8080"
var INDIA_CODE_FILTER = "tt"
var INDIA_CODE = "ind"
var STATUS = "status"
var CREATED = "created"
