package adapter

import "github.com/tsongpon/ginraidee/model"

type GeoCodeAdapter interface {
	GetLocation(string) (*model.Location, error)
}
