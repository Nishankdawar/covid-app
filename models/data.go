package models

type JSONResponse struct {
	StateWise []StateWiseArray `json:"statewise"`
}

type JSONResponseReverseGeoCoding struct {
	Data []Address `json:"data"`
}

type Address struct {
	Region     string `json:"region"`
	RegionCode string `json:"region_code"`
}

type StateWiseArray struct {
	State           string `json:"state"`
	StateCode       string `json:"statecode"`
	ActiveCases     string `json:"active"`
	ConfirmedCases  string `json:"confirmed"`
	LastUpdatedTime string `json:"lastupdatedtime"`
}

type CovidData struct {
	State           string `json:"state" bson:"state"`
	StateCode       string `json:"statecode" bson:"statecode"`
	ActiveCases     int64  `json:"active" bson:"active"`
	ConfirmedCases  int64  `json:"confirmed" bson:"confirmed"`
	LastUpdatedTime int64  `json:"lastupdatedtime" bson:"lastupdatedtime"`
}

type UserStats struct {
	CountryCases    int64  `json:"countrycases"`
	StateCases      int64  `json:"statecases"`
	StateName       string `json:"statename"`
	LastUpdatedTime string `json:"lastupdatedtime"`
}
