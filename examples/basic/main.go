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
	posts, err := loadPosts("example_posts.csv")
	if err != nil {
		slog.Error("Error loading posts", "error", err)
		os.Exit(1)
	}

	// Load comments and associate them with posts
	if err := loadComments("example_comments.csv", posts); err != nil {
		slog.Error("Error loading comments", "error", err)
		os.Exit(1)
	}

	// Score the posts
	var postSlice []reddit.Post
	for _, p := range posts {
		postSlice = append(postSlice, *p)
	}
	scoredPosts, err := s.ScorePosts(context.Background(), postSlice)
	if err != nil {
		slog.Error("Failed to score posts", "error", err)
		os.Exit(1)
	}

	// Print results
	for _, post := range scoredPosts {
		fmt.Printf("Post: %s\nScore: %d\nReason: %s\n\n", post.Post.Title, post.Score, post.Reason)
	}
}

// loadPosts reads posts from a CSV file
func loadPosts(filename string) (map[string]*reddit.Post, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening posts file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading CSV header: %w", err)
	}

	posts := make(map[string]*reddit.Post)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV records: %w", err)
	}

	for _, record := range records {
		posts[record[0]] = &reddit.Post{
			ID:       record[0],
			Title:    record[1],
			SelfText: record[2],
		}
	}

	return posts, nil
}

// loadComments reads comments from a CSV file and associates them with posts
func loadComments(filename string, posts map[string]*reddit.Post) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening comments file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("reading CSV header: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reading CSV records: %w", err)
	}

	for _, record := range records {
		postID, text := record[0], record[1]
		if post, exists := posts[postID]; exists {
			post.Comments = append(post.Comments, reddit.Comment{
				Body: text,
			})
		}
	}

	return nil
}
