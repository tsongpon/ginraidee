package transport

type GooglePlaceTransport struct {
	NextPageToken string   `json:"next_page_token"`
	Results       []PlaceResult `json:"results"`
}

type PlaceResult struct {
	PlaceID string `json:"place_id"`
	Name    string `json:"name"`
	Rating  float32 `json:"rating"`
}
