package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/JohnPlummer/post-scorer/scorer"
	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/joho/godotenv"
)

func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError // Default to ERROR for safety
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	// Initialize logger
	opts := &slog.HandlerOptions{
		Level: getLogLevel(os.Getenv("LOG_LEVEL")),
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	// Read custom prompt
	promptText, err := os.ReadFile("custom_prompt.txt")
	if err != nil {
		slog.Error("Error reading prompt file", "error", err)
		os.Exit(1)
	}

	// Initialize the scorer
	cfg := scorer.Config{
		OpenAIKey:  os.Getenv("OPENAI_API_KEY"),
		PromptText: string(promptText),
	}

	s, err := scorer.New(cfg)
	if err != nil {
		slog.Error("Failed to create scorer", "error", err)
		os.Exit(1)
	}

	// Read posts from CSV file
	file, err := os.Open("example_posts.csv")
	if err != nil {
		slog.Error("Error opening posts file", "error", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		slog.Error("Error reading CSV header", "error", err)
		os.Exit(1)
	}

	var posts []reddit.Post
	records, err := reader.ReadAll()
	if err != nil {
		slog.Error("Error reading CSV records", "error", err)
		os.Exit(1)
	}

	for _, record := range records {
		posts = append(posts, reddit.Post{
			ID:       record[0],
			Title:    record[1],
			SelfText: record[2],
		})
	}

	// Score the posts
	scoredPosts, err := s.ScorePosts(context.Background(), posts)
	if err != nil {
		slog.Error("Failed to score posts", "error", err)
		os.Exit(1)
	}

	// Print results
	for _, post := range scoredPosts {
		fmt.Printf("Post: %s\nScore: %.2f\n\n", post.Post.Title, post.Score)
	}
}
