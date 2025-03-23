import './Login.css';
import { useNavigate } from 'react-router-dom';

// Routes
import { login } from '../../utils/api';

function Login() {

    const navigate = useNavigate();

    const goToIndex = () => {
        navigate('/');
    };

    const handleLogin = async (event) => {
        event.preventDefault();
        const email = event.target.email.value;
        const password = event.target.password.value;
        const response = await login(email, password);

        if (response.message === 'Login successful') {
            navigate('/home');
        }

    }

    return (
        <div className="pageLogin">
            <form onSubmit={handleLogin}>
                <img src="./logo.png" onClick={goToIndex} alt="Logo da empresa" />
                <input type="text" placeholder="Email" name="email" />
                <input type="password" placeholder="Password" name="password" />
                <button type="submit">Login</button>

                <a href="#">Não possui conta? Faça o Registro aqui!</a>
            </form>
        </div>
    );
}

export default Login;
