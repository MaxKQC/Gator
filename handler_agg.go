package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MaxKQC/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid time_between_reqs: %w", err)
	}
	if timeBetweenRequests <= 0 {
		return fmt.Errorf("time_between_reqs must be greater than zero")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	fmt.Printf("Fetching feed: %s\n", feed.Name)
	if err := s.db.MarkFeedFetched(ctx, feed.ID); err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed %s: %w", feed.Name, err)
	}

	for _, item := range rssFeed.Channel.Items {
		if item.Link == "" {
			fmt.Printf("Skipping post %q with no URL\n", item.Title)
			continue
		}

		publishedAt, err := parsePublishedAt(item.PubDate)
		if err != nil {
			fmt.Printf("Couldn't parse published time for post %q: %v\n", item.Title, err)
		}

		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			fmt.Printf("Couldn't save post %q: %v\n", item.Title, err)
		}
	}
	return nil
}

func parsePublishedAt(value string) (sql.NullTime, error) {
	if value == "" {
		return sql.NullTime{}, nil
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
	}

	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return sql.NullTime{Time: parsed.UTC(), Valid: true}, nil
		}
		lastErr = err
	}
	return sql.NullTime{}, lastErr
}
