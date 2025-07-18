package scorer

import (
	"fmt"
)

const maxBatchSize = 10

var batchScorePrompt string
var batchPromptError error

func init() {
	// Load batch prompt during package initialization
	promptBytes, err := promptFS.ReadFile("prompts/batch_prompt.txt")
	if err != nil {
		batchPromptError = fmt.Errorf("failed to load batch prompt: %w", err)
		return
	}
	batchScorePrompt = string(promptBytes)
}
