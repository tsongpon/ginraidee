package adapter

import "github.com/tsongpon/ginraidee/model"

type PlaceAdapter interface {
	GetPlaces(string, float32, float32, string) ([]model.Place, string, error)
}
