package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/SergioFloresCorrea/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("the user must be registered before login")
	}

	fmt.Printf("The username %s has been set.\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the name")
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		fmt.Println("a user with the same name has already been registered")
		os.Exit(1)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	_, err = s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	if err = s.cfg.SetUser(name); err != nil {
		return err
	}
	fmt.Printf("A user with the name %s was created\n", name)
	fmt.Printf("User data: %+v\n", params)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteAll(context.Background()); err != nil {
		fmt.Printf("%v, the database couldn't be deleted.", err)
		os.Exit(1)
		return err
	}
	fmt.Println("Successful deletion.")
	os.Exit(0)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUserName := s.cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expecting a single argument, the time between requests")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	userID := user.ID

	if len(cmd.args) != 2 {
		return fmt.Errorf("adding a feed expects a name and a url")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	params := database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    userID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User data: %+v\n", feed)

	paramsFeedFollow := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), paramsFeedFollow)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for idx, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed %d:\n", idx+1)
		fmt.Printf(" * Name: %s\n", feed.Name)
		fmt.Printf(" * URL: %s\n", feed.Url)
		fmt.Printf(" * Uploaded by: %s\n", user.Name)
		fmt.Println("")
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the follow command expects a single argument")
	}
	url := cmd.args[0]

	userID := user.ID

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	feedID := feed.ID

	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feedID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("User data: %+v\n", feedFollow)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	userID := user.ID
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), userID)
	if err != nil {
		return err
	}

	for idx, feed := range feeds {
		fmt.Printf("Feed #%d: %s\n", idx+1, feed.FeedName)
	}
	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("to unfollow a feed, an url is expected")
	}
	url := cmd.args[0]
	userID := user.ID
	params := database.UnFollowFeedParams{
		Url:    url,
		UserID: userID,
	}
	if err := s.db.UnFollowFeed(context.Background(), params); err != nil {
		return err
	}
	return nil
}

func handlerBrowse(s *state, cmd command) error {
	var limit int32
	if len(cmd.args) == 0 {
		limit = 2
	} else {
		limit64, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(limit64)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), limit)
	if err != nil {
		return err
	}
	for _, post := range posts {
		data, err := json.MarshalIndent(post, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	if err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	for _, feedItem := range feed.Channel.Item {
		pubDateTime, err := parsePubDate(feedItem.PubDate)
		if err != nil {
			return err
		}
		params := database.CreatePostParams{
			CreatedAt:   nextFeed.CreatedAt,
			UpdatedAt:   nextFeed.UpdatedAt,
			Title:       stringToNull(feedItem.Title),
			Url:         feedItem.Link,
			Description: stringToNull(feedItem.Description),
			PublishedAt: pubDateTime,
			FeedID:      nextFeed.ID,
		}
		err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			return err
		}
	}
	return nil
}

func stringToNull(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func parsePubDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,                     // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,                      // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,                      // "02 Jan 06 15:04 -0700"
		time.RFC822,                       // "02 Jan 06 15:04 MST"
		time.RFC3339,                      // "2006-01-02T15:04:05Z07:00"
		"Mon, 02 Jan 2006 15:04:05 -0700", // Explicit custom (some feeds use odd spacing)
		"02 Jan 2006 15:04:05 MST",        // Custom variations
	}

	var err error
	for _, layout := range formats {
		var t time.Time
		t, err = time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported pubDate format: %q", dateStr)
}
