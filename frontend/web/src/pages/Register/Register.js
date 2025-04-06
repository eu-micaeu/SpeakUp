import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import styles from './Register.module.css';

// Routes
import { register } from '../../utils/api';

function Register() {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const [selectedLanguage, setSelectedLanguage] = useState('');
    const [levels, setLevels] = useState([]);

    const goToIndex = () => {
        navigate('/');
    };

    const handleLanguageChange = (e) => {
        const language = e.target.value;
        setSelectedLanguage(language);
        
        if (language === 'english') {
            setLevels(['A1', 'A2', 'B1', 'B2', 'C1', 'C2']);
        } else if (language === 'japanese') {
            setLevels(['N5', 'N4', 'N3', 'N2', 'N1']);
        } else {
            setLevels([]);
        }
    };

    const handleRegister = async (event) => {
        event.preventDefault();
        setIsLoading(true);
        setError('');

        try {
            const userData = {
                first_name: event.target.first_name.value,
                last_name: event.target.last_name.value,
                email: event.target.email.value,
                password: event.target.password.value,
                language: event.target.language.value,
                level: event.target.level.value,
            };

            if (userData.password !== event.target['confirm-password'].value) {
                setError('As senhas não coincidem');
                setIsLoading(false);
                return;
            }

            const response = await register(userData);

            if (response.message === 'User created successfully') {
                navigate('/login');
            } else {
                setError('Erro ao criar usuário');
            }
        } catch (err) {
            setError('Ocorreu um erro ao fazer login. Tente novamente.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className={styles.pageRegister}>
            <form className={styles.form} onSubmit={handleRegister}>
                <div className={styles.card}>
                    <div className={styles.textCenter}>
                        <img
                            src="./logo.png"
                            onClick={goToIndex}
                            alt="SpeakUp Logo"
                            className={styles.logo}
                        />
                        <h2 className={styles.title}>Crie sua conta</h2>
                        <p className={styles.subtitle}>Entre para continuar a conversa</p>
                    </div>

                    {error && (
                        <div className={styles.errorMessage}>
                            {error}
                        </div>
                    )}

                    <div className={styles.inputRow}>
                        <div className={styles.inputContainer}>
                            <label htmlFor="first_name" className={styles.inputLabel}>
                                Primeiro nome
                            </label>
                            <input
                                type="text"
                                id="first_name"
                                name="first_name"
                                placeholder="Digite seu nome"
                                required
                                className={styles.input}
                            />
                        </div>

                        <div className={styles.inputContainer}>
                            <label htmlFor="last_name" className={styles.inputLabel}>
                                Último nome
                            </label>
                            <input
                                type="text"
                                id="last_name"
                                name="last_name"
                                placeholder="Digite seu último sobrenome"
                                required
                                className={styles.input}
                            />
                        </div>
                    </div>

                    <div className={styles.inputRow}>
                        <div className={styles.inputContainer}>
                            <label htmlFor="language" className={styles.inputLabel}>
                                Idioma que deseja aprender
                            </label>
                            <select
                                id="language"
                                name="language"
                                onChange={handleLanguageChange}
                                required
                                className={styles.select}
                            >
                                <option value="" disabled selected>Selecione</option>
                                <option value="english">Inglês</option>
                                <option value="japanese">Japonês</option>
                            </select>
                        </div>

                        <div className={styles.inputContainer}>
                            <label htmlFor="level" className={styles.inputLabel}>
                                Nível atual no idioma
                            </label>
                            <select
                                id="level"
                                name="level"
                                required
                                disabled={!selectedLanguage}
                                className={styles.select}
                            >
                                <option value="" disabled selected>Selecione</option>
                                {levels.map((level) => (
                                    <option key={level} value={level}>{level}</option>
                                ))}
                            </select>
                        </div>
                    </div>

                    <div className={styles.inputContainerFull}>
                        <label htmlFor="email" className={styles.inputLabel}>
                            Email
                        </label>
                        <input
                            type="email"
                            id="email"
                            name="email"
                            placeholder="Digite seu email"
                            required
                            className={styles.input}
                        />
                    </div>

                    <div className={styles.inputRow}>
                        <div className={styles.inputContainer}>
                            <label htmlFor="password" className={styles.inputLabel}>
                                Senha
                            </label>
                            <input
                                type="password"
                                id="password"
                                name="password"
                                placeholder="Digite sua senha"
                                required
                                className={styles.input}
                            />
                        </div>

                        <div className={styles.inputContainer}>
                            <label htmlFor="confirm-password" className={styles.inputLabel}>
                                Confirmar Senha
                            </label>
                            <input
                                type="password"
                                id="confirm-password"
                                name="confirm-password"
                                placeholder="Confirme sua senha"
                                required
                                className={styles.input}
                            />
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        className={styles.button}
                    >
                        {isLoading ? 'Carregando...' : 'Cadastrar'}
                    </button>

                    <Link to="/login" className={styles.styledLink}>
                        Já tem uma conta? <span>Entre agora!</span>
                    </Link>
                </div>
            </form>
        </div>
    );
}

export default Register;