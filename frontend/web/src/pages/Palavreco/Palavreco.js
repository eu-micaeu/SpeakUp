import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { generateRandomWord } from '../../utils/api';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import styles from './Palavreco.module.css';

function Palavreco() {
    const navigate = useNavigate();
    const [words, setWords] = useState([]);
    const [loading, setLoading] = useState(false);
    const [currentWordIndex, setCurrentWordIndex] = useState(0);
    const containerRef = useRef(null);
    const isLoadingMore = useRef(false);

    const loadMultipleWords = async (count) => {
        if (isLoadingMore.current) return;
        try {
            isLoadingMore.current = true;
            setLoading(true);
            const promises = Array(count).fill().map(() => generateRandomWord());
            const newWords = await Promise.all(promises);
            setWords(prevWords => [...prevWords, ...newWords]);
        } catch (error) {
            console.error('Erro ao carregar palavras:', error);
        } finally {
            setLoading(false);
            isLoadingMore.current = false;
        }
    };

    useEffect(() => {
        // Carregar 3 palavras iniciais
        loadMultipleWords(3);
    }, []);

    const handleScroll = (e) => {
        if (loading || !containerRef.current) return;

        const container = containerRef.current;
        const scrollTop = container.scrollTop;
        const totalHeight = container.scrollHeight;
        
        // Calcula a porcentagem do scroll
        const scrollPercentage = scrollTop / (totalHeight - container.clientHeight);
        
        // Calcula o índice baseado na porcentagem do scroll e número total de palavras
        const currentIndex = Math.round(scrollPercentage * (words.length - 1));
        
        console.log({
            scrollTop,
            totalHeight,
            scrollPercentage,
            currentIndex,
            totalWords: words.length
        });

        // Atualiza o índice atual
        setCurrentWordIndex(currentIndex);

        // Se estiver na penúltima palavra, carrega mais
        if (currentIndex === words.length - 2) {
            console.log('Carregando mais palavras na posição:', currentIndex);
            loadMultipleWords(3);
        }
    };

    const handleAddToDeck = (word) => {
        // Funcionalidade será implementada posteriormente
        console.log('Adicionar ao deck:', word);
    };

    return (
        <div className={styles.pageContainer}>
            <ArrowBackIcon
                className={styles.backArrow}
                onClick={() => navigate('/home')}
            />
            <div className={styles.progressIndicator}>
                {words.length > 0 && `Palavra ${currentWordIndex + 1} de ${words.length}`}
            </div>
            <div
                className={styles.container}
                ref={containerRef}
                onScroll={handleScroll}
            >
                {words.map((word, index) => (
                    <div
                        key={word.id || index}
                        className={`${styles.wordCard} ${index === currentWordIndex ? styles.currentWord : ''}`}
                    >
                        <h2>{word.word}</h2>
                        <div className={styles.translation}>{word.word_translated}</div>
                        <div className={styles.context}>{word.context_in_phrase}</div>
                        <button
                            className={styles.addButton}
                            onClick={() => handleAddToDeck(word)}
                        >
                            Adicionar ao Baralho
                        </button>
                    </div>
                ))}
            </div>
            {loading && (
                <div className={styles.loading}>
                    <div className={styles.loadingSpinner}></div>
                </div>
            )}
        </div>
    );
}

export default Palavreco;
