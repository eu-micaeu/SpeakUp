<script lang="ts">
  import { onMount } from "svelte";
  import * as THREE from "three";
  import FeatureCard from "../components/FeatureCard.svelte";
  import Footer from "../components/Footer.svelte";
  import { isAuthTokenValid } from "../utils/cookies";
  import { goto } from "$app/navigation";
  import { env } from "$env/dynamic/public";
  import { getAuth, GoogleAuthProvider, signInWithPopup } from "firebase/auth";
  import { initializeApp } from "firebase/app";
  import { login as loginApi, register as registerApi } from "../utils/api";

  const cardData = [
    {
      imgSrc: "./chat.png",
      alt: "Ícone de chat representando conversação",
      title:
        "Plataforma inovadora que utiliza IA para auxiliar no aprendizado de idiomas.",
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
      imgSrc: "./mic.png",
      alt: "Ícone de microfone para reconhecimento de voz",
      title:
        "Reconhecimento de voz avançado com análise de pronúncia em tempo real.",
      content:
        "Transforme sua fala em texto e receba feedback instantâneo sobre pronúncia, entonação e fluência. Pratique conversação de forma natural e aprimore suas habilidades de comunicação oral com suporte completo da IA.",
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

  // Firebase config
  const isFirebaseConfigured =
    env.PUBLIC_FIREBASE_API_KEY &&
    env.PUBLIC_FIREBASE_API_KEY !== "your_firebase_api_key_here";

  const firebaseConfig = {
    apiKey: env.PUBLIC_FIREBASE_API_KEY || "",
    authDomain: env.PUBLIC_FIREBASE_AUTH_DOMAIN || "",
    projectId: env.PUBLIC_FIREBASE_PROJECT_ID || "",
    storageBucket: env.PUBLIC_FIREBASE_STORAGE_BUCKET || "",
    messagingSenderId: env.PUBLIC_FIREBASE_MESSAGING_SENDER_ID || "",
    appId: env.PUBLIC_FIREBASE_APP_ID || "",
    measurementId: env.PUBLIC_FIREBASE_MEASUREMENT_ID || "",
  };

  let app: any = null;
  let auth: any = null;

  if (isFirebaseConfigured) {
    try {
      app = initializeApp(firebaseConfig);
      auth = getAuth(app);
    } catch (error) {
      console.warn("Firebase initialization failed:", error);
    }
  }

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
          goto("/chat");
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

  async function handleGoogleLogin() {
    if (!isFirebaseConfigured || !auth) {
      showToast("Login com Google não está configurado", "error");
      return;
    }

    const provider = new GoogleAuthProvider();
    try {
      const result = await signInWithPopup(auth, provider);
      const user = result.user;
      const email = user.email!;
      const password = user.uid;
      const name = user.displayName!;

      let isFirstTime = false;

      try {
        await loginApi(email, password);
      } catch (err) {
        isFirstTime = true;
        const userData = {
          name,
          email,
          password,
          language: "en",
          level: "beginner",
        };
        await registerApi(userData);
        await loginApi(email, password);
      }

      showToast(`Bem-vindo, ${name}!`, "success");
      closeAuthModal();

      setTimeout(() => {
        goto("/chat");
      }, 1500);
    } catch (error) {
      console.error(error);
      showToast("Erro ao fazer login com Google", "error");
    }
  }

  let mountRef: HTMLDivElement;

  onMount(() => {
    if (isAuthTokenValid()) {
      goto("/chat");
      return;
    }

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(
      75,
      window.innerWidth / window.innerHeight,
      0.1,
      1000,
    );
    const renderer = new THREE.WebGLRenderer({
      alpha: true,
      antialias: true,
      powerPreference: "high-performance",
    });

    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2)); // Limita pixel ratio
    renderer.setClearColor(0x000000, 0);

    if (mountRef) {
      mountRef.appendChild(renderer.domElement);
      renderer.domElement.style.position = "absolute";
      renderer.domElement.style.top = "0";
      renderer.domElement.style.left = "0";
      renderer.domElement.style.zIndex = "0";
      renderer.domElement.style.pointerEvents = "none";
    }

    camera.position.z = 30;

    const particleCount = 200; // Reduzido de 400 para 200
    const particles = new THREE.BufferGeometry();
    const positions = new Float32Array(particleCount * 3);
    const colors = new Float32Array(particleCount * 3);
    const sizes = new Float32Array(particleCount);

    for (let i = 0; i < particleCount; i++) {
      positions[i * 3] = (Math.random() - 0.5) * 200;
      positions[i * 3 + 1] = (Math.random() - 0.5) * 200;
      positions[i * 3 + 2] = (Math.random() - 0.5) * 200;

      colors[i * 3] = 1.0;
      colors[i * 3 + 1] = 1.0;
      colors[i * 3 + 2] = 1.0;

      sizes[i] = 0.5 + Math.random() * 1.5;
    }

    particles.setAttribute("position", new THREE.BufferAttribute(positions, 3));
    particles.setAttribute("color", new THREE.BufferAttribute(colors, 3));
    particles.setAttribute("size", new THREE.BufferAttribute(sizes, 1));

    const particleMaterial = new THREE.PointsMaterial({
      size: 0.5,
      vertexColors: true,
      transparent: true,
      opacity: 0.8,
      blending: THREE.AdditiveBlending,
      sizeAttenuation: true,
    });

    const particleSystem = new THREE.Points(particles, particleMaterial);
    scene.add(particleSystem);

    let animationId: number;
    let isVisible = true;

    // Detectar visibilidade da página
    const handleVisibilityChange = () => {
      isVisible = !document.hidden;
      if (isVisible && !animationId) {
        animate();
      }
    };

    // Detectar scroll para pausar quando fora da hero section
    const handleScroll = () => {
      const heroSection = document.querySelector(".hero");
      if (heroSection) {
        const rect = heroSection.getBoundingClientRect();
        isVisible = rect.bottom > 0 && rect.top < window.innerHeight;
      }
    };

    const animate = () => {
      if (!isVisible) {
        animationId = 0;
        return;
      }

      animationId = requestAnimationFrame(animate);
      particleSystem.rotation.x += 0.0005;
      particleSystem.rotation.y += 0.001;
      renderer.render(scene, camera);
    };

    animate();

    const handleResize = () => {
      camera.aspect = window.innerWidth / window.innerHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(window.innerWidth, window.innerHeight);
      renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    };

    window.addEventListener("resize", handleResize);
    document.addEventListener("visibilitychange", handleVisibilityChange);
    window.addEventListener("scroll", handleScroll, { passive: true });

    return () => {
      // Limpeza adequada
      if (animationId) {
        cancelAnimationFrame(animationId);
      }

      window.removeEventListener("resize", handleResize);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("scroll", handleScroll);

      if (mountRef && mountRef.contains(renderer.domElement)) {
        mountRef.removeChild(renderer.domElement);
      }

      // Liberar recursos do Three.js
      particles.dispose();
      particleMaterial.dispose();
      renderer.dispose();
      scene.clear();
    };
  });
</script>

<main>
  <div bind:this={mountRef} class="particles"></div>

  {#if toastMessage}
    <div class="toast {toastType}">
      {toastMessage}
    </div>
  {/if}

  <section class="hero">
    <img src="logo.png" alt="Logo" width={60} />
    <h1 class="heroTitle">SpeakUp</h1>
    <p class="heroSubtitle">
      Aprendizado personalizado com correções, explicações e evolução contínua.
    </p>
    <button on:click={() => openAuthModal("register")} class="cta"
      >Vamos começar!</button
    >
    <button on:click={() => openAuthModal("login")} class="ctaSecondary">
      Já possui uma conta? Faça o login aqui!
    </button>
  </section>
  <section class="featureSection">
    <div class="sectionHeader">
      <h1>Recursos do SpeakUp</h1>
      <p class="sectionSubtitle">Descubra tudo o que você pode fazer.</p>
    </div>
    <div class="featuresGrid">
      {#each cardData as card, i}
        <FeatureCard {...card} index={i} />
      {/each}
    </div>
  </section>

  <section class="apiSection">
    <div class="apiContainer">
      <h1>Potencializado pelas melhores IAs do mercado</h1>
      <div class="apiLogos">
        <div class="apiCard">
          <img src="./openai-logo.png" alt="OpenAI Logo" />
          <h3>OpenAI GPT</h3>
          <p>
            Tecnologia de ponta para compreensão e geração de linguagem natural
          </p>
        </div>
        <div class="apiCard">
          <img src="./gemini-logo.png" alt="Google Gemini Logo" />
          <h3>Google Gemini</h3>
          <p>
            IA multimodal avançada do Google para análise contextual profunda
          </p>
        </div>
      </div>
      <p class="apiDescription">
        Combinamos as APIs mais avançadas de inteligência artificial disponíveis
        para oferecer correções precisas e contextualizadas detalhadas que se
        adaptam ao seu nível de conhecimento. Cada interação é otimizada para
        maximizar seu aprendizado.
      </p>
    </div>
  </section>

  <!-- Auth Modal -->
  {#if showAuthModal}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="modalOverlay" on:click={closeAuthModal}>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
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
          <img src="./logo.png" alt="SpeakUp Logo" class="modalLogo" />
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

            {#if isFirebaseConfigured}
              <button
                type="button"
                on:click={handleGoogleLogin}
                disabled={isLoading}
                class="googleButton"
              >
                <img src="./google.png" alt="Google" />
                Continuar com Google
              </button>
            {/if}
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
  }

  .ctaSecondary:hover {
    transform: translateY(-2px);
  }

  .particles {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 0;
    pointer-events: none;
  }

  .featureSection {
    padding: 6rem 2rem;
    background: linear-gradient(180deg, #0a0a0a 0%, #121212 50%, #0a0a0a 100%);
    position: relative;
  }

  .featureSection::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.2),
      transparent
    );
  }

  .sectionHeader {
    text-align: center;
    margin-bottom: 4rem;
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
    font-size: 3rem;
    font-weight: 800;
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

  .sectionSubtitle {
    font-size: 1.2rem;
    opacity: 0.7;
    color: #e0e0e0;
  }

  .featuresGrid {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .apiSection {
    background: linear-gradient(135deg, #1a1a1a 0%, #0a0a0a 100%);
    padding: 4rem 2rem;
    border-top: 1px solid #333;
  }

  .apiContainer {
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
  }

  .apiContainer h1 {
    font-size: 2.5rem;
    margin-bottom: 3rem;
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

  .apiLogos {
    display: flex;
    justify-content: center;
    gap: 4rem;
    margin-bottom: 3rem;
    flex-wrap: wrap;
  }

  .apiCard {
    background: rgba(255, 255, 255, 0.05);
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    transition: all 0.3s ease;
    width: 280px;
    backdrop-filter: blur(10px);
  }

  .apiCard:hover {
    transform: translateY(-10px);
    border-color: rgba(255, 255, 255, 0.3);
    box-shadow: 0 10px 30px rgba(255, 255, 255, 0.1);
  }

  .apiCard img {
    width: 80px;
    height: 80px;
    margin-bottom: 1.5rem;
    filter: brightness(0) invert(1);
  }

  .apiCard h3 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
    color: #ffffff;
  }

  .apiCard p {
    font-size: 0.95rem;
    opacity: 0.7;
    line-height: 1.6;
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
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modalContent {
    background: linear-gradient(135deg, #1f1f1f 0%, #2a2a2a 100%);
    border-radius: 24px;
    padding: 2.5rem;
    max-width: 550px;
    width: 90%;
    max-height: 90vh;
    overflow-y: auto;
    position: relative;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  }

  .closeButton {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: rgba(255, 255, 255, 0.1);
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: white;
  }

  .closeButton:hover {
    background: rgba(255, 255, 255, 0.2);
  }

  .modalHeader {
    text-align: center;
    margin-bottom: 2rem;
  }

  .modalLogo {
    width: 50px;
    height: 50px;
    margin-bottom: 1rem;
  }

  .modalHeader h2 {
    font-size: 1.8rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
    color: white;
  }

  .modalSubtitle {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.95rem;
  }

  .authTabs {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
    background: rgba(255, 255, 255, 0.05);
    padding: 0.4rem;
    border-radius: 12px;
  }

  .authTab {
    flex: 1;
    padding: 0.75rem;
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.6);
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.3s;
  }

  .authTab.active {
    background: rgba(255, 255, 255, 0.1);
    color: white;
  }

  .authTab:hover:not(.active) {
    color: rgba(255, 255, 255, 0.8);
  }

  .inputGroup {
    margin-bottom: 1.25rem;
  }

  .inputGroup label {
    display: block;
    color: rgba(255, 255, 255, 0.8);
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }

  .inputGroup input,
  .inputGroup select {
    width: 100%;
    padding: 0.875rem;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 10px;
    color: white;
    font-size: 1rem;
    transition: all 0.3s;
    box-sizing: border-box;
  }

  .inputGroup input:focus,
  .inputGroup select:focus {
    outline: none;
    border-color: rgba(255, 255, 255, 0.3);
    background: rgba(255, 255, 255, 0.08);
  }

  .inputGroup input::placeholder {
    color: rgba(255, 255, 255, 0.4);
  }

  .inputGroup select {
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23ffffff' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.875rem center;
    background-size: 12px;
    padding-right: 2.5rem;
    cursor: pointer;
  }

  .inputGroup select:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .inputGroup select option {
    background: #1a1a1a;
    color: white;
    padding: 0.5rem;
  }

  .inputGroup select option:disabled {
    color: rgba(255, 255, 255, 0.5);
  }

  .inputRow {
    display: flex;
    gap: 1rem;
  }

  .inputRow .inputGroup {
    flex: 1;
  }

  .submitButton {
    width: 100%;
    padding: 1rem;
    background: linear-gradient(135deg, #4a4a4a 0%, #353535 100%);
    border: none;
    border-radius: 10px;
    color: white;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    margin-bottom: 1rem;
  }

  .submitButton:hover:not(:disabled) {
    background: linear-gradient(135deg, #5a5a5a 0%, #454545 100%);
    transform: translateY(-2px);
    box-shadow: 0 10px 30px rgba(255, 255, 255, 0.1);
  }

  .submitButton:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .googleButton {
    width: 100%;
    padding: 1rem;
    background: white;
    border: none;
    border-radius: 10px;
    color: #333;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
  }

  .googleButton img {
    width: 20px;
    height: 20px;
  }

  .googleButton:hover:not(:disabled) {
    background: #f5f5f5;
    transform: translateY(-2px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  }

  .googleButton:disabled {
    opacity: 0.6;
    cursor: not-allowed;
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

  @media screen and (max-width: 768px) {
    .particles {
      display: none;
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
  }
</style>
