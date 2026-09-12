package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func fetchLocationAreas(cfg *config, pageURL *string) (LocationAreasResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	// 1. Check cache
	if data, ok := cfg.cache.Get(url); ok {
		locationAreas := LocationAreasResponse{}
		err := json.Unmarshal(data, &locationAreas)
		if err != nil {
			return LocationAreasResponse{}, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return locationAreas, nil
	}

	// 2. Fetch over network
	resp, err := http.Get(url)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return LocationAreasResponse{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// 3. Store in cache
	cfg.cache.Add(url, data)

	locationAreas := LocationAreasResponse{}
	err = json.Unmarshal(data, &locationAreas)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return locationAreas, nil
}
