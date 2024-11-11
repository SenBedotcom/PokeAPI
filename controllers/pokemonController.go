package controllers

import (
	"net/http"
	"pokemon-api/services"

	"github.com/gin-gonic/gin"
)

// PokemonController struct
type PokemonController struct {
	PokemonService services.PokemonService
	CacheService   services.CacheService
}

// NewPokemonController is a constructor for PokemonController
func NewPokemonController(pokemonService services.PokemonService, cacheService services.CacheService) *PokemonController {
	return &PokemonController{
		PokemonService: pokemonService,
		CacheService:   cacheService,
	}
}

// GetPokemonByName handles the GET /pokemon/:name route
func (pc *PokemonController) GetPokemonByName(c *gin.Context) {
	name := c.Param("name")
	if data, found := pc.CacheService.Get(name); found {
		c.JSON(http.StatusOK, data)
		return
	}

	data, err := pc.PokemonService.GetPokemonByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon not found"})
		return
	}

	pc.CacheService.Set(name, data)
	c.JSON(http.StatusOK, data)
}

// GetPokemonAbility handles the GET /pokemon/:name/ability route
func (pc *PokemonController) GetPokemonAbility(c *gin.Context) {
	name := c.Param("name")
	if data, found := pc.CacheService.Get(name + "_ability"); found {
		c.JSON(http.StatusOK, data)
		return
	}

	data, err := pc.PokemonService.GetPokemonAbility(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ability not found"})
		return
	}

	pc.CacheService.Set(name+"_ability", data)
	c.JSON(http.StatusOK, data)
}

// GetRandomPokemon handles the GET /pokemon/random route
func (pc *PokemonController) GetRandomPokemon(c *gin.Context) {
	data, err := pc.PokemonService.GetRandomPokemon()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching random Pokemon"})
		return
	}
	c.JSON(http.StatusOK, data)
}
