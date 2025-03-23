import './Login.css';

import {login} from '../../utils/api';

function Login() {

    const handleLogin = async (event) => {
        event.preventDefault();
        const email = event.target[0].value;
        const password = event.target[1].value;
        const response = await login(email, password);

        if (response.message === 'Login successful') {
            window.location.href = '/home';
        }

    }
    
    return (

        <div className="pageLogin">

            <form onSubmit={handleLogin}>

                <img src="./logo.png"></img>
                <input type="text" placeholder="Email" />
                <input type="password" placeholder="Password" />
                <button type="submit">Login</button>

                <a href="#">Não possui conta? Faça o Registro aqui!</a>

            </form>

        </div>

    );
}

export default Login;