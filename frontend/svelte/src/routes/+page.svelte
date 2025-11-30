<script>
  import { onMount } from "svelte";
  import * as THREE from "three";
  import FeatureCard from "../components/FeatureCard.svelte";
  import Footer from "../components/Footer.svelte";
  import { isAuthTokenValid } from "../utils/cookies";
  import { goto } from "$app/navigation";

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

  let mountRef;

  onMount(() => {
    if (isAuthTokenValid()) {
      goto("/chat");
    }

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(
      75,
      window.innerWidth / window.innerHeight,
      0.1,
      1000,
    );
    const renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true });

    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setPixelRatio(window.devicePixelRatio || 1);
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

    const particleCount = 400;
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

    const animate = () => {
      requestAnimationFrame(animate);
      particleSystem.rotation.x += 0.0005;
      particleSystem.rotation.y += 0.001;
      renderer.render(scene, camera);
    };

    animate();

    const handleResize = () => {
      camera.aspect = window.innerWidth / window.innerHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(window.innerWidth, window.innerHeight);
    };

    window.addEventListener("resize", handleResize);

    return () => {
      if (mountRef && mountRef.contains(renderer.domElement)) {
        mountRef.removeChild(renderer.domElement);
      }
      window.removeEventListener("resize", handleResize);
    };
  });
</script>

<main>
  <div bind:this={mountRef} class="particles"></div>
  <section class="hero">
    <img src="logo.png" alt="Logo" width={60} />
    <h1 class="heroTitle">SpeakUp</h1>
    <p class="heroSubtitle">
      Aprendizado personalizado com correções, explicações e evolução contínua.
    </p>
    <a href="/register" class="cta">Vamos começar!</a>
    <a href="/login" class="ctaSecondary"
      >Já possui uma conta? Faça o login aqui!</a
    >
  </section>
  <section class="featureSection">
    <h1>Recursos</h1>
    {#each cardData as card}
      <FeatureCard {...card} />
    {/each}
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
        para oferecer correções precisas e contextualizadas
        detalhadas que se adaptam ao seu nível de conhecimento. Cada interação é
        otimizada para maximizar seu aprendizado.
      </p>
    </div>
  </section>
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
    text-decoration: none;
    color: rgb(0, 0, 0);
    font-weight: bold;
    transition: 0.3s;
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
    text-decoration: none;
    font-weight: bold;
    transition: 0.3s;
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
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
    background-color: #121212;
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
  }
</style>
