// Components
import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';

// Icons
import CreateIcon from '@mui/icons-material/Create';
import BookIcon from '@mui/icons-material/Book';

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

                        <CreateIcon 
                        
                            sx={{

                                fontSize: "5rem",

                            }}
                        
                        ></CreateIcon>

                        <h2>Praticar.</h2>

                    </Card>

                    <Card>

                        <BookIcon

                            sx={{

                                fontSize: "5rem",

                            }}

                        ></BookIcon>

                        <h2>Plano de Ensino.</h2>

                    </Card>

                </CardsContainer>

            </PageHome>

            <Footer />

        </>

    );
}

export default Home;