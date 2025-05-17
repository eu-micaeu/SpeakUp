import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { getAuth, GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { initializeApp } from "firebase/app";

import { useAuth } from '../../contexts/Auth';
import { login as loginApi, register as registerApi } from '../../utils/api';
import styles from './Login.module.css';

// Configuração do Firebase: depois dos imports
const firebaseConfig = {
    apiKey: process.env.REACT_APP_FIREBASE_API_KEY,
    authDomain: process.env.REACT_APP_FIREBASE_AUTH_DOMAIN,
    projectId: process.env.REACT_APP_FIREBASE_PROJECT_ID,
    storageBucket: process.env.REACT_APP_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: process.env.REACT_APP_FIREBASE_MESSAGING_SENDER_ID,
    appId: process.env.REACT_APP_FIREBASE_APP_ID,
    measurementId: process.env.REACT_APP_FIREBASE_MEASUREMENT_ID
};

// Inicialização do Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

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

            if (response && response.user) {
                login(response.user);
                toast.success('Login realizado com sucesso!');
                setTimeout(() => {
                    goToHome();
                }, 2000);
            } else {
                toast.error('Email ou senha inválidos');
            }
        } catch (err) {
            console.error(err);
            toast.error('Email ou senha inválidos');
        } finally {
            setIsLoading(false);
        }
    };

    const handleGoogleLogin = async () => {
        const provider = new GoogleAuthProvider();
        try {
            const result = await signInWithPopup(auth, provider);
            const user = result.user;

            const email = user.email;
            const password = user.uid;
            const name = user.displayName;

            let isFirstTime = false;

            try {
                await loginApi(email, password);
            } catch (err) {
                isFirstTime = true;
                await registerApi({ name, email, password });
                await loginApi(email, password);
            }

            login({ email, name }); // Armazene mais dados se quiser

            toast.success(`Bem-vindo, ${name}!`);

            setTimeout(() => {
                if (isFirstTime) {
                    navigate('/onboarding'); // redireciona para tela de seleção de level/language
                } else {
                    goToHome();
                }
            }, 2000);

        } catch (error) {
            console.error(error);
            toast.error('Erro ao fazer login com Google');
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

                    <button
                        type="button"
                        onClick={handleGoogleLogin}
                        disabled={isLoading}
                        className={styles.buttonGoogle}
                    >
                        <img
                            src="./google.png"
                            alt="Google Icon"
                            className={styles.googleIcon}
                        />
                    </button>

                    <div className={styles.textCenter}>
                        <Link to="/register" className={styles.styledLink}>
                            Não tem uma conta ainda? <span>Registre-se</span>
                        </Link>
                    </div>
                </div>
            </form>
        </div>
    );
}

export default Login;
