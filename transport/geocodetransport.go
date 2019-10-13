package transport

type GoogleGeoCodeTransport struct {
	NextPageToken string          `json:"next_page_token"`
	Results       []GeoCodeResult `json:"results"`
}

type GeoCodeResult struct {
	Geometry Geometry `json:"geometry"`
	Name     string   `json:"formatted_address"`
	PlaceID  string  `json:"place_id"`
}

type Geometry struct {
	Location struct {
		Lat float32 `json:"Lat"`
		Lng float32 `json:"lng"`
	} `json:"location"`
}
