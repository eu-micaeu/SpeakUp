package prompts

import (
	"embed"
	"fmt"
)

//go:embed *.txt
var FS embed.FS

// GetPrompt reads an embedded prompt file by name.
func GetPrompt(name string) (string, error) {
	b, err := FS.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("failed to load embedded prompt %s: %w", name, err)
	}
	return string(b), nil
}
