import { apiClient } from './api';
import type { AIResponse } from '../types';

export const aiService = {
    async generateDialog(message: string): Promise<AIResponse> {
        const response = await apiClient.post<AIResponse>('/api/ai/generate-response-dialog', { message });
        return response.data;
    },

    async generateCorrection(message: string): Promise<AIResponse> {
        const response = await apiClient.post<AIResponse>('/api/ai/generate-response-correction', { message });
        return response.data;
    },

    async generateTranslation(message: string): Promise<AIResponse> {
        const response = await apiClient.post<AIResponse>('/api/ai/generate-response-translation', { message });
        return response.data;
    },

    async generateTopic(message: string): Promise<AIResponse> {
        const response = await apiClient.post<AIResponse>('/api/ai/generate-response-topic', { message });
        return response.data;
    },

    async generateRandomWord(): Promise<AIResponse> {
        const response = await apiClient.post<AIResponse>('/api/ai/generate-random-word', {});
        return response.data;
    },
};
