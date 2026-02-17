package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	name := flag.String("name", "World", "The name to greet")
	flag.Parse()

	msg, err := Greet(*name)
	if err != nil {
		slog.Error("failed to greet", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Println(msg)
	slog.Info("Greeting sent", slog.String("name", *name))
}

var ErrEmptyName = errors.New("name cannot be empty")

// Greet returns a greeting message for the given name.
func Greet(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("open %s: %w", "greeting", ErrEmptyName) // ERR-1 context wrap
	}
	return fmt.Sprintf("Hello, %s!", name), nil
}
