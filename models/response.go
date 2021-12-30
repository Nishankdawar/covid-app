package models

type Address struct {
	Region      string `json:"region"`
	RegionCode  string `json:"region_code"`
	CountryCode string `json:"country_code"`
}

type StateWiseArray struct {
	State           string `json:"state"`
	StateCode       string `json:"statecode"`
	ActiveCases     string `json:"active"`
	ConfirmedCases  string `json:"confirmed"`
	LastUpdatedTime string `json:"lastupdatedtime"`
}

type JSONResponse struct {
	StateWise []StateWiseArray `json:"statewise"`
}

type JSONResponseReverseGeoCoding struct {
	Data []Address `json:"data"`
}

type UserStatsResponse struct {
	CountryCases    int64  `json:"countrycases"`
	StateCases      int64  `json:"statecases"`
	StateName       string `json:"statename"`
	LastUpdatedTime string `json:"lastupdatedtime"`
}
