package model

type Place struct {
	Name    string  `json:"name"`
	Rating  float32 `json:"rating"`
	MapLink string  `json:"link"`
}
