export interface User {
    id: string;
    name: string;
    email: string;
    language: string;
    level: string;
}

export interface Chat {
    id: string;
    user_id?: string;
    topic: string;
    start_time?: string;
}

export interface Message {
    id: string;
    chat_id: string;
    content: string;
    sender: 'user' | 'ai';
    type: 'request' | 'response';
    timestamp?: string;
}

export interface Word {
    id: string;
    word: string;
    translation: string;
    user_id: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    message: string;
}

export interface RegisterRequest {
    name: string;
    email: string;
    password: string;
    language: string;
    level: string;
}

export interface RegisterResponse {
    message: string;
}

export interface AIResponse {
    message?: string;
    response: string;
}
