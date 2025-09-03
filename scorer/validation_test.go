package scorer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/JohnPlummer/llm-client/scorer"
)

var _ = Describe("Validation", func() {
	Describe("ValidateContent", func() {
		var opts scorer.ValidationOptions

		BeforeEach(func() {
			opts = scorer.DefaultValidationOptions()
		})

		Context("with empty content", func() {
			It("should fail validation by default", func() {
				result := scorer.ValidateContent("", opts)
				Expect(result.Valid).To(BeFalse())
				Expect(result.Issues).To(ContainElement("content is empty"))
			})

			It("should pass when AllowEmpty is true", func() {
				opts.AllowEmpty = true
				result := scorer.ValidateContent("", opts)
				Expect(result.Valid).To(BeTrue())
			})
		})

		Context("with whitespace-only content", func() {
			It("should fail validation by default", func() {
				result := scorer.ValidateContent("   \t\n  ", opts)
				Expect(result.Valid).To(BeFalse())
				Expect(result.Issues).To(ContainElement("content contains only whitespace"))
			})

			It("should pass when AllowWhitespace is true", func() {
				opts.AllowWhitespace = true
				result := scorer.ValidateContent("   ", opts)
				Expect(result.Valid).To(BeTrue())
			})
		})

		Context("with content length validation", func() {
			It("should fail when content is too short", func() {
				opts.MinLength = 10
				result := scorer.ValidateContent("short", opts)
				Expect(result.Valid).To(BeFalse())
				Expect(result.Issues[0]).To(ContainSubstring("content too short"))
			})

			It("should fail when content is too long", func() {
				opts.MaxLength = 10
				result := scorer.ValidateContent("this is a very long content", opts)
				Expect(result.Valid).To(BeFalse())
				Expect(result.Issues[0]).To(ContainSubstring("content too long"))
			})

			It("should pass with valid length content", func() {
				opts.MinLength = 5
				opts.MaxLength = 20
				result := scorer.ValidateContent("valid content", opts)
				Expect(result.Valid).To(BeTrue())
				Expect(result.Issues).To(BeEmpty())
			})
		})

		Context("with TrimWhitespace option", func() {
			It("should validate trimmed length when enabled", func() {
				opts.MaxLength = 10
				opts.TrimWhitespace = true
				result := scorer.ValidateContent("  content  ", opts)
				Expect(result.Valid).To(BeTrue()) // "content" is 7 chars after trim
			})

			It("should validate full length when disabled", func() {
				opts.MaxLength = 10
				opts.TrimWhitespace = false
				result := scorer.ValidateContent("  content  ", opts)
				Expect(result.Valid).To(BeFalse()) // "  content  " is 11 chars
			})
		})
	})

	Describe("ValidateTextItems", func() {
		var opts scorer.ValidationOptions

		BeforeEach(func() {
			opts = scorer.DefaultValidationOptions()
		})

		It("should return error for empty items", func() {
			_, err := scorer.ValidateTextItems([]scorer.TextItem{}, opts)
			Expect(err).To(Equal(scorer.ErrEmptyInput))
		})

		It("should validate all items in batch", func() {
			items := []scorer.TextItem{
				{ID: "1", Content: "valid content"},
				{ID: "2", Content: ""},
				{ID: "3", Content: "another valid"},
			}

			results, err := scorer.ValidateTextItems(items, opts)
			Expect(err).To(HaveOccurred())
			Expect(results).To(HaveLen(3))
			Expect(results[0].Valid).To(BeTrue())
			Expect(results[1].Valid).To(BeFalse())
			Expect(results[2].Valid).To(BeTrue())
		})

		It("should fail items with empty IDs", func() {
			items := []scorer.TextItem{
				{ID: "", Content: "content"},
			}

			results, err := scorer.ValidateTextItems(items, opts)
			Expect(err).To(HaveOccurred())
			Expect(results[0].Valid).To(BeFalse())
			Expect(results[0].Issues).To(ContainElement("item ID is empty"))
		})
	})

	Describe("SanitizeContent", func() {
		It("should trim whitespace", func() {
			result := scorer.SanitizeContent("  content  ")
			Expect(result).To(Equal("content"))
		})

		It("should normalize multiple spaces", func() {
			result := scorer.SanitizeContent("too    many     spaces")
			Expect(result).To(Equal("too many spaces"))
		})

		It("should remove non-printable characters", func() {
			result := scorer.SanitizeContent("content\x00with\x01control\x02chars")
			Expect(result).To(Equal("contentwithcontrolchars"))
		})

		It("should preserve newlines and tabs", func() {
			result := scorer.SanitizeContent("line1\nline2\ttabbed")
			Expect(result).To(Equal("line1\nline2\ttabbed"))
		})
	})

	Describe("SanitizeTextItems", func() {
		It("should sanitize all items", func() {
			items := []scorer.TextItem{
				{ID: "1", Content: "  spaced  "},
				{ID: "2", Content: "multiple    spaces"},
			}

			sanitized := scorer.SanitizeTextItems(items)
			Expect(sanitized[0].Content).To(Equal("spaced"))
			Expect(sanitized[1].Content).To(Equal("multiple spaces"))
		})

		It("should preserve metadata", func() {
			items := []scorer.TextItem{
				{ID: "1", Content: "  content  ", Metadata: map[string]interface{}{"key": "value"}},
			}

			sanitized := scorer.SanitizeTextItems(items)
			Expect(sanitized[0].Metadata).To(HaveKeyWithValue("key", "value"))
		})
	})

	Describe("ValidateAndSanitize", func() {
		It("should sanitize then validate", func() {
			opts := scorer.DefaultValidationOptions()
			opts.MaxLength = 10

			items := []scorer.TextItem{
				{ID: "1", Content: "  short  "}, // Will be "short" after sanitization
			}

			sanitized, results, err := scorer.ValidateAndSanitize(items, opts)
			Expect(err).ToNot(HaveOccurred())
			Expect(sanitized[0].Content).To(Equal("short"))
			Expect(results[0].Valid).To(BeTrue())
		})
	})
})
