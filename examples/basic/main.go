package main

import (
	"context"
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
	}

	s, err := scorer.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Example posts
	posts := []reddit.Post{
		{
			ID:       "1",
			Title:    "Best coffee shops in downtown",
			SelfText: "Looking for recommendations for coffee shops with good wifi for working.",
		},
		{
			ID:       "2",
			Title:    "Cat picture",
			SelfText: "Here's my cat sleeping.",
		},
		{
			ID:       "3",
			Title:    "Tonight: Live music at The Basement Bar",
			SelfText: "Local band playing at 8pm, $10 cover. Great venue for indie music!",
		},
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
