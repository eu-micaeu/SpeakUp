import React, { useState, useEffect, useRef } from 'react';
import {
    View,
    Text,
    StyleSheet,
    FlatList,
    TextInput,
    TouchableOpacity,
    KeyboardAvoidingView,
    Platform,
    Alert,
} from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { chatService } from '../../src/services/chat.service';
import { aiService } from '../../src/services/ai.service';
import { MessageBubble } from '../../src/components/MessageBubble';
import { AudioRecorder } from '../../src/components/AudioRecorder';
import { Loading } from '../../src/components/Loading';
import { Colors, Spacing, FontSizes, BorderRadius } from '../../src/constants/theme';
import type { Message } from '../../src/types';

export default function ChatScreen() {
    const params = useLocalSearchParams();
    const chatId = params.id as string | undefined;

    const [messages, setMessages] = useState<Message[]>([]);
    const [inputText, setInputText] = useState('');
    const [loading, setLoading] = useState(false);
    const [sending, setSending] = useState(false);
    const [currentChatId, setCurrentChatId] = useState<string | null>(chatId || null);
    const [chatTitle, setChatTitle] = useState<string>('');
    const flatListRef = useRef<FlatList>(null);

    useEffect(() => {
        if (chatId) {
            loadChat(chatId);
        }
    }, [chatId]);

    const loadChat = async (id: string) => {
        setLoading(true);
        try {
            const [chatData, messagesData] = await Promise.all([
                chatService.getChatById(id),
                chatService.getMessagesByChatId(id)
            ]);
            setChatTitle(chatData.topic || 'Conversa');
            setMessages(messagesData);
        } catch (error) {
            Alert.alert('Erro', 'Não foi possível carregar a conversa');
        } finally {
            setLoading(false);
        }
    };

    const handleSend = async () => {
        if (!inputText.trim() || sending) return;

        const userMessage = inputText.trim();
        setInputText('');
        setSending(true);

        try {
            let chatIdToUse = currentChatId;

            // Create new chat if needed
            if (!chatIdToUse) {
                const topicRes = await aiService.generateTopic(userMessage);
                const newChat = await chatService.createChat(topicRes.response);
                chatIdToUse = newChat.id;
                setCurrentChatId(chatIdToUse);
                setChatTitle(topicRes.response);
            }

            // Get correction
            const correction = await aiService.generateCorrection(userMessage);
            const fullUserMessage = `${userMessage}\n\nCorreção: ${correction.response}`;

            // Save user message
            const savedUserMsg = await chatService.addMessage(
                chatIdToUse,
                fullUserMessage,
                'user',
                'request'
            );

            setMessages((prev) => [...prev, savedUserMsg]);

            // Get AI response
            const dialogRes = await aiService.generateDialog(userMessage);

            // Get translation
            const translation = await aiService.generateTranslation(dialogRes.response);
            const fullAIResponse = `${dialogRes.response}\n\n[TRANSLATION]: ${translation.response}`;

            // Save AI response
            const savedAiMsg = await chatService.addMessage(
                chatIdToUse,
                fullAIResponse,
                'ai',
                'response'
            );

            setMessages((prev) => [...prev, savedAiMsg]);

            // Scroll to bottom
            setTimeout(() => {
                flatListRef.current?.scrollToEnd({ animated: true });
            }, 100);
        } catch (error: any) {
            Alert.alert('Erro', error.response?.data?.error || 'Erro ao enviar mensagem');
        } finally {
            setSending(false);
        }
    };

    const handleAudioTranscription = (text: string) => {
        if (text) {
            setInputText(text);
        }
    };

    if (loading) {
        return <Loading text="Carregando conversa..." />;
    }

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
            keyboardVerticalOffset={100}
        >
            <View style={styles.header}>
                <Text style={styles.headerTitle}>
                    {chatTitle || (currentChatId ? 'Conversa' : 'Nova Conversa')}
                </Text>
            </View>

            {messages.length === 0 ? (
                <View style={styles.emptyState}>
                    <Ionicons name="chatbubble-outline" size={64} color={Colors.textSecondary} />
                    <Text style={styles.emptyText}>Comece uma nova conversa</Text>
                    <Text style={styles.emptySubtext}>
                        Digite uma mensagem ou grave um áudio para praticar
                    </Text>
                </View>
            ) : (
                <FlatList
                    ref={flatListRef}
                    data={messages}
                    keyExtractor={(item) => item.id}
                    renderItem={({ item }) => <MessageBubble message={item} />}
                    contentContainerStyle={styles.messagesList}
                    onContentSizeChange={() => flatListRef.current?.scrollToEnd()}
                />
            )}

            <View style={styles.inputContainer}>
                <AudioRecorder onTranscription={handleAudioTranscription} disabled={sending} />

                <TextInput
                    style={styles.input}
                    value={inputText}
                    onChangeText={setInputText}
                    placeholder="Digite sua mensagem..."
                    placeholderTextColor={Colors.textSecondary}
                    multiline
                    maxLength={500}
                    editable={!sending}
                />

                <TouchableOpacity
                    style={[styles.sendButton, !inputText.trim() && styles.sendButtonDisabled]}
                    onPress={handleSend}
                    disabled={!inputText.trim() || sending}
                >
                    <Ionicons name="send" size={20} color={Colors.background} />
                </TouchableOpacity>
            </View>
        </KeyboardAvoidingView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: Colors.background,
    },
    header: {
        paddingHorizontal: Spacing.lg,
        paddingTop: Spacing.xxl,
        paddingBottom: Spacing.md,
        borderBottomWidth: 1,
        borderBottomColor: Colors.border,
        backgroundColor: Colors.backgroundSecondary,
    },
    headerTitle: {
        fontSize: FontSizes.xl,
        fontWeight: 'bold',
        color: Colors.text,
        textAlign: 'center',
    },
    emptyState: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: Spacing.xl,
    },
    emptyText: {
        fontSize: FontSizes.lg,
        fontWeight: '600',
        color: Colors.text,
        marginTop: Spacing.md,
    },
    emptySubtext: {
        fontSize: FontSizes.md,
        color: Colors.textSecondary,
        textAlign: 'center',
        marginTop: Spacing.xs,
    },
    messagesList: {
        paddingVertical: Spacing.md,
    },
    inputContainer: {
        flexDirection: 'row',
        alignItems: 'flex-end',
        padding: Spacing.md,
        borderTopWidth: 1,
        borderTopColor: Colors.border,
        backgroundColor: Colors.surface,
        gap: Spacing.sm,
    },
    input: {
        flex: 1,
        backgroundColor: Colors.backgroundSecondary,
        borderRadius: BorderRadius.lg,
        paddingHorizontal: Spacing.md,
        paddingVertical: Spacing.sm,
        fontSize: FontSizes.md,
        color: Colors.text,
        maxHeight: 100,
    },
    sendButton: {
        backgroundColor: Colors.primary,
        width: 40,
        height: 40,
        borderRadius: BorderRadius.full,
        justifyContent: 'center',
        alignItems: 'center',
    },
    sendButtonDisabled: {
        opacity: 0.5,
    },
});
