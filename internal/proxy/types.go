package proxy

import (
	"time"
)

type ProxyResponse struct {
	Count    int           `json:"count"`
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []ProxyDetail `json:"results"`
}

type ProxyDetail struct {
	ID                    string    `json:"id"`
	Username              string    `json:"username"`
	Password              string    `json:"password"`
	ProxyAddress          string    `json:"proxy_address"`
	Port                  int       `json:"port"`
	Valid                 bool      `json:"valid"`
	LastVerification      time.Time `json:"last_verification"`
	CountryCode           string    `json:"country_code"`
	CityName              string    `json:"city_name"`
	AsnName               string    `json:"asn_name"`
	AsnNumber             int       `json:"asn_number"`
	HighCountryConfidence bool      `json:"high_country_confidence"`
	CreatedAt             time.Time `json:"created_at"`
}
