package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

type GetLocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(url *string) (*GetLocationAreasResponse, error)  {
	var fullURL string

	if url != nil {
		fullURL = *url
	} else {
		fullURL = baseURL + "/location-area"
	}

	res, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("error getting request: status code is %d", res.StatusCode)
	}
	
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var response GetLocationAreasResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &response, nil
}
