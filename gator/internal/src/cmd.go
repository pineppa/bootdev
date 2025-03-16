package src

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

type CliState struct {
	DbQueries *database.Queries
	Cfg       config.Config
}

type CliCommand struct {
	Name string
	Args []string
}

type CliCommands struct {
	Commands map[string]func(*CliState, CliCommand) error
}

// This method registers a new handler function for a command name.
func (c *CliCommands) Register(name string, f func(*CliState, CliCommand) error) {
	c.Commands[name] = f
}

// This method runs a given command with the provided CliState if it exists.
func (c *CliCommands) Run(s *CliState, cmd CliCommand) error {
	f := c.Commands[cmd.Name]
	err := f(s, cmd)
	return err
}

func RegisterCommands() CliCommands {
	cmds := CliCommands{
		Commands: make(map[string]func(*CliState, CliCommand) error),
	}

	cmds.Register("login", HandlerLogin)
	cmds.Register("register", HandlerRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerUsers)
	cmds.Register("agg", HandlerAgg)
	cmds.Register("addfeed", HandlerAddFeed)
	// cmds.Register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	cmds.Register("feeds", HandlerFeeds)
	cmds.Register("follow", HandlerFollow)
	cmds.Register("following", HandlerFollowing)
	cmds.Register("unfollow", HandlerUnfollow)
	//	cmds.Register("unfollow", middlewareLoggedIn(HandleUnfollow))
	return cmds
}

// Handles the Login command with a username
func HandlerLogin(s *CliState, cmd CliCommand) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("the function accept one single argument")
	}
	if _, err := s.DbQueries.GetUser(
		context.Background(), cmd.Args[0],
	); err != nil {
		return fmt.Errorf("you can't login to an account that doesn't exist")
	}
	s.Cfg.CurrentUserName = cmd.Args[0]
	err := s.Cfg.SetUser()
	if err != nil {
		return err
	}
	fmt.Printf("The user %s has been set\n", s.Cfg.CurrentUserName)
	return nil
}

func HandlerRegister(s *CliState, cmd CliCommand) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("the function accept one single argument")
	}
	if _, err := s.DbQueries.GetUser(
		context.Background(), cmd.Args[0],
	); err == nil {
		return fmt.Errorf("user is already registered")
	}

	db_user, err := s.DbQueries.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      sql.NullString{String: cmd.Args[0], Valid: true},
		},
	)
	s.Cfg.CurrentUserName = db_user.Name.String
	if s.Cfg.SetUser() != nil {
		return err
	}
	fmt.Printf("The newly registered user %s has been set\n", s.Cfg.CurrentUserName)
	return err
}

func HandlerReset(s *CliState, cmd CliCommand) error {
	err := s.DbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}
	fmt.Println("All users deleted successfully!")
	if err := s.DbQueries.DeleteAllFeeds(context.Background()); err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}
	fmt.Println("All feeds deleted successfully!")
	if err := s.DbQueries.DeleteAllFeedFollows(context.Background()); err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}
	fmt.Println("All feeds_follow deleted successfully!")
	return nil
}

func HandlerUsers(s *CliState, cmd CliCommand) error {
	users, err := s.DbQueries.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}
	for _, user := range users {
		name := user.Name.String
		fmt.Printf("* %s", name)
		if name == s.Cfg.CurrentUserName {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}
	return nil
}

func HandlerAgg(s *CliState, cmd CliCommand) error {
	ctx := context.Background()
	rssFeed, err := fetchFeed(ctx, rssUrl)
	if err != nil {
		fmt.Printf("Error fetching the feed: %v\n", err)
		return nil
	}
	printRssFeed(rssFeed)
	return nil
}

// // wrapper function
// func middlewareLoggedIn(handler func(s *CliState, cmd CliCommand, user database.User) error) (func(*CliState, CliCommand), error) {

// }

func HandlerAddFeed(s *CliState, cmd CliCommand) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed takes exactly 2 arguments: name, url")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	ctx := context.Background()
	userEntry, err := s.DbQueries.GetUser(ctx, s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	feed_id := uuid.New()
	s.DbQueries.CreateFeed(ctx, database.CreateFeedParams{
		ID:        feed_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    userEntry.ID,
	})
	s.DbQueries.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userEntry.ID,
			FeedID:    feed_id,
		},
	)

	rssFeed, err := fetchFeed(ctx, url)
	if err != nil {
		fmt.Printf("Error fetching the feed: %v\n", err)
		return nil
	}
	printRssFeed(rssFeed)
	return nil
}

func HandlerFeeds(s *CliState, cmd CliCommand) error {
	feeds, err := s.DbQueries.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error in retrieving the feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Println("New feed retrieved:")
		fmt.Println("Name: ", feed.Name)
		fmt.Println("Url: ", feed.Url)
		fmt.Println("Username: ", feed.Username.String)
		fmt.Println("")
	}
	return nil
}

func HandlerFollow(s *CliState, cmd CliCommand) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("wrong amount of arguments")
	}
	fmt.Println(s.Cfg.CurrentUserName)
	url := cmd.Args[0]
	ctx := context.Background()

	userEntry, err := s.DbQueries.GetUser(ctx, s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	fmt.Println(s.Cfg.CurrentUserName)

	feedDb, err := s.DbQueries.GetFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}

	_, err = s.DbQueries.GetFeedFollowsForUserFeedPair(
		ctx, database.GetFeedFollowsForUserFeedPairParams{
			Username: s.Cfg.CurrentUserName,
			Feedurl:  feedDb.Url,
		})

	if err == nil {
		return fmt.Errorf("the User, Feed pair already exists")
	}

	_, err = s.DbQueries.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userEntry.ID,
			FeedID:    feedDb.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error in retrieving the feeds: %w", err)
	}
	fmt.Printf("Feed name: %s\n", feedDb.Name)
	fmt.Printf("Current User: %s\n", userEntry.Name.String)
	return err
}

func HandlerFollowing(s *CliState, cmd CliCommand) error {
	userFeeds, err := s.DbQueries.GetFeedFollowsForUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error in retrieving the feeds: %w", err)
	}

	for _, feed := range userFeeds {
		fmt.Println(" - ", feed.Username.String, feed.Feedname)
	}
	return nil
}

func HandlerUnfollow(s *CliState, cmd CliCommand) error {
	if err := s.DbQueries.UnfollowFeedFollow(context.Background(), database.UnfollowFeedFollowParams{
		Username: s.Cfg.CurrentUserName,
		Feedurl:  cmd.Args[0],
	}); err != nil {
		return fmt.Errorf("failed to unfollow the feed - %v", err)
	}
	return nil
}
