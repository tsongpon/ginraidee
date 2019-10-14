package model

type Place struct {
	PlaceID string  `json:"placeId"`
	Name    string  `json:"name"`
	Rating  float32 `json:"rating"`
	MapLink string  `json:"link"`
}
