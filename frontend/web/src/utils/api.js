import axios from 'axios';
import Cookies from 'js-cookie';

// Login
export const login = async (email, password) => {
    
    const response = await axios.post('http://localhost:8080/user/login', {
        email,
        password
    });

    const authToken = response.data.token;
    Cookies.set('authToken', authToken);

    return response.data;

}

