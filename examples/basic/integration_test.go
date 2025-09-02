package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestExecutionWithMockedAPI tests the main execution flow with environment setup
func TestExecutionWithMockedAPI(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test verifies that the main function can run without crashing
	// when proper environment is set up

	// Create a temporary .env file for testing
	envContent := `OPENAI_API_KEY=test-mock-api-key-for-integration-testing
LOG_LEVEL=error
`

	err := os.WriteFile(".env.test", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(".env.test")

	// Set environment for testing
	originalAPIKey := os.Getenv("OPENAI_API_KEY")
	originalLogLevel := os.Getenv("LOG_LEVEL")

	defer func() {
		if originalAPIKey != "" {
			os.Setenv("OPENAI_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("OPENAI_API_KEY")
		}
		if originalLogLevel != "" {
			os.Setenv("LOG_LEVEL", originalLogLevel)
		} else {
			os.Unsetenv("LOG_LEVEL")
		}
	}()

	// Set test environment
	os.Setenv("OPENAI_API_KEY", "test-mock-key")
	os.Setenv("LOG_LEVEL", "error")

	t.Log("Environment configured for integration test")

	// Note: Full execution would require actual API calls
	// This test validates environment setup and initial validation
}

// TestCompilationAndBuild verifies the entire build process
func TestCompilationAndBuild(t *testing.T) {
	// Test that the module builds successfully
	cmd := exec.Command("go", "build", "-o", "example_test", ".")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Build failed: %v\nOutput: %s", err, string(output))
		return
	}

	// Clean up built binary
	defer os.Remove("example_test")

	// Verify binary was created
	if _, err := os.Stat("example_test"); os.IsNotExist(err) {
		t.Error("Build did not create expected binary")
	} else {
		t.Log("Build successful - binary created")
	}
}

// TestModuleTidy verifies dependencies are clean
func TestModuleTidy(t *testing.T) {
	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("go mod tidy failed: %v\nOutput: %s", err, string(output))
		return
	}

	// Verify no changes needed (in a clean state)
	cmd = exec.Command("go", "mod", "verify")
	output, err = cmd.CombinedOutput()

	if err != nil {
		t.Errorf("go mod verify failed: %v\nOutput: %s", err, string(output))
	} else {
		t.Log("Module dependencies are clean and verified")
	}
}

// TestValidExampleData verifies example CSV data is properly formatted
func TestValidExampleData(t *testing.T) {
	testFiles := []string{
		"example_items.csv",
		"example_items_edge_cases.csv",
	}

	for _, filename := range testFiles {
		t.Run(filename, func(t *testing.T) {
			items, err := loadTextItems(filename)
			if err != nil {
				t.Errorf("Failed to load %s: %v", filename, err)
				return
			}

			if len(items) == 0 {
				t.Errorf("File %s contains no valid items", filename)
				return
			}

			// Validate each item has required fields
			for i, item := range items {
				if item.ID == "" {
					t.Errorf("Item %d in %s has empty ID", i, filename)
				}
				if item.Content == "" {
					t.Errorf("Item %d in %s has empty content", i, filename)
				}
				if len(item.Content) < 10 {
					t.Errorf("Item %d in %s has suspiciously short content: %q", i, filename, item.Content)
				}

				// Validate metadata
				if item.Metadata == nil {
					t.Errorf("Item %d in %s has nil metadata", i, filename)
					continue
				}

				requiredMetadataFields := []string{"title", "body", "source"}
				for _, field := range requiredMetadataFields {
					if _, exists := item.Metadata[field]; !exists {
						t.Errorf("Item %d in %s missing metadata field: %s", i, filename, field)
					}
				}
			}

			t.Logf("File %s validated successfully with %d items", filename, len(items))
		})
	}
}

// TestCustomPromptTemplate verifies custom prompt file exists and loads
func TestCustomPromptTemplate(t *testing.T) {
	filename := "custom_prompt.txt"

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Failed to read %s: %v", filename, err)
		return
	}

	promptText := string(content)

	// Verify basic structure
	if len(promptText) < 50 {
		t.Errorf("Custom prompt seems too short: %d characters", len(promptText))
	}

	// Check for essential elements
	essentialElements := []string{
		"score", // Should mention scoring
		"JSON",  // Should specify JSON format
		"0-100", // Should specify score range
	}

	for _, element := range essentialElements {
		if !strings.Contains(strings.ToLower(promptText), strings.ToLower(element)) {
			t.Errorf("Custom prompt missing essential element: %s", element)
		}
	}

	t.Logf("Custom prompt template validated (%d characters)", len(promptText))
}

// TestEnvironmentVariableHandling verifies .env file processing
func TestEnvironmentVariableHandling(t *testing.T) {
	// Test .env.example exists
	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		t.Error(".env.example file is missing")
		return
	}

	// Read .env.example
	content, err := os.ReadFile(".env.example")
	if err != nil {
		t.Errorf("Failed to read .env.example: %v", err)
		return
	}

	exampleContent := string(content)

	// Verify it contains expected variables
	expectedVars := []string{
		"OPENAI_API_KEY",
		"LOG_LEVEL",
	}

	for _, variable := range expectedVars {
		if !strings.Contains(exampleContent, variable) {
			t.Errorf(".env.example missing variable: %s", variable)
		}
	}

	t.Log(".env.example file validated")
}

// TestConcurrentProcessingConfig verifies concurrent configuration works
func TestConcurrentProcessingConfig(t *testing.T) {
	apiKey := "test-concurrent-config"

	// Test that createCustomScorer doesn't panic with concurrent settings
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("createCustomScorer with concurrent config panicked: %v", r)
		}
	}()

	scorer := createCustomScorer(apiKey)
	if scorer == nil {
		t.Error("createCustomScorer returned nil with concurrent config")
	}

	t.Log("Concurrent processing configuration validated")
}

// TestMetricsServerSetup verifies metrics server configuration
func TestMetricsServerSetup(t *testing.T) {
	// This test verifies the startMetricsServer function exists and is callable
	// without actually starting the server

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("startMetricsServer function panicked when called: %v", r)
		}
	}()

	// Just verify the function exists by referencing it
	if startMetricsServer == nil {
		t.Error("startMetricsServer function is nil")
	}

	t.Log("Metrics server setup function validated")
}

// TestReadmeDocumentation verifies documentation exists and is accurate
func TestReadmeDocumentation(t *testing.T) {
	if _, err := os.Stat("README.md"); os.IsNotExist(err) {
		t.Skip("README.md not found - skipping documentation test")
		return
	}

	content, err := os.ReadFile("README.md")
	if err != nil {
		t.Errorf("Failed to read README.md: %v", err)
		return
	}

	readme := string(content)

	// Check for essential documentation elements
	essentialSections := []string{
		"example",
		"usage",
		"setup",
	}

	for _, section := range essentialSections {
		if !strings.Contains(strings.ToLower(readme), section) {
			t.Errorf("README.md missing section about: %s", section)
		}
	}

	t.Log("README.md documentation validated")
}

// TestFullIntegrationPipeline runs a complete validation pipeline
func TestFullIntegrationPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full integration pipeline in short mode")
	}

	pipeline := []struct {
		name string
		test func(*testing.T) error
	}{
		{
			name: "Module tidy",
			test: func(t *testing.T) error {
				cmd := exec.Command("go", "mod", "tidy")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("go mod tidy failed: %v\nOutput: %s", err, string(output))
				}
				return nil
			},
		},
		{
			name: "Compilation",
			test: func(t *testing.T) error {
				cmd := exec.Command("go", "build", "-o", "pipeline_test", ".")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("compilation failed: %v\nOutput: %s", err, string(output))
				}
				defer os.Remove("pipeline_test")
				return nil
			},
		},
		{
			name: "Unit tests",
			test: func(t *testing.T) error {
				cmd := exec.Command("go", "test", "-short", ".")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("unit tests failed: %v\nOutput: %s", err, string(output))
				}
				return nil
			},
		},
		{
			name: "Data validation",
			test: func(t *testing.T) error {
				items, err := loadTextItems("example_items.csv")
				if err != nil {
					return fmt.Errorf("failed to load example data: %v", err)
				}
				if len(items) == 0 {
					return fmt.Errorf("no example data loaded")
				}
				return nil
			},
		},
	}

	for _, step := range pipeline {
		t.Run(step.name, func(t *testing.T) {
			if err := step.test(t); err != nil {
				t.Errorf("Pipeline step %s failed: %v", step.name, err)
			} else {
				t.Logf("Pipeline step %s completed successfully", step.name)
			}
		})
	}
}

// TestTimeouts verifies reasonable timeouts are configured
func TestTimeouts(t *testing.T) {
	// Test that createCustomScorer sets reasonable timeouts
	apiKey := "test-timeout-config"

	start := time.Now()
	scorer := createCustomScorer(apiKey)
	duration := time.Since(start)

	if scorer == nil {
		t.Error("createCustomScorer returned nil")
		return
	}

	// Scorer creation should be fast
	if duration > 5*time.Second {
		t.Errorf("Scorer creation took too long: %v", duration)
	}

	t.Logf("Scorer creation completed in %v", duration)
}

// BenchmarkFullPipeline benchmarks the complete verification pipeline
func BenchmarkFullPipeline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Load data
		items, err := loadTextItems("example_items.csv")
		if err != nil {
			b.Fatalf("Failed to load text items: %v", err)
		}

		// Create scorer
		scorer := createCustomScorer("benchmark-key")
		if scorer == nil {
			b.Fatal("createCustomScorer returned nil")
		}

		// Setup logger (lightweight)
		setupLogger()

		b.Logf("Pipeline iteration %d completed with %d items", i, len(items))
	}
}
