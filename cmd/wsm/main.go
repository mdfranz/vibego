package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Config holds our environment configuration.
type Config struct {
	WsRoot string
	GitDir string
	Log    bool
}

var (
	ErrInvalidSelection = errors.New("invalid selection")
	ErrNoBranches       = errors.New("no branches found")
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	setupLogging(cfg.Log)

	// Ensure we have a title for the folder
	if flag.NArg() < 1 {
		fmt.Println("Usage: wsm [options] <workspace-title>")
		flag.Usage()
		os.Exit(1)
	}
	inputTitle := strings.Join(flag.Args(), " ")

	ctx := context.Background()

	// 2. Git Operations: List and Select
	selectedBranch, err := selectBranch(ctx, cfg.GitDir)
	if err != nil {
		slog.Error("failed to select branch", slog.Any("error", err))
		os.Exit(1)
	}

	// 3. Git Operations: Checkout
	slog.Info("Switching git repo", slog.String("branch", selectedBranch))
	if err := checkoutBranch(ctx, cfg.GitDir, selectedBranch); err != nil {
		slog.Error("failed to checkout branch", slog.Any("error", err))
		os.Exit(1)
	}

	// 4. Workspace Creation
	folderName := generateFolderName(inputTitle)
	fullPath := filepath.Join(cfg.WsRoot, folderName)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		slog.Error("failed to create directory", slog.Any("error", err), slog.String("path", fullPath))
		os.Exit(1)
	}

	fmt.Println("✅ Workspace ready!")
	fmt.Printf("   📂 Directory: %s\n", fullPath)
	fmt.Printf("   🌿 Git Branch: %s\n", selectedBranch)
	slog.Info("workspace created", slog.String("path", fullPath), slog.String("branch", selectedBranch))
}

func setupLogging(enabled bool) {
	var logOut io.Writer = io.Discard
	if enabled {
		logOut = os.Stdout
	}
	logger := slog.New(slog.NewTextHandler(logOut, nil))
	slog.SetDefault(logger)
}

// loadConfig fetches config from flags and environment variables.
func loadConfig() (Config, error) {
	rootEnv := os.Getenv("WS_ROOT")
	gitDirEnv := os.Getenv("WS_GIT_DIR")
	logEnv := os.Getenv("LOG_ENABLED") == "true"

	root := flag.String("root", rootEnv, "Workspace root directory (env: WS_ROOT)")
	gitDir := flag.String("git-dir", gitDirEnv, "Git repository directory (env: WS_GIT_DIR)")
	logFlag := flag.Bool("log", logEnv, "Enable logging (env: LOG_ENABLED)")

	flag.Parse()

	if *root == "" {
		return Config{}, errors.New("WS_ROOT or -root flag is required")
	}
	if *gitDir == "" {
		return Config{}, errors.New("WS_GIT_DIR or -git-dir flag is required")
	}

	return Config{
		WsRoot: *root,
		GitDir: *gitDir,
		Log:    *logFlag,
	}, nil
}

// selectBranch lists branches and prompts the user.
func selectBranch(ctx context.Context, gitDir string) (string, error) {
	// We use --sort=-committerdate to show the most recently worked-on branches first
	cmd := exec.CommandContext(ctx, "git", "-C", gitDir, "branch", "--sort=-committerdate", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git branch command failed: %w", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(branches) == 0 || (len(branches) == 1 && branches[0] == "") {
		return "", fmt.Errorf("listing branches: %w", ErrNoBranches)
	}

	// Display top 15 branches to avoid clutter
	limit := 15
	if len(branches) < limit {
		limit = len(branches)
	}

	fmt.Println("👇 Select a branch to switch to:")
	for i := 0; i < limit; i++ {
		fmt.Printf("[%d] %s\n", i+1, branches[i])
	}

	// Prompt user
	fmt.Print("\nEnter number: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	selection, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil || selection < 1 || selection > limit {
		return "", fmt.Errorf("selecting branch: %w", ErrInvalidSelection)
	}

	return branches[selection-1], nil
}

// checkoutBranch executes git checkout for the specified branch.
func checkoutBranch(ctx context.Context, gitDir, branch string) error {
	cmd := exec.CommandContext(ctx, "git", "-C", gitDir, "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git checkout %s: %w", branch, err)
	}
	return nil
}

// generateFolderName creates "YYYY-MM-DD_slug-title".
func generateFolderName(input string) string {
	date := time.Now().Format("2006-01-02")

	// Simple slugify: lowercase, trim, replace non-alphanumeric with hyphen, collapse hyphens
	slug := strings.ToLower(strings.TrimSpace(input))
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return fmt.Sprintf("%s_%s", date, slug)
}
