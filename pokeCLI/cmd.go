package main

import (
	"fmt"
	"os"
	"strconv"

	"pokeCLI/internal/pokeAPI"
)

type cliCommand struct {
	Name        string
	Description string
	Callback    func(conf *config, texts []string)
}

func (c cliCommand) GetDescription() string {
	return c.Description
}

func HelpCommand() {
	fmt.Println("Welcome to the pokeCLI! Usage:")
	for name, cmd := range cliCom {
		fmt.Printf("> %s: %s\n", name, cmd.GetDescription())
	}
}

func ExitCommand(conf *config, texts []string) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

const NResPerPage = 20

func MapCommand(conf *config, texts []string) {
	fmt.Println("Fetching locations data...")
	cmdName := texts[0]
	if cmdName == "mapb" && conf.CurrentId < NResPerPage {
		fmt.Println("You are in the first page! There is no way backward")
		return
	}
	if cmdName == "map" {
		conf.CurrentId = conf.NextId
		conf.NextId += NResPerPage
		conf.PreviousId = conf.CurrentId
	} else if cmdName == "mapb" {
		conf.NextId -= NResPerPage
		conf.CurrentId -= NResPerPage
		conf.PreviousId -= NResPerPage
	}
	Map(conf.CurrentId)
}

func Map(id int) {
	for i := range int(NResPerPage) {
		url := UrlBasis + strconv.Itoa(id+i)
		location, err := pokeAPI.FetchLocationArea(url)
		if err != nil {
			fmt.Println("Error fetching location:", err)
			fmt.Println("Error fetching location data:", location)
			break
		}
		fmt.Println(location.Name)
	}
}

func ExploreCommand(conf *config, texts []string) {
	if len(texts[0]) == 1 {
		fmt.Println("Please insert the right location to be visited... Check the help for more info")
	}
	url := UrlBasis + texts[1]
	location, err := pokeAPI.FetchLocationArea(url)
	if err != nil {
		fmt.Println("Error fetching the mentioned location! Retry")
		return
	}
	fmt.Printf("Exploring %v...\n", location.Name)
	pokemons := pokeAPI.CheckLocationForPokemons(location)
	for _, poke := range pokemons {
		fmt.Println(poke)
	}
}

func CatchCommand(conf *config, texts []string) {
	pokemon, err := pokeAPI.FetchPokemon(texts[1])
	if err != nil {
		fmt.Println("Error fetching the pokemon; Check again the pokemon name")
		return
	}
	pokeAPI.CatchPokemon(pokemon)
}

func InspectCommand(conf *config, texts []string) {
	pokeAPI.InspectPokemon(texts[1])
}

func pokedexCommand(conf *config, texts []string) {
	pokeAPI.ShowPokedex()
}

func NewCLICommand(name, description string, callback func(*config, []string)) cliCommand {
	return cliCommand{Name: name, Description: description, Callback: callback}
}

var cliCom = map[string]cliCommand{
	"exit":    NewCLICommand("exit", "Exit the Pokedex", ExitCommand),
	"help":    NewCLICommand("help", "Provides information about the possible commands", func(conf *config, texts []string) {}),
	"map":     NewCLICommand("map", "Displays the next 20 locations", MapCommand),
	"mapb":    NewCLICommand("mapb", "Displays the previous 20 locations", MapCommand),
	"remap":   NewCLICommand("remap", "Re-displays the current 20 locations", MapCommand),
	"explore": NewCLICommand("explore", "Explore available Pokémon in an area", ExploreCommand),
	"catch":   NewCLICommand("catch", "Attempts to catch a Pokémon", CatchCommand),
	"inspect": NewCLICommand("inspect", "Inspect a Pokémon's stats and types", InspectCommand),
	"pokedex": NewCLICommand("pokedex", "Displays the list of Pokémon captured in the Pokédex", pokedexCommand),
}
