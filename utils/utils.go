package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Nishankdawar/covid-app/literals"
	"github.com/Nishankdawar/covid-app/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Logger(level string, message string, file string, time time.Time) {
	log_string := fmt.Sprintf("level:%[1]v, message: %[2]v, file: %[3]v, time:%[4]v", level, message, file, time)
	log.Println(log_string)
}

func ErrorLogger(err error) {
	var error_string string = fmt.Sprintf("%s", err)
	Logger("ERROR", error_string, "utils.go", time.Now())
}

func GetMongoClient() *mongo.Client {
	Logger("INFO", "Creating mongo client in GetMongoClient:", "utils.go", time.Now())
	client, err := mongo.NewClient(options.Client().ApplyURI(literals.MONGO_ENDPOINT))
	if err != nil {
		ErrorLogger(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	Logger("INFO", "Setting up client connection in GetMongoClient:", "utils.go", time.Now())
	err = client.Connect(ctx)
	if err != nil {
		ErrorLogger(err)
	}
	return client
}

func GetMongoCollection(client *mongo.Client) *mongo.Collection {
	db := client.Database(literals.DB_NAME)
	covidcollection := db.Collection(literals.COLLECTION_NAME)
	Logger("INFO", "Getting collection in GetMongoCollection:", "utils.go", time.Now())
	return covidcollection
}

func ConvertStringToInteger(str string) int64 {

	params_string := fmt.Sprintf("Inside ConvertStringToInteger with params %s", str)
	Logger("INFO", params_string, "utils.go", time.Now())
	cases, err := strconv.Atoi(str)
	if err != nil {
		ErrorLogger(err)
	}
	return int64(cases)
}

func MakeData(state_code string, state string, active_cases int64, confirmed_cases int64, last_updated_time int64) models.CovidData {
	params_string := fmt.Sprintf("Inside MakeData with params %s, %s, %v, %v, %v", state_code, state, active_cases, confirmed_cases, last_updated_time)
	Logger("INFO", params_string, "utils.go", time.Now())
	var data = models.CovidData{
		State:           strings.ToLower(state),
		StateCode:       strings.ToLower(state_code),
		ActiveCases:     active_cases,
		ConfirmedCases:  confirmed_cases,
		LastUpdatedTime: last_updated_time,
	}
	return data
}

func ReadIOResponseBody(response *http.Response) []byte {
	Logger("INFO", "Inside ReadIOResponseBody", "utils.go", time.Now())
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ErrorLogger(err)
	}
	return body
}
