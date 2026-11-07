package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("* %s\n  URL: %s\n  User: %s\n", feed.FeedName, feed.Url, feed.UserName)
	}
	return nil
}
