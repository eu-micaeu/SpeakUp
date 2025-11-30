import { apiClient } from './api';
import type { Word } from '../types';

export const wordService = {
    async getWordsByUserId(): Promise<Word[]> {
        const response = await apiClient.get<Word[]>('/api/word/user');
        return response.data || [];
    },

    async addWord(word: string, translation: string): Promise<Word> {
        const response = await apiClient.post<Word>('/api/word', {
            word,
            translation,
        });
        return response.data;
    },

    async deleteWord(wordId: string): Promise<void> {
        await apiClient.delete(`/api/word/${wordId}`);
    },

    async updateWord(wordId: string, word: string, translation: string): Promise<Word> {
        const response = await apiClient.put<Word>(`/api/word/${wordId}`, {
            word,
            translation,
        });
        return response.data;
    },
};
