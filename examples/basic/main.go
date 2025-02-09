package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/JohnPlummer/post-scorer/scorer"
	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read custom prompt
	promptText, err := os.ReadFile("custom_prompt.txt")
	if err != nil {
		log.Fatal("Error reading prompt file:", err)
	}

	// Initialize the scorer
	cfg := scorer.Config{
		OpenAIKey:  os.Getenv("OPENAI_API_KEY"),
		PromptText: string(promptText),
		Debug:      false,
	}

	s, err := scorer.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Read posts from CSV file
	file, err := os.Open("example_posts.csv")
	if err != nil {
		log.Fatal("Error opening posts file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		log.Fatal("Error reading CSV header:", err)
	}

	var posts []reddit.Post
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV records:", err)
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
		log.Fatal(err)
	}

	// Print results
	for _, post := range scoredPosts {
		fmt.Printf("Post: %s\nScore: %.2f\n\n", post.Post.Title, post.Score)
	}
}
