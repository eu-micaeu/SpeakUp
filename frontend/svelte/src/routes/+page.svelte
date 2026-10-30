<script lang="ts">
  import { onMount } from "svelte";
  import FeatureCard from "../components/FeatureCard.svelte";
  import Footer from "../components/Footer.svelte";
  import { isAuthTokenValid } from "../utils/cookies";
  import { goto } from "$app/navigation";
  import { login as loginApi, register as registerApi } from "../utils/api";

  const cardData = [
    {
      svg: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>`,
      tag: "Conversação",
      title: "Chat Adaptativo por Nível",
      content:
        "Diálogos naturais ajustados automaticamente ao seu nível CEFR (A1 ao C2) para praticar no seu ritmo.",
    },
    {
      svg: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.9 5.8a2 2 0 0 1-1.3 1.3L3 12l5.8 1.9a2 2 0 0 1 1.3 1.3L12 21l1.9-5.8a2 2 0 0 1 1.3-1.3L21 12l-5.8-1.9a2 2 0 0 1-1.3-1.3Z"/></svg>`,
      tag: "Gramática",
      title: "Correção em Tempo Real",
      content:
        "Identifica erros gramaticais e de digitação instantaneamente, fornecendo explicações claras em português.",
    },
    {
      svg: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="5" width="18" height="14" rx="2.5"/><path d="M3 10h18M8 15h4"/></svg>`,
      tag: "Memória",
      title: "Flashcards SuperMemo-2",
      content:
        "Salve termos diretamente do chat e revise no momento ideal com repetição espaçada e pronúncia em áudio.",
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

  // Live preview animation state in Hero
  let previewStep = 0;
  let previewInterval: ReturnType<typeof setInterval> | null = null;

  function reveal(node: HTMLElement, delay = 0) {
    if (delay > 0) {
      node.style.transitionDelay = `${delay}ms`;
    }
    if (typeof IntersectionObserver === "undefined") {
      node.classList.add("revealed");
      return;
    }
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            node.classList.add("revealed");
            observer.unobserve(node);
          }
        });
      },
      { threshold: 0.15, rootMargin: "0px 0px -40px 0px" }
    );
    observer.observe(node);
    return {
      destroy() {
        observer.disconnect();
      },
    };
  }

  onMount(() => {
    if (isAuthTokenValid()) {
      goto("/dashboard");
      return;
    }

    if (typeof window !== "undefined") {
      const params = new URLSearchParams(window.location.search);
      if (params.get("login") === "true") {
        openAuthModal("login");
        showToast("Faça login para acessar sua conta.", "error");
        window.history.replaceState({}, "", "/");
      }
    }

    // Cycle through the 3 demo stages gently
    previewStep = 1;
    previewInterval = setInterval(() => {
      previewStep = previewStep >= 3 ? 1 : previewStep + 1;
    }, 3400);

    return () => {
      if (previewInterval) clearInterval(previewInterval);
    };
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
  <!-- Top Navigation Translúcida -->
  <header class="landingNav">
    <div class="navInner">
      <a href="/" class="navBrand" aria-label="SpeakUp Início">
        <img src="logo.png" alt="SpeakUp Logo" width="32" height="32" />
        <span>SpeakUp</span>
      </a>
      <div class="navActions">
        <button type="button" on:click={() => openAuthModal("login")} class="navLoginBtn">
          Entrar
        </button>
        <button type="button" on:click={() => openAuthModal("register")} class="navRegisterBtn">
          Começar grátis
        </button>
      </div>
    </div>
  </header>

  {#if toastMessage}
    <div class="toast {toastType}">
      {toastMessage}
    </div>
  {/if}

  <section class="hero">
    <div class="heroContainer">
      <!-- Coluna Esquerda: Copy + CTAs -->
      <div class="heroContent">

        <h1 class="heroTitle heroAnim2">
          Aprenda idiomas conversando em tempo real com IA.
        </h1>

        <p class="heroSubtitle heroAnim3">
          Pratique diálogos adaptados ao seu nível CEFR, receba correções gramaticais instantâneas explicadas em português e memorize vocabulário com flashcards inteligentes.
        </p>

        <div class="heroActions heroAnim4">
          <button on:click={() => openAuthModal("register")} class="cta">
            <span>Vamos começar!</span>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M5 12h14" />
              <path d="m12 5 7 7-7 7" />
            </svg>
          </button>
          <button on:click={() => openAuthModal("login")} class="ctaSecondary">
            Já possui uma conta? Entrar
          </button>
        </div>
      </div>

      <!-- Coluna Direita: Demonstração Animada do Chat -->
      <div class="heroPreview heroAnim5" aria-label="Demonstração interativa do tutor de IA">
        <div class="previewWindow">
          <div class="previewHeader">
            <div class="windowDots">
              <span></span><span></span><span></span>
            </div>
            <div class="previewTitle">
              <span class="onlineDot"></span>
              SpeakUp • llama3.2:latest
            </div>
            <div class="stepPills">
              {#each [1, 2, 3] as step}
                <button
                  type="button"
                  class="stepPill"
                  class:active={previewStep >= step}
                  on:click={() => (previewStep = step)}
                  aria-label="Etapa {step} da demonstração"
                ></button>
              {/each}
            </div>
          </div>

          <div class="previewBody">
            <!-- Etapa 1: Mensagem do Usuário -->
            <div class="demoMsg userDemo" class:visible={previewStep >= 1}>
              <div class="demoMsgHeader">Você escreveu</div>
              <p class="demoMsgText">Yesterday I <span class="errHighlight">go</span> to the cinema with my friends.</p>

              <!-- Etapa 2: Correção + Explicação -->
              <div class="demoCorrectionBox" class:expanded={previewStep >= 2}>
                <div class="demoCorrLine">
                  <span class="corrBadge">✓ Correção</span>
                  <span>Yesterday I <strong>went</strong> to the cinema with my friends.</span>
                </div>
                <p class="demoExplLine">
                  <strong>Explicação:</strong> Para ações concluídas no passado (<em>"Yesterday"</em>), usamos o passado simples (<strong>went</strong>) em vez de "go".
                </p>
              </div>
            </div>

            <!-- Etapa 3: Resposta do Tutor IA + Tradução -->
            <div class="demoMsg aiDemo" class:visible={previewStep >= 3}>
              <div class="demoMsgHeader">
                <span>SpeakUp IA</span>
                <span class="cefrTag">Nível B1</span>
              </div>
              <p class="demoMsgText">That sounds fun! Which movie did you watch together?</p>
              <p class="demoTranslation">
                <strong>Tradução:</strong> Isso parece divertido! Qual filme vocês assistiram juntos?
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <a href="#recursos" class="scrollHint" aria-label="Rolar para ver recursos">
      <span>Explorar recursos</span>
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="m6 9 6 6 6-6" />
      </svg>
    </a>
  </section>

  <!-- Seção 1: Recursos -->
  <section id="recursos" class="featureSection">
    <div class="sectionContainer">
      <div class="sectionHeader revealItem" use:reveal>
        <span class="sectionEyebrow">Recursos</span>
        <h2>Recursos do SpeakUp</h2>
        <p class="sectionSubtitle">Ferramentas essenciais para acelerar sua fluência todos os dias.</p>
      </div>
      <div class="featuresGrid">
        {#each cardData as card, i}
          <FeatureCard {...card} index={i} />
        {/each}
      </div>
    </div>
  </section>

  <!-- Seção 2: 100% Gratuito -->
  <section class="freeSection">
    <div class="sectionContainer">
      <div class="sectionHeader revealItem" use:reveal>
        <span class="sectionEyebrow">Acesso Livre</span>
        <h2>Projeto gratuito!</h2>
        <p class="sectionSubtitle">
          Educação de idiomas acessível e sem barreiras financeiras.
        </p>
      </div>

      <div class="freeGrid">
        <div class="freeCard revealItem" use:reveal={0}>
          <span class="freeIndex">01</span>
          <h3>100% Gratuito</h3>
          <p>Acesso completo a todas as funcionalidades sem mensalidades ou taxas ocultas.</p>
        </div>
        <div class="freeCard revealItem" use:reveal={90}>
          <span class="freeIndex">02</span>
          <h3>Uso Ilimitado</h3>
          <p>Pratique conversas, correções gramaticais e traduções sem limite diário.</p>
        </div>
        <div class="freeCard revealItem" use:reveal={180}>
          <span class="freeIndex">03</span>
          <h3>Para Todos</h3>
          <p>Aprenda no seu ritmo com suporte de IA adaptado ao seu nível de proficiência.</p>
        </div>
      </div>
    </div>
  </section>

  <!-- Seção 3: Inteligência Artificial -->
  <section class="apiSection">
    <div class="sectionContainer">
      <div class="apiBanner revealItem" use:reveal>
        <div class="apiInfo">
          <span class="sectionEyebrow">Tecnologia</span>
          <h2>Potencializado por Inteligência Artificial</h2>
          <p class="sectionSubtitle">
            Executamos o modelo <strong>llama3.2:latest</strong> via Ollama diretamente em nosso servidor local — entregando correções instantâneas e conversação adaptativa sem exigir nada da sua máquina.
          </p>
        </div>

        <div class="apiCard">
          <div class="apiCardHeader">
            <div class="ollamaLogoBox">
              <img src="/Ollama Logo - White - zonalogo.com.svg" alt="Ollama Logo" />
            </div>
            <div>
              <h3>Ollama • llama3.2:latest</h3>
              <span class="apiMeta">Hospedado no Servidor SpeakUp</span>
            </div>
          </div>
          <p>
            Toda a inferência do <code>llama3.2:latest</code> roda no nosso servidor próprio, garantindo respostas rápidas e privacidade direto no seu navegador, sem pesar no seu computador ou celular.
          </p>
        </div>
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
  :global(html) {
    scroll-behavior: smooth;
  }

  main {
    position: relative;
    overflow-x: hidden;
    background-color: #0a0a0a;
    color: white;
    min-height: 100vh;
  }

  /* Navbar superior translúcida */
  .landingNav {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 900;
    background: rgba(10, 10, 12, 0.72);
    backdrop-filter: blur(14px);
    -webkit-backdrop-filter: blur(14px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
    animation: navFadeIn 0.55s ease forwards;
  }

  @keyframes navFadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .navInner {
    max-width: 1180px;
    margin: 0 auto;
    padding: 0.85rem 2rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .navBrand {
    display: inline-flex;
    align-items: center;
    gap: 0.65rem;
    text-decoration: none;
    color: #ffffff;
    font-weight: 800;
    font-size: 1.2rem;
    letter-spacing: -0.02em;
    transition: opacity 0.2s ease;
  }

  .navBrand:hover {
    opacity: 0.88;
  }

  .navActions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .navLoginBtn {
    background: transparent;
    border: 1px solid transparent;
    color: rgba(255, 255, 255, 0.78);
    padding: 0.5rem 1rem;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .navLoginBtn:hover {
    color: #ffffff;
    background: rgba(255, 255, 255, 0.06);
  }

  .navRegisterBtn {
    background: #ffffff;
    color: #09090b;
    border: none;
    padding: 0.52rem 1.1rem;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s ease;
  }

  .navRegisterBtn:hover {
    background: #f2f2f5;
  }

  /* Hero Section */
  .hero {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    z-index: 1;
    position: relative;
    padding: 6.5rem 2rem 3.5rem;
    box-sizing: border-box;
  }

  .heroContainer {
    width: 100%;
    max-width: 1160px;
    margin: auto;
    display: grid;
    grid-template-columns: 1.08fr 0.92fr;
    gap: 3.5rem;
    align-items: center;
  }

  .heroContent {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    text-align: left;
  }

  /* Animação escalonada de entrada no Hero (apenas fade suave sem deslocar) */
  @keyframes heroFade {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .heroAnim1 {
    opacity: 0;
    animation: heroFade 0.65s ease 0.05s forwards;
  }
  .heroAnim2 {
    opacity: 0;
    animation: heroFade 0.65s ease 0.14s forwards;
  }
  .heroAnim3 {
    opacity: 0;
    animation: heroFade 0.65s ease 0.23s forwards;
  }
  .heroAnim4 {
    opacity: 0;
    animation: heroFade 0.65s ease 0.32s forwards;
  }
  .heroAnim5 {
    opacity: 0;
    animation: heroFade 0.75s ease 0.22s forwards;
  }

  /* Pill Badge sutil substituindo a antiga faixa amarela */
  .heroBadge {
    display: inline-flex;
    align-items: center;
    gap: 0.55rem;
    padding: 0.42rem 0.9rem;
    border-radius: 999px;
    background: rgba(92, 109, 255, 0.1);
    border: 1px solid rgba(92, 109, 255, 0.28);
    color: #c7d2fe;
    font-size: 0.82rem;
    font-weight: 600;
    margin-bottom: 1.35rem;
    letter-spacing: 0.01em;
  }

  .badgeDot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: #5c6dff;
    animation: dotPulse 2.4s ease-in-out infinite;
  }

  @keyframes dotPulse {
    0%,
    100% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.35);
      opacity: 0.65;
    }
  }

  .heroTitle {
    font-size: clamp(2.25rem, 4.2vw, 3.35rem);
    font-weight: 800;
    line-height: 1.12;
    letter-spacing: -0.03em;
    margin: 0 0 1.25rem;
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
    animation: chromeShine 7s ease-in-out infinite;
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
    font-size: 1.06rem;
    line-height: 1.65;
    max-width: 540px;
    margin: 0 0 2rem;
    color: rgba(255, 255, 255, 0.68);
  }

  .heroActions {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.9rem;
  }

  .cta {
    display: inline-flex;
    align-items: center;
    gap: 0.55rem;
    background: #ffffff;
    padding: 0.9rem 1.65rem;
    border-radius: 10px;
    font-size: 1rem;
    border: none;
    color: #09090b;
    font-weight: 700;
    cursor: pointer;
  }

  .cta svg {
    transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .cta:hover {
    background: #f4f4f6;
  }

  .cta:hover svg {
    transform: translateX(3px);
  }

  .ctaSecondary {
    background: rgba(255, 255, 255, 0.04);
    color: rgba(255, 255, 255, 0.88);
    padding: 0.9rem 1.35rem;
    font-size: 0.96rem;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 10px;
    font-weight: 600;
    transition: all 0.25s ease;
    cursor: pointer;
  }

  .ctaSecondary:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.22);
    color: #ffffff;
  }

  /* Janela de Demonstração Interativa do Chat no Hero */
  .heroPreview {
    width: 100%;
  }

  .previewWindow {
    background: rgba(18, 18, 22, 0.78);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 14px;
    overflow: hidden;
    transition: border-color 0.35s ease;
  }

  .previewWindow:hover {
    border-color: rgba(92, 109, 255, 0.32);
  }

  .previewHeader {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.8rem 1.1rem;
    background: rgba(255, 255, 255, 0.03);
    border-bottom: 1px solid rgba(255, 255, 255, 0.07);
  }

  .windowDots {
    display: flex;
    gap: 6px;
  }

  .windowDots span {
    width: 9px;
    height: 9px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.18);
  }

  .previewTitle {
    display: inline-flex;
    align-items: center;
    gap: 0.45rem;
    font-size: 0.78rem;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.65);
  }

  .onlineDot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #22c55e;
  }

  .stepPills {
    display: flex;
    gap: 5px;
  }

  .stepPill {
    width: 16px;
    height: 5px;
    border-radius: 99px;
    border: none;
    padding: 0;
    background: rgba(255, 255, 255, 0.16);
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .stepPill.active {
    background: #5c6dff;
    width: 22px;
  }

  .previewBody {
    padding: 1.35rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    min-height: 290px;
  }

  .demoMsg {
    border-radius: 11px;
    padding: 0.95rem 1.1rem;
    text-align: left;
    opacity: 0;
    transition: opacity 0.45s ease;
  }

  .demoMsg.visible {
    opacity: 1;
  }

  .userDemo {
    background: rgba(255, 255, 255, 0.045);
    border: 1px solid rgba(255, 255, 255, 0.09);
  }

  .aiDemo {
    background: rgba(92, 109, 255, 0.08);
    border: 1px solid rgba(92, 109, 255, 0.24);
  }

  .demoMsgHeader {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: rgba(255, 255, 255, 0.48);
    margin-bottom: 0.4rem;
  }

  .cefrTag {
    background: rgba(92, 109, 255, 0.22);
    color: #c7d2fe;
    padding: 0.12rem 0.45rem;
    border-radius: 4px;
    font-size: 0.68rem;
  }

  .demoMsgText {
    margin: 0;
    font-size: 0.94rem;
    line-height: 1.5;
    color: #f4f4f5;
  }

  .errHighlight {
    text-decoration: underline wavy #f87171;
    text-underline-offset: 3px;
    color: #fca5a5;
  }

  .demoCorrectionBox {
    margin-top: 0;
    max-height: 0;
    opacity: 0;
    overflow: hidden;
    border-top: 1px solid transparent;
    transition: max-height 0.5s cubic-bezier(0.16, 1, 0.3, 1),
      opacity 0.4s ease,
      margin-top 0.4s ease,
      padding-top 0.4s ease,
      border-color 0.4s ease;
  }

  .demoCorrectionBox.expanded {
    margin-top: 0.75rem;
    padding-top: 0.75rem;
    max-height: 160px;
    opacity: 1;
    border-top-color: rgba(255, 255, 255, 0.09);
  }

  .demoCorrLine {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.45rem;
    font-size: 0.84rem;
    color: #e4e4e7;
    margin-bottom: 0.35rem;
  }

  .corrBadge {
    background: rgba(34, 197, 94, 0.16);
    color: #86efac;
    border: 1px solid rgba(34, 197, 94, 0.3);
    padding: 0.1rem 0.45rem;
    border-radius: 4px;
    font-size: 0.72rem;
    font-weight: 700;
  }

  .demoExplLine {
    margin: 0;
    font-size: 0.78rem;
    line-height: 1.45;
    color: rgba(255, 255, 255, 0.62);
  }

  .demoTranslation {
    margin: 0.55rem 0 0;
    padding-top: 0.5rem;
    border-top: 1px solid rgba(92, 109, 255, 0.2);
    font-size: 0.8rem;
    color: rgba(255, 255, 255, 0.68);
  }

  /* Indicador de Scroll sutil */
  .scrollHint {
    margin-top: 2.25rem;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    color: rgba(255, 255, 255, 0.45);
    text-decoration: none;
    font-size: 0.82rem;
    font-weight: 500;
    transition: color 0.25s ease;
  }

  .scrollHint:hover {
    color: rgba(255, 255, 255, 0.85);
  }

  .scrollHint svg {
    transition: opacity 0.25s ease;
  }

  /* Scroll Reveal genérico acionado por IntersectionObserver (sem deslocamento vertical) */
  .revealItem {
    opacity: 0;
    transition: opacity 0.65s ease,
      background 0.3s ease,
      border-color 0.3s ease;
  }

  :global(.revealItem.revealed) {
    opacity: 1;
  }

  /* Seções Minimalistas da Tela Inicial */
  .featureSection,
  .freeSection,
  .apiSection {
    padding: 4.5rem 2rem;
    background: #0a0a0a;
    position: relative;
    border-top: 1px solid #181818;
  }

  .sectionContainer {
    max-width: 1080px;
    margin: 0 auto;
    width: 100%;
  }

  .sectionHeader {
    text-align: left;
    margin-bottom: 2.25rem;
  }

  .sectionEyebrow {
    display: inline-block;
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #777777;
    background-color: #141414;
    border: 1px solid #222222;
    padding: 0.22rem 0.6rem;
    border-radius: 4px;
    margin-bottom: 0.75rem;
  }

  .sectionHeader h2,
  .apiInfo h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 0.45rem 0;
    letter-spacing: -0.02em;
    line-height: 1.25;
  }

  .sectionSubtitle {
    font-size: 0.95rem;
    color: #888888;
    max-width: 560px;
    margin: 0;
    line-height: 1.6;
  }

  .sectionSubtitle strong,
  .apiCard code {
    color: #e4e4e7;
    font-weight: 600;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    font-size: 0.86em;
    background-color: #181818;
    border: 1px solid #262626;
    padding: 0.1rem 0.4rem;
    border-radius: 4px;
  }

  .featuresGrid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1.25rem;
  }

  .freeGrid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1.25rem;
  }

  .freeCard {
    background-color: #111111;
    border: 1px solid #222222;
    border-radius: 12px;
    padding: 1.65rem;
    text-align: left;
  }

  :global(.freeCard.revealed:hover) {
    background-color: #151515;
    border-color: #333333;
  }

  .freeIndex {
    display: inline-block;
    font-size: 0.72rem;
    font-weight: 600;
    letter-spacing: 0.05em;
    color: #666666;
    background-color: #181818;
    border: 1px solid #242424;
    padding: 0.18rem 0.5rem;
    border-radius: 4px;
    margin-bottom: 1.1rem;
  }

  .freeCard h3 {
    font-size: 1.1rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 0.45rem 0;
    letter-spacing: -0.01em;
  }

  .freeCard p {
    color: #888888;
    font-size: 0.88rem;
    line-height: 1.6;
    margin: 0;
  }

  /* Seção IA Minimalista (Split Layout) */
  .apiBanner {
    display: grid;
    grid-template-columns: 1.1fr 0.9fr;
    gap: 2.5rem;
    align-items: center;
  }

  .apiInfo {
    text-align: left;
  }

  .apiCard {
    background-color: #111111;
    padding: 1.75rem;
    border-radius: 12px;
    border: 1px solid #222222;
    text-align: left;
    transition: border-color 0.25s ease, background-color 0.25s ease;
  }

  .apiCard:hover {
    background-color: #151515;
    border-color: #333333;
  }

  .apiCardHeader {
    display: flex;
    align-items: center;
    gap: 0.9rem;
    margin-bottom: 1rem;
  }

  .ollamaLogoBox {
    width: 44px;
    height: 44px;
    border-radius: 10px;
    background-color: #181818;
    border: 1px solid #262626;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .ollamaLogoBox img {
    width: 26px;
    height: 26px;
    object-fit: contain;
  }

  .apiCardHeader h3 {
    font-size: 1.15rem;
    font-weight: 700;
    margin: 0 0 0.15rem 0;
    color: #ffffff;
  }

  .apiMeta {
    font-size: 0.76rem;
    color: #777777;
    font-weight: 500;
  }

  .apiCard p {
    font-size: 0.88rem;
    color: #888888;
    line-height: 1.6;
    margin: 0 0 1.25rem 0;
  }

  .apiPills {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    padding-top: 1rem;
    border-top: 1px solid #1c1c1c;
  }

  .apiPills span {
    font-size: 0.72rem;
    font-weight: 600;
    color: #999999;
    background-color: #181818;
    border: 1px solid #242424;
    padding: 0.22rem 0.6rem;
    border-radius: 4px;
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
    from { opacity: 0; transform: scale(0.97); }
    to { opacity: 1; transform: scale(1); }
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

  @media screen and (max-width: 960px) {
    .heroContainer {
      grid-template-columns: 1fr;
      gap: 2.5rem;
      text-align: center;
    }

    .heroContent {
      align-items: center;
      text-align: center;
    }

    .heroActions {
      justify-content: center;
    }

    .featuresGrid,
    .freeGrid,
    .apiBanner {
      grid-template-columns: 1fr;
    }

    .apiBanner {
      gap: 1.75rem;
    }
  }

  @media screen and (max-width: 768px) {
    .navInner {
      padding: 0.75rem 1.15rem;
    }

    .navRegisterBtn {
      padding: 0.45rem 0.85rem;
      font-size: 0.85rem;
    }

    .hero {
      padding: 5.5rem 1.25rem 2.5rem;
    }

    .heroTitle {
      font-size: 2.05rem;
      text-align: center;
    }

    .heroSubtitle {
      font-size: 0.95rem;
      max-width: 95%;
      margin-bottom: 1.75rem;
    }

    .cta {
      font-size: 0.95rem;
      padding: 0.85rem 1.4rem;
    }

    .ctaSecondary {
      font-size: 0.9rem;
      padding: 0.85rem 1.2rem;
    }

    .previewBody {
      padding: 1rem;
      min-height: 260px;
    }

    .featureSection,
    .freeSection,
    .apiSection {
      padding: 3.5rem 1.25rem;
    }

    .sectionHeader h2,
    .apiInfo h2 {
      font-size: 1.5rem;
    }

    .sectionSubtitle {
      font-size: 0.9rem;
    }

    .featuresGrid,
    .freeGrid {
      gap: 1rem;
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

  @media (prefers-reduced-motion: reduce) {
    .heroAnim1,
    .heroAnim2,
    .heroAnim3,
    .heroAnim4,
    .heroAnim5,
    .revealItem {
      opacity: 1 !important;
      transform: none !important;
      animation: none !important;
      transition: none !important;
    }
  }
</style>
