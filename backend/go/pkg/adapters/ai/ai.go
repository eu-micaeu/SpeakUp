package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
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
	default:
		return connectors.NewOllamaConnector()
	}
}

const maxDialogResponseChars = 128
const maxDialogRewriteAttempts = 3
const maxChatHistoryMessages = 10

func cleanStoredMessageContent(msg models.Message) string {
	content := strings.TrimSpace(msg.Content)
	if content == "" {
		return ""
	}

	if strings.EqualFold(msg.Sender, "user") {
		// In MongoDB, user messages are stored as "<raw>\n\nCorreção: <corrected>"
		if idx := strings.Index(content, "\n\nCorreção:"); idx != -1 {
			corrected := strings.TrimSpace(content[idx+len("\n\nCorreção:"):])
			if corrected != "" {
				return corrected
			}
			return strings.TrimSpace(content[:idx])
		}
		return content
	}

	// In MongoDB, AI messages are stored as "<dialog>\n\n[TRANSLATION]: <pt>"
	if idx := strings.Index(content, "\n\n[TRANSLATION]:"); idx != -1 {
		return strings.TrimSpace(content[:idx])
	}
	return content
}

func formatChatHistory(messages []models.Message, currentMessage string) string {
	if len(messages) == 0 {
		return "(No previous messages)"
	}

	trimmedCurrent := strings.TrimSpace(currentMessage)

	// Exclude the last stored user message if it duplicates currentMessage
	endIdx := len(messages)
	if endIdx > 0 && strings.EqualFold(messages[endIdx-1].Sender, "user") {
		lastClean := cleanStoredMessageContent(messages[endIdx-1])
		if strings.EqualFold(lastClean, trimmedCurrent) || strings.Contains(messages[endIdx-1].Content, trimmedCurrent) {
			endIdx--
		}
	}

	relevant := messages[:endIdx]
	if len(relevant) == 0 {
		return "(No previous messages)"
	}
	if len(relevant) > maxChatHistoryMessages {
		relevant = relevant[len(relevant)-maxChatHistoryMessages:]
	}

	var b strings.Builder
	for _, msg := range relevant {
		clean := cleanStoredMessageContent(msg)
		if clean == "" {
			continue
		}
		role := "User"
		if strings.EqualFold(msg.Sender, "ai") || strings.EqualFold(msg.Sender, "assistant") {
			role = "Output"
		}
		b.WriteString(fmt.Sprintf("%s: %s\n", role, clean))
	}

	res := strings.TrimSpace(b.String())
	if res == "" {
		return "(No previous messages)"
	}
	return res
}

func GetDialogResponse(ctx context.Context, message string, messages []models.Message, language string, level string) (string, error) {
	if strings.TrimSpace(language) == "" {
		language = "English"
	}
	if strings.TrimSpace(level) == "" {
		level = "B1"
	}
	prePrompt, err := prompts.GetPrompt("promptDialog.txt")
	if err != nil {
		return "", fmt.Errorf("failed to load prompt: %w", err)
	}

	connector := getConnector(ctx)
	historyText := formatChatHistory(messages, message)

	fullPrompt := fmt.Sprintf(
		"%s\n\nCONTEXT:\n- Target Language: %s\n- Learner CEFR Level: %s\n- Max Length: %d characters\n\nCHAT HISTORY:\n%s\n\nCURRENT USER MESSAGE:\nUser: %s\nOutput:",
		prePrompt,
		language,
		level,
		maxDialogResponseChars,
		historyText,
		strings.TrimSpace(message),
	)

	return generateDialogResponseWithinLimit(ctx, connector, fullPrompt, maxDialogResponseChars, level)
}

func GetCorrectionResponse(ctx context.Context, message string) (string, string, error) {
	prompt, err := prompts.GetPrompt("promptCorrection.txt")
	if err != nil {
		return "", "", fmt.Errorf("failed to load prompt: %w", err)
	}

	connector := getConnector(ctx)
	fullPrompt := fmt.Sprintf("%s\n\nINPUT:\n%s\n\nOUTPUT:", prompt, strings.TrimSpace(message))

	var correctionResp string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are a strict language proofreader. Do NOT answer questions or continue conversations. Correct the INPUT text and explain the corrections in Brazilian Portuguese using the exact prefixes 'Corrected:' and 'Explanation:'."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 384,
			"stop": []string{
				"\nINPUT:",
				"\nInput:",
			},
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
	fullPrompt := fmt.Sprintf("%s\n\nINPUT:\n%s\n\nOUTPUT:", prompt, strings.TrimSpace(message))

	var response string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are a strict translation engine. Translate the INPUT text into natural Brazilian Portuguese and return ONLY the translated text. Never answer questions or add labels."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 256,
			"stop": []string{
				"\n\n",
				"\nINPUT:",
				"\nInput:",
				"Explanation:",
				"Explicação:",
				"Note:",
				"Nota:",
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
	strictPrompt := fmt.Sprintf("%s\n\nInput: %s\nOutput:", prompt, strings.TrimSpace(message))

	var topicResp string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are a strict topic labeler. Return ONLY two words in English Title Case (e.g., 'Travel Memories') with no punctuation or labels."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 12,
			"stop": []string{
				"\n",
				"Input:",
				"INPUT:",
				"Output:",
				"OUTPUT:",
			},
		}
		topicResp, err = optConnector.GenerateResponseWithOptions(ctx, strictPrompt, systemPrompt, options)
	} else {
		topicResp, err = connector.GenerateResponse(ctx, strictPrompt)
	}

	if err != nil {
		return "", err
	}

	return sanitizeTopicResponse(topicResp), nil
}

func sanitizeTopicResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return "New Topic"
	}

	// Keep only the first non-empty line
	for _, line := range strings.Split(cleaned, "\n") {
		if t := strings.TrimSpace(line); t != "" {
			cleaned = t
			break
		}
	}

	cleaned = strings.TrimLeft(cleaned, "*#-_ \t")
	for _, prefix := range []string{"output:", "topic:", "category:", "tópico:", "theme:"} {
		if strings.HasPrefix(strings.ToLower(cleaned), prefix) {
			cleaned = strings.TrimSpace(cleaned[len(prefix):])
		}
	}

	// Keep only letters, digits, and spaces
	var filtered strings.Builder
	for _, r := range cleaned {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			filtered.WriteRune(r)
		} else if r == '-' || r == '/' {
			filtered.WriteRune(' ')
		}
	}

	words := strings.Fields(filtered.String())
	if len(words) == 0 {
		return "New Topic"
	}
	if len(words) > 3 {
		words = words[:2]
	}

	for i, w := range words {
		runes := []rune(strings.ToLower(w))
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[i] = string(runes)
	}

	res := strings.Join(words, " ")
	if utf8.RuneCountInString(res) > 24 && len(words) > 1 {
		res = strings.Join(words[:2], " ")
	}
	if strings.TrimSpace(res) == "" {
		return "New Topic"
	}
	return res
}

func normalizeDialogResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return ""
	}

	// Strip leaked translation blocks or subsequent user turns if any
	for _, cutMarker := range []string{"\n[TRANSLATION]", "[TRANSLATION]:", "\nUser:", "\nInput:"} {
		if idx := strings.Index(cleaned, cutMarker); idx != -1 {
			cleaned = strings.TrimSpace(cleaned[:idx])
		}
	}

	// Remove common speaker prefixes
	cleaned = strings.TrimLeft(cleaned, "*_ \t")
	for _, prefix := range []string{"output:", "you:", "speakup:", "assistant:", "ai:", "response:"} {
		if strings.HasPrefix(strings.ToLower(cleaned), prefix) {
			cleaned = strings.TrimSpace(cleaned[len(prefix):])
			break
		}
	}

	cleaned = strings.TrimSpace(strings.Trim(cleaned, "\"'`"))
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}

func truncateDialogResponse(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	runes := []rune(text)
	window := string(runes[:maxChars])

	// 1. Prefer cutting at the last complete sentence punctuation (. ! ?)
	lastPunct := -1
	for i, r := range window {
		if r == '.' || r == '!' || r == '?' {
			lastPunct = i
		}
	}
	if lastPunct >= maxChars/3 {
		return strings.TrimSpace(window[:lastPunct+1])
	}

	// 2. Otherwise cut cleanly at the last word boundary and ensure ending punctuation
	lastSpace := strings.LastIndex(window, " ")
	if lastSpace > 0 {
		trimmed := strings.TrimRight(strings.TrimSpace(window[:lastSpace]), ",;:-")
		if trimmed != "" {
			lastRune, _ := utf8.DecodeLastRuneInString(trimmed)
			if lastRune != '.' && lastRune != '!' && lastRune != '?' {
				if utf8.RuneCountInString(trimmed) < maxChars {
					trimmed += "."
				}
			}
			return trimmed
		}
	}

	return strings.TrimSpace(window)
}

func generateDialogResponseWithinLimit(ctx context.Context, connector connectors.AIConnector, prompt string, maxChars int, level string) (string, error) {
	if maxChars <= 0 {
		return "", fmt.Errorf("invalid dialog response limit")
	}

	currentPrompt := prompt
	options := map[string]any{
		"temperature": 0.3,
		"top_p":       0.9,
		"num_predict": 80,
		"stop": []string{
			"\nUser:",
			"\nYou:",
			"\nInput:",
			"\n[TRANSLATION]",
		},
	}
	if level == "" {
		level = "B1"
	}
	systemPrompt := fmt.Sprintf("You are SpeakUp, a friendly language exchange partner. The user's CEFR level is %s. Reply in the target language only, using 1-2 short sentences (at most %d characters total). Output only your reply without labels.", level, maxChars)

	var lastCleaned string
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
		lastCleaned = cleaned

		if utf8.RuneCountInString(cleaned) <= maxChars {
			return cleaned, nil
		}

		currentPrompt = fmt.Sprintf(
			"%s\n\nIMPORTANT: Your previous draft was %d characters, which exceeds the %d-character limit. Rewrite it into ONE short sentence + ONE short question under %d characters.",
			prompt,
			utf8.RuneCountInString(cleaned),
			maxChars,
			maxChars,
		)
	}

	if lastCleaned != "" {
		return truncateDialogResponse(lastCleaned, maxChars), nil
	}

	return "", fmt.Errorf("failed to generate response within %d characters", maxChars)
}

func sanitizeTranslationResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return ""
	}

	cleaned = strings.TrimLeft(cleaned, "*#-_ \t")
	labels := []string{
		"Translation:",
		"Tradução:",
		"Traducao:",
		"Translated text:",
		"OUTPUT:",
		"Output:",
	}

	lower := strings.ToLower(cleaned)
	for _, label := range labels {
		labelLower := strings.ToLower(label)
		if strings.HasPrefix(lower, labelLower) {
			cleaned = strings.TrimSpace(cleaned[len(label):])
			cleaned = strings.TrimLeft(cleaned, "*_ \t")
			lower = strings.ToLower(cleaned)
			break
		}
	}

	for _, marker := range []string{
		"\ntranslation:",
		"\ntradução:",
		"\nnota:",
		"\nnote:",
		"\nobs:",
		"\nexplicação:",
		"\nexplanation:",
		" translation:",
		" tradução:",
	} {
		if idx := strings.Index(strings.ToLower(cleaned), marker); idx > 0 {
			cleaned = strings.TrimSpace(cleaned[:idx])
			break
		}
	}

	cleaned = strings.TrimSpace(strings.Trim(cleaned, "\"'`"))
	return cleaned
}

func stripMarkdownLabel(line string) string {
	s := strings.TrimSpace(line)
	s = strings.TrimLeft(s, "*-#> \t")
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "__", "")
	return strings.TrimSpace(s)
}

func parseCorrectionResponse(raw string) (string, string) {
	lines := strings.Split(raw, "\n")
	var corrected, explanation []string
	isExplanation := false

	correctedPrefixes := []string{"corrected:", "correction:", "corrigido:", "correção:", "correcao:"}
	explanationPrefixes := []string{"explanation:", "explicação:", "explicacao:"}

	for _, line := range lines {
		normalized := stripMarkdownLabel(line)
		if normalized == "" {
			continue
		}

		lower := strings.ToLower(normalized)

		// Check if both Corrected: and Explanation: appear on the same line
		matchedBoth := false
		for _, cp := range correctedPrefixes {
			if strings.HasPrefix(lower, cp) {
				rest := strings.TrimSpace(normalized[len(cp):])
				restLower := strings.ToLower(rest)
				for _, ep := range explanationPrefixes {
					if idx := strings.Index(restLower, ep); idx != -1 {
						cPart := strings.TrimSpace(rest[:idx])
						ePart := strings.TrimSpace(rest[idx+len(ep):])
						if cPart != "" {
							corrected = append(corrected, cPart)
						}
						if ePart != "" {
							explanation = append(explanation, ePart)
						}
						isExplanation = true
						matchedBoth = true
						break
					}
				}
				if !matchedBoth {
					if rest != "" {
						corrected = append(corrected, rest)
					}
					isExplanation = false
					matchedBoth = true
				}
				break
			}
		}
		if matchedBoth {
			continue
		}

		matchedExpl := false
		for _, ep := range explanationPrefixes {
			if strings.HasPrefix(lower, ep) {
				rest := strings.TrimSpace(normalized[len(ep):])
				if rest != "" {
					explanation = append(explanation, rest)
				}
				isExplanation = true
				matchedExpl = true
				break
			}
		}
		if matchedExpl {
			continue
		}

		if isExplanation {
			explanation = append(explanation, normalized)
		} else {
			corrected = append(corrected, normalized)
		}
	}

	// Join corrected with spaces so it never introduces "\n\n" which would break frontend splitting
	correctedText := strings.TrimSpace(strings.Join(corrected, " "))
	correctedText = strings.TrimSpace(strings.Trim(correctedText, "\"'`"))
	explanationText := strings.TrimSpace(strings.Join(explanation, " "))

	// Fallback if no tags were matched
	if correctedText == "" && explanationText == "" {
		correctedText = strings.TrimSpace(strings.Join(strings.Fields(raw), " "))
	}

	return correctedText, explanationText
}

type GeneratedFlashcardData struct {
	Back            string `json:"back"`
	ContextSentence string `json:"context_sentence"`
	Explanation     string `json:"explanation"`
}

func extractJSONObject(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	start := strings.Index(cleaned, "{")
	end := strings.LastIndex(cleaned, "}")
	if start != -1 && end != -1 && end > start {
		return cleaned[start : end+1]
	}
	return cleaned
}

func GetFlashcardGenerateResponse(ctx context.Context, term string, contextSentence string) (GeneratedFlashcardData, error) {
	var result GeneratedFlashcardData

	prompt, err := prompts.GetPrompt("promptFlashcard.txt")
	if err != nil {
		return result, fmt.Errorf("failed to load prompt: %w", err)
	}

	connector := getConnector(ctx)
	fullPrompt := fmt.Sprintf("%s\n\nWORD/PHRASE:\n%s\n\nCONTEXT SENTENCE:\n%s\n\nOUTPUT (JSON ONLY):", prompt, strings.TrimSpace(term), strings.TrimSpace(contextSentence))

	var resp string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are a strict language flashcard generator. Output a raw JSON object only with keys 'back', 'context_sentence', and 'explanation'."
		options := map[string]any{
			"format":      "json",
			"temperature": 0.2,
			"top_p":       0.9,
			"num_predict": 256,
		}
		resp, err = optConnector.GenerateResponseWithOptions(ctx, fullPrompt, systemPrompt, options)
	} else {
		resp, err = connector.GenerateResponse(ctx, fullPrompt)
	}

	if err != nil {
		return result, err
	}

	cleaned := extractJSONObject(resp)
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return result, fmt.Errorf("failed to parse flashcard JSON: %w", err)
	}

	result.Back = FormatSingleWordTranslation(result.Back)
	if result.Back == "" {
		return result, fmt.Errorf("empty flashcard back translation")
	}

	return result, nil
}

func FormatSingleWordTranslation(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, " \t\n\r\"'.,;:!?()-")
	if s == "" {
		return ""
	}
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '/' || r == ',' || r == ';'
	})
	if len(fields) == 0 {
		return ""
	}

	targetWord := fields[0]
	// Skip common articles/particles if there is a second word
	lowerFirst := strings.ToLower(strings.Trim(fields[0], " \t\n\r\"'.,;:!?()-"))
	if len(fields) > 1 && (lowerFirst == "a" || lowerFirst == "o" || lowerFirst == "um" || lowerFirst == "uma" || lowerFirst == "to") {
		targetWord = fields[1]
	}

	word := strings.TrimSpace(targetWord)
	word = strings.Trim(word, " \t\n\r\"'.,;:!?()-")
	if word == "" {
		return ""
	}

	runes := []rune(strings.ToLower(word))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

type BatchFlashcardItem struct {
	Front           string `json:"front"`
	Back            string `json:"back"`
	ContextSentence string `json:"context_sentence"`
	Explanation     string `json:"explanation"`
}

var fallbackFlashcardsPool = []BatchFlashcardItem{
	{Front: "Resilient", Back: "Resiliente", ContextSentence: "She remained resilient despite facing multiple challenges.", Explanation: "Descreve a capacidade de se recuperar rapidamente de dificuldades."},
	{Front: "Empathy", Back: "Empatia", ContextSentence: "Showing empathy helps build stronger relationships.", Explanation: "A capacidade de se colocar no lugar de outra pessoa e entender seus sentimentos."},
	{Front: "Overcome", Back: "Superar", ContextSentence: "He worked hard to overcome his fear of public speaking.", Explanation: "Verbo que significa vencer ou ultrapassar uma dificuldade ou obstáculo."},
	{Front: "Thrive", Back: "Prosperar", ContextSentence: "Plants thrive when they receive enough sunlight and water.", Explanation: "Significa crescer, desenvolver-se bem ou ter muito sucesso."},
	{Front: "Insight", Back: "Discernimento", ContextSentence: "Her research provided valuable insights into human behavior.", Explanation: "Compreensão profunda ou intuitiva sobre algo complexo."},
	{Front: "Endeavor", Back: "Esforço", ContextSentence: "Starting a new business is a challenging endeavor.", Explanation: "Um esforço consciente, empreendimento ou tentativa séria de alcançar um objetivo."},
	{Front: "Reluctant", Back: "Hesitante", ContextSentence: "He was reluctant to share his opinion in the meeting.", Explanation: "Pessoa que reluta ou hesita em fazer algo por dúvida ou receio."},
	{Front: "Plausible", Back: "Plausível", ContextSentence: "Her explanation sounded plausible and well reasoned.", Explanation: "Algo que parece provável, aceitável ou razoável."},
	{Front: "Meticulous", Back: "Meticuloso", ContextSentence: "He paid meticulous attention to every detail of the project.", Explanation: "Pessoa cuidadosa, detalhista e extremamente precisa."},
	{Front: "Proactive", Back: "Proativo", ContextSentence: "Being proactive helps you solve problems before they happen.", Explanation: "Tomar a iniciativa para antecipar problemas e buscar soluções."},
	{Front: "Perseverance", Back: "Perseverança", ContextSentence: "Success comes to those who show great perseverance.", Explanation: "Persistência em alcançar um objetivo apesar de dificuldades."},
	{Front: "Candid", Back: "Sincero", ContextSentence: "She gave a candid feedback during the performance review.", Explanation: "Franco, direto e honesto sem segundas intenções."},
	{Front: "Versatile", Back: "Versátil", ContextSentence: "This tool is versatile and can be used for many tasks.", Explanation: "Capaz de se adaptar facilmente a diferentes funções ou situações."},
	{Front: "Inevitable", Back: "Inevitável", ContextSentence: "Change is an inevitable part of life and growth.", Explanation: "Algo que certamente vai acontecer e não pode ser evitado."},
	{Front: "Pragmatic", Back: "Pragmático", ContextSentence: "We need a pragmatic approach to solve this issue quickly.", Explanation: "Focado em soluções práticas e no realismo, mais do que em teorias."},
	{Front: "Fostering", Back: "Fomentando", ContextSentence: "The teacher is dedicated to fostering creativity in students.", Explanation: "Incentivar, nutrir ou promover o desenvolvimento de algo."},
	{Front: "Scrutinize", Back: "Examinar", ContextSentence: "Auditors will scrutinize the company financial records.", Explanation: "Examinar ou inspecionar detalhadamente e criticamente."},
	{Front: "Enhance", Back: "Aprimorar", ContextSentence: "Regular exercise can enhance your overall health.", Explanation: "Melhorar, aumentar a qualidade ou valor de algo."},
	{Front: "Comprehensive", Back: "Abrangente", ContextSentence: "The guide offers a comprehensive overview of the subject.", Explanation: "Amplo, completo, que inclui todos ou quase todos os detalhes."},
	{Front: "Ambiguous", Back: "Ambíguo", ContextSentence: "His answer was ambiguous and left room for interpretation.", Explanation: "Com duplo sentido, que não é claro ou pode ser interpretado de várias formas."},
	{Front: "Diligent", Back: "Diligente", ContextSentence: "The diligent student completed all assignments ahead of time.", Explanation: "Trabalhador, dedicado e cuidadoso nas suas tarefas."},
	{Front: "Feasible", Back: "Viável", ContextSentence: "It is not feasible to complete this project in one day.", Explanation: "Possível de ser feito ou realizado com sucesso."},
	{Front: "Substantial", Back: "Substancial", ContextSentence: "They made substantial progress on the new software update.", Explanation: "De grande valor, tamanho ou importância."},
	{Front: "Coherent", Back: "Coerente", ContextSentence: "She presented a coherent argument during the debate.", Explanation: "Lógico, bem estruturado e consistente."},
	{Front: "Nurture", Back: "Nutrir", ContextSentence: "Parents work hard to nurture their children talents.", Explanation: "Alimentar, cuidar e apoiar o crescimento de algo ou alguém."},
	{Front: "Cohesive", Back: "Coeso", ContextSentence: "The team worked well together to form a cohesive unit.", Explanation: "Unido, que se mantém firme e integrado."},
	{Front: "Elaborate", Back: "Detalhado", ContextSentence: "He gave an elaborate presentation about the project scope.", Explanation: "Bem trabalhado, com muitos detalhes e explicações."},
	{Front: "Reconcile", Back: "Reconciliar", ContextSentence: "They managed to reconcile their differences and work together.", Explanation: "Resolver conflitos ou fazer duas coisas opostas se harmonizarem."},
	{Front: "Spontaneous", Back: "Espontâneo", ContextSentence: "They took a spontaneous trip to the mountains over the weekend.", Explanation: "Feito de forma natural, sem planejamento prévio."},
	{Front: "Paramount", Back: "Fundamental", ContextSentence: "Safety is of paramount importance during construction.", Explanation: "De máxima prioridade, mais importante do que qualquer outra coisa."},
	{Front: "Acknowledge", Back: "Reconhecer", ContextSentence: "It is important to acknowledge your mistakes and learn from them.", Explanation: "Admitir a existência, verdade ou validade de algo."},
	{Front: "Advocate", Back: "Defender", ContextSentence: "She continues to advocate for better healthcare services.", Explanation: "Apoiar publicamente ou defender uma causa ou ideia."},
	{Front: "Assess", Back: "Avaliar", ContextSentence: "We need to assess the risks before making a final decision.", Explanation: "Analisar e julgar o valor, importância ou qualidade de algo."},
	{Front: "Collaborate", Back: "Colaborar", ContextSentence: "The two teams will collaborate to complete the project on time.", Explanation: "Trabalhar junto com outras pessoas para alcançar um objetivo comum."},
	{Front: "Formulate", Back: "Formular", ContextSentence: "They worked hard to formulate a new strategy for growth.", Explanation: "Desenvolver ou expressar uma ideia de forma clara e estruturada."},
	{Front: "Implement", Back: "Implementar", ContextSentence: "The company plans to implement new security policies next month.", Explanation: "Colocar um plano, decisão ou sistema em prática."},
	{Front: "Integrate", Back: "Integrar", ContextSentence: "The software allows you to integrate different applications easily.", Explanation: "Combinar duas ou mais coisas para funcionarem juntas."},
	{Front: "Mitigate", Back: "Atenuar", ContextSentence: "Steps were taken to mitigate the environmental impact.", Explanation: "Tornar algo menos severo, sério ou doloroso."},
	{Front: "Navigate", Back: "Navegar", ContextSentence: "He managed to navigate through the complex legal rules.", Explanation: "Guiar-se ou encontrar o caminho através de uma situação difícil."},
	{Front: "Optimize", Back: "Otimizar", ContextSentence: "We must optimize our workflow to save time and resources.", Explanation: "Fazer o melhor uso possível de algo para obter o melhor resultado."},
	{Front: "Perceive", Back: "Perceber", ContextSentence: "How people perceive a situation depends on their experience.", Explanation: "Notar, entender ou interpretar algo através dos sentidos ou da mente."},
	{Front: "Prioritize", Back: "Priorizar", ContextSentence: "You should prioritize urgent tasks before moving on to others.", Explanation: "Tratar algo como mais importante do que outras coisas."},
	{Front: "Reiterate", Back: "Reiterar", ContextSentence: "I want to reiterate that attendance is mandatory.", Explanation: "Dizer algo novamente para enfatizar ou esclarecer."},
	{Front: "Strengthen", Back: "Fortalecer", ContextSentence: "Daily practice will strengthen your language skills.", Explanation: "Tornar algo mais forte, resistente ou eficaz."},
	{Front: "Sustain", Back: "Sustentar", ContextSentence: "The country needs new policies to sustain economic growth.", Explanation: "Manter algo funcionando ou existindo ao longo do tempo."},
	{Front: "Transform", Back: "Transformar", ContextSentence: "Technology has the power to transform how we work.", Explanation: "Mudar completamente a forma, aparência ou caráter de algo."},
	{Front: "Validate", Back: "Validar", ContextSentence: "The results of the test validate our original theory.", Explanation: "Confirmar ou provar que algo é verdadeiro e correto."},
	{Front: "Vulnerable", Back: "Vulnerável", ContextSentence: "Without protection, the system remains vulnerable to attacks.", Explanation: "Suscetível a danos, críticas ou ataques."},
	{Front: "Yield", Back: "Render", ContextSentence: "The new investment strategy is expected to yield high returns.", Explanation: "Produzir ou fornecer um resultado, lucro ou cultivo."},
	{Front: "Zealous", Back: "Zeloso", ContextSentence: "The zealous volunteer worked tirelessly for the community.", Explanation: "Cheio de entusiasmo, paixão ou dedicação a uma causa."},
}

func GetBatchFlashcardsResponse(ctx context.Context, existingTerms []string, count int) ([]BatchFlashcardItem, error) {
	if count <= 0 {
		count = 20
	}

	existingMap := make(map[string]bool)
	for _, term := range existingTerms {
		existingMap[strings.ToLower(strings.TrimSpace(term))] = true
	}

	var results []BatchFlashcardItem

	// 1. Pick non-duplicate words from curated pool
	for _, item := range fallbackFlashcardsPool {
		termLower := strings.ToLower(strings.TrimSpace(item.Front))
		if !existingMap[termLower] {
			existingMap[termLower] = true
			item.Back = FormatSingleWordTranslation(item.Back)
			results = append(results, item)
		}
		if len(results) >= count {
			break
		}
	}

	// 2. If pool is exhausted, call AI with a fast 5-second timeout
	if len(results) < count {
		needed := count - len(results)
		aiCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		connector := getConnector(aiCtx)
		systemPrompt := "You are a strict language flashcard generator. Output raw JSON array only."
		fullPrompt := fmt.Sprintf(`Generate %d common English vocabulary words/phrases for a Portuguese student.
DO NOT use any of these words: [%s].

Return a JSON array containing %d objects:
[
  {"front": "Word", "back": "Tradução", "context_sentence": "Sentence.", "explanation": "Explicação."}
]`, needed, strings.Join(existingTerms, ", "), needed)

		var resp string
		var err error
		if optConnector, ok := connector.(connectors.OptionableConnector); ok {
			options := map[string]any{
				"temperature": 0.7,
				"num_predict": 512,
			}
			resp, err = optConnector.GenerateResponseWithOptions(aiCtx, fullPrompt, systemPrompt, options)
		} else {
			resp, err = connector.GenerateResponse(aiCtx, fullPrompt)
		}

		if err == nil {
			cleaned := strings.TrimSpace(resp)
			cleaned = strings.TrimPrefix(cleaned, "```json")
			cleaned = strings.TrimPrefix(cleaned, "```")
			cleaned = strings.TrimSuffix(cleaned, "```")
			cleaned = strings.TrimSpace(cleaned)

			start := strings.Index(cleaned, "[")
			end := strings.LastIndex(cleaned, "]")
			if start != -1 && end != -1 && end > start {
				cleaned = cleaned[start : end+1]
			}

			var aiItems []BatchFlashcardItem
			if errUnmarshal := json.Unmarshal([]byte(cleaned), &aiItems); errUnmarshal == nil {
				for _, item := range aiItems {
					termClean := strings.TrimSpace(item.Front)
					if termClean == "" {
						continue
					}
					termLower := strings.ToLower(termClean)
					if !existingMap[termLower] {
						existingMap[termLower] = true
						item.Back = FormatSingleWordTranslation(item.Back)
						results = append(results, item)
					}
					if len(results) >= count {
						break
					}
				}
			}
		}
	}

	return results, nil
}

