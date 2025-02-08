package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ArturM94/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

type GetLocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type GetLocationAreaDetailsResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocationAreas(cache *pokecache.Cache, url *string) (*GetLocationAreasResponse, error) {
	var fullURL string

	if url != nil {
		fullURL = *url
	} else {
		fullURL = baseURL + "/location-area"
	}

	var data []byte

	cachedData, ok := cache.Get(fullURL)
	if ok {
		data = cachedData
	} else {
		res, err := http.Get(fullURL)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return nil, fmt.Errorf("error getting request: status code is %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %w", err)
		}

		cache.Add(fullURL, data)
	}

	var response GetLocationAreasResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &response, nil
}

func GetLocationAreaDetails(cache *pokecache.Cache, idOrName string) (*GetLocationAreaDetailsResponse, error) {
	fullURL := baseURL + "/location-area/" + idOrName

	var data []byte

	cachedData, ok := cache.Get(fullURL)
	if ok {
		data = cachedData
	} else {
		res, err := http.Get(fullURL)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return nil, fmt.Errorf("error getting request: status code is %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %w", err)
		}
	}

	var response GetLocationAreaDetailsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &response, nil
}
