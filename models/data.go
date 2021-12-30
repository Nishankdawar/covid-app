package models

type CovidData struct {
	State           string `json:"state" bson:"state"`
	StateCode       string `json:"statecode" bson:"statecode"`
	ActiveCases     int64  `json:"active" bson:"active"`
	ConfirmedCases  int64  `json:"confirmed" bson:"confirmed"`
	LastUpdatedTime int64  `json:"lastupdatedtime" bson:"lastupdatedtime"`
}
