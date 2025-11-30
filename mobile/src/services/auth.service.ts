import { apiClient } from './api';
import { storage } from '../utils/storage';
import type {
    LoginRequest,
    LoginResponse,
    RegisterRequest,
    RegisterResponse,
    User,
} from '../types';

export const authService = {
    async login(data: LoginRequest): Promise<LoginResponse> {
        const response = await apiClient.post<LoginResponse>('/api/user/login', data);
        if (response.data.token) {
            await storage.setToken(response.data.token);

            // Decodificar token e extrair dados do usuário
            const userData = this.decodeToken(response.data.token);
            if (userData) {
                await storage.setUserData(userData);
            }
        }
        return response.data;
    },

    decodeToken(token: string): User | null {
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(
                atob(base64)
                    .split('')
                    .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                    .join('')
            );
            const decoded = JSON.parse(jsonPayload);

            return {
                id: decoded.user_id,
                name: decoded.name,
                email: decoded.email,
                language: decoded.language,
                level: decoded.level,
            };
        } catch (error) {
            console.error('❌ [Auth] Error decoding token:', error);
            return null;
        }
    },

    async register(data: RegisterRequest): Promise<RegisterResponse> {
        const response = await apiClient.post<RegisterResponse>('/api/user/', data);
        return response.data;
    },

    async logout(): Promise<void> {
        await storage.clear();
    },

    async getUserProfile(): Promise<User> {
        // Primeiro tentar obter do storage
        const cachedUser = await storage.getUserData();
        if (cachedUser) {
            return cachedUser;
        }

        // Se não tiver no storage, tentar decodificar do token
        const token = await storage.getToken();
        if (token) {
            const userData = this.decodeToken(token);
            if (userData) {
                await storage.setUserData(userData);
                return userData;
            }
        }

        // Se ainda não tiver, buscar do backend (usando o ID do token)
        const response = await apiClient.get<User>('/api/user/me');
        await storage.setUserData(response.data);
        return response.data;
    },

    async updateUserProfile(data: Partial<User>): Promise<User> {
        // Obter ID do usuário atual
        const currentUser = await this.getUserProfile();
        const response = await apiClient.put<User>(`/api/user/${currentUser.id}`, data);
        await storage.setUserData(response.data);
        return response.data;
    },

    async deleteAccount(userId: string): Promise<void> {
        await apiClient.delete(`/api/user/${userId}`);
        await storage.clear();
    },

    async isAuthenticated(): Promise<boolean> {
        const token = await storage.getToken();
        return !!token;
    },
};
