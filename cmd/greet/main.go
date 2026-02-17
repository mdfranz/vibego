package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
)

func main() {
	envLog := os.Getenv("LOG_ENABLED") == "false"
	logFlag := flag.Bool("log", envLog, "Enable logging (env: LOG_ENABLED)")
	defaultName := os.Getenv("GREET_NAME")
	if defaultName == "" {
		defaultName = "World"
	}

	name := flag.String("name", defaultName, "The name to greet (env: GREET_NAME)")
	flag.Parse()

	var logOut io.Writer = io.Discard
	if *logFlag {
		logOut = os.Stdout
	}
	logger := slog.New(slog.NewTextHandler(logOut, nil))
	slog.SetDefault(logger)

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
