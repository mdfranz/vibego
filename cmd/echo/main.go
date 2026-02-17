package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	args := os.Args[1:]
	if len(args) == 0 {
		slog.Warn("no arguments provided")
		return
	}

	msg := strings.Join(args, " ")
	fmt.Println(msg)
	slog.Info("echoed message", slog.Int("arg_count", len(args)))
}
