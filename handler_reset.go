package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't reset database: %w", err)
	}

	fmt.Println("Database reset successfully!")
	return nil
}
