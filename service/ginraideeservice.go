package service

import (
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/model"
)

type GinRaiDeeService struct {
	placeAdapter   adapter.PlaceAdapter
	geoCodeAdapter adapter.GeoCodeAdapter
}

func NewGinRaiDeeService(placeAdapter adapter.PlaceAdapter, geoCodeAdapter adapter.GeoCodeAdapter) *GinRaiDeeService {
	service := new(GinRaiDeeService)
	service.placeAdapter = placeAdapter
	service.geoCodeAdapter = geoCodeAdapter
	return service
}

func (s *GinRaiDeeService) GetRestaurants(address string) ([]model.Place, error) {
	var err error
	location, err := s.geoCodeAdapter.GetLocation(address)
	if err != nil {
		return nil, err
	}
	restaurants, err :=  s.placeAdapter.GetPlaces("restaurant", location.Lat, location.Lng)

	return restaurants, nil
}