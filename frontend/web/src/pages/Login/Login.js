import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

// Routes
import { login } from '../../utils/api';

function Login() {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    const goToIndex = () => {
        navigate('/');
    };

    const goToHome = () => {
        navigate('/home');
    }

    const handleLogin = async (event) => {
        event.preventDefault();
        setError('');
        setIsLoading(true);
        const email = event.target.email.value;
        const password = event.target.password.value;
        try {
            const response = await login(email, password);

            if (response.message === 'Login successful') {
                toast.success('Login successful!');
                setTimeout(() => {
                    goToHome();
                }, 3000);
            } else {
                setError('Login failed. Please check your credentials.');
                toast.error('Login failed. Please check your credentials.');
            }
        } catch (err) {
            setError('An error occurred. Please try again.');
            toast.error('An error occurred. Please try again.');
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div className="pageLogin">
            <ToastContainer />
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

