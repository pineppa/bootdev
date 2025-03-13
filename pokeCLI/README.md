# POKE CLI

## A pokemon game, played in a command line interface written in Go

This is a project created as part of the Boot.dev training to learn how to work in Golang. It was a fun project to execture to reinforce the learning on the basics of the language and the `http` module on top of how to setup our own project.

Note: This project won't be improved in the near future

# How to install:

## Linux

- Download the repository in your preferred folder.
- Open the terminal in the same folder and add run ```go build -o pokeCLI```
- Play the game with `./pokeCLI`

# How to play:

 - exit: Exits the Pokedex,
 - help: Provides information about the possible commands,
 - map: Displays the next 20 locations,
 - mapb: Displays the previous 20 locations,
 - remap: Re-displays the current 20 locations,
 - explore: Explore available Pokémon in an area,
 - catch: Attempts to catch a Pokémon,
 - inspect: Inspect a Pokémon's stats and types,
 - pokedex: Displays the list of Pokémon captured in the Pokédex.

# Current status:

- Pokémon are not limited to they specific area in the map
- The rate of capture is 50%, independently of levels

## Possible suggested improvements

- Update the CLI to support the "up" arrow to cycle through previous commands
- Simulate battles between pokemon
- Add more unit tests (Kept in the local environment for now)
- Refactor your code to organize it better and make it more testable
- Keep pokemon in a "party" and allow them to level up
- Allow for pokemon that are caught to evolve after a set amount of time
- Persist a user's Pokedex to disk so they can save progress between sessions
- Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe you are given choices of areas and just type "left" or "right"
- Random encounters with wild pokemon
- Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching pokemon
