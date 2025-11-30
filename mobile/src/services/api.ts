import axios, { AxiosInstance } from 'axios';
import { API_URL } from '../constants';
import { storage } from '../utils/storage';

class ApiClient {
    private client: AxiosInstance;

    constructor() {
        console.log('🌐 [API] Initializing API client with baseURL:', API_URL);

        this.client = axios.create({
            baseURL: API_URL,
            timeout: 30000,
            headers: {
                'Content-Type': 'application/json',
            },
        });

        // Request interceptor to add token
        this.client.interceptors.request.use(
            async (config) => {
                console.log('📤 [API] Request:', config.method?.toUpperCase(), config.url);
                const token = await storage.getToken();
                if (token) {
                    config.headers.Authorization = `Bearer ${token}`;
                    console.log('🔑 [API] Token attached');
                }
                return config;
            },
            (error) => {
                console.error('❌ [API] Request error:', error);
                return Promise.reject(error);
            }
        );

        // Response interceptor for error handling
        this.client.interceptors.response.use(
            (response) => {
                console.log('📥 [API] Response:', response.status, response.config.url);
                return response;
            },
            async (error) => {
                console.error('❌ [API] Response error:', error.message);
                if (error.response) {
                    console.error('❌ [API] Status:', error.response.status);
                    console.error('❌ [API] Data:', error.response.data);
                }
                if (error.response?.status === 401) {
                    console.log('🔓 [API] Unauthorized, clearing storage');
                    await storage.clear();
                    // Redirect to login handled by the app
                }
                return Promise.reject(error);
            }
        );
    }

    getClient(): AxiosInstance {
        return this.client;
    }
}

export const apiClient = new ApiClient().getClient();
