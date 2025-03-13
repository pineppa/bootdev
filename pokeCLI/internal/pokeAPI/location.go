package pokeAPI

import "encoding/json"

type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocVersionDetail struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}

type EncounterMethodRate struct {
	EncounterMethod EncounterMethod    `json:"encounter_method"`
	VersionDetails  []LocVersionDetail `json:"version_details"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Name struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}

type LocPokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ConditionValue struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterDetail struct {
	MinLevel        int              `json:"min_level"`
	MaxLevel        int              `json:"max_level"`
	ConditionValues []ConditionValue `json:"condition_values"`
	Chance          int              `json:"chance"`
	Method          EncounterMethod  `json:"method"`
}

type PokemonVersionDetail struct {
	Version          Version           `json:"version"`
	MaxChance        int               `json:"max_chance"`
	EncounterDetails []EncounterDetail `json:"encounter_details"`
}

type PokemonEncounter struct {
	Pokemon               LocPokemon             `json:"pokemon"`
	PokemonVersionDetails []PokemonVersionDetail `json:"version_details"`
}

type LocationArea struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             Location              `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

func CheckLocationForPokemons(location *LocationArea) []string {
	var pokemons []string
	for _, encounter := range location.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon.Name)
	}
	return pokemons
}

func FetchLocationArea(url string) (*LocationArea, error) {
	body, err := fetchFromURL(url)
	if err != nil {
		return nil, err
	}

	var locationArea LocationArea
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return nil, err
	}

	return &locationArea, nil
}
