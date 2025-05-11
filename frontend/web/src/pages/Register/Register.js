import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast, ToastContainer } from 'react-toastify';
import { register } from '../../utils/api';
import styles from './Register.module.css';

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

    const isValidPassword = (password) => {
        const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
        return regex.test(password);
    };

    const handleRegister = async (event) => {
        event.preventDefault();
        setIsLoading(true);
        setError('');

        try {
            const userData = {
                name: event.target.name.value,
                email: event.target.email.value,
                password: event.target.password.value,
                language: event.target.language.value,
                level: event.target.level.value,
            };

            const confirmPassword = event.target['confirm-password'].value;

            if (userData.password !== confirmPassword) {
                setError('As senhas não coincidem');
                setIsLoading(false);
                return;
            }

            if (!isValidPassword(userData.password)) {
                setError('A senha deve ter pelo menos 8 caracteres, uma letra maiúscula, uma minúscula, um número e um caractere especial.');
                setIsLoading(false);
                return;
            }

            const response = await register(userData);

            if (response.message === 'User created successfully') {
                toast.success('Usuário criado com sucesso!');
                setTimeout(() => {
                    navigate('/login');
                }, 2000);
            } else {
                setError('Erro ao criar usuário.');
            }
        } catch (err) {

            if (err.response && err.response.status === 409) {
                setError('Email já cadastrado.');
            } else {
                setError('Erro ao criar usuário.');
            }

        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className={styles.pageRegister}>
            <ToastContainer />
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

                    <div className={styles.inputContainerFull}>
                        <label htmlFor="name" className={styles.inputLabel}>
                            Nome completo
                        </label>
                        <input
                            type="text"
                            id="name"
                            name="name"
                            placeholder="Digite seu nome"
                            required
                            className={styles.input}
                        />
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
