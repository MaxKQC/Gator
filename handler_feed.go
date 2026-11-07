package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MaxKQC/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	if _, err := createFeedFollow(ctx, s, user, feed); err != nil {
		return fmt.Errorf("couldn't follow newly created feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:        %v\n", feed.ID)
	fmt.Printf(" * Created:   %v\n", feed.CreatedAt)
	fmt.Printf(" * Updated:   %v\n", feed.UpdatedAt)
	fmt.Printf(" * Name:      %v\n", feed.Name)
	fmt.Printf(" * URL:       %v\n", feed.Url)
	fmt.Printf(" * User ID:   %v\n", feed.UserID)
}
