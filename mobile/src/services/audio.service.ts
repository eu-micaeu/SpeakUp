import { apiClient } from './api';

export const audioService = {
    async transcribeAudio(audioBlob: Blob): Promise<string> {
        const formData = new FormData();
        formData.append('audio', audioBlob, 'audio.webm');

        const response = await apiClient.post<{ text: string }>('/api/ai/transcribe-audio', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });

        return response.data.text;
    },
};
