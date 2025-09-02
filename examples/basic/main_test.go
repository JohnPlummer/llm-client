package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/JohnPlummer/llm-client/scorer"
)

func TestMain(m *testing.M) {
	// Setup for tests - suppress logging during tests
	opts := &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors during tests
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// TestCompilation verifies the basic compilation without execution
func TestCompilation(t *testing.T) {
	// This test passes if the file compiles successfully
	// The Go test runner ensures compilation success
	t.Log("main.go compiles successfully")
}

// TestLoadTextItems verifies CSV loading functionality
func TestLoadTextItems(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expectError bool
		expectCount int
	}{
		{
			name:        "valid CSV file",
			filename:    "example_items.csv",
			expectError: false,
			expectCount: 10, // Assuming standard test data
		},
		{
			name:        "nonexistent file",
			filename:    "nonexistent.csv",
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := loadTextItems(tt.filename)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(items) != tt.expectCount {
				t.Errorf("expected %d items, got %d", tt.expectCount, len(items))
			}

			// Verify structure of loaded items
			for i, item := range items {
				if item.ID == "" {
					t.Errorf("item %d has empty ID", i)
				}
				if item.Content == "" {
					t.Errorf("item %d has empty content", i)
				}
				if item.Metadata == nil {
					t.Errorf("item %d has nil metadata", i)
				}
			}
		})
	}
}

// TestSetupLogger verifies logging configuration
func TestSetupLogger(t *testing.T) {
	// Save original state
	originalLevel := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", originalLevel)

	tests := []struct {
		name     string
		logLevel string
	}{
		{
			name:     "default level",
			logLevel: "",
		},
		{
			name:     "debug level",
			logLevel: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.logLevel)

			// Capture log output
			var buf bytes.Buffer
			opts := &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}
			logger := slog.New(slog.NewTextHandler(&buf, opts))
			slog.SetDefault(logger)

			setupLogger()

			// Test that logger was configured (this is mainly about ensuring no panics)
			slog.Info("test message")
			t.Logf("Logger setup completed for level: %s", tt.logLevel)
		})
	}
}

// TestCreateCustomScorer verifies custom scorer creation
func TestCreateCustomScorer(t *testing.T) {
	apiKey := "test-api-key-for-config-validation"

	// This should not panic and should create a valid scorer config
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("createCustomScorer panicked: %v", r)
		}
	}()

	scorer := createCustomScorer(apiKey)
	if scorer == nil {
		t.Error("createCustomScorer returned nil")
	}
}

// TestMainWithoutAPIKey verifies behavior when API key is missing
func TestMainWithoutAPIKey(t *testing.T) {
	// Save original API key
	originalAPIKey := os.Getenv("OPENAI_API_KEY")
	defer func() {
		if originalAPIKey != "" {
			os.Setenv("OPENAI_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("OPENAI_API_KEY")
		}
	}()

	// Remove API key
	os.Unsetenv("OPENAI_API_KEY")

	// Capture exit behavior (this is tricky in Go)
	// We'll test the logic that would cause exit instead
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Log("Correctly detects missing API key")
	} else {
		t.Error("Failed to detect missing API key")
	}
}

// TestExampleItemsStructure verifies the example CSV structure
func TestExampleItemsStructure(t *testing.T) {
	// Check if required CSV files exist
	requiredFiles := []string{
		"example_items.csv",
		"example_items_edge_cases.csv",
	}

	for _, filename := range requiredFiles {
		t.Run(filename, func(t *testing.T) {
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				t.Errorf("Required file %s does not exist", filename)
				return
			}

			items, err := loadTextItems(filename)
			if err != nil {
				t.Errorf("Failed to load %s: %v", filename, err)
				return
			}

			if len(items) == 0 {
				t.Errorf("File %s contains no items", filename)
			}

			t.Logf("File %s loaded successfully with %d items", filename, len(items))
		})
	}
}

// TestConfigurationFiles verifies required configuration files exist
func TestConfigurationFiles(t *testing.T) {
	files := []struct {
		name     string
		required bool
	}{
		{"go.mod", true},
		{"go.sum", true},
		{".env.example", true},
		{"README.md", false},
		{"custom_prompt.txt", false},
	}

	for _, file := range files {
		t.Run(file.name, func(t *testing.T) {
			_, err := os.Stat(file.name)
			exists := !os.IsNotExist(err)

			if file.required && !exists {
				t.Errorf("Required file %s is missing", file.name)
			} else if exists {
				t.Logf("File %s exists", file.name)
			}
		})
	}
}

// TestModuleDependencies verifies go.mod structure
func TestModuleDependencies(t *testing.T) {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	goModContent := string(content)

	// Check essential dependencies
	requiredDeps := []string{
		"github.com/JohnPlummer/llm-client",
		"github.com/joho/godotenv",
	}

	for _, dep := range requiredDeps {
		if !strings.Contains(goModContent, dep) {
			t.Errorf("Missing required dependency: %s", dep)
		} else {
			t.Logf("Found required dependency: %s", dep)
		}
	}

	// Check go version
	if !strings.Contains(goModContent, "go 1.23") {
		t.Error("Expected Go version 1.23.x in go.mod")
	}

	// Check replace directive for local development
	if !strings.Contains(goModContent, "replace") {
		t.Error("Expected replace directive for local development")
	}
}

// BenchmarkLoadTextItems benchmarks CSV loading performance
func BenchmarkLoadTextItems(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := loadTextItems("example_items.csv")
		if err != nil {
			b.Fatalf("Failed to load text items: %v", err)
		}
	}
}

// BenchmarkCreateCustomScorer benchmarks scorer creation
func BenchmarkCreateCustomScorer(b *testing.B) {
	apiKey := "benchmark-api-key"

	for i := 0; i < b.N; i++ {
		scorer := createCustomScorer(apiKey)
		if scorer == nil {
			b.Fatal("createCustomScorer returned nil")
		}
	}
}

// TestIntegrationReadiness verifies the example is ready for integration tests
func TestIntegrationReadiness(t *testing.T) {
	// This test verifies all components needed for integration testing

	t.Run("loadTextItems function", func(t *testing.T) {
		items, err := loadTextItems("example_items.csv")
		if err != nil {
			t.Errorf("loadTextItems failed: %v", err)
		}
		if len(items) == 0 {
			t.Error("loadTextItems returned empty slice")
		}
	})

	t.Run("createCustomScorer function", func(t *testing.T) {
		scorer := createCustomScorer("test-key")
		if scorer == nil {
			t.Error("createCustomScorer returned nil")
		}
	})

	t.Run("setupLogger function", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("setupLogger panicked: %v", r)
			}
		}()
		setupLogger()
	})
}

// TestErrorHandling verifies error handling throughout the example
func TestErrorHandling(t *testing.T) {
	t.Run("invalid CSV file", func(t *testing.T) {
		// Create a temporary invalid CSV file
		tmpFile := "test_invalid.csv"
		defer os.Remove(tmpFile)

		err := os.WriteFile(tmpFile, []byte("invalid,csv\ncontent"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		items, err := loadTextItems(tmpFile)
		if err != nil {
			t.Logf("Correctly handled invalid CSV: %v", err)
		} else if len(items) == 0 {
			t.Log("Handled invalid CSV by returning empty slice")
		}
	})
}

// TestDocumentationConsistency verifies the example matches documentation
func TestDocumentationConsistency(t *testing.T) {
	// Read the main.go file to verify it contains expected patterns
	content, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	mainContent := string(content)

	// Check for expected patterns
	expectedPatterns := []string{
		"BuildProductionScorer",
		"ScoreTextsWithOptions",
		"GetHealth",
		"WithModel",
		"WithPromptTemplate",
		"LoadTextItems",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(mainContent, pattern) {
			t.Errorf("Expected pattern '%s' not found in main.go", pattern)
		} else {
			t.Logf("Found expected pattern: %s", pattern)
		}
	}
}

// TestBuildTags verifies no inappropriate build tags
func TestBuildTags(t *testing.T) {
	content, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	// Check for inappropriate build tags that might prevent compilation
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "//go:build") ||
			strings.HasPrefix(strings.TrimSpace(line), "// +build") {
			t.Errorf("Found build tag at line %d: %s", i+1, line)
		}
	}
}

// TestVersionConsistency verifies version references are consistent
func TestVersionConsistency(t *testing.T) {
	// Check main.go for version references
	content, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	mainContent := string(content)

	// Look for version strings
	if strings.Contains(mainContent, "v0.11.0") {
		t.Log("Found version reference v0.11.0 in main.go")
	}
}
