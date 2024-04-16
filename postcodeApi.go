package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type PostcodeResponse struct {
	Status int `json:"status"`
	Result struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"result"`
}

func getLatLong(postcode string) (float64, float64, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.postcodes.io/postcodes/%s", postcode))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data PostcodeResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, 0, err
	}

	if data.Status != 200 {
		return 0, 0, fmt.Errorf("API response status: %d", data.Status)
	}

	return data.Result.Latitude, data.Result.Longitude, nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of earth in KM
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
