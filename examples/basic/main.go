package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/JohnPlummer/post-scorer/scorer"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Note: .env file not found, using environment variables")
	}

	// Setup logging
	setupLogger()

	// Get API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		slog.Error("OPENAI_API_KEY environment variable not set")
		os.Exit(1)
	}

	fmt.Println("=== Post-Scorer v0.10.0 Example ===\n")

	// Example 1: Production-ready scorer with all resilience features
	fmt.Println("1. Production Configuration (with resilience patterns)")
	fmt.Println("   - Circuit breaker: prevents cascade failures")
	fmt.Println("   - Retry logic: handles transient errors")
	fmt.Println("   - Prometheus metrics: production monitoring")
	fmt.Println("   - Health checks: service status monitoring\n")

	productionScorer, err := scorer.BuildProductionScorer(apiKey)
	if err != nil {
		slog.Error("Failed to create production scorer", "error", err)
		os.Exit(1)
	}

	// Load test data
	items, err := loadTextItems("example_posts.csv")
	if err != nil {
		slog.Error("Failed to load test data", "error", err)
		os.Exit(1)
	}

	// Score with production scorer
	ctx := context.Background()
	results, err := productionScorer.ScoreTexts(ctx, items)
	if err != nil {
		slog.Error("Failed to score texts", "error", err)
		
		// Check health status when scoring fails
		health := productionScorer.GetHealth(ctx)
		if !health.Healthy {
			slog.Error("Scorer is unhealthy", "status", health.Status, "details", health.Details)
		}
		os.Exit(1)
	}

	// Display results
	fmt.Println("Results:")
	for _, result := range results {
		fmt.Printf("  [Score: %3d] %s\n", result.Score, result.Item.Content)
		fmt.Printf("    Reason: %s\n\n", result.Reason)
	}

	// Example 2: Custom configuration with specific resilience settings
	fmt.Println("\n2. Custom Configuration Example")
	customScorer := createCustomScorer(apiKey)
	
	// Score with custom prompt template
	customResults, err := customScorer.ScoreTextsWithOptions(ctx, items,
		scorer.WithModel("gpt-4o-mini"),
		scorer.WithPromptTemplate("Rate this for relevance to local events: {{.Content}}"),
	)
	if err != nil {
		slog.Warn("Custom scoring failed", "error", err)
	} else {
		fmt.Printf("  Scored %d items with custom configuration\n", len(customResults))
	}

	// Example 3: Health monitoring
	fmt.Println("\n3. Health Check")
	health := productionScorer.GetHealth(ctx)
	fmt.Printf("  Status: %s\n", health.Status)
	fmt.Printf("  Healthy: %v\n", health.Healthy)
	if health.Details != nil {
		if integration, ok := health.Details["integration"].(map[string]interface{}); ok {
			fmt.Printf("  Circuit Breaker: %v\n", integration["circuit_breaker_enabled"])
			fmt.Printf("  Retry Enabled: %v\n", integration["retry_enabled"])
			fmt.Printf("  Metrics: %v\n", integration["metrics_enabled"])
		}
	}

	// Example 4: Start metrics server
	fmt.Println("\n4. Starting Metrics Server")
	go startMetricsServer()
	fmt.Println("  Prometheus metrics available at http://localhost:8080/metrics")
	fmt.Println("  Health check available at http://localhost:8080/health")
	
	fmt.Println("\nPress Ctrl+C to exit...")
	select {} // Keep running to serve metrics
}

func setupLogger() {
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	}
	
	opts := &slog.HandlerOptions{
		Level: level,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func createCustomScorer(apiKey string) scorer.TextScorer {
	// Create a custom configuration with specific settings
	cfg := scorer.NewDefaultConfig(apiKey)
	cfg = cfg.WithCircuitBreaker()
	cfg = cfg.WithRetry()
	cfg = cfg.WithMaxConcurrent(5)
	cfg = cfg.WithTimeout(30 * time.Second)
	
	// Customize retry strategy
	cfg.RetryConfig.Strategy = scorer.RetryStrategyExponential
	cfg.RetryConfig.MaxAttempts = 3
	cfg.RetryConfig.InitialDelay = 100 * time.Millisecond
	cfg.RetryConfig.MaxDelay = 2 * time.Second
	
	// Customize circuit breaker
	cfg.CircuitBreakerConfig.MaxRequests = 5
	cfg.CircuitBreakerConfig.Interval = 10 * time.Second
	cfg.CircuitBreakerConfig.Timeout = 30 * time.Second
	
	s, err := scorer.NewIntegratedScorer(cfg)
	if err != nil {
		slog.Error("Failed to create custom scorer", "error", err)
		os.Exit(1)
	}
	
	return s
}

func loadTextItems(filename string) ([]scorer.TextItem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading records: %w", err)
	}

	items := make([]scorer.TextItem, 0, len(records))
	for _, record := range records {
		if len(record) >= 3 {
			// Combine title and body into content
			content := record[1] + " - " + record[2]
			items = append(items, scorer.TextItem{
				ID:      record[0],
				Content: content,
				Metadata: map[string]interface{}{
					"title":  record[1],
					"body":   record[2],
					"source": "csv_file",
				},
			})
		}
	}

	return items, nil
}

func startMetricsServer() {
	mux := http.NewServeMux()
	
	// Metrics endpoint
	mux.Handle("/metrics", scorer.GetMetricsHandler())
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Would need access to the scorer instance for real health check
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Metrics server failed", "error", err)
	}
}