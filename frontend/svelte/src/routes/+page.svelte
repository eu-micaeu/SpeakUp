<script lang="ts">
  import { onMount } from "svelte";
  import FeatureCard from "../components/FeatureCard.svelte";
  import Footer from "../components/Footer.svelte";
  import { isAuthTokenValid } from "../utils/cookies";
  import { goto } from "$app/navigation";
  import { login as loginApi, register as registerApi } from "../utils/api";

  const cardData = [
    {
      imgSrc: "./chat.png",
      alt: "Ícone de chat representando conversação",
      title:
        "Plataforma que utiliza IA para auxiliar no aprendizado de idiomas.",
      content:
        "Analisa textos em tempo real, identificando e corrigindo erros com explicações detalhadas para um aprendizado contínuo.",
    },
    {
      imgSrc: "./ia.png",
      alt: "Ícone representando Inteligência Artificial",
      title:
        "Oferece experiência personalizada com análise em tempo real de erros gramaticais e de digitação.",
      content:
        "Fornece correções com explicações para melhor compreensão das regras linguísticas.",
    },
    {
      svg: `<svg width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="5" width="18" height="14" rx="3"/><path d="M3 10h18M8 15h3" stroke-linecap="round"/></svg>`,
      alt: "Ícone representando Flashcards e Memória",
      title:
        "Sistema de Flashcards com Repetição Espaçada (SuperMemo-2).",
      content:
        "Memorize vocabulário e expressões salvas do chat com áudio em tempo real e revisões diárias otimizadas para retenção de longo prazo.",
    },
  ];

  // Auth modal state
  let showAuthModal = false;
  let authMode = "login"; // "login" or "register"
  let isLoading = false;
  let toastMessage = "";
  let toastType = "";
  let error = "";
  let selectedLanguage = "";
  let levels: string[] = [];
  let showLevelTest = false;
  let currentQuestion = 0;
  let testAnswers: number[] = [];

  const levelingQuestions = [
    {
      question: "What ___ you doing right now?",
      options: ["is", "are", "am", "be"],
      correct: 1,
      level: "A1",
    },
    {
      question: "I ___ to the store yesterday.",
      options: ["go", "went", "gone", "going"],
      correct: 1,
      level: "A2",
    },
    {
      question: "If I ___ rich, I would travel the world.",
      options: ["am", "was", "were", "be"],
      correct: 2,
      level: "B1",
    },
    {
      question: "By the time you arrive, I ___ the project.",
      options: ["finish", "finished", "will finish", "will have finished"],
      correct: 3,
      level: "B2",
    },
    {
      question: "Had I known about the meeting, I ___ attended it.",
      options: ["would have", "will have", "would", "will"],
      correct: 0,
      level: "C1",
    },
  ];

  function openAuthModal(mode: string) {
    authMode = mode;
    showAuthModal = true;
    error = "";
  }

  function closeAuthModal() {
    showAuthModal = false;
    error = "";
    selectedLanguage = "";
    levels = [];
  }

  function showToast(message: string, type: string) {
    toastMessage = message;
    toastType = type;
    setTimeout(() => {
      toastMessage = "";
      toastType = "";
    }, 3000);
  }

  const levelDescriptions: Record<string, Record<string, string>> = {
    english: {
      A1: "A1 - Básico",
      A2: "A2 - Elementar",
      B1: "B1 - Intermediário",
      B2: "B2 - Intermediário Superior",
      C1: "C1 - Avançado",
      C2: "C2 - Proficiente",
    },
  };

  function handleLanguageChange(e: Event) {
    const language = (e.target as HTMLSelectElement).value;
    selectedLanguage = language;

    if (language === "english") {
      levels = ["A1", "A2", "B1", "B2", "C1", "C2"];
    } else {
      levels = [];
    }
  }

  function startLevelTest() {
    if (!selectedLanguage) {
      return;
    }
    showLevelTest = true;
    currentQuestion = 0;
    testAnswers = [];
  }

  function answerQuestion(answerIndex: number) {
    testAnswers[currentQuestion] = answerIndex;

    if (currentQuestion < levelingQuestions.length - 1) {
      currentQuestion++;
    } else {
      calculateLevel();
    }
  }

  function calculateLevel() {
    let correctAnswers = 0;
    testAnswers.forEach((answer, index) => {
      if (answer === levelingQuestions[index].correct) {
        correctAnswers++;
      }
    });

    let detectedLevel = "A1";
    if (correctAnswers >= 4) detectedLevel = "C1";
    else if (correctAnswers >= 3) detectedLevel = "B2";
    else if (correctAnswers >= 2) detectedLevel = "B1";
    else if (correctAnswers >= 1) detectedLevel = "A2";

    // Definir o nível no select
    const levelSelect = document.getElementById(
      "register-level",
    ) as HTMLSelectElement;
    if (levelSelect) {
      levelSelect.value = detectedLevel;
    }

    showLevelTest = false;
  }

  function closeLevelTest() {
    showLevelTest = false;
    currentQuestion = 0;
    testAnswers = [];
  }

  function isValidPassword(password: string): boolean {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
    return regex.test(password);
  }

  async function handleLogin(event: Event) {
    event.preventDefault();
    isLoading = true;
    error = "";

    const formData = new FormData(event.target as HTMLFormElement);
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
      const response = await loginApi(email, password);

      if (response) {
        showToast("Login realizado com sucesso!", "success");
        closeAuthModal();
        setTimeout(() => {
          goto("/dashboard");
        }, 1500);
      } else {
        error = "Email ou senha inválidos";
      }
    } catch (err) {
      console.error(err);
      error = "Email ou senha inválidos";
    } finally {
      isLoading = false;
    }
  }

  async function handleRegister(event: Event) {
    event.preventDefault();
    isLoading = true;
    error = "";

    try {
      const formData = new FormData(event.target as HTMLFormElement);

      const userData = {
        name: formData.get("name") as string,
        email: formData.get("email") as string,
        password: formData.get("password") as string,
        language: formData.get("language") as string,
        level: formData.get("level") as string,
      };

      const confirmPassword = formData.get("confirm-password") as string;

      if (userData.password !== confirmPassword) {
        error = "As senhas não coincidem";
        isLoading = false;
        return;
      }

      if (!isValidPassword(userData.password)) {
        error =
          "A senha deve ter pelo menos 8 caracteres, uma letra maiúscula, uma minúscula, um número e um caractere especial.";
        isLoading = false;
        return;
      }

      const response = await registerApi(userData);

      if (response.message === "User created successfully") {
        showToast("Usuário criado com sucesso!", "success");
        authMode = "login";
        error = "";
      } else {
        error = "Erro ao criar usuário.";
      }
    } catch (err: any) {
      if (err?.response?.status === 409) {
        error = "Email já cadastrado.";
      } else {
        error = "Erro ao criar usuário.";
      }
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    if (isAuthTokenValid()) {
      goto("/dashboard");
      return;
    }
  });
</script>

<svelte:head>
    <title>SpeakUp</title>
    <meta name="description" content="Melhore sua pronúncia, escrita e gramática em tempo real com a SpeakUp. A plataforma definitiva de aprendizado de idiomas baseada em Inteligência Artificial. Pratique agora mesmo de forma prática e inovadora!" />
    <meta property="og:title" content="SpeakUp" />
    <meta property="og:description" content="Melhore sua pronúncia, escrita e gramática em tempo real com a SpeakUp. A plataforma definitiva de aprendizado de idiomas baseada em Inteligência Artificial. Pratique agora mesmo de forma prática e inovadora!" />
    <meta name="twitter:title" content="SpeakUp" />
    <meta name="twitter:description" content="Melhore sua pronúncia, escrita e gramática em tempo real com a SpeakUp. A plataforma definitiva de aprendizado de idiomas baseada em Inteligência Artificial. Pratique agora mesmo de forma prática e inovadora!" />
</svelte:head>

<main>
  <div class="betaBanner">
    <div class="betaBannerContent">
      <div class="betaText">
         O SpeakUp está em desenvolvimento.
      </div>
    </div>
  </div>

  {#if toastMessage}
    <div class="toast {toastType}">
      {toastMessage}
    </div>
  {/if}

  <section class="hero">
    <img src="logo.png" alt="Logo" width={60} />
    <h1 class="heroTitle">SpeakUp</h1>
    <p class="heroSubtitle">
      Aprenda qualquer idioma conversando com uma inteligência artificial.
    </p>
    <button on:click={() => openAuthModal("register")} class="cta"
      >Vamos começar!</button
    >
    <button on:click={() => openAuthModal("login")} class="ctaSecondary">
      Já possui uma conta? Aperte aqui!
    </button>
  </section>

  <!-- Seção 1: Recursos -->
  <section class="featureSection">
    <div class="sectionHeader">
      <h1>Recursos do SpeakUp</h1>
      <p class="sectionSubtitle">Descubra tudo o que você pode fazer para alcançar a fluência.</p>
    </div>
    <div class="featuresGrid">
      {#each cardData as card, i}
        <FeatureCard {...card} index={i} />
      {/each}
    </div>
  </section>

  <!-- Seção 2: 100% Gratuito -->
  <section class="freeSection">
    <div class="sectionHeader">
      <h1>Projeto gratuito!</h1>
      <p class="sectionSubtitle">
        Um projeto dedicado à educação de idiomas sem barreiras financeiras.
      </p>
    </div>

    <div class="freeGrid">
      <div class="freeCard">
        <h3>100% Gratuito</h3>
        <p>Acesso livre a todas as ferramentas sem mensalidades ou taxas ocultas.</p>
      </div>
      <div class="freeCard">
        <h3>Uso Ilimitado</h3>
        <p>Pratique conversas, correções gramaticais e traduções sem limite diário.</p>
      </div>
      <div class="freeCard">
        <h3>Para Todos</h3>
        <p>Aprenda no seu ritmo com suporte de IA adaptado ao seu nível de proficiência.</p>
      </div>
    </div>
  </section>

  <!-- Seção 3: Inteligência Artificial -->
  <section class="apiSection">
    <div class="sectionHeader">
      <h1>Potencializado por Inteligência Artificial</h1>
      <p class="sectionSubtitle">
        Utilizamos modelos locais e open-source do Ollama para respostas rápidas, privacidade total, correções instantâneas e aprendizado adaptativo.
      </p>
    </div>
    <div class="apiLogos">
      <div class="apiCard">
        <img src="/Ollama Logo - White - zonalogo.com.svg" alt="Ollama Logo" />
        <h3>Ollama</h3>
        <p>
          Modelos de linguagem open-source rodando localmente para análise gramatical precisa, conversação fluida e privacidade completa dos seus dados.
        </p>
      </div>
    </div>
  </section>

  {#if showAuthModal}
    <div class="modalOverlay" on:click={closeAuthModal}>
      <div class="modalContent" on:click|stopPropagation>
        <button
          class="closeButton"
          on:click={closeAuthModal}
          aria-label="Fechar modal"
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
          >
            <path
              d="M18 6L6 18M6 6l12 12"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>

        <div class="modalHeader">
          <h2>
            {authMode === "login" ? "Bem-vindo de volta!" : "Crie sua conta"}
          </h2>
          <p class="modalSubtitle">
            {authMode === "login"
              ? "Entre para continuar a conversa"
              : "Junte-se a nós e comece a aprender"}
          </p>
        </div>

        <div class="authTabs">
          <button
            class="authTab"
            class:active={authMode === "login"}
            on:click={() => {
              authMode = "login";
              error = "";
            }}
          >
            Login
          </button>
          <button
            class="authTab"
            class:active={authMode === "register"}
            on:click={() => {
              authMode = "register";
              error = "";
            }}
          >
            Registrar
          </button>
        </div>

        {#if error}
          <div class="errorMessage">
            {error}
          </div>
        {/if}

        {#if authMode === "login"}
          <form on:submit={handleLogin}>
            <div class="inputGroup">
              <label for="login-email">Email</label>
              <input
                id="login-email"
                type="email"
                name="email"
                placeholder="seu@email.com"
                required
              />
            </div>

            <div class="inputGroup">
              <label for="login-password">Senha</label>
              <input
                id="login-password"
                type="password"
                name="password"
                placeholder="••••••••"
                required
              />
            </div>

            <button type="submit" disabled={isLoading} class="submitButton">
              {isLoading ? "Entrando..." : "Entrar"}
            </button>
          </form>
        {:else}
          <form on:submit={handleRegister}>
            <div class="inputGroup">
              <label for="register-name">Nome completo</label>
              <input
                id="register-name"
                type="text"
                name="name"
                placeholder="Digite seu nome"
                required
              />
            </div>

            <div class="inputRow">
              <div class="inputGroup">
                <label for="register-language">Idioma que deseja aprender</label
                >
                <select
                  id="register-language"
                  name="language"
                  bind:value={selectedLanguage}
                  on:change={handleLanguageChange}
                  required
                >
                  <option value="" disabled selected>Escolha um idioma</option>
                  <option value="english">🇺🇸 Inglês</option>
                </select>
              </div>

              <div class="inputGroup">
                <label for="register-level">Nível atual</label>
                <select
                  id="register-level"
                  name="level"
                  required
                  disabled={!selectedLanguage}
                  class:disabled={!selectedLanguage}
                >
                  <option value="" disabled selected>
                    {selectedLanguage
                      ? "Escolha seu nível"
                      : "Selecione um idioma primeiro"}
                  </option>
                  {#each levels as level}
                    <option value={level}>
                      {levelDescriptions[selectedLanguage]?.[level] || level}
                    </option>
                  {/each}
                </select>
              </div>
            </div>

            <div class="inputGroup">
              <button
                type="button"
                class="levelTestButton"
                on:click={startLevelTest}
                disabled={!selectedLanguage}
              >
                Não sabe o seu nível? - Fazer teste rápido
              </button>
            </div>

            <div class="inputGroup">
              <label for="register-email">Email</label>
              <input
                id="register-email"
                type="email"
                name="email"
                placeholder="Digite seu email"
                required
              />
            </div>

            <div class="inputRow">
              <div class="inputGroup">
                <label for="register-password">Senha</label>
                <input
                  id="register-password"
                  type="password"
                  name="password"
                  placeholder="Digite sua senha"
                  required
                />
              </div>

              <div class="inputGroup">
                <label for="register-confirm">Confirmar Senha</label>
                <input
                  id="register-confirm"
                  type="password"
                  name="confirm-password"
                  placeholder="Confirme sua senha"
                  required
                />
              </div>
            </div>

            <button type="submit" disabled={isLoading} class="submitButton">
              {isLoading ? "Carregando..." : "Cadastrar"}
            </button>
          </form>
        {/if}
      </div>
    </div>
  {/if}

  <!-- Level Test Modal -->
  {#if showLevelTest}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="modalOverlay testOverlay" on:click={closeLevelTest}>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <div class="modalContent testModal" on:click|stopPropagation>
        <button
          class="closeButton"
          on:click={closeLevelTest}
          aria-label="Fechar teste"
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
          >
            <path
              d="M18 6L6 18M6 6l12 12"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>

        <div class="testHeader">
          <h2>Teste de Nivelamento</h2>
          <p class="testProgress">
            Questão {currentQuestion + 1} de {levelingQuestions.length}
          </p>
        </div>

        <div class="questionContainer">
          <p class="question">{levelingQuestions[currentQuestion].question}</p>
          <div class="optionsGrid">
            {#each levelingQuestions[currentQuestion].options as option, index}
              <button
                type="button"
                class="optionButton"
                on:click={() => answerQuestion(index)}
              >
                {option}
              </button>
            {/each}
          </div>
        </div>
      </div>
    </div>
  {/if}
</main>

<Footer />

<style>
  main {
    position: relative;
    overflow-x: hidden;
    background-color: #0a0a0a;
    color: white;
    min-height: 100vh;
  }

  .hero {
    height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
    z-index: 1;
    position: relative;
    padding: 2rem;
  }

  .heroTitle {
    font-size: 3rem;
    font-weight: 800;
    margin: 0;
    margin-bottom: 1rem;
    background: linear-gradient(
      270deg,
      #ffffff,
      #bfbfbf,
      #a0a0a0,
      #8a8a8a,
      #bfbfbf,
      #ffffff
    );
    background-size: 200% auto;
    background-clip: text;
    -webkit-background-clip: text;
    color: transparent;
    -webkit-text-fill-color: transparent;
    animation: chromeShine 6s ease-in-out infinite;
  }

  @keyframes chromeShine {
    0% {
      background-position: 0% center;
    }

    100% {
      background-position: -200% center;
    }
  }

  .heroSubtitle {
    font-size: 1rem;
    max-width: 600px;
    margin-bottom: 2rem;
    opacity: 0.8;
  }

  .cta {
    background: #ffffff;
    padding: 1rem 2rem;
    border-radius: 8px;
    font-size: 1.2rem;
    border: none;
    color: rgb(0, 0, 0);
    font-weight: bold;
    transition: 0.3s;
    cursor: pointer;
  }

  .cta:hover {
    background: #f0f0f0;
    transform: translateY(-2px);
  }

  .ctaSecondary {
    background: transparent;
    color: #ffffff;
    padding: 0.8rem 1.5rem;
    font-size: 1rem;
    border: none;
    font-weight: bold;
    transition: 0.3s;
    cursor: pointer;
    margin: 15px 0;
  }

  /* Seções Padronizadas da Tela Inicial */
  .featureSection,
  .freeSection,
  .apiSection {
    padding: 5rem 2rem;
    background: #0a0a0a;
    position: relative;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
  }

  .freeSection {
    background: rgba(255, 255, 255, 0.015);
  }

  .sectionHeader {
    text-align: center;
    margin-bottom: 3.5rem;
    animation: fadeIn 1s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .sectionHeader h1 {
    font-size: 2.4rem;
    font-weight: 800;
    margin-bottom: 0.75rem;
    background: linear-gradient(
      270deg,
      #ffffff,
      #bfbfbf,
      #a0a0a0,
      #8a8a8a,
      #bfbfbf,
      #ffffff
    );
    background-size: 200% auto;
    background-clip: text;
    -webkit-background-clip: text;
    color: transparent;
    -webkit-text-fill-color: transparent;
    animation: chromeShine 6s ease-in-out infinite;
    letter-spacing: -0.02em;
  }

  .sectionSubtitle {
    font-size: 1.05rem;
    color: rgba(255, 255, 255, 0.6);
    max-width: 600px;
    margin: 0 auto;
    line-height: 1.6;
  }

  .featuresGrid {
    max-width: 1100px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .freeGrid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 1.5rem;
    max-width: 1100px;
    margin: 0 auto;
  }

  .freeCard {
    background: rgba(255, 255, 255, 0.025);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    padding: 2.25rem 1.75rem;
    text-align: center;
    transition: all 0.3s ease;
  }

  .freeCard:hover {
    background: rgba(255, 255, 255, 0.04);
    border-color: rgba(255, 255, 255, 0.18);
    transform: translateY(-3px);
  }

  .freeCard h3 {
    font-size: 1.25rem;
    font-weight: 700;
    color: #ffffff;
    margin-bottom: 0.75rem;
  }

  .freeCard p {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.95rem;
    line-height: 1.6;
    margin: 0;
  }

  .apiLogos {
    display: flex;
    justify-content: center;
    gap: 2rem;
    max-width: 1100px;
    margin: 0 auto;
    flex-wrap: wrap;
  }

  .apiCard {
    background: rgba(255, 255, 255, 0.025);
    padding: 2.25rem 2rem;
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.08);
    transition: all 0.3s ease;
    max-width: 450px;
    width: 100%;
    text-align: center;
  }

  .apiCard:hover {
    background: rgba(255, 255, 255, 0.04);
    border-color: rgba(255, 255, 255, 0.18);
    transform: translateY(-3px);
  }

  .apiCard img {
    width: 64px;
    height: 64px;
    margin-bottom: 1rem;
  }

  .apiCard h3 {
    font-size: 1.3rem;
    font-weight: 700;
    margin-bottom: 0.75rem;
    color: #ffffff;
  }

  .apiCard p {
    font-size: 0.95rem;
    color: rgba(255, 255, 255, 0.6);
    line-height: 1.6;
    margin: 0;
  }

  .apiDescription {
    font-size: 1.1rem;
    max-width: 800px;
    margin: 0 auto;
    opacity: 0.8;
    line-height: 1.8;
  }

  .modalOverlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.75);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: modalFadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  @keyframes modalFadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes modalScaleUp {
    from { opacity: 0; transform: scale(0.96) translateY(6px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
  }

  .modalContent {
    background: #111113;
    border-radius: 8px;
    padding: 2rem 2.25rem;
    max-width: 440px;
    width: 90%;
    overflow-y: auto;
    max-height: 90vh;
    position: relative;
    border: 1px solid rgba(255, 255, 255, 0.08);
    animation: modalScaleUp 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .closeButton {
    position: absolute;
    top: 1.25rem;
    right: 1.25rem;
    background: transparent;
    border: none;
    border-radius: 4px;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: rgba(255, 255, 255, 0.4);
    transition: all 0.2s ease;
  }

  .closeButton:hover {
    background: rgba(255, 255, 255, 0.08);
    color: #ffffff;
  }

  .modalHeader {
    text-align: center;
    margin-bottom: 1.25rem;
  }

  .modalLogo {
    width: 42px;
    height: 42px;
    margin-bottom: 0.75rem;
  }

  .modalHeader h2 {
    font-size: 1.4rem;
    font-weight: 700;
    margin-bottom: 0.25rem;
    color: #ffffff;
    letter-spacing: -0.02em;
  }

  .modalSubtitle {
    color: rgba(255, 255, 255, 0.45);
    font-size: 0.85rem;
  }

  .authTabs {
    display: flex;
    background: rgba(255, 255, 255, 0.035);
    border: 1px solid rgba(255, 255, 255, 0.06);
    padding: 3px;
    border-radius: 6px;
    margin-bottom: 1.5rem;
  }

  .authTab {
    flex: 1;
    padding: 0.55rem;
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.88rem;
    font-weight: 500;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .authTab.active {
    background: rgba(255, 255, 255, 0.12);
    color: #ffffff;
    font-weight: 600;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }

  .authTab:hover:not(.active) {
    color: rgba(255, 255, 255, 0.8);
  }

  .errorMessage {
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.2);
    color: #fca5a5;
    border-radius: 4px;
    padding: 0.6rem 0.9rem;
    font-size: 0.85rem;
    margin-bottom: 1.25rem;
    text-align: center;
  }

  .inputGroup {
    margin-bottom: 1rem;
  }

  .inputGroup label {
    display: block;
    color: rgba(255, 255, 255, 0.65);
    font-size: 0.78rem;
    margin-bottom: 0.35rem;
    font-weight: 500;
  }

  .inputGroup input,
  .inputGroup select {
    width: 100%;
    padding: 0.7rem 0.85rem;
    background: rgba(255, 255, 255, 0.035);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 6px;
    color: #ffffff;
    font-size: 0.9rem;
    transition: all 0.2s ease;
    box-sizing: border-box;
  }

  .inputGroup input:focus,
  .inputGroup select:focus {
    outline: none;
    border-color: rgba(92, 109, 255, 0.6);
    background: rgba(255, 255, 255, 0.06);
    box-shadow: 0 0 0 3px rgba(92, 109, 255, 0.15);
  }

  .inputGroup input::placeholder {
    color: rgba(255, 255, 255, 0.3);
  }

  .inputGroup select {
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23ffffff' opacity='0.5' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.85rem center;
    background-size: 12px;
    padding-right: 2.25rem;
    cursor: pointer;
  }

  .inputGroup select:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .inputGroup select option {
    background: #18181b;
    color: white;
    padding: 0.5rem;
  }

  .inputGroup select option:disabled {
    color: rgba(255, 255, 255, 0.4);
  }

  .inputRow {
    display: flex;
    gap: 0.75rem;
  }

  .inputRow .inputGroup {
    flex: 1;
  }

  .submitButton {
    width: 100%;
    padding: 0.8rem;
    background: #ffffff;
    border: none;
    border-radius: 6px;
    color: #09090b;
    font-size: 0.92rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .submitButton:hover:not(:disabled) {
    background: #f4f4f5;
    transform: translateY(-1px);
    box-shadow: 0 4px 14px rgba(255, 255, 255, 0.15);
  }

  .submitButton:active:not(:disabled) {
    transform: translateY(0);
  }

  .submitButton:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .levelTestButton {
    width: 100%;
    padding: 0.4rem 0;
    background: transparent;
    border: none;
    color: #8c9eff;
    font-size: 0.8rem;
    font-weight: 500;
    cursor: pointer;
    text-align: right;
    transition: all 0.2s ease;
  }

  .levelTestButton:hover:not(:disabled) {
    color: #b3c0ff;
    text-decoration: underline;
  }

  .levelTestButton:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .testOverlay {
    z-index: 1100;
  }

  .testModal {
    max-width: 600px;
    min-height: 400px;
    border-radius: 8px;
  }

  .testHeader {
    text-align: center;
    margin-bottom: 2rem;
  }

  .testHeader h2 {
    font-size: 1.8rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
    color: white;
  }

  .testProgress {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.9rem;
  }

  .questionContainer {
    margin-top: 2rem;
  }

  .question {
    font-size: 1.3rem;
    font-weight: 600;
    color: white;
    margin-bottom: 2rem;
    line-height: 1.6;
    text-align: center;
  }

  .optionsGrid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .optionButton {
    padding: 1.25rem;
    background: rgba(255, 255, 255, 0.05);
    border: 2px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    color: white;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s;
  }

  .optionButton:hover {
    background: rgba(255, 255, 255, 0.1);
    border-color: rgba(255, 255, 255, 0.3);
    transform: translateY(-2px);
  }

  .errorMessage {
    background: rgba(239, 68, 68, 0.15);
    border: 1px solid rgba(239, 68, 68, 0.3);
    color: #ff6b6b;
    padding: 1rem;
    border-radius: 10px;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }

  .toast {
    position: fixed;
    top: 20px;
    right: 20px;
    padding: 1rem 1.5rem;
    border-radius: 12px;
    color: white;
    z-index: 2000;
    animation: slideIn 0.3s ease-in;
    font-weight: 500;
  }

  .toast.success {
    background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  }

  .toast.error {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  }

  @keyframes slideIn {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  .betaBanner {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    padding: 1rem 2rem;
    border-bottom: 2px solid #b45309;
    position: sticky;
    top: 0;
    z-index: 999;
  }

  .betaBannerContent {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    gap: 1rem;
    font-size: 0.95rem;
  }

  .betaText {
    color: #1f2937;
    font-weight: bold;
    line-height: 1.5;
    text-align: center;
    width: 100%;
  }

  .betaText strong {
    font-weight: 700;
  }

  @media screen and (max-width: 768px) {
    .particles {
      display: none;
    }

    .betaBanner {
      padding: 0.875rem 1rem;
    }

    .betaBannerContent {
      font-size: 0.85rem;
      gap: 0.75rem;
    }

    .hero {
      padding: 0;
    }

    .heroTitle {
      font-size: 2rem;
      text-align: center;
    }

    .heroSubtitle {
      font-size: 0.8rem;
      max-width: 90%;
      margin-bottom: 2rem;
    }

    .cta {
      font-size: 1rem;
      padding: 0.8rem 1.5rem;
    }

    .featureSection {
      padding: 4rem 1rem;
    }

    .sectionHeader h1 {
      font-size: 2rem;
    }

    .sectionSubtitle {
      font-size: 1rem;
    }

    .featuresGrid {
      gap: 1.5rem;
    }

    .apiSection {
      padding: 3rem 1rem;
    }

    .apiContainer h1 {
      font-size: 1.8rem;
      margin-bottom: 2rem;
    }

    .apiLogos {
      gap: 2rem;
      flex-direction: column;
      align-items: center;
    }

    .apiCard {
      width: 90%;
      max-width: 300px;
    }

    .apiDescription {
      font-size: 0.95rem;
      padding: 0 1rem;
    }

    .modalContent {
      padding: 2rem 1.5rem;
      width: 95%;
    }

    .modalHeader h2 {
      font-size: 1.5rem;
    }

    .inputRow {
      flex-direction: column;
      gap: 0;
    }

    .authTab {
      font-size: 0.9rem;
      padding: 0.6rem;
    }

    .optionsGrid {
      grid-template-columns: 1fr;
    }

    .question {
      font-size: 1.1rem;
    }

    .testModal {
      min-height: 350px;
    }
  }
</style>
