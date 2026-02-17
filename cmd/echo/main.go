package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func main() {
	envLog := os.Getenv("LOG_ENABLED") == "false"
	logFlag := flag.Bool("log", envLog, "Enable logging (env: LOG_ENABLED)")
	flag.Parse()

	var logOut io.Writer = io.Discard
	if *logFlag {
		logOut = os.Stdout
	}
	logger := slog.New(slog.NewTextHandler(logOut, nil))
	slog.SetDefault(logger)

	prefix := os.Getenv("ECHO_PREFIX")
	args := flag.Args()
	if len(args) == 0 {
		slog.Warn("no arguments provided")
		return
	}

	msg := strings.Join(args, " ")
	if prefix != "" {
		msg = fmt.Sprintf("%s: %s", prefix, msg)
	}
	fmt.Println(msg)
	slog.Info("echoed message", slog.Int("arg_count", len(args)), slog.String("prefix", prefix))
}
