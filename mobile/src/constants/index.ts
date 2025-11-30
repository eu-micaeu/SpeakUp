export const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8082';

console.log('🔧 [Config] API_URL:', API_URL);

export const ROUTES = {
    LOGIN: '/(auth)/login',
    REGISTER: '/(auth)/register',
    HOME: '/(tabs)/home',
    CHAT: '/(tabs)/chat',
    PROFILE: '/(tabs)/profile',
    WORDS: '/(tabs)/words',
} as const;

export const STORAGE_KEYS = {
    AUTH_TOKEN: 'authToken',
    USER_DATA: 'userData',
} as const;
