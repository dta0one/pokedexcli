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

	if data, ok := cfg.cache.Get(url); ok {
		locationAreas := LocationAreasResponse{}
		err := json.Unmarshal(data, &locationAreas)
		if err != nil {
			return LocationAreasResponse{}, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return locationAreas, nil
	}

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

	cfg.cache.Add(url, data)

	locationAreas := LocationAreasResponse{}
	err = json.Unmarshal(data, &locationAreas)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return locationAreas, nil
}

type LocationAreaResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func fetchLocationArea(cfg *config, areaName string) (LocationAreaResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	if data, ok := cfg.cache.Get(url); ok {
		locationArea := LocationAreaResponse{}
		err := json.Unmarshal(data, &locationArea)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return locationArea, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to fetch location area: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return LocationAreaResponse{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	cfg.cache.Add(url, data)

	locationArea := LocationAreaResponse{}
	err = json.Unmarshal(data, &locationArea)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return locationArea, nil
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

func fetchPokemon(cfg *config, pokemonName string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	if data, ok := cfg.cache.Get(url); ok {
		pokemon := Pokemon{}
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return Pokemon{}, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return pokemon, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to fetch pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to read response body: %w", err)
	}

	cfg.cache.Add(url, data)

	pokemon := Pokemon{}
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return pokemon, nil
}
