import styled from 'styled-components';
import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';
import { useState, useEffect } from 'react';

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

// Styled components
const Main = styled.main`
  background-color: #313131;
  padding: 0;
  height: 90vh;
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const Title = styled.h1`
  text-align: center;
  font-size: 3em;
  padding: 30px 0;
  margin: 0;
  color: #fff;
`;

const CardsContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  padding: 0 20px;
  width: 80%;
  position: relative;
  height: 400px;
`;

const Card = styled.div`
  width: 300px;
  height: 60%;
  padding: 40px;
  background-color: #000;
  color: #fff;
  border-radius: 10px;
  transition: all 0.5s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: absolute;
  opacity: 0;
  transform: scale(0.9);
  
  &.active {
    opacity: 1;
    transform: scale(1);
  }
  
  &:hover {
    transform: scale(1.05);
  }
`;

const CardContent = styled.p`
  font-size: 1em;
  margin: 20px;
  text-align: justify;
  line-height: 1.5em;
`;

const DotsContainer = styled.div`
  display: flex;
  justify-content: center;
  gap: 10px;
  margin: 10px 0;
`;

const Dot = styled.div`
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background-color: ${props => props.active ? '#fff' : '#666'};
  cursor: pointer;
  transition: background-color 0.3s;
`;

function Index() {
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentIndex((prevIndex) => (prevIndex + 1) % cardData.length);
    }, 5000);
    
    return () => clearInterval(interval);
  }, []);

  return (
    <>
      <Header />
      
      <Main>
        <Title>SpeakUp</Title>
        
        <CardsContainer>
          {cardData.map((card, index) => (
            <Card 
              key={index}
              className={index === currentIndex ? 'active' : ''}
            >
              <img src={card.imgSrc} alt={card.alt} width={75} height={75} />
              <CardContent>{card.content}</CardContent>
            </Card>
          ))}
        </CardsContainer>

        <DotsContainer>
          {cardData.map((_, index) => (
            <Dot 
              key={index} 
              active={index === currentIndex}
              onClick={() => setCurrentIndex(index)}
            />
          ))}
        </DotsContainer>
      </Main>
      
      <Footer />
    </>
  );
}

export default Index;