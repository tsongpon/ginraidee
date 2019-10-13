package model

type Location struct {
	Name string  `json:"name"`
	Lat  float32 `json:"lat"`
	Lng  float32 `json:"lng"`
}
