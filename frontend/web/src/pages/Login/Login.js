import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import { useAuth } from '../../contexts/Auth';

// Routes
import { login as loginApi } from '../../utils/api';

// CSS Module
import styles from './Login.module.css';

function Login() {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);
    const { login } = useAuth();

    const goToIndex = () => {
        navigate('/');
    };

    const goToHome = () => {
        navigate('/home');
    };

    const handleLogin = async (event) => {
        event.preventDefault();
        setIsLoading(true);
        const email = event.target.email.value;
        const password = event.target.password.value;
        try {
            const response = await loginApi(email, password);

            if (response.message === 'Login successful') {
                login();
                toast.success('Login realizado com sucesso!');
                setTimeout(() => {
                    goToHome();
                }, 2000);
            } else {
                toast.error('Email ou senha inválidos');
            }
        } catch (err) {
            console.error(err);
            toast.error('Email ou senha inválidos'); // Adicione esta linha
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className={styles.pageLogin}>
            <ToastContainer />
            <form className={styles.form} onSubmit={handleLogin}>
                <div className={styles.card}>
                    <div className={styles.textCenter}>
                        <img
                            src="./logo.png"
                            onClick={goToIndex}
                            alt="SpeakUp Logo"
                            className={styles.logo}
                        />
                        <h2 className={styles.h2}>Bem-vindo de volta!</h2>
                        <p className={styles.p}>Entre para continuar a conversa</p>
                    </div>

                    <div className={styles.inputContainer}>
                        <label htmlFor="email" className={styles.inputLabel}>
                            Email
                        </label>
                        <input
                            id="email"
                            type="email"
                            name="email"
                            placeholder="seu@email.com"
                            required
                            className={styles.input}
                        />
                    </div>

                    <div className={styles.inputContainer}>
                        <label htmlFor="password" className={styles.inputLabel}>
                            Senha
                        </label>
                        <input
                            id="password"
                            type="password"
                            name="password"
                            placeholder="••••••••"
                            required
                            className={styles.input}
                        />
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        className={styles.button}
                    >
                        {isLoading ? 'Entrando...' : 'Entrar'}
                    </button>

                    <div className={styles.textCenter}>
                        <Link to="/register" className={styles.styledLink}>
                            Não tem uma conta ainda?{' '}
                            <span>Registre-se</span>
                        </Link>
                    </div>
                </div>
            </form>
        </div>
    );
}

export default Login;