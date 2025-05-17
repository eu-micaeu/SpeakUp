import { useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import Cookies from 'js-cookie';

import { updateUser } from '../../utils/api'; // ajuste o caminho conforme necessário
import styles from './OnBoarding.module.css';

function OnBoarding() {
    const navigate = useNavigate();
    const [language, setLanguage] = useState('');
    const [level, setLevel] = useState('');
    const [levels, setLevels] = useState([]);
    const [userId, setUserId] = useState(null);

    useEffect(() => {
        const storedUserId = Cookies.get('userId');
        if (storedUserId) {
            setUserId(storedUserId);
        } else {
            console.warn('Usuário não autenticado!');
            navigate('/login');
        }
    }, [navigate]);

    const handleLanguageChange = (e) => {
        const selected = e.target.value;
        setLanguage(selected);
        setLevel('');

        if (selected === 'english') {
            setLevels(['A1', 'A2', 'B1', 'B2', 'C1', 'C2']);
        } else if (selected === 'japanese') {
            setLevels(['N5', 'N4', 'N3', 'N2', 'N1']);
        } else {
            setLevels([]);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            await updateUser(userId, { language, level });
            navigate('/home');
        } catch (error) {
            alert('Erro ao salvar suas preferências.');
        }
    };

    return (
        <div className={styles.pageOnBoarding}>
            <form onSubmit={handleSubmit} className={styles.onboardingForm}>
                <h2>Vamos começar!</h2>

                <label>
                    Idioma que deseja aprender:
                    <select value={language} onChange={handleLanguageChange} required>
                        <option value="">Selecione...</option>
                        <option value="english">Inglês</option>
                        <option value="japanese">Japonês</option>
                    </select>
                </label>

                <label>
                    Seu nível:
                    <select
                        value={level}
                        onChange={(e) => setLevel(e.target.value)}
                        required
                        disabled={!language}
                    >
                        <option value="">Selecione...</option>
                        {levels.map((lvl) => (
                            <option key={lvl} value={lvl}>{lvl}</option>
                        ))}
                    </select>
                </label>

                <button type="submit" disabled={!language || !level}>
                    Continuar
                </button>
            </form>
        </div>
    );
}

export default OnBoarding;
