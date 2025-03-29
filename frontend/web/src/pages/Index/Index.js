import styled from 'styled-components';
import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';

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
  padding: 50px 0;
  margin: 0;
  color: #fff;
`;

const CardsContainer = styled.div`
  display: flex;
  justify-content: center;
  gap: 20px;
  padding: 0 20px;
  width: 60%;
`;

const Card = styled.div`
  width: 30%;
  padding: 30px;
  background-color: #000;
  color: #fff;
  border-radius: 15px;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  align-items: center;

  &:hover {
    transform: scale(1.05);
  }
`;

const CardContent = styled.p`
  font-size: 1em;
  margin: 20px 10px;
  text-align: justify;
  line-height: 1.6em;
`;

function Index() {
  return (
    <>
      <Header />
      
      <Main>
        <Title>SpeakUp</Title>
        
        <CardsContainer>
          {cardData.map((card, index) => (
            <Card key={index}>
              <img src={card.imgSrc} alt={card.alt} width={75} height={75} />
              <CardContent>{card.content}</CardContent>
            </Card>
          ))}
        </CardsContainer>

      </Main>
      
      <Footer />
    </>
  );
}

export default Index;