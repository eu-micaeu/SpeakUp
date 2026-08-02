package prompts_test

import (
	"testing"

	"speakup/pkg/prompts"
)

func TestGetPrompt(t *testing.T) {
	files := []string{
		"promptCorrection.txt",
		"promptDialog.txt",
		"promptTopic.txt",
		"promptTranslate.txt",
	}

	for _, file := range files {
		content, err := prompts.GetPrompt(file)
		if err != nil {
			t.Errorf("expected no error loading %s, got: %v", file, err)
		}
		if len(content) == 0 {
			t.Errorf("expected content for %s, got empty string", file)
		}
	}
}
