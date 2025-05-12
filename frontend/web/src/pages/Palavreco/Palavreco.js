import { useState, useEffect, useRef } from 'react';
import { generateRandomWord } from '../../utils/api';
import Header from '../../components/Header/Header';
import styles from './Palavreco.module.css';

function Palavreco() {
    const [words, setWords] = useState([]);
    const [loading, setLoading] = useState(false);
    const containerRef = useRef(null);
    const lastScrollPos = useRef(0);

    const loadNewWord = async () => {
        try {
            setLoading(true);
            const newWord = await generateRandomWord();
            setWords(prevWords => [...prevWords, newWord]);
        } catch (error) {
            console.error('Erro ao carregar palavra:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        // Carregar primeira palavra ao montar o componente
        loadNewWord();
    }, []);

    const handleScroll = (e) => {
        if (loading) return;

        const container = containerRef.current;
        const currentScrollPos = container.scrollTop;
        const isScrollingDown = currentScrollPos > lastScrollPos.current;
        const scrollHeight = container.scrollHeight;
        const clientHeight = container.clientHeight;
        
        // Se chegou perto do final, carrega nova palavra
        if (isScrollingDown &&
            currentScrollPos + clientHeight > scrollHeight - 100) {
            loadNewWord();
        }

        lastScrollPos.current = currentScrollPos;
    };

    const handleAddToDeck = (word) => {
        // Funcionalidade será implementada posteriormente
        console.log('Adicionar ao deck:', word);
    };

    return (
        <>
            <Header />
            <div
                className={styles.container}
                ref={containerRef}
                onScroll={handleScroll}
            >
                {words.map((word, index) => (
                    <div key={word.id || index} className={styles.wordCard}>
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
                {loading && (
                    <div className={styles.loading}>
                        <div className={styles.loadingSpinner}></div>
                    </div>
                )}
            </div>
        </>
    );
}

export default Palavreco;
