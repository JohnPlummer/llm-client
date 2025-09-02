package scorer

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidationResult contains the results of content validation
type ValidationResult struct {
	Valid   bool
	Issues  []string
	Suggestions []string
}

// ValidationOptions configures content validation behavior
type ValidationOptions struct {
	MaxLength        int
	MinLength        int
	AllowEmpty       bool
	AllowWhitespace  bool
	TrimWhitespace   bool
}

// DefaultValidationOptions returns sensible defaults for content validation
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		MaxLength:       DefaultMaxContentLength,
		MinLength:       MinContentLength,
		AllowEmpty:      false,
		AllowWhitespace: false,
		TrimWhitespace:  true,
	}
}

// ValidateContent validates a single text item's content
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

// ValidateTextItems validates a batch of text items
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
				Valid:  false,
				Issues: []string{"item ID is empty"},
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

// SanitizeContent cleans and normalizes text content
func SanitizeContent(content string) string {
	// Trim leading and trailing whitespace
	content = strings.TrimSpace(content)
	
	// Normalize whitespace (replace multiple spaces with single space)
	content = normalizeWhitespace(content)
	
	// Remove non-printable characters except newlines and tabs
	content = removeNonPrintable(content)
	
	return content
}

// SanitizeTextItems sanitizes content for a batch of text items
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

// ValidateAndSanitize performs both validation and sanitization
func ValidateAndSanitize(items []TextItem, opts ValidationOptions) ([]TextItem, []ValidationResult, error) {
	// First sanitize
	sanitized := SanitizeTextItems(items)
	
	// Then validate the sanitized content
	results, err := ValidateTextItems(sanitized, opts)
	
	return sanitized, results, err
}