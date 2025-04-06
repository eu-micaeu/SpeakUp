// Components
import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';

// Icons
import CreateIcon from '@mui/icons-material/Create';
import BookIcon from '@mui/icons-material/Book';

import { useNavigate } from 'react-router-dom';

// CSS Module
import styles from './Home.module.css';

function Home() {

    const navigate = useNavigate();

    const goToChat = () => {
        navigate('/chat');
    }
    
    console.log("Home component is rendering");

    return (

        <>

            <Header />

            <div className={styles.pageHome}>

                <h1>Menu</h1>

                <div className={styles.cardsContainer}>

                    <div className={styles.card} onClick={() => goToChat()}>

                        <CreateIcon 
                        
                            sx={{

                                fontSize: "5rem",

                            }}
                        
                        ></CreateIcon>

                        <h2>Praticar</h2>

                    </div>

                    <div className={styles.card}>

                        <BookIcon

                            sx={{

                                fontSize: "5rem",

                            }}

                        ></BookIcon>

                        <h2>Plano de Ensino</h2>

                    </div>

                </div>

            </div>

            <Footer />

        </>

    );
}

export default Home;