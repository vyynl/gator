package main

import (
	"database/sql"
	"log"
	"os"
	"vyynl/gator/internal/config"
	"vyynl/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	/* Creating our state that the program will interact with */
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	/* Registering command list */
	c := commands{commandList: map[string]func(*state, command) error{}}
	c.register("dburl", handlerDBUrl)
	c.register("login", handlerLogin)
	c.register("users", handlerGetUsers)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("follow", handlerFollow)
	c.register("following", handlerFollowing)
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("browse", middlewareLoggedIn(handlerBrowse))
	c.register("agg", handlerAgg)

	/* Grabbing and parsing out user input to pass through to "c" */
	input := os.Args
	if len(input) < 2 {
		log.Fatalf("No command provided")
	}
	cmd := command{name: input[1]}
	if len(input) >= 3 {
		cmd.args = input[2:]
	}

	/* Executing the valid user input command with args */
	err = c.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
