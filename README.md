# SpeakUp 🗣️

<p align="center">
  <img src="./frontend/svelte/static/logo.png" width="120" alt="SpeakUp Logo">
</p>

<p align="center">
  <strong>Uma plataforma inovadora impulsionada por inteligência artificial para aprender novos idiomas de forma personalizada e contínua.</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Svelte-FF3E00?style=for-the-badge&logo=svelte&logoColor=white" alt="Svelte" />
  <img src="https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white" alt="MongoDB" />
  <img src="https://img.shields.io/badge/Stripe-5433FF?style=for-the-badge&logo=stripe&logoColor=white" alt="Stripe" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
</p>

---

## 📋 Descrição do Projeto

O **SpeakUp** é um tutor de idiomas inteligente alimentado por IA que acompanha você em sua jornada de fluência. A plataforma analisa seus textos em tempo real, identificando erros gramaticais e de digitação, fornecendo correções e explicações detalhadas em português sobre cada alteração. 

Além disso, o sistema conta com chat adaptativo baseado nos níveis do quadro europeu comum (CEFR), teste de nivelamento automático e a possibilidade de gravar áudio para praticar a fala e conversação.

---

## ✨ Funcionalidades Principais

*   **Tutor de IA em Tempo Real:** Correção gramatical instantânea e contextualizada das mensagens digitadas.
*   **Explicações Didáticas:** Para cada erro corrigido, a IA gera explicações em português detalhando a regra gramatical aplicada.
*   **Chat Inteligente Adaptativo:** Converse com a IA no idioma de destino. A IA ajusta automaticamente a complexidade e vocabulário ao seu nível de proficiência.
*   **Nivelamento Automático:** Questionário rápido de nivelamento (A1 a C1) durante o cadastro para definir o ponto de partida do aprendizado.
*   **Gravação de Voz (Audio Transcription):** Funcionalidade para converter áudio gravado em texto e responder no chat.
*   **Controle de Uso & Assinaturas (Stripe):**
    *   **Plano Gratuito:** Limite configurável de requisições de IA diárias (padrão de 10 interações).
    *   **Plano Pro:** Acesso ilimitado às ferramentas com planos de assinatura mensais e anuais integrados ao Stripe.
*   **Documentação Interativa:** API documentada com Swagger disponível localmente.

---

## 🛠️ Tecnologias Utilizadas

### Backend
*   **Linguagem:** Go (Golang 1.22+)
*   **Framework:** Gin Gonic
*   **Banco de Dados:** MongoDB (com drivers nativos em Go)
*   **Autenticação:** JWT (JSON Web Tokens) e Criptografia com bcrypt
*   **Hot-reloads:** Air (para desenvolvimento rápido)
*   **APIs de IA:** Google Gemini API (via SDK)

### Frontend
*   **Framework:** SvelteKit (Svelte 5)
*   **Estilização:** CSS Vanilla com variáveis globais e animações
*   **Comunicação:** Axios para API REST
*   **Efeitos Visuais:** Three.js (gráficos na Home)
*   **Componentes Adicionais:** Svelte Sonner (notificações) e Firebase SDK

---

## 📂 Estrutura do Projeto

```text
SpeakUp/
├── backend/
│   └── go/
│       ├── cmd/
│       │   └── api/          # Ponto de entrada (main.go)
│       └── pkg/
│           ├── adapters/     # Conectores externos (Gemini, IA)
│           ├── config/       # Conexão com MongoDB e Configurações gerais
│           ├── handlers/     # Controladores das rotas
│           ├── middlewares/  # Middlewares de Auth (JWT) e CORS
│           ├── models/       # Modelos estruturais e esquemas do MongoDB
│           ├── planlimits/   # Lógica de controle de limites de requisições
│           ├── prompts/      # Arquivos de prompts da IA (.txt)
│           ├── repositories/ # Camada de dados do MongoDB (User, Chat, Message)
│           ├── routes/       # Definição das rotas HTTP
│           └── utils/        # Funções utilitárias
├── frontend/
│   └── svelte/
│       ├── src/
│       │   ├── components/   # Componentes reutilizáveis
│       │   ├── routes/       # Páginas do SvelteKit (Landing Page, Chat, Perfil)
│       │   └── utils/        # APIs, cookies e gerenciadores auxiliares
│       └── static/           # Ativos e imagens estáticas
├── docker-compose.yml        # Orquestração do ambiente
├── Makefile                  # Comandos utilitários simplificados
└── LICENSE                   # Licença do repositório
```

---

## ⚙️ Configuração do Ambiente

Para iniciar o projeto localmente, você precisa configurar os arquivos de ambiente `.env`.

### 1. Backend (`/backend/go/.env`)
Crie uma cópia do arquivo `/backend/go/.env.example` com o nome `.env` e preencha com as suas credenciais:

```env
# MongoDB
MONGO_CONNECTION_STRING=mongodb://host.docker.internal:27017/speakup
DB_NAME=speakup

# Autenticação
JWT_KEY=insira-um-segredo-longo-e-seguro

# Provedor de IA (Gemini)
GEMINI_API_KEY=sua_chave_gemini_aqui
GEMINI_MODEL=gemini-1.5-flash

# Limites do Plano Gratuito
FREE_DAILY_AI_LIMIT=10

# Configurações do Stripe (Assinaturas)
STRIPE_SECRET_KEY=sk_test_...
STRIPE_SUCCESS_URL=http://localhost:5175/checkout/success
STRIPE_CANCEL_URL=http://localhost:5175/checkout/cancel
STRIPE_PORTAL_RETURN_URL=http://localhost:5175/perfil
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_PRICE_MONTHLY=price_...
STRIPE_PRICE_ANNUAL=price_...
```

### 2. Frontend (`/frontend/svelte/.env`)
Crie uma cópia do arquivo `/frontend/svelte/.env.example` com o nome `.env` e ajuste as variáveis:

```env
VITE_API_URL=http://host.docker.internal:8082/api
VITE_STRIPE_PRICE_MONTHLY=price_...
VITE_STRIPE_PRICE_ANNUAL=price_...

# Firebase (Auth & Analytics)
PUBLIC_FIREBASE_API_KEY=sua_chave_firebase
PUBLIC_FIREBASE_AUTH_DOMAIN=seu-projeto.firebaseapp.com
PUBLIC_FIREBASE_PROJECT_ID=seu-projeto-id
PUBLIC_FIREBASE_STORAGE_BUCKET=seu-projeto.firebasestorage.app
PUBLIC_FIREBASE_MESSAGING_SENDER_ID=sender_id
PUBLIC_FIREBASE_APP_ID=app_id
PUBLIC_FIREBASE_MEASUREMENT_ID=measurement_id
```

---

## 🚀 Como Executar

### Pré-requisitos
*   [Docker](https://www.docker.com/) e Docker Compose instalados.
*   [MongoDB](https://www.mongodb.com/) rodando no host local (na porta padrão `27017`) ou na nuvem.
*   *Opcional para execução local manual:* [Go 1.22+](https://go.dev/) e [Node.js](https://nodejs.org/).

### Executando com Docker (Recomendado)
A forma mais simples de subir toda a aplicação é utilizando o Docker Compose via `Makefile` disponibilizado no projeto:

1.  Inicie os containers da aplicação:
    ```bash
    make up
    ```
    Isso subirá os serviços:
    *   **Frontend (SvelteKit):** `http://localhost:5175`
    *   **Backend (Go/Gin):** `http://localhost:8082`

2.  Para derrubar e limpar os containers:
    ```bash
    make down
    ```

3.  Para reiniciar os containers rapidamente:
    ```bash
    make restart
    ```

### Executando Manualmente (Desenvolvimento)

#### Iniciar o Backend
```bash
cd backend/go
# Certifique-se de que o MongoDB esteja ativo localmente
go run cmd/api/main.go
# Ou utilizando o Air para recarregamento em tempo real (caso possua instalado):
air
```

#### Iniciar o Frontend
```bash
cd frontend/svelte
npm install
npm run dev -- --port 5175
```

---

## 🧪 Testes no Backend

Para rodar os testes unitários e de integração no backend com relatório de cobertura de código:

```bash
make tests
```

---

## 📡 Endpoints da API

O backend oferece uma documentação interativa baseada no Swagger. Para acessá-la, com o servidor backend rodando:
🔗 **http://localhost:8082/swagger/index.html**

Principais grupos de rotas (todas as rotas protegidas exigem o header `Authorization: Bearer <token>`):

| Método | Rota | Descrição | Requer Auth |
| :--- | :--- | :--- | :---: |
| **POST** | `/api/user/` | Criação de Novo Usuário (Cadastro) | Não |
| **POST** | `/api/user/login` | Login e Geração do Token JWT | Não |
| **POST** | `/api/chat` | Cria um Novo Canal de Chat de Aprendizado | Sim |
| **GET** | `/api/chat/user` | Busca todos os chats do usuário logado | Sim |
| **POST** | `/api/message` | Envia e registra uma mensagem | Sim |
| **POST** | `/api/ai/generate-response-dialog` | Gera a resposta do bot baseada no nível | Sim |
| **POST** | `/api/ai/generate-response-correction` | Analisa e corrige um texto enviado | Sim |
| **POST** | `/api/ai/generate-response-translation` | Traduz o input de entrada para português | Sim |
| **GET** | `/api/ai/usage` | Retorna o consumo atual do limite diário do usuário | Sim |
| **POST** | `/api/billing/checkout` | Cria sessão de pagamento no Stripe | Sim |
| **POST** | `/api/billing/webhook` | Recebe notificações de status de pagamento do Stripe | Não |

---

## 🧠 Integração com IA (Modelos)

O backend possui uma arquitetura desacoplada baseada no padrão de projeto Adapter para gerenciar a integração com a inteligência artificial:

*   **Google Gemini:** É altamente eficiente e rápido. O modelo padrão configurado é o `gemini-1.5-flash`.

---

## 📄 Licença

Este projeto está licenciado sob a licença MIT - consulte o arquivo [LICENSE](file:///home/micael/GitHub/SpeakUp/LICENSE) para obter mais detalhes.
