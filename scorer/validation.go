// Package scorer provides content validation and sanitization functionality
// for text items before they are processed by the LLM scoring engine.
//
// The validation system ensures content meets quality standards through
// configurable rules for length limits, whitespace handling, and content
// cleanliness. It provides detailed feedback with specific issues and
// actionable suggestions for improvement.
//
// Key components:
// - ValidationResult: Comprehensive validation feedback with issues and suggestions
// - ValidationOptions: Flexible configuration for validation rules
// - Content sanitization: Automatic cleanup of problematic text content
// - Batch processing: Efficient validation of multiple text items
//
// Example usage:
//
//	opts := DefaultValidationOptions()
//	result := ValidateContent("Sample text", opts)
//	if !result.Valid {
//	    log.Printf("Issues: %v", result.Issues)
//	}
package scorer

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidationResult contains comprehensive results of content validation including
// validity status, specific issues found, and actionable suggestions for improvement.
// This structured approach enables callers to provide detailed feedback to users
// about why validation failed and how to fix the problems.
type ValidationResult struct {
	Valid       bool     // Whether the content passes all validation rules
	Issues      []string // Specific problems identified in the content
	Suggestions []string // Actionable recommendations to fix the issues
}

// ValidationOptions configures content validation behavior with flexible rules
// for length limits, whitespace handling, and content requirements. This allows
// different use cases to apply appropriate validation strictness while maintaining
// consistent validation logic across the system.
type ValidationOptions struct {
	MaxLength       int  // Maximum allowed content length in characters
	MinLength       int  // Minimum required content length in characters
	AllowEmpty      bool // Whether empty content is acceptable
	AllowWhitespace bool // Whether whitespace-only content is acceptable
	TrimWhitespace  bool // Whether to trim whitespace before length checks
}

// DefaultValidationOptions returns production-ready validation settings optimized
// for LLM text scoring use cases. These defaults prevent common issues like empty
// content, excessive length that increases API costs, and whitespace-only submissions
// that provide no meaningful content for scoring.
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		MaxLength:       DefaultMaxContentLength,
		MinLength:       MinContentLength,
		AllowEmpty:      false,
		AllowWhitespace: false,
		TrimWhitespace:  true,
	}
}

// ValidateContent performs comprehensive validation of text content against
// configurable rules including length limits, emptiness checks, and whitespace
// handling. It returns detailed feedback with specific issues and actionable
// suggestions, enabling callers to provide meaningful error messages to users.
//
// The validation process follows a structured approach: empty content checks,
// whitespace-only detection, and length validation using either original or
// trimmed content based on configuration. This ensures consistent validation
// behavior across different use cases while maintaining flexibility.
func ValidateContent(content string, opts ValidationOptions) ValidationResult {
	result := ValidationResult{Valid: true}

	// Check for empty content
	if content == "" {
		if !opts.AllowEmpty {
			result.Valid = false
			result.Issues = append(result.Issues, "content is empty")
			result.Suggestions = append(result.Suggestions, "provide meaningful text content")
		}
		return result
	}

	// Check for whitespace-only content
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		if !opts.AllowWhitespace {
			result.Valid = false
			result.Issues = append(result.Issues, "content contains only whitespace")
			result.Suggestions = append(result.Suggestions, "provide non-whitespace content")
		}
		return result
	}

	// Use trimmed content for length checks if TrimWhitespace is enabled
	checkContent := content
	if opts.TrimWhitespace {
		checkContent = trimmed
	}

	// Check minimum length
	if len(checkContent) < opts.MinLength {
		result.Valid = false
		result.Issues = append(result.Issues, fmt.Sprintf("content too short (%d chars, minimum %d)",
			len(checkContent), opts.MinLength))
		result.Suggestions = append(result.Suggestions, "provide more detailed content")
	}

	// Check maximum length
	if len(checkContent) > opts.MaxLength {
		result.Valid = false
		result.Issues = append(result.Issues, fmt.Sprintf("content too long (%d chars, maximum %d)",
			len(checkContent), opts.MaxLength))
		result.Suggestions = append(result.Suggestions, fmt.Sprintf("reduce content to under %d characters", opts.MaxLength))
	}

	return result
}

// ValidateTextItems performs efficient batch validation of multiple text items,
// checking both ID requirements and content validity against the provided options.
// It returns individual validation results for each item along with an aggregate
// error if any items fail validation.
//
// This batch approach enables processing large datasets while maintaining detailed
// per-item feedback. The function continues validation even when individual items
// fail, allowing callers to address all issues in a single pass rather than
// requiring iterative fixes.
func ValidateTextItems(items []TextItem, opts ValidationOptions) ([]ValidationResult, error) {
	if len(items) == 0 {
		return nil, ErrEmptyInput
	}

	results := make([]ValidationResult, len(items))
	hasErrors := false

	for i, item := range items {
		// Validate ID
		if item.ID == "" {
			results[i] = ValidationResult{
				Valid:       false,
				Issues:      []string{"item ID is empty"},
				Suggestions: []string{fmt.Sprintf("provide unique ID for item at index %d", i)},
			}
			hasErrors = true
			continue
		}

		// Validate content
		results[i] = ValidateContent(item.Content, opts)
		if !results[i].Valid {
			hasErrors = true
		}
	}

	if hasErrors {
		// Return results even with errors so caller can see what failed
		return results, fmt.Errorf("validation failed for one or more text items")
	}

	return results, nil
}

// SanitizeContent performs comprehensive text cleaning and normalization to
// prepare content for LLM processing. It removes problematic characters,
// normalizes whitespace, and ensures consistent formatting while preserving
// meaningful content structure like newlines and tabs.
//
// The sanitization process applies three transformations: whitespace trimming,
// whitespace normalization (multiple spaces become single spaces), and removal
// of non-printable characters. This improves LLM processing reliability and
// reduces token consumption from redundant whitespace.
func SanitizeContent(content string) string {
	// Trim leading and trailing whitespace
	content = strings.TrimSpace(content)

	// Normalize whitespace (replace multiple spaces with single space)
	content = normalizeWhitespace(content)

	// Remove non-printable characters except newlines and tabs
	content = removeNonPrintable(content)

	return content
}

// SanitizeTextItems applies content sanitization across a batch of text items
// while preserving item structure and metadata. This batch operation is more
// efficient than individual sanitization calls and ensures consistent processing
// across all items in a dataset.
func SanitizeTextItems(items []TextItem) []TextItem {
	sanitized := make([]TextItem, len(items))
	for i, item := range items {
		sanitized[i] = TextItem{
			ID:       item.ID,
			Content:  SanitizeContent(item.Content),
			Metadata: item.Metadata,
		}
	}
	return sanitized
}

// normalizeWhitespace replaces multiple consecutive spaces with a single space
// but preserves newlines and tabs
func normalizeWhitespace(s string) string {
	var result strings.Builder
	wasSpace := false

	for _, r := range s {
		if r == '\n' || r == '\t' {
			// Preserve newlines and tabs
			result.WriteRune(r)
			wasSpace = false
		} else if unicode.IsSpace(r) {
			if !wasSpace {
				result.WriteRune(' ')
				wasSpace = true
			}
		} else {
			result.WriteRune(r)
			wasSpace = false
		}
	}

	return result.String()
}

// removeNonPrintable removes non-printable characters except newlines and tabs
func removeNonPrintable(s string) string {
	var result strings.Builder

	for _, r := range s {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ValidateAndSanitize provides a complete content preparation pipeline that
// combines sanitization and validation in the optimal order. It first cleans
// the content through sanitization, then validates the cleaned content against
// the provided rules.
//
// This two-phase approach ensures that validation occurs on the actual content
// that will be processed by the LLM, preventing false validation failures due
// to formatting issues that would be resolved during sanitization. The function
// returns both the sanitized content and detailed validation results.
func ValidateAndSanitize(items []TextItem, opts ValidationOptions) ([]TextItem, []ValidationResult, error) {
	// First sanitize
	sanitized := SanitizeTextItems(items)

	// Then validate the sanitized content
	results, err := ValidateTextItems(sanitized, opts)

	return sanitized, results, err
}
