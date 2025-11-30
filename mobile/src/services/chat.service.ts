import { apiClient } from './api';
import type { Chat, Message, AIResponse } from '../types';

export const chatService = {
    async createChat(topic: string): Promise<Chat> {
        const response = await apiClient.post<Chat>('/api/chat', { topic });
        return response.data;
    },

    async getChatsByUserId(): Promise<Chat[]> {
        const response = await apiClient.get<{ chats: Chat[] }>('/api/chat/user');
        return response.data.chats || [];
    },

    async getChatById(chatId: string): Promise<Chat> {
        const response = await apiClient.get<Chat>(`/api/chat/${chatId}`);
        return response.data;
    },

    async deleteChat(chatId: string): Promise<void> {
        await apiClient.delete(`/api/chat/${chatId}`);
    },

    async getMessagesByChatId(chatId: string): Promise<Message[]> {
        console.log('📨 [Chat] Fetching messages for chat:', chatId);
        const response = await apiClient.get<Message[]>(`/api/message/chat/${chatId}`);
        console.log('📨 [Chat] Messages received:', response.data?.length || 0);
        return response.data || [];
    },

    async addMessage(
        chatId: string,
        content: string,
        sender: 'user' | 'ai',
        type: 'request' | 'response'
    ): Promise<Message> {
        const response = await apiClient.post<Message>('/api/message', {
            chat_id: chatId,
            content,
            sender,
            type,
        });
        return response.data;
    },
};
