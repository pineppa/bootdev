package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

const UrlBasis = "https://pokeapi.co/api/v2/location-area/"

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	conf := initConfig()
	// Infinite loop asking for input until exit is triggered
	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")

		// Trigger exit
		if !scanner.Scan() {
			break
		}

		// Retrieve, clean and process input
		text := scanner.Text()

		texts := cleanInput(text)
		if texts[0] == "help" {
			HelpCommand()
		} else {
			cmd, ok := cliCom[texts[0]]
			if ok {
				cmd.Callback(&conf, texts)
			} else {
				fmt.Println("Unknown command! Use help to get a list of commands")
			}
		}
	}
	// callback Exit
	os.Exit(0)
}
