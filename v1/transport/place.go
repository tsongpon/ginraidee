package transport

type PlaceTransport struct {
	PlaceID string `json:"placeId"`
	Name    string `json:"name"`
	Link    string `json:"link"`
}

type PlaceResponseTransport struct {
	Result   []PlaceTransport `json:"results"`
	NextPage *string          `json:"next"`
}
