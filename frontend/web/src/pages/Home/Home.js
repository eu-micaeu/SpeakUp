import { useNavigate } from 'react-router-dom';
import CreateIcon from '@mui/icons-material/Create';
import BookIcon from '@mui/icons-material/Book';
import AutoStoriesIcon from '@mui/icons-material/AutoStories';

import Header from '../../components/Header/Header';
import { Footer } from '../../components/Footer/Footer';

import styles from './Home.module.css';

function Home() {
    
    const navigate = useNavigate();

    const handleNavigate = (path) => {
        navigate(path);
    };

    return (
        <>
            <Header />
            <div className={styles.pageHome}>

                <h1>SpeakUp</h1>

                <h2>Aprenda e pratique inglês de uma inteligência artificial!</h2>

                <p>O SpeakUp é um aplicativo de aprendizado de inglês que utiliza inteligência artificial para ajudar os usuários a praticar conversação e melhorar suas habilidades linguísticas. Com uma interface amigável e recursos interativos, o SpeakUp torna o aprendizado de inglês mais acessível e envolvente.</p>

                <h3>Menu</h3>

                <p>Escolha uma das opções abaixo para começar a praticar seu inglês!</p>

                {/* Cards */}
                <div className={styles.cardsContainer}>
                    <div className={styles.card} onClick={() => handleNavigate('/chat')}>
                        <CreateIcon sx={{ fontSize: '3rem' }} />
                        <h2>Praticar</h2>
                    </div>

                    <div className={styles.card} onClick={() => handleNavigate('/teaching-plan')}>
                        <BookIcon sx={{ fontSize: '3rem' }} />
                        <h2>Plano de Ensino</h2>
                    </div>

                    {/* <div className={styles.card} onClick={() => handleNavigate('/palavreco')}>
                        <AutoStoriesIcon sx={{ fontSize: '3rem' }} />
                        <h2>Palavreco</h2>
                    </div> */}

                </div>
            </div>
            <Footer />
        </>
    );
}

export default Home;
