package main

import (
	"github.com/SergioFloresCorrea/gator/internal/config"
	"github.com/SergioFloresCorrea/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func createState(db *database.Queries, cfg *config.Config) *state {
	return &state{db: db, cfg: cfg}
}

func createCommand(name string, args []string) command {
	return command{name: name, args: args}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	err := c.handlers[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}
