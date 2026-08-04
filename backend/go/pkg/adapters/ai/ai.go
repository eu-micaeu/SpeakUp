package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"speakup/pkg/adapters/connectors"
	"speakup/pkg/models"
	"speakup/pkg/prompts"
)

func getConnector(ctx context.Context) connectors.AIConnector {
	var provider string
	if ctx != nil {
		if val, ok := ctx.Value("aiProvider").(string); ok && val != "" {
			provider = val
		}
	}
	if provider == "" {
		provider = os.Getenv("AI_PROVIDER")
	}

	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "ollama":
		return connectors.NewOllamaConnector()
	case "gemini":
		return connectors.NewGeminiConnector()
	default:
		if os.Getenv("GEMINI_API_KEY") != "" && strings.EqualFold(os.Getenv("AI_PROVIDER"), "gemini") {
			return connectors.NewGeminiConnector()
		}
		return connectors.NewOllamaConnector()
	}
}

const maxDialogResponseChars = 128
const maxDialogRewriteAttempts = 3

func GetDialogResponse(ctx context.Context, message string, messages []models.Message, language string, level string) (string, error) {
	if level == "" {
		level = "B1"
	}
	prePrompt, err := prompts.GetPrompt("promptDialog.txt")
	if err != nil {
		return "", fmt.Errorf("failed to load prompt: %w", err)
	}

	var chatHistory strings.Builder
	for _, msg := range messages {
		chatHistory.WriteString(fmt.Sprintf("%s: %s\n", msg.Sender, msg.Content))
	}

	connector := getConnector(ctx)

	resumeHist, err := connector.GenerateResponse(ctx, "Format the following chat history to only show user reponse and AI response: "+chatHistory.String())
	if err != nil {
		return "", err
	}

	fullPrompt := fmt.Sprintf("%s\nChat history:\n%s\nATENTION! All Before this point is system instructions and chat history, to generate your response consider the current user message -> = %s\nAnswer me in this language: %s\nUSER LEVEL: %s (Adjust your vocabulary and grammar complexity to match this CEFR level).",
		prePrompt,
		resumeHist,
		message,
		language,
		level)
	fullPrompt += fmt.Sprintf("\nIMPORTANT: Your final answer must be complete, natural, and at most %d characters. Never cut words or end mid-sentence.", maxDialogResponseChars)

	return generateDialogResponseWithinLimit(ctx, connector, fullPrompt, maxDialogResponseChars, level)
}

func GetCorrectionResponse(ctx context.Context, message string) (string, string, error) {
	prompt, err := prompts.GetPrompt("promptCorrection.txt")
	if err != nil {
		return "", "", fmt.Errorf("failed to load prompt: %w", err)
	}

	connector := getConnector(ctx)
	fullPrompt := fmt.Sprintf("%s\n\nINPUT:\n%s\n\nOUTPUT:", prompt, message)

	var correctionResp string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are an English language tutor. Correct the input text and provide a short, clear explanation of the errors in Portuguese. You must format your response exactly with 'Corrected:' and 'Explanation:' prefixes."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 384,
		}
		correctionResp, err = optConnector.GenerateResponseWithOptions(ctx, fullPrompt, systemPrompt, options)
	} else {
		correctionResp, err = connector.GenerateResponse(ctx, fullPrompt)
	}

	if err != nil {
		return "", "", err
	}

	correctionResp = strings.TrimSpace(correctionResp)
	if correctionResp == "" {
		return message, "Nenhuma correção necessária.", nil
	}

	correctedText, explanationText := parseCorrectionResponse(correctionResp)
	if correctedText == "" {
		correctedText = message
	}
	if explanationText == "" {
		explanationText = "Nenhuma explicação adicional fornecida."
	}

	return correctedText, explanationText, nil
}

func GetTranslationResponse(ctx context.Context, message string) (string, error) {
	prompt, err := prompts.GetPrompt("promptTranslate.txt")
	if err != nil {
		return "", fmt.Errorf("failed to load prompt: %w", err)
	}

	connector := getConnector(ctx)
	fullPrompt := fmt.Sprintf("%s\n\nINPUT:\n%s\n\nOUTPUT:", prompt, message)

	var response string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are a strict translation engine. Translate the INPUT text into Brazilian Portuguese and return ONLY the translated text."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 256,
			"stop": []string{
				"\n\n",
				"Explanation:",
				"Explicação:",
			},
		}
		response, err = optConnector.GenerateResponseWithOptions(ctx, fullPrompt, systemPrompt, options)
	} else {
		response, err = connector.GenerateResponse(ctx, fullPrompt)
	}

	if err != nil {
		return "", err
	}

	response = sanitizeTranslationResponse(response)
	if response == "" {
		return "", fmt.Errorf("empty translation response")
	}

	return response, nil
}

func GetTopicResponse(ctx context.Context, message string) (string, error) {
	prompt, err := prompts.GetPrompt("promptTopic.txt")
	if err != nil {
		return "", fmt.Errorf("failed to load prompt: %w", err)
	}

	connector := getConnector(ctx)
	strictPrompt := fmt.Sprintf("%s\n\nInput: %s\nOutput:", prompt, message)

	var topicResp string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are a strict topic labeler. Return exactly two words in Title Case with only letters and one space."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 12,
		}
		topicResp, err = optConnector.GenerateResponseWithOptions(ctx, strictPrompt, systemPrompt, options)
	} else {
		topicResp, err = connector.GenerateResponse(ctx, strictPrompt)
	}

	if err != nil {
		return "", err
	}

	topicResp = strings.TrimSpace(topicResp)
	if topicResp == "" {
		topicResp = "New Topic"
	}

	return topicResp, nil
}

func normalizeDialogResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return ""
	}

	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}

func generateDialogResponseWithinLimit(ctx context.Context, connector connectors.AIConnector, prompt string, maxChars int, level string) (string, error) {
	if maxChars <= 0 {
		return "", fmt.Errorf("invalid dialog response limit")
	}

	currentPrompt := prompt
	options := map[string]any{
		"temperature": 0.3,
		"top_p":       0.9,
		"num_predict": 96,
	}
	if level == "" {
		level = "B1"
	}
	systemPrompt := fmt.Sprintf("You are a natural language exchange partner. The user has a language proficiency of %s. Keep coherence with chat context and user message. Return exactly one complete answer with at most %d characters. Adjust your language to be appropriate for a %s learner.", level, maxChars, level)

	for attempt := 0; attempt < maxDialogRewriteAttempts; attempt++ {
		var (
			rawResp string
			err     error
		)

		if optConnector, ok := connector.(connectors.OptionableConnector); ok {
			rawResp, err = optConnector.GenerateResponseWithOptions(ctx, currentPrompt, systemPrompt, options)
		} else {
			rawResp, err = connector.GenerateResponse(ctx, currentPrompt)
		}
		if err != nil {
			return "", err
		}

		cleaned := normalizeDialogResponse(rawResp)
		if cleaned == "" {
			currentPrompt = prompt
			continue
		}

		if utf8.RuneCountInString(cleaned) <= maxChars {
			return cleaned, nil
		}

		currentPrompt = fmt.Sprintf(
			"%s\n\nIMPORTANT: Your previous answer exceeded %d characters. Answer again with one complete, natural sentence up to %d characters, keeping the same context and language.",
			prompt,
			maxChars,
			maxChars,
		)
	}

	return "", fmt.Errorf("failed to generate response within %d characters", maxChars)
}

func sanitizeTranslationResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return ""
	}

	labels := []string{
		"Translation:",
		"Tradução:",
		"Translated text:",
		"OUTPUT:",
		"Output:",
	}

	lower := strings.ToLower(cleaned)
	for _, label := range labels {
		labelLower := strings.ToLower(label)
		if strings.HasPrefix(lower, labelLower) {
			cleaned = strings.TrimSpace(cleaned[len(label):])
			lower = strings.ToLower(cleaned)
			break
		}
	}

	for _, marker := range []string{"\ntranslation:", "\ntradução:", " translation:", " tradução:"} {
		if idx := strings.Index(strings.ToLower(cleaned), marker); idx > 0 {
			cleaned = strings.TrimSpace(cleaned[:idx])
			break
		}
	}

	cleaned = strings.TrimSpace(strings.Trim(cleaned, "\"'`"))
	return cleaned
}

func parseCorrectionResponse(raw string) (string, string) {
	lines := strings.Split(raw, "\n")
	var corrected, explanation []string
	isExplanation := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(trimmed), "corrected:") {
			corrected = append(corrected, strings.TrimSpace(line[len("corrected:"):]))
			isExplanation = false
		} else if strings.HasPrefix(strings.ToLower(trimmed), "explanation:") {
			explanation = append(explanation, strings.TrimSpace(line[len("explanation:"):]))
			isExplanation = true
		} else if strings.HasPrefix(strings.ToLower(trimmed), "explicação:") {
			explanation = append(explanation, strings.TrimSpace(line[len("explicação:"):]))
			isExplanation = true
		} else {
			if isExplanation {
				explanation = append(explanation, trimmed)
			} else {
				corrected = append(corrected, trimmed)
			}
		}
	}

	correctedText := strings.TrimSpace(strings.Join(corrected, "\n"))
	explanationText := strings.TrimSpace(strings.Join(explanation, "\n"))

	// Fallback if no tags were matched
	if correctedText == "" && explanationText == "" {
		correctedText = strings.TrimSpace(raw)
	}

	return correctedText, explanationText
}
