package pokeAPI

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
)

type Pokemon struct {
	ID                     int                `json:"id"`
	Name                   string             `json:"name"`
	BaseExperience         int                `json:"base_experience"`
	Height                 int                `json:"height"`
	IsDefault              bool               `json:"is_default"`
	Order                  int                `json:"order"`
	Weight                 int                `json:"weight"`
	Abilities              []Ability          `json:"abilities"`
	Forms                  []NamedAPIResource `json:"forms"`
	GameIndices            []GameIndex        `json:"game_indices"`
	HeldItems              []HeldItem         `json:"held_items"`
	LocationAreaEncounters string             `json:"location_area_encounters"`
	Moves                  []Move             `json:"moves"`
	Species                NamedAPIResource   `json:"species"`
	Sprites                Sprites            `json:"sprites"`
	Cries                  Cries              `json:"cries"`
	Stats                  []Stat             `json:"stats"`
	Types                  []Type             `json:"types"`
	PastTypes              []PastType         `json:"past_types"`
}

type Ability struct {
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
	Ability  NamedAPIResource `json:"ability"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GameIndex struct {
	GameIndex int              `json:"game_index"`
	Version   NamedAPIResource `json:"version"`
}

type HeldItem struct {
	Item               NamedAPIResource    `json:"item"`
	PokeVersionDetails []PokeVersionDetail `json:"version_details"`
}

type PokeVersionDetail struct {
	Rarity  int              `json:"rarity"`
	Version NamedAPIResource `json:"version"`
}

type Move struct {
	Move                NamedAPIResource     `json:"move"`
	VersionGroupDetails []VersionGroupDetail `json:"version_group_details"`
}

type VersionGroupDetail struct {
	LevelLearnedAt  int              `json:"level_learned_at"`
	VersionGroup    NamedAPIResource `json:"version_group"`
	MoveLearnMethod NamedAPIResource `json:"move_learn_method"`
}

type Sprites struct {
	BackDefault      string       `json:"back_default"`
	BackFemale       *string      `json:"back_female"`
	BackShiny        string       `json:"back_shiny"`
	BackShinyFemale  *string      `json:"back_shiny_female"`
	FrontDefault     string       `json:"front_default"`
	FrontFemale      *string      `json:"front_female"`
	FrontShiny       string       `json:"front_shiny"`
	FrontShinyFemale *string      `json:"front_shiny_female"`
	Other            OtherSprites `json:"other"`
}

type OtherSprites struct {
	DreamWorld      ImageSet `json:"dream_world"`
	Home            ImageSet `json:"home"`
	OfficialArtwork ImageSet `json:"official-artwork"`
	Showdown        ImageSet `json:"showdown"`
}

type ImageSet struct {
	FrontDefault     string  `json:"front_default"`
	FrontFemale      *string `json:"front_female"`
	FrontShiny       *string `json:"front_shiny"`
	FrontShinyFemale *string `json:"front_shiny_female"`
}

type Cries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type Stat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

type Type struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

type PastType struct {
	Generation NamedAPIResource `json:"generation"`
	Types      []Type           `json:"types"`
}

var pokedex = make(map[string]Pokemon)

func CatchPokemon(pokemon *Pokemon) bool {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	// Calculate the chances of catching that pokemon
	if rand.IntN(100) > 50 {
		fmt.Println("Pokemon was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
		pokedex[pokemon.Name] = *pokemon
		return true
	}
	fmt.Println("Pokemon has not been captured")
	return false
}

func FetchPokemon(pokemonName string) (*Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	body, err := fetchFromURL(url)
	if err != nil {
		return nil, err
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return nil, err
	}
	return &pokemon, nil
}

func printPokeStats(pokemon Pokemon) {
	fmt.Println("Name: ", pokemon.Name)
	fmt.Println("Height: ", pokemon.Height)
	fmt.Println("Weight: ", pokemon.Weight)
	fmt.Println("Stats: ")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types: ")
	for _, typeVal := range pokemon.Types {
		fmt.Printf("  - %s\n", typeVal.Type.Name)
	}
}

func InspectPokemon(pokemonName string) {
	pokemon, ok := pokedex[pokemonName]
	if !ok {
		fmt.Println("Error fetching the pokemon; Check again the pokemon name")
		return
	}
	printPokeStats(pokemon)
}

func ShowPokedex() {
	fmt.Println("Your Pokedex:")
	for _, val := range pokedex {
		fmt.Println(" - ", val.Name)
	}
}
