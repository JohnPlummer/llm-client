# Post Scorer

A Go package that scores Reddit posts using ChatGPT to evaluate their likelihood of containing information about events, activities, or venue recommendations.

## Installation

```bash
go get github.com/JohnPlummer/post-scorer
```

## Usage

### Basic Usage

```golang
package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/JohnPlummer/post-scorer/scorer"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    godotenv.Load()

    // Initialize logger
    opts := &slog.HandlerOptions{
        Level: getLogLevel(os.Getenv("LOG_LEVEL")),
    }
    logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
    slog.SetDefault(logger)

    cfg := scorer.Config{
        OpenAIKey: os.Getenv("OPENAI_API_KEY"),
    }

    s, err := scorer.New(cfg)
    if err != nil {
        slog.Error("Failed to create scorer", "error", err)
        os.Exit(1)
    }

    posts := []reddit.Post{
        // Your Reddit posts here
    }

    scoredPosts, err := s.ScorePosts(context.Background(), posts)
    if err != nil {
        slog.Error("Failed to score posts", "error", err)
        os.Exit(1)
    }

    for _, post := range scoredPosts {
        slog.Info("Scored post", "title", post.Post.Title, "score", post.Score)
    }
}
```

### Using Custom Prompt

You can provide your own scoring prompt:

```golang
promptText, err := os.ReadFile("custom_prompt.txt")
if err != nil {
    log.Fatal(err)
}

cfg := scorer.Config{
    OpenAIKey:  os.Getenv("OPENAI_API_KEY"),
    PromptText: string(promptText),
    Debug:      false,
}
```

### Reading Posts from CSV

The package includes an example of reading posts from a CSV file:

```golang
file, err := os.Open("posts.csv")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

reader := csv.NewReader(file)
reader.TrimLeadingSpace = true
reader.LazyQuotes = true

// Skip header row
_, err = reader.Read()
if err != nil {
    log.Fatal(err)
}

var posts []reddit.Post
records, err := reader.ReadAll()
if err != nil {
    log.Fatal(err)
}

for _, record := range records {
    posts = append(posts, reddit.Post{
        ID:       record[0],
        Title:    record[1],
        SelfText: record[2],
    })
}
```

## Configuration Options

The `Config` struct supports the following options:

```golang
type Config struct {
    OpenAIKey     string  // Required: Your OpenAI API key
    PromptText    string  // Optional: Custom scoring prompt
    MaxConcurrent int     // Optional: For rate limiting
}
```

### Environment Variables

```env
OPENAI_API_KEY=your-openai-api-key-here
LOG_LEVEL=debug  # Options: debug, info, warn, error (defaults to error)
```

## Project Structure

```text
.
├── scorer/            # Main package code
│   ├── scorer.go     # Core scoring logic
│   └── scorer_test.go # Tests
├── examples/          # Example usage
│   └── basic/        # Basic usage example
│       ├── main.go   # Example implementation
│       └── custom_prompt.txt # Example custom prompt
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Example Files

The `examples/basic` directory contains a complete working example including:

- `main.go`: Example implementation
- `custom_prompt.txt`: Example custom scoring prompt
- `example_posts.csv`: Example CSV input file
