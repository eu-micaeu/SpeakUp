import axios from 'axios';
import Cookies from 'js-cookie';

const API_URL = process.env.REACT_APP_API_URL

// Login
export const login = async (email, password) => {

    const response = await axios.post(API_URL + '/user/login', {
        email,
        password
    });

    const authToken = response.data.token;
    Cookies.set('authToken', authToken);

    return response.data;

}

export const register = async (userData) => {
     
    const response = await axios.post(API_URL + '/user/', userData);

    console.log(response)

    return response.data;

}
// Create chat
export const createChat = async (topic) => {

    const response = await axios.post(API_URL + '/chat', {
        topic
    }, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    return response.data;
}

// Delete chat
export const deleteChat = async (chatId) => {

    const response = await axios.delete(API_URL + `/chat/${chatId}`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });
    return response.data;
}

// Get chats by user id
export const getChatsByUserId = async () => {

    const response = await axios.get(API_URL + `/chat/user`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    return response.data;

}

// Get messages by chat id
export const getMessagesByChatId = async (chatId) => {

    const response = await axios.get(API_URL + `/message/chat/${chatId}`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    return response.data;

}

// Add message to chat
export const addMessageToChat = async (chat_id, content, sender, type) => {

    const response = await axios.post(API_URL + `/message`, {
        chat_id,
        content,
        sender, 
        type
    }, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });

    return response.data;

}

// generate AI response dialog
export const generateAIResponseDialog = async (message) => {
    try {
        const response = await axios.post(API_URL + `/ai/generate-response-dialog`, {
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

// generate AI response correction
export const generateAIResponseCorrection = async (message) => {
    try {
        const response = await axios.post(API_URL + `/ai/generate-response-correction`, {
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

// generate AI response translation
export const generateAIResponseTranslation = async (message) => {
    try {
        const response = await axios.post(API_URL + `/ai/generate-response-translation`, {
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

// generate AI response topic
export const generateAIResponseTopic = async (message) => {
    try {
        const response = await axios.post(API_URL + `/ai/generate-response-topic`, {
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