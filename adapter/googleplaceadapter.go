package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tsongpon/ginraidee/model"
	"github.com/tsongpon/ginraidee/transport"
	"log"
)

type GooglePlaceAdapter struct {
}

func NewGooglePlaceAdapter() *GooglePlaceAdapter {
	return new(GooglePlaceAdapter)
}

func (a *GooglePlaceAdapter) GetPlaces(placeType string, lat float32, lon float32) []model.Place {
	client := resty.New()

	location := fmt.Sprintf("%f", lat) + "," + fmt.Sprintf("%f", lon)
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"location": location,
			"type":     placeType,
			"radius":   "1000",
			"key":      "AIzaSyBfGD0y888DZ8FUfpBjDCRVRhFimnG0z78",
		}).
		SetHeader("Accept", "application/json").
		Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json")

	if err != nil {
		log.Printf("get error %s", err.Error())
	}
	result := transport.GooglePlaceTransport{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Printf("Unmarshal result error %s", err.Error())
	}
	var palces []model.Place
	for _, each := range result.Results {
		place := model.Place{each.Name, each.Rating}
		palces = append(palces, place)
	}
	return palces
}
