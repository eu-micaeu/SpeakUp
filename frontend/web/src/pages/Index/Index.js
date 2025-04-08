import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';
import { useState, useEffect } from 'react';
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
  const [activeIndex, setActiveIndex] = useState(1); // Carta do meio como padrão
  const { user } = useAuth();

  useEffect(() => {
    if (user) {
      window.location.href = '/home';
    }
  }, [user]);

  useEffect(() => {
    const particlesContainer = document.querySelector(`.${styles.particles}`);
    const numberOfParticles = 50;

    for (let i = 0; i < numberOfParticles; i++) {
      const particle = document.createElement('div');
      particle.className = styles.particle;

      const size = Math.random() * 5 + 1;
      const left = Math.random() * 100;
      const animationDuration = Math.random() * 5 + 3;
      const delay = Math.random() * 5;

      particle.style.width = `${size}px`;
      particle.style.height = `${size}px`;
      particle.style.left = `${left}%`;
      particle.style.animationDuration = `${animationDuration}s`;
      particle.style.animationDelay = `${delay}s`;

      particlesContainer.appendChild(particle);
    }
  }, []);

  return (
    <>
      <Header />
      <main className={styles.main}>
        <div className={styles.particles}></div>
        <h1 className={styles.title}>SpeakUp</h1>
        <div className={styles.cardsContainer}>
          {cardData.map((card, index) => (
            <div
              key={index}
              className={`${styles.card} ${index === activeIndex ? styles.active : ''}`}
              onClick={() => setActiveIndex(index)}
            >
              <img
                src={card.imgSrc}
                alt={card.alt}
                width={75}
                height={75}
                className={styles.cardImage}
              />
              <p className={styles.cardContent}>{card.content}</p>
            </div>
          ))}
        </div>
        <div className={styles.dotsContainer}>
          {cardData.map((_, index) => (
            <div
              key={index}
              className={`${styles.dot} ${index === activeIndex ? styles.active : ''}`}
              onClick={() => setActiveIndex(index)}
            />
          ))}
        </div>
      </main>
      <Footer />
    </>
  );
}

export default Index;