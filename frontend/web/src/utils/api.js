import axios from 'axios';

// Login
export const login = async (email, password) => {
    
    const response = await axios.post('http://localhost:8080/user/login', {
        email,
        password
    });

    const authToken = response.data.token;
    localStorage.setItem('token', authToken);

    return response.data;

}

//Register
export const register = async (userData) => {
    
    const response = await axios.post('http://localhost:8080/user', userData);

    return response.data;

}