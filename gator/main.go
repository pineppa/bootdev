package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/src"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	state := src.CliState{
		Cfg: config.Read(),
	}
	db, err := sql.Open("postgres", state.Cfg.DBUrl)
	if err != nil {
		fmt.Println("Error opening the SQL server")
		os.Exit(1)
	}
	state.DbQueries = database.New(db)
	cmds := src.RegisterCommands()

	if len(os.Args) < 2 {
		fmt.Println("Error, add something here and you... Just call a command!")
		os.Exit(1)
	}

	cmd := src.CliCommand{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cmds.Run(&state, cmd); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
