import { useState, useEffect, useRef } from 'react';
import * as THREE from 'three';
import styles from './Index.module.css';
import { useAuth } from '../../contexts/Auth';
import { FeatureCard } from '../../components/FeatureCard/FeatureCard';
import { Footer } from '../../components/Footer/Footer';

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
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
    const renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true });

    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setPixelRatio(window.devicePixelRatio || 1);
    renderer.setClearColor(0x000000, 0); // Fundo transparente

    if (mountRef.current) {
      mountRef.current.appendChild(renderer.domElement);
      renderer.domElement.style.position = 'absolute';
      renderer.domElement.style.top = '0';
      renderer.domElement.style.left = '0';
      renderer.domElement.style.zIndex = '0';
      renderer.domElement.style.pointerEvents = 'none';
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

    particles.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    particles.setAttribute('color', new THREE.BufferAttribute(colors, 3));
    particles.setAttribute('size', new THREE.BufferAttribute(sizes, 1));

    const particleMaterial = new THREE.PointsMaterial({
      size: 0.5,
      vertexColors: true,
      transparent: true,
      opacity: 0.8,
      blending: THREE.AdditiveBlending,
      sizeAttenuation: true
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

    window.addEventListener('resize', handleResize);

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
          {cardData.map((card, index) => (
            <FeatureCard
              key={index}
              imgSrc={card.imgSrc}
              alt={card.alt}
              title={card.content.split('.')[0]}
              content={card.content}
            />
          ))}
        </section>
      </main>

      <Footer />

    </>
  );
}

export default Index;