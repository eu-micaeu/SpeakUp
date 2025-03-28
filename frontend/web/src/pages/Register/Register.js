import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './Register.css';

// Routes
import { register } from '../../utils/api';

function Register() {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    
    const goToIndex = () => {
        navigate('/');
    };

    const handleRegister = async (event) => {
        event.preventDefault();
        setIsLoading(true);
        setError('');

        try {
            const userData = {
                name: event.target.name.value,
                email: event.target.email.value,
                password: event.target.password.value
            };

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
    }

    return (
        <div className="pageRegister">
            <form onSubmit={handleRegister}>
                <div className="card">
                    <div className="text-center">
                        <img 
                            src="./logo.png" 
                            onClick={goToIndex} 
                            alt="SpeakUp Logo" 
                        />
                        <h2>Crie sua conta</h2>
                        <p>Entre para continuar a conversa</p>
                    </div>

                    {error && (
                        <div className="error-message">
                            {error}
                        </div>
                    )}

                    <div className="input-container">
                        <label htmlFor="name" className="input-label">
                            Nome
                        </label>
                        <input 
                            type="text" 
                            id="name" 
                            name="name" 
                            placeholder="Digite seu nome" 
                            required 
                        />
                    </div>

                    <div className="input-container">
                        <label htmlFor="email" className="input-label">
                            Email
                        </label>
                        <input 
                            type="email" 
                            id="email" 
                            name="email" 
                            placeholder="Digite seu email" 
                            required 
                        />
                    </div>

                    <div className="input-container">
                        <label htmlFor="password" className="input-label">
                            Senha
                        </label>
                        <input 
                            type="password" 
                            id="password" 
                            name="password" 
                            placeholder="Digite sua senha" 
                            required 
                        />
                    </div>

                    <button 
                        type="submit" 
                        className="btn-primary"
                        disabled={isLoading}
                    >
                        {isLoading ? 'Carregando...' : 'Cadastrar'}
                    </button>

                    <Link to="/login" className="link">Já tem uma conta? Entre agora!</Link>
                </div>
            </form>

        </div>
    );


}

export default Register;