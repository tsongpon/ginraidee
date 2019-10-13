package adapter

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/tsongpon/ginraidee/model"
	"github.com/tsongpon/ginraidee/transport"
	"log"
	"os"
)

var googleGeoCodeAPIKey = os.Getenv("GOOGLE_API_KEY")

type GoogleGeoCodeAdapter struct {
}

func NewGoogleGeoCodeAdapter() *GoogleGeoCodeAdapter {
	return new(GoogleGeoCodeAdapter)
}

func (g *GoogleGeoCodeAdapter) GetLocation(address string) (*model.Location, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"address": address,
			"key":     googleGeoCodeAPIKey,
		}).
		SetHeader("Accept", "application/json").
		Get("https://maps.googleapis.com/maps/api/geocode/json")

	if err != nil {
		log.Printf("get error %s", err.Error())
	}
	result := transport.GoogleGeoCodeTransport{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Printf("Unmarshal result error %s", err.Error())
	}
	log.Printf("result %v", result)
	if len(result.Results) == 0 {
		return nil, errors.New("location not found")
	}
	location := model.Location{}
	location.Name = result.Results[0].Name
	location.Lat = result.Results[0].Geometry.Location.Lat
	location.Lng = result.Results[0].Geometry.Location.Lng

	return &location, nil
}
