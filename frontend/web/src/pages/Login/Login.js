import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './Login.css';

// Routes
import { login } from '../../utils/api';

function Login() {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    const goToIndex = () => {
        navigate('/');
    };



    const handleLogin = async (event) => {
        event.preventDefault();
        setIsLoading(true);
        setError('');

        try {
            const email = event.target.email.value;
            const password = event.target.password.value;
            const response = await login(email, password);

            if (response.message === 'Login successful') {
                navigate('/home');
            } else {
                setError('Email ou senha inválidos');
            }
        } catch (err) {
            setError('Ocorreu um erro ao fazer login. Tente novamente.');
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div className="pageLogin">
            <form onSubmit={handleLogin}>
                <div className="card">
                    <div className="text-center">
                        <img 
                            src="./logo.png" 
                            onClick={goToIndex} 
                            alt="SpeakUp Logo" 
                        />
                        <h2>Bem-vindo de volta!</h2>
                        <p>Entre para continuar a conversa</p>
                    </div>

                    {error && (
                        <div className="error-message">
                            {error}
                        </div>
                    )}

                    <div className="input-container">
                        <label htmlFor="email" className="input-label">
                            Email
                        </label>
                        <input
                            id="email"
                            type="email"
                            name="email"
                            placeholder="seu@email.com"
                            required
                        />
                    </div>

                    <div className="input-container">
                        <label htmlFor="password" className="input-label">
                            Senha
                        </label>
                        <input
                            id="password"
                            type="password"
                            name="password"
                            placeholder="••••••••"
                            required
                        />
                    </div>

                    <button 
                        type="submit"
                        disabled={isLoading}
                    >
                        {isLoading ? 'Entrando...' : 'Entrar'}
                    </button>

                    <div className="text-center">
                        <Link to="/register">
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
