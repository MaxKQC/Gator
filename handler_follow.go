package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MaxKQC/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	follow, err := createFeedFollow(ctx, s, user, feed)
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)
	return nil
}

func createFeedFollow(ctx context.Context, s *state, user database.User, feed database.Feed) (database.CreateFeedFollowRow, error) {
	return s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	ctx := context.Background()
	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get followed feeds: %w", err)
	}

	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	if err := s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Printf("%s unfollowed %s\n", user.Name, feed.Name)
	return nil
}
