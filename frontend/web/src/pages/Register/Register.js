import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import styled from 'styled-components';

// Routes
import { register } from '../../utils/api';

const PageRegister = styled.div`

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
`;

const Title = styled.h2`
    font-size: 1.875rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
`;

const Subtitle = styled.p`
    color: rgb(156, 163, 175);
    margin-bottom: 2rem;
`;

const ErrorMessage = styled.div`
    background-color: rgb(239, 68, 68);
    color: white;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
`;

const InputRow = styled.div`
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
`;

const InputContainer = styled.div`
    margin-bottom: 1.5rem;
    width: 48%;
`;

const InputContainerFull = styled.div`
    margin-bottom: 1.5rem;
    width: 100%;
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
    box-sizing: border-box; /* Garante que padding e border não aumentem a largura */

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

const Select = styled.select`
    width: 100%;
    padding: 0.75rem;
    border-radius: 0.375rem;
    background-color: rgb(55, 65, 81);
    border: 1px solid rgb(75, 85, 99);
    color: white;
    box-sizing: border-box;
`;

const Option = styled.option`
    color: white;
`;

const StyledLink = styled(Link)`
    display: block;
    margin-top: 10px;
    text-align: center;
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
        
        // Define os níveis com base no idioma selecionado
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

            console.log(userData);

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
    }

    return (
        <PageRegister>
            <Form onSubmit={handleRegister}>
                <Card>
                    <TextCenter>
                        <Logo
                            src="./logo.png"
                            onClick={goToIndex}
                            alt="SpeakUp Logo"
                        />
                        <Title>Crie sua conta</Title>
                        <Subtitle>Entre para continuar a conversa</Subtitle>
                    </TextCenter>

                    {error && (
                        <ErrorMessage>
                            {error}
                        </ErrorMessage>
                    )}

                    <InputRow>
                        <InputContainer>
                            <InputLabel htmlFor="first_name">
                                Primeiro nome
                            </InputLabel>
                            <Input
                                type="first_name"
                                id="first_name"
                                name="first_name"
                                placeholder="Digite seu nome"
                                required
                            />
                        </InputContainer>

                        <InputContainer>
                            <InputLabel htmlFor="last_name">
                                Último nome
                            </InputLabel>
                            <Input
                                type="last_name"
                                id="last_name"
                                name="last_name"
                                placeholder="Digite seu último sobrenome"
                                required
                            />
                        </InputContainer>
                    </InputRow>

                    <InputRow>
                        <InputContainer>
                            <InputLabel htmlFor="language">
                                Idioma que deseja aprender
                            </InputLabel>
                            <Select
                                id="language"
                                name="language"
                                onChange={handleLanguageChange}
                                required
                            >
                                <Option value="" disabled selected>Selecione</Option>
                                <Option value="english">Inglês</Option>
                                <Option value="japanese">Japonês</Option>
                            </Select>
                        </InputContainer>

                        <InputContainer>
                            <InputLabel htmlFor="level">
                                Nível atual no idioma
                            </InputLabel>
                            <Select
                                id="level"
                                name="level"
                                required
                                disabled={!selectedLanguage}
                            >
                                <Option value="" disabled selected>Selecione</Option>
                                {levels.map((level) => (
                                    <Option key={level} value={level}>{level}</Option>
                                ))}
                            </Select>
                        </InputContainer>
                    </InputRow>

                    {/* Restante do seu formulário permanece o mesmo */}
                    <InputContainerFull>
                        <InputLabel htmlFor="email">
                            Email
                        </InputLabel>
                        <Input
                            type="email"
                            id="email"
                            name="email"
                            placeholder="Digite seu email"
                            required
                        />
                    </InputContainerFull>

                    <InputRow>
                        <InputContainer>
                            <InputLabel htmlFor="password">
                                Senha
                            </InputLabel>
                            <Input
                                type="password"
                                id="password"
                                name="password"
                                placeholder="Digite sua senha"
                                required
                            />
                        </InputContainer>

                        <InputContainer>
                            <InputLabel htmlFor="confirm-password">
                                Confirmar Senha
                            </InputLabel>
                            <Input
                                type="password"
                                id="confirm-password"
                                name="confirm-password"
                                placeholder="Confirme sua senha"
                                required
                            />
                        </InputContainer>
                    </InputRow>

                    <Button
                        type="submit"
                        disabled={isLoading}
                    >
                        {isLoading ? 'Carregando...' : 'Cadastrar'}
                    </Button>

                    <StyledLink to="/login">Já tem uma conta? Entre agora!</StyledLink>
                </Card>
            </Form>
        </PageRegister>
    );
}

export default Register;