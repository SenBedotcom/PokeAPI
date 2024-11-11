package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

// PokemonService handles data retrieval from PokeAPI
type PokemonService struct {
	HttpClient   *http.Client
	CacheService *CacheService
}

// NewPokemonService initializes and returns a new instance of PokemonService
func NewPokemonService(cacheService *CacheService) *PokemonService {
	return &PokemonService{
		HttpClient:   &http.Client{},
		CacheService: cacheService,
	}
}

var pokeCache = cache.New(10*time.Minute, 15*time.Minute)

func FetchPokemon(name string) (map[string]interface{}, error) {
	if data, found := pokeCache.Get(name); found {
		return data.(map[string]interface{}), nil
	}

	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	pokeCache.Set(name, result, cache.DefaultExpiration)
	return result, nil
}

func (s *PokemonService) GetPokemonByName(name string) (interface{}, error) {
	cacheKey := fmt.Sprintf("pokemon_%s", name)
	if cachedData, found := s.CacheService.Get(cacheKey); found {
		return cachedData, nil
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	resp, err := s.HttpClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve pokemon data")
	}
	defer resp.Body.Close()

	var pokemonData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pokemonData); err != nil {
		return nil, err
	}

	// Cache the result for 10 minutes
	s.CacheService.Set(cacheKey, pokemonData)

	return pokemonData, nil
}

func (s *PokemonService) GetPokemonAbility(name string) (interface{}, error) {
	cacheKey := fmt.Sprintf("pokemon_ability_%s", name)
	if cachedAbilities, found := s.CacheService.Get(cacheKey); found {
		return cachedAbilities, nil
	}

	pokemonData, err := s.GetPokemonByName(name)
	if err != nil {
		return nil, err
	}

	abilities := extractAbilities(pokemonData)
	s.CacheService.Set(cacheKey, abilities)

	return abilities, nil
}

func (s *PokemonService) GetRandomPokemon() (interface{}, error) {
	randomID := getRandomPokemonID()
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", randomID)
	resp, err := s.HttpClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve random pokemon")
	}
	defer resp.Body.Close()

	var pokemonData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pokemonData); err != nil {
		return nil, err
	}

	return pokemonData, nil
}

func extractAbilities(pokemonData interface{}) []string {
	var abilities []string
	data, ok := pokemonData.(map[string]interface{})
	if !ok {
		return abilities
	}

	if abilitiesData, exists := data["abilities"].([]interface{}); exists {
		for _, ability := range abilitiesData {
			abilityMap, ok := ability.(map[string]interface{})
			if ok && abilityMap["ability"] != nil {
				if abilityDetail, ok := abilityMap["ability"].(map[string]interface{}); ok {
					if abilityName, exists := abilityDetail["name"].(string); exists {
						abilities = append(abilities, abilityName)
					}
				}
			}
		}
	}
	return abilities
}

func getRandomPokemonID() int {
	return int(1 + (time.Now().UnixNano() % 898))
}
