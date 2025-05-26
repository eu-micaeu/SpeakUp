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

                <h2>Aprenda e pratique inglês com uma inteligência artificial!</h2>

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
