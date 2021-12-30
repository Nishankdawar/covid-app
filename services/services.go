package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Nishankdawar/covid-app/literals"
	"github.com/Nishankdawar/covid-app/models"
	"github.com/Nishankdawar/covid-app/utils"
)

func GetUserRegionCode(latitude string, longitude string) string {

	utils.Logger("INFO", "Inside GetUserRegionCode:", "services.go", time.Now())
	reverse_geocoding := fmt.Sprintf(literals.REVERSE_GEOCODING_ENDPOINT, literals.TOKEN, latitude, longitude)
	utils.Logger("INFO", "Reverse Geocoding URL: "+reverse_geocoding, "services.go", time.Now())
	response, err := http.Get(reverse_geocoding)

	if err != nil {
		utils.ErrorLogger(err)
	}

	if response.StatusCode != 200 {
		utils.Logger("WARNING", "Reverse geocoding status code is not 200", "services.go", time.Now())
	}

	body := utils.ReadIOResponseBody(response)

	var jsonResponseReverseGeoCoding models.JSONResponseReverseGeoCoding
	err1 := json.Unmarshal([]byte(body), &jsonResponseReverseGeoCoding)

	if err1 != nil {
		utils.ErrorLogger(err1)
	}

	return strings.ToLower(jsonResponseReverseGeoCoding.Data[0].RegionCode)
}

func GetCovidData() models.JSONResponse {
	utils.Logger("INFO", "Inside GetCovidData:", "services.go", time.Now())
	response, err := http.Get(literals.COVID_DATA_ENDPOINT)

	if err != nil {
		utils.ErrorLogger(err)
	}

	body := utils.ReadIOResponseBody(response)

	var jsonResponse models.JSONResponse
	err1 := json.Unmarshal([]byte(body), &jsonResponse)

	if err1 != nil {
		utils.ErrorLogger(err1)
	}
	return jsonResponse
}
