import { useState, useEffect, useRef } from 'react';
import * as THREE from 'three';
import styles from './Index.module.css';
import { useAuth } from '../../contexts/Auth';

const cardData = [
  {
    imgSrc: './chat.png',
    alt: 'Ícone de chat representando conversação',
    content: 'Plataforma inovadora que utiliza IA para auxiliar no aprendizado de idiomas. Analisa textos em tempo real, identificando e corrigindo erros com explicações detalhadas para um aprendizado contínuo.'
  },
  {
    imgSrc: './ia.png',
    alt: 'Ícone representando Inteligência Artificial',
    content: 'Oferece experiência personalizada com análise em tempo real de erros gramaticais e de digitação. Fornece correções com explicações para melhor compreensão das regras linguísticas.'
  },
  {
    imgSrc: './mic.png',
    alt: 'Ícone de microfone para reconhecimento de voz',
    content: 'Em breve: funcionalidade de reconhecimento de voz para análise de pronúncia em tempo real, com feedback detalhado para aprimoramento dinâmico.'
  }
];

function Index() {
  const [activeIndex, setActiveIndex] = useState(1);
  const { user } = useAuth();
  const mountRef = useRef(null);

  useEffect(() => {
    if (user) {
      window.location.href = '/home';
    }
  }, [user]);

  useEffect(() => {
    const interval = setInterval(() => {
      setActiveIndex((prevIndex) => (prevIndex + 1) % cardData.length);
    }, 5000);

    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    // Configuração da cena Three.js
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
    const renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true });

    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setPixelRatio(window.devicePixelRatio || 1);
    renderer.setClearColor(0x000000, 0); // Fundo transparente

    // Configura o container das partículas
    if (mountRef.current) {
      mountRef.current.appendChild(renderer.domElement);
      renderer.domElement.style.position = 'absolute';
      renderer.domElement.style.top = '0';
      renderer.domElement.style.left = '0';
      renderer.domElement.style.zIndex = '0';
      renderer.domElement.style.pointerEvents = 'none';
    }

    // Posiciona a câmera
    camera.position.z = 30;

    // Criação das partículas
    const particleCount = 400;
    const particles = new THREE.BufferGeometry();
    const positions = new Float32Array(particleCount * 3);
    const colors = new Float32Array(particleCount * 3);
    const sizes = new Float32Array(particleCount);

    // Preenche os arrays com valores aleatórios
    for (let i = 0; i < particleCount; i++) {
      // Posições aleatórias em um espaço 3D
      positions[i * 3] = (Math.random() - 0.5) * 200;
      positions[i * 3 + 1] = (Math.random() - 0.5) * 200;
      positions[i * 3 + 2] = (Math.random() - 0.5) * 200;

      // Cores brancas
      colors[i * 3] = 1.0; // R
      colors[i * 3 + 1] = 1.0; // G
      colors[i * 3 + 2] = 1.0; // B

      // Tamanhos aleatórios
      sizes[i] = 0.5 + Math.random() * 1.5;
    }

    // Atribui os arrays à geometria
    particles.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    particles.setAttribute('color', new THREE.BufferAttribute(colors, 3));
    particles.setAttribute('size', new THREE.BufferAttribute(sizes, 1));

    // Material das partículas
    const particleMaterial = new THREE.PointsMaterial({
      size: 0.5,
      vertexColors: true,
      transparent: true,
      opacity: 0.8,
      blending: THREE.AdditiveBlending,
      sizeAttenuation: true
    });

    // Cria o sistema de partículas
    const particleSystem = new THREE.Points(particles, particleMaterial);
    scene.add(particleSystem);

    // Animação
    const animate = () => {
      requestAnimationFrame(animate);

      // Rotação suave das partículas
      particleSystem.rotation.x += 0.0005;
      particleSystem.rotation.y += 0.001;

      renderer.render(scene, camera);
    };

    animate();

    // Resize handler
    const handleResize = () => {
      camera.aspect = window.innerWidth / window.innerHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(window.innerWidth, window.innerHeight);
    };

    window.addEventListener('resize', handleResize);

    // Limpeza
    return () => {
      if (mountRef.current && mountRef.current.contains(renderer.domElement)) {
        mountRef.current.removeChild(renderer.domElement);
      }
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  return (
    <>
      <main className={styles.main}>
        <div ref={mountRef} className={styles.particles} />

        <section className={styles.hero}>
          <img src="logo.png" alt="Logo" className={styles.logo} />
          <h1 className={styles.heroTitle}>SpeakUp</h1>
          <p className={styles.heroSubtitle}>
            Aprendizado personalizado com correções, explicações e evolução contínua.
          </p>
          <a href="/login" className={styles.cta}>Vamos começar!</a>
        </section>

        <section className={styles.featureSection}>
          <div className={styles.feature}>
            <img src="./chat.png" alt="Chat" />
            <div>
              <h2>Correções inteligentes em tempo real</h2>
              <p>A plataforma analisa seus textos e corrige com explicações didáticas.</p>
            </div>
          </div>

          <div className={styles.feature}>
            <img src="./ia.png" alt="IA" />
            <div>
              <h2>Experiência personalizada</h2>
              <p>Entende seus erros mais comuns e ajuda você a dominá-los com contexto.</p>
            </div>
          </div>

          <div className={styles.feature}>
            <img src="./mic.png" alt="Voz" />
            <div>
              <h2>Reconhecimento de fala</h2>
              <p>Em breve: pratique a pronúncia com feedback em tempo real.</p>
            </div>
          </div>
        </section>
      </main>

    </>
  );
}

export default Index;