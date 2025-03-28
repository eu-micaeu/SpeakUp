import './Login.css';
import { useNavigate } from 'react-router-dom';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

// Routes
import { login } from '../../utils/api';

function Login() {

    const navigate = useNavigate();

    const goToIndex = () => {
        navigate('/');
    };

    const goToHome = () => {
        navigate('/home');
    }

    const handleLogin = async (event) => {
        event.preventDefault();
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
                toast.error('Login failed. Please check your credentials.');
            }
        } catch (error) {
            toast.error('An error occurred. Please try again.');
        }
    }

    return (
        <div className="pageLogin">
            <ToastContainer />
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

