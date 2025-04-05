// Components
import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';

import {useNavigate} from 'react-router-dom';

import styled from 'styled-components';

const PageHome = styled.div`
    display: flex;
    background-color: rgb(0, 0, 0);
    color: white;
    height: 100vh;
    text-align: center;
    margin: 10px 0;
    padding: 40px 0;
    flex-direction: column;
`;

const CardsContainer = styled.div`
    display: flex;
    justify-content: center;
    align-items: center;
    flex-wrap: wrap;
    gap: 20px;
    margin: 20px;
    padding: 20px;
    background-color: #000;
    border-radius: 10px;
`;

const Card = styled.div`
    background-color: #1a1a1a;
    border-radius: 10px;
    width: 250px;
    height: 200px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    transition: ease 0.3s all;

    &:hover {
        cursor: pointer;
        transform: scale(1.05);
    }
`;

function Home() {

    const navigate = useNavigate();

    const goToChat = () => {
        navigate('/chat');
    }
    
    console.log("Home component is rendering");

    return (

        <>

            <Header />

            <PageHome>

                <h1>Menu</h1>

                <CardsContainer>

                    <Card onClick={() => goToChat()}>

                        <h3>IA</h3>

                    </Card>

                    <Card>

                        <h3>Plano de ensino</h3>

                    </Card>

                    <Card>

                        <h3>Card 3</h3>

                    </Card>

                </CardsContainer>

            </PageHome>

            <Footer />

        </>

    );
}

export default Home;