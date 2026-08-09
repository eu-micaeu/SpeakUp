import axios, { AxiosError } from 'axios';
import Cookies from 'js-cookie';

// For Vite / SvelteKit projects, use environment variable VITE_API_URL or PUBLIC_API_URL
export const API_URL = import.meta.env.VITE_API_URL || import.meta.env.PUBLIC_API_URL || '/api';

// Types
interface LoginResponse {
    token: string;
    message: string;
}

interface RegisterResponse {
    message: string;
}

interface UserData {
    name: string;
    email: string;
    password: string;
    language: string;
    level: string;
}

interface Chat {
    id: string;
    topic: string;
}

interface Message {
    id: string;
    chat_id: string;
    content: string;
    sender: string;
    type: string;
}

interface AIResponse {
    message: string;
    response: string;
    explanation?: string;
}

interface User {
    id: string;
    name: string;
    email: string;
    language: string;
    level: string;
}

export interface BillingStatus {
    stripe_customer_id: string;
    stripe_subscription_id: string;
    stripe_price_id: string;
    stripe_status: string;
    stripe_current_period_end: number;
}

export interface AIUsageStatus {
    is_pro: boolean;
    daily_limit: number;
    used_today: number;
    remaining: number;
}

interface Word {
    id: string;
    word: string;
    translation: string;
    user_id: string;
}

// Login
export const login = async (email: string, password: string): Promise<{ token: string }> => {
    try {
        const response = await axios.post<LoginResponse>(API_URL + '/user/login', {
            email,
            password
        });

        const { token, message } = response.data;

        if (token) {
            Cookies.set('authToken', token);
            return { token };
        }

        throw new Error('Dados de usuário não encontrados na resposta');
    } catch (error) {
        console.error('Erro no login:', error);
        throw error;
    }
}

// Register
export const register = async (userData: UserData): Promise<RegisterResponse> => {
    const response = await axios.post<RegisterResponse>(API_URL + '/user/', userData);
    return response.data;
}

// Create chat
export const createChat = async (topic: string): Promise<Chat> => {
    const response = await axios.post<Chat>(API_URL + '/chat', {
        topic
    }, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });
    return response.data;
}

// Delete chat
export const deleteChat = async (chatId: string): Promise<void> => {
    const response = await axios.delete(API_URL + `/chat/${chatId}`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });
    return response.data;
}

// Get chats by user id
export const getChatsByUserId = async (): Promise<Chat[]> => {
    const response = await axios.get<Chat[]>(API_URL + `/chat/user`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });
    return response.data;
}

// Get messages by chat id
export const getMessagesByChatId = async (chatId: string): Promise<Message[]> => {
    const response = await axios.get<Message[]>(API_URL + `/message/chat/${chatId}`, {
        headers: {
            Authorization: `Bearer ${Cookies.get('authToken')}`
        }
    });
    return response.data;
}

// Add message to chat
export const addMessageToChat = async (chat_id: string, content: string, sender: string, type: string): Promise<Message> => {
    const response = await axios.post<Message>(API_URL + `/message`, {
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

const getAIProvider = (): string => {
    if (typeof window === 'undefined') return 'ollama';
    const stored = localStorage.getItem('aiProvider');
    if (!stored || stored === 'gemini') {
        localStorage.setItem('aiProvider', 'ollama');
        return 'ollama';
    }
    return stored;
};

// generate AI response dialog
export const generateAIResponseDialog = async (message: string, chatId: string): Promise<AIResponse | undefined> => {
    try {
        const response = await axios.post<AIResponse>(API_URL + `/ai/generate-response-dialog`, {
            message,
            chat_id: chatId
        }, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`,
                'X-AI-Provider': getAIProvider()
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        console.log(error);
        return undefined;
    }
}

// generate AI response correction
export const generateAIResponseCorrection = async (message: string): Promise<AIResponse | undefined> => {
    try {
        const response = await axios.post<AIResponse>(API_URL + `/ai/generate-response-correction`, {
            message
        }, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`,
                'X-AI-Provider': getAIProvider()
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        console.log(error);
        return undefined;
    }
}

// generate AI response translation
export const generateAIResponseTranslation = async (message: string): Promise<AIResponse | undefined> => {
    try {
        const response = await axios.post<AIResponse>(API_URL + `/ai/generate-response-translation`, {
            message
        }, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`,
                'X-AI-Provider': getAIProvider()
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        console.log(error);
        return undefined;
    }
}

// generate AI response topic
export const generateAIResponseTopic = async (message: string): Promise<AIResponse | undefined> => {
    try {
        const response = await axios.post<AIResponse>(API_URL + `/ai/generate-response-topic`, {
            message
        }, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`,
                'X-AI-Provider': getAIProvider()
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        console.log(error);
        return undefined;
    }
}

// generate random word
export const generateRandomWord = async (): Promise<Word> => {
    try {
        const response = await axios.post<Word>(API_URL + `/ai/generate-random-word`, {}, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        console.error('Erro ao gerar palavra:', error);
        throw error;
    }
}

// get user words
export const getUserWords = async (): Promise<Word[]> => {
    try {
        const response = await axios.get<Word[]>(API_URL + `/word/user`, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        console.error('Erro ao buscar palavras:', error);
        throw error;
    }
}

// get user by id
export const getUserById = async (userId: string): Promise<User> => {
    try {
        if (!userId) {
            throw new Error('ID do usuário é obrigatório');
        }

        const response = await axios.get<{ user: User }>(API_URL + `/user/${userId}`, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });

        return response.data.user;
    }
    catch (error: unknown) {
        const axiosError = error as AxiosError;
        console.error('Erro ao buscar usuário:', axiosError.response || axiosError);
        throw error;
    }
}

// update user 
export const updateUser = async (userId: string, userData: Partial<UserData>): Promise<User> => {
    try {
        if (!userId) {
            throw new Error('ID do usuário é obrigatório');
        }

        const response = await axios.put<User>(API_URL + `/user/${userId}`, userData, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        const axiosError = error as AxiosError;
        console.error('Erro ao atualizar usuário:', axiosError.response || axiosError);
        throw error;
    }
}

// delete user
export const deleteUser = async (userId: string): Promise<void> => {
    try {
        if (!userId) {
            throw new Error('ID do usuário é obrigatório');
        }

        const response = await axios.delete(API_URL + `/user/${userId}`, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`
            }
        });

        return response.data;
    }
    catch (error: unknown) {
        const axiosError = error as AxiosError;
        console.error('Erro ao deletar usuário:', axiosError.response || axiosError);
        throw error;
    }
}

// Transcribe audio using Whisper API
export const transcribeAudio = async (audioBlob: Blob): Promise<string> => {
    try {
        const formData = new FormData();
        formData.append('audio', audioBlob, 'audio.webm');

        const response = await axios.post<{ text: string }>(API_URL + '/ai/transcribe-audio', formData, {
            headers: {
                Authorization: `Bearer ${Cookies.get('authToken')}`,
                'Content-Type': 'multipart/form-data'
            }
        });

        return response.data.text;
    }
    catch (error: unknown) {
        console.error('Erro ao transcrever áudio:', error);
        throw error;
    }
}

// Billing status
export const getBillingStatus = async (): Promise<BillingStatus> => {
    return {
        stripe_customer_id: 'free_user',
        stripe_subscription_id: 'free_subscription',
        stripe_price_id: 'free_price',
        stripe_status: 'active',
        stripe_current_period_end: 9999999999
    };
}

// Create Stripe checkout session
export const createCheckoutSession = async (_plan: 'monthly' | 'annual'): Promise<{ url: string }> => {
    throw new Error('O SpeakUp é 100% gratuito!');
}

// Create Stripe customer portal session
export const createPortalSession = async (): Promise<{ url: string }> => {
    throw new Error('O SpeakUp é 100% gratuito!');
}

// AI usage status
export const getAIUsageStatus = async (): Promise<AIUsageStatus> => {
    return {
        is_pro: true,
        daily_limit: -1,
        used_today: 0,
        remaining: -1
    };
}

