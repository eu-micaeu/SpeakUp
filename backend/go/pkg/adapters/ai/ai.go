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

type GeneratedFlashcardData struct {
	Back            string `json:"back"`
	ContextSentence string `json:"context_sentence"`
	Explanation     string `json:"explanation"`
}

func GetFlashcardGenerateResponse(ctx context.Context, term string, contextSentence string) (GeneratedFlashcardData, error) {
	var result GeneratedFlashcardData

	prompt, err := prompts.GetPrompt("promptFlashcard.txt")
	if err != nil {
		return result, fmt.Errorf("failed to load prompt: %w", err)
	}

	connector := getConnector(ctx)
	fullPrompt := fmt.Sprintf("%s\n\nWORD/PHRASE:\n%s\n\nCONTEXT SENTENCE:\n%s\n\nOUTPUT (JSON ONLY):", prompt, term, contextSentence)

	var resp string
	if optConnector, ok := connector.(connectors.OptionableConnector); ok {
		systemPrompt := "You are a strict language flashcard generator. Output raw JSON only."
		options := map[string]any{
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

	cleaned := strings.TrimSpace(resp)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		// Fallback if parsing fails
		result.Back = term
		result.ContextSentence = contextSentence
		result.Explanation = "Tradução/Explicação automática."
	}

	result.Back = FormatSingleWordTranslation(result.Back)

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

