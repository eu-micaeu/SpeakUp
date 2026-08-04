# SpeakUp

<p align="center">
  <img src="./frontend/svelte/static/logo.png" width="100" alt="SpeakUp Logo">
</p>

<p align="center">
  <strong>Uma plataforma de inteligência artificial para aprender novos idiomas de forma personalizada e contínua.</strong>
</p>

---

## Descrição do Projeto

O **SpeakUp** é um tutor de idiomas inteligente alimentado por IA que acompanha você em sua jornada de fluência. A plataforma analisa seus textos em tempo real, identificando erros gramaticais e de digitação, fornecendo correções e explicações detalhadas em português sobre cada alteração. 

Além disso, o sistema conta com chat adaptativo baseado nos níveis do quadro europeu comum (CEFR), teste de nivelamento automático e a possibilidade de gravar áudio para praticar a fala e conversação.

---

## Funcionalidades

*   **Tutor de IA em Tempo Real:** Correção gramatical instantânea e contextualizada das mensagens digitadas.
*   **Explicações Didáticas:** Para cada erro corrigido, a IA gera explicações em português detalhando a regra gramatical aplicada.
*   **Chat Inteligente Adaptativo:** Converse com a IA no idioma de destino. A IA ajusta automaticamente a complexidade e vocabulário ao seu nível de proficiência.
*   **Provedores de IA Flexíveis (Ollama & Gemini):** Suporte nativo para executar localmente com **Ollama** via Docker (dev/prod) ou alternar para **Gemini**.
*   **Nivelamento Automático:** Questionário rápido de nivelamento (A1 a C1) durante o cadastro para definir o ponto de partida do aprendizado.
*   **Gravação de Voz (Audio Transcription):** Funcionalidade para converter áudio gravado em texto e responder no chat.
*   **Plataforma 100% Gratuita e Ilimitada:** Acesso total a todas as ferramentas de IA sem limites de requisições diárias ou necessidade de assinatura.
*   **Documentação Interativa:** API documentada com Swagger disponível localmente.

---

## Configuração de IA (Ollama & Gemini)

Você pode alternar facilmente entre o **Ollama** (execução local totalmente offline/open-source via Docker) e o **Gemini**:

No arquivo `.env` (ou no Docker Compose):

```env
# Define o provedor ativo: "ollama" ou "gemini"
AI_PROVIDER=ollama

# Configurações do Ollama
OLLAMA_HOST=http://ollama:11434
OLLAMA_MODEL=llama3.2:latest

# Configurações do Gemini (opcional)
GEMINI_API_KEY=sua_chave_aqui
GEMINI_MODEL=gemini-2.5-flash
```

Ao rodar com `docker-compose -f docker-compose.dev.yml up -d` ou `docker-compose -f docker-compose.prod.yml up -d`, o serviço do Ollama será inicializado e baixará automaticamente o modelo especificado em `OLLAMA_MODEL`.
