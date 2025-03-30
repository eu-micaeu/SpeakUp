import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast, ToastContainer } from 'react-toastify';
import styled from 'styled-components';
import 'react-toastify/dist/ReactToastify.css';

// Routes
import { login } from '../../utils/api';

const PageLogin = styled.div`
    
`;

const Form = styled.form`
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: rgb(0, 0, 0);
    color: white;
    height: 100vh;
`;

const Card = styled.div`
    width: 100%;
    max-width: 28rem;
    padding: 2rem;
    background-color: rgb(31, 31, 31);
    border-radius: 0.5rem;
`;

const TextCenter = styled.div`
    text-align: center;
`;

const Logo = styled.img`
    height: 3rem;
    margin: 0 auto 1rem;
    cursor: pointer;
    transition: 1s ease all;

    &:hover {
        transform: scale(1.1);
        cursor: pointer;
    }
`;

const H2 = styled.h2`
    font-size: 1.875rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
`;

const P = styled.p`
    color: rgb(156, 163, 175);
    margin-bottom: 2rem;
`;

const InputContainer = styled.div`
    margin-bottom: 1.5rem;
`;

const InputLabel = styled.label`
    display: block;
    font-size: 0.875rem;
    color: rgb(209, 213, 219);
    margin-bottom: 0.25rem;
`;

const Input = styled.input`
    width: 100%;
    padding: 0.75rem;
    border-radius: 0.375rem;
    background-color: rgb(55, 65, 81);
    border: 1px solid rgb(75, 85, 99);
    color: white;
    box-sizing: border-box;

    &::placeholder {
        color: rgb(156, 163, 175);
    }
`;

const Button = styled.button`
    width: 100%;
    padding: 0.75rem;
    margin: 0 0 1.5rem 0;
    border-radius: 0.375rem;
    border: none;
    background-color: rgb(53, 53, 53);
    color: white;
    cursor: pointer;
    transition: background-color 0.2s;
    box-sizing: border-box;

    &:hover:not(:disabled) {
        background-color: rgb(37, 99, 235);
    }

    &:disabled {
        opacity: 0.75;
        cursor: not-allowed;
    }
`;

const ErrorMessage = styled.div`
    background-color: rgb(239, 68, 68);
    color: white;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
`;

const StyledLink = styled(Link)`
    color: rgb(156, 163, 175);
    text-decoration: none;
    font-size: 0.875rem;

    span {
        color: rgb(59, 130, 246);
        font-weight: 500;
    }

    span:hover {
        color: rgb(37, 99, 235);
    }
`;

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
        <PageLogin>
            <ToastContainer />
            <Form onSubmit={handleLogin}>
                <Card>
                    <TextCenter>
                        <Logo
                            src="./logo.png"
                            onClick={goToIndex}
                            alt="SpeakUp Logo"
                        />
                        <H2>Bem-vindo de volta!</H2>
                        <P>Entre para continuar a conversa</P>
                    </TextCenter>

                    {error && (
                        <ErrorMessage>
                            {error}
                        </ErrorMessage>
                    )}

                    <InputContainer>
                        <InputLabel htmlFor="email">
                            Email
                        </InputLabel>
                        <Input
                            id="email"
                            type="email"
                            name="email"
                            placeholder="seu@email.com"
                            required
                        />
                    </InputContainer>

                    <InputContainer>
                        <InputLabel htmlFor="password">
                            Senha
                        </InputLabel>
                        <Input
                            id="password"
                            type="password"
                            name="password"
                            placeholder="••••••••"
                            required
                        />
                    </InputContainer>

                    <Button
                        type="submit"
                        disabled={isLoading}
                    >
                        {isLoading ? 'Entrando...' : 'Entrar'}
                    </Button>

                    <TextCenter>
                        <StyledLink to="/register">
                            Não tem uma conta ainda?{' '}
                            <span>Registre-se</span>
                        </StyledLink>
                    </TextCenter>
                </Card>
            </Form>
        </PageLogin>
    );
}

export default Login;


