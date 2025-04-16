package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/SergioFloresCorrea/gator/internal/config"
	"github.com/SergioFloresCorrea/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	var cmds commands
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Printf("%v", err)
	}
	dbQueries := database.New(db)
	s := createState(dbQueries, &cfg)
	cmds.handlers = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnFollow))
	cmds.register("browse", handlerBrowse)
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Println("at least one argument is expected")
		os.Exit(1)
	}

	commandName := argsWithProg[1]
	commandArgs := argsWithProg[2:]
	cmd := createCommand(commandName, commandArgs)
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
