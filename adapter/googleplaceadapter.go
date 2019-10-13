package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tsongpon/ginraidee/model"
	"github.com/tsongpon/ginraidee/transport"
	"log"
)

const defaultRedis = "5000"
const mapLinkBasURL = "https://www.google.com/maps/place/?q=place_id:"

type GooglePlaceAdapter struct {
}

func NewGooglePlaceAdapter() *GooglePlaceAdapter {
	return new(GooglePlaceAdapter)
}

func (a *GooglePlaceAdapter) GetPlaces(placeType string, lat float32, lng float32) ([]model.Place, error) {
	var err error
	client := resty.New()
	location := fmt.Sprintf("%f", lat) + "," + fmt.Sprintf("%f", lng)
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"location": location,
			"type":     placeType,
			"radius":   defaultRedis,
			"key":      "AIzaSyBfGD0y888DZ8FUfpBjDCRVRhFimnG0z78",
		}).
		SetHeader("Accept", "application/json").
		Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json")

	if err != nil {
		log.Printf("get error %s", err.Error())
		return nil, err
	}
	result := transport.GooglePlaceTransport{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Printf("Unmarshal result error %s", err.Error())
		return nil, err
	}
	var places []model.Place
	for _, each := range result.Results {
		mapLink := mapLinkBasURL + each.PlaceID
		place := model.Place{Name: each.Name, Ratting: each.Rating, MapLink: mapLink}
		places = append(places, place)
	}
	return places, nil
}
