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
// Create chat
export const createChat = async (topic) => {

    console.log(topic)

    const response = await axios.post('http://localhost:8080/chat', {
        topic
    }, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    return response.data;
}

// Get chats by user id
export const getChatsByUserId = async () => {

    const response = await axios.get(`http://localhost:8080/chat/user`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    return response.data;

}

// Get messages by chat id
export const getMessagesByChatId = async (chatId) => {

    const response = await axios.get(`http://localhost:8080/message/chat/${chatId}`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    console.log(response.data)

    return response.data;

}

// Add message to chat
export const addMessageToChat = async (chat_id, content, sender) => {

    const response = await axios.post(`http://localhost:8080/message`, {
        chat_id,
        content,
        sender
    }, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    return response.data;

}

// generate AI response
export const generateAIResponse = async (message) => {
    try {
        const response = await axios.post(`http://localhost:8080/ai/generate-response`, {
            message
        }, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });

        return response.data;
    }
    catch (error) {
        console.log(error)
    }
}