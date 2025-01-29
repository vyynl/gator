package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"vyynl/gator/internal/database"

	"github.com/google/uuid"
)

/* Command handling */
type commands struct {
	commandList map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	name := cmd.name
	f, exists := c.commandList[name]
	if exists {
		return f(s, cmd)
	} else {
		return fmt.Errorf("ERROR command not registered")
	}
}

func handlerDBUrl(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("ERROR dburl handler expects one argument: postgres db url")
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("ERROR login handler expects one argument: username")
	}

	user, err := s.db.GetUser(
		context.Background(),
		cmd.args[0],
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("ERROR user not registered")
		}
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("Username successfully set")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("ERROR reset handler expects no additional arguments")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.Current_user_name {
			fmt.Printf("%s (current)", user.Name)
			continue
		}
		fmt.Println(user.Name)
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("ERROR register handler expects one argument: new username")
	}

	/* Checking if user already exists prior to creating duplicate record */
	_, err := s.db.GetUser(
		context.Background(),
		cmd.args[0],
	)
	if err == nil {
		return fmt.Errorf("ERROR user already registered")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}

	/* Creating new user */
	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()
	name := cmd.args[0]

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        id,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
			Name:      name,
		},
	)
	if err != nil {
		return fmt.Errorf("ERROR registering user: %v", err)
	}

	/* Updating config to match new user */
	s.cfg.SetUser(user.Name)
	fmt.Println("User successfully registered:")
	fmt.Printf("ID: %v\nCreatedAt: %v\nUpdatedAt: %v\nName: %v\n", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("ERROR reset handler expects no additional arguments")
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}
	err = s.db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}
	err = s.db.ResetFeedFollows(context.Background())
	if err != nil {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}
	err = s.db.ResetPosts(context.Background())
	if err != nil {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}
	fmt.Println("Reset successful")
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("ERROR addFeed handler expects 2 args: name, url")
	}
	if strings.Contains(cmd.args[0], "https://") || strings.Contains(cmd.args[0], "www.") {
		return fmt.Errorf("ERROR addFeed handler expects name then url as args")
	}

	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()
	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return err
	}

	feed, err := s.db.AddFeed(
		context.Background(),
		database.AddFeedParams{
			ID:        id,
			UserID:    user.ID,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
			Name:      name,
			Url:       url,
		},
	)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR adding feed: %v", err)
	}
	fmt.Println(feed)

	follow_id := uuid.New()
	_, err = s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        follow_id,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR adding feed follow: %v", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("ERROR feeds handler expects no additional arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching feeds: %v", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("DATABASE ERROR fetching user name: %v", err)
		}
		fmt.Printf("Name: %s - URL: %s - Username: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("ERROR follow handler expects only 1 argument: follow url")
	}

	id := uuid.New()
	url := cmd.args[0]
	created_at := time.Now()
	updated_at := time.Now()

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching user: %v", err)
	}
	user_id := user.ID

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching feed: %v", err)
	}
	feed_id := feed.ID

	follow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        id,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
			UserID:    user_id,
			FeedID:    feed_id,
		},
	)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR adding follow: %v", err)
	}

	fmt.Println(follow)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("ERROR followers handler expect no additional arguments")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching user: %v", err)
	}

	following, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching following: %v", err)
	}

	for _, feed := range following {
		fmt.Printf("Feed: %s\n", feed.FeedName)
	}

	fmt.Println(following)
	return nil
}

func handlerUnfollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("ERROR unfollow handler expects only 1 argument: unfollow url")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching current user: %v", err)
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching feed: %v", err)
	}

	err = s.db.RemoveFeedFollow(
		context.Background(),
		database.RemoveFeedFollowParams{
			FeedID: feed.ID,
			UserID: user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR unfollowing feed: %v", err)
	}

	return nil
}

func handlerBrowse(s *state, cmd command) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("ERROR browse handler expects 1 argument max: limit")
	}

	user, _ := s.db.GetUser(context.Background(), s.cfg.Current_user_name)

	var limit int = 2
	if len(cmd.args) > 0 {
		num, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("ERROR invalid argument: %v", err)
		}
		limit = num
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching posts: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\nDescription: %s\nLink: %v\n\n", post.Title, post.Description.String, post.Url)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("ERROR agg handler expects only 1 argument: time between reqs")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("ERROR parsing time between reqs: %v", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

/* Middleware */
func middlewareLoggedIn(handler func(s *state, cmd command) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		_, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			return fmt.Errorf("DATABASE ERROR user not logged in: %v", err)
		}
		return handler(s, cmd)
	}
}
