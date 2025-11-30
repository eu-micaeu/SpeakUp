import * as SecureStore from 'expo-secure-store';
import { STORAGE_KEYS } from '../constants';

export const storage = {
    async setToken(token: string): Promise<void> {
        await SecureStore.setItemAsync(STORAGE_KEYS.AUTH_TOKEN, token);
    },

    async getToken(): Promise<string | null> {
        return await SecureStore.getItemAsync(STORAGE_KEYS.AUTH_TOKEN);
    },

    async removeToken(): Promise<void> {
        await SecureStore.deleteItemAsync(STORAGE_KEYS.AUTH_TOKEN);
    },

    async setUserData(data: any): Promise<void> {
        await SecureStore.setItemAsync(STORAGE_KEYS.USER_DATA, JSON.stringify(data));
    },

    async getUserData(): Promise<any | null> {
        const data = await SecureStore.getItemAsync(STORAGE_KEYS.USER_DATA);
        return data ? JSON.parse(data) : null;
    },

    async removeUserData(): Promise<void> {
        await SecureStore.deleteItemAsync(STORAGE_KEYS.USER_DATA);
    },

    async clear(): Promise<void> {
        await this.removeToken();
        await this.removeUserData();
    },
};
