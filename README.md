# Post Scorer

A Go package that scores Reddit posts using ChatGPT to evaluate their likelihood of containing information about events, activities, or venue recommendations.

## Installation

```bash
go get github.com/JohnPlummer/post-scorer
```

## Usage

```golang
package main

import (
    "context"
    "log"
    "os"
    "github.com/JohnPlummer/post-scorer/scorer"
)

func main() {
    cfg := scorer.Config{
        OpenAIKey: os.Getenv("OPENAI_API_KEY"),
    }

    s, err := scorer.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    posts := []reddit.Post{
        // Your Reddit posts here
    }

    scoredPosts, err := s.ScorePosts(context.Background(), posts)
    if err != nil {
        log.Fatal(err)
    }

    for _, post := range scoredPosts {
        log.Printf("Post: %s, Score: %.2f\n", post.Post.Title, post.Score)
    }
}
```

## Project Structure

```text
.
├── scorer/            # Main package code
│   ├── scorer.go     # Core scoring logic
│   └── scorer_test.go # Tests
├── examples/          # Example usage
│   └── basic/        # Basic usage example
│       └── main.go   # Example implementation
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
