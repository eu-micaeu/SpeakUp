import axios from 'axios';
import Cookies from 'js-cookie';

const API_URL = process.env.REACT_APP_API_URL

// Login
export const login = async (email, password) => {
    try {
        const response = await axios.post(API_URL + '/user/login', {
            email,
            password
        });
        
        const { token, message } = response.data;
        
        if (token) {
            Cookies.set('authToken', token);
            
            // Extrair o user_id do token JWT
            try {
                const tokenParts = token.split('.');
                const tokenPayload = JSON.parse(atob(tokenParts[1]));
                
                if (tokenPayload.user_id) {
                    Cookies.set('userId', tokenPayload.user_id);
                    return {
                        token,
                        user: {
                            id: tokenPayload.user_id
                        }
                    };
                }
            } catch (tokenError) {
                console.error('Erro ao decodificar token:', tokenError);
                throw new Error('Erro ao processar token');
            }
        }
        
        throw new Error('Dados de usuário não encontrados na resposta');
    } catch (error) {
        console.error('Erro no login:', error);
        throw error;
    }
}

// Register
export const register = async (userData) => {
    const response = await axios.post(API_URL + '/user/', userData);
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
export const generateAIResponseDialog = async (message, chatId) => {
    try {
        const response = await axios.post(API_URL + `/ai/generate-response-dialog`, {
            message,
            chat_id: chatId
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

// generate random word
export const generateRandomWord = async (previousWords = []) => {
    try {
        const response = await axios.post(API_URL + `/ai/generate-random-word`, {
            previousWords: previousWords.slice(-10) // Envia apenas as últimas 10 palavras
        }, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });

        return response.data;
    }
    catch (error) {
        console.error('Erro ao gerar palavra:', error);
        throw error;
    }
}

// get user by id
export const getUserById = async (userId) => {
    try {
        if (!userId) {
            throw new Error('ID do usuário é obrigatório');
        }

        const response = await axios.get(API_URL + `/user/${userId}`, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });
        
        return response.data;
    }
    catch (error) {
        console.error('Erro ao buscar usuário:', error.response || error);
        throw error;
    }
}

// update user 
export const updateUser = async (userId, userData) => {
    try {
        if (!userId) {
            throw new Error('ID do usuário é obrigatório');
        }

        const response = await axios.put(API_URL + `/user/${userId}`, userData, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });
        
        return response.data;
    }
    catch (error) {
        console.error('Erro ao atualizar usuário:', error.response || error);
        throw error;
    }
}