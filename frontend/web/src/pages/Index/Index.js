import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';
import { useState, useEffect } from 'react';
import styles from './Index.module.css';

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
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentIndex((prevIndex) => (prevIndex + 1) % cardData.length);
    }, 2500);
    
    return () => clearInterval(interval);
  }, []);

  return (
    <>
      <Header />
      
      <main className={styles.main}>

        <h1 className={styles.title}>SpeakUp</h1>
        
        <div className={styles.cardsContainer}>
          {cardData.map((card, index) => (
            <div 
              key={index}
              className={`${styles.card} ${index === currentIndex ? styles.active : ''}`}
            >
              <img src={card.imgSrc} alt={card.alt} width={75} height={75} />
              <p className={styles.cardContent}>{card.content}</p>
            </div>
          ))}
        </div>

        <div className={styles.dotsContainer}>
          {cardData.map((_, index) => (
            <div 
              key={index} 
              className={`${styles.dot} ${index === currentIndex ? styles.active : ''}`}
              onClick={() => setCurrentIndex(index)}
            />
          ))}
        </div>
      </main>
      
      <Footer />
    </>
  );
}

export default Index;