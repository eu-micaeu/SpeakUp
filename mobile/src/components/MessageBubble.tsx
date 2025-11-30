import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Colors, BorderRadius, Spacing, FontSizes } from '../constants/theme';
import type { Message } from '../types';

interface MessageBubbleProps {
    message: Message;
}

export const MessageBubble: React.FC<MessageBubbleProps> = ({ message }) => {
    const isUser = message.sender === 'user';
    const isRequest = message.type === 'request';
    const isResponse = message.type === 'response';

    // Extrair texto principal e extras
    const parts = message.content.split('\n\n');
    const mainText = parts[0];

    // Para mensagens do usuário, extrair correção
    let correction = '';
    if (isRequest && parts[1]) {
        correction = parts[1].replace('Correção: ', '');
    }

    // Para respostas da IA, extrair tradução
    let translation = '';
    if (isResponse && message.content.includes('[TRANSLATION]: ')) {
        const splitTranslation = message.content.split('[TRANSLATION]: ');
        translation = splitTranslation[1] || '';
    }

    return (
        <View style={[styles.container, isUser ? styles.userContainer : styles.aiContainer]}>
            <View style={[styles.bubble, isUser ? styles.userBubble : styles.aiBubble]}>
                <Text style={[styles.text, isUser ? styles.userText : styles.aiText]}>
                    {mainText}
                </Text>

                {correction && (
                    <View style={styles.extras}>
                        <Text style={styles.extrasLabel}>Correção:</Text>
                        <Text style={styles.extrasText}>{correction}</Text>
                    </View>
                )}

                {translation && (
                    <View style={styles.extras}>
                        <Text style={styles.extrasLabel}>Tradução:</Text>
                        <Text style={styles.extrasText}>{translation}</Text>
                    </View>
                )}
            </View>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        marginVertical: Spacing.xs,
        paddingHorizontal: Spacing.md,
    },
    userContainer: {
        alignItems: 'flex-end',
    },
    aiContainer: {
        alignItems: 'flex-start',
    },
    bubble: {
        maxWidth: '80%',
        paddingHorizontal: Spacing.md,
        paddingVertical: Spacing.sm,
        borderRadius: BorderRadius.lg,
    },
    userBubble: {
        backgroundColor: Colors.userMessage,
        borderBottomRightRadius: Spacing.xs,
    },
    aiBubble: {
        backgroundColor: Colors.aiMessage,
        borderBottomLeftRadius: Spacing.xs,
    },
    text: {
        fontSize: FontSizes.md,
        lineHeight: 22,
    },
    userText: {
        color: Colors.userMessageText,
    },
    aiText: {
        color: Colors.aiMessageText,
    },
    extras: {
        marginTop: Spacing.sm,
        paddingTop: Spacing.sm,
        borderTopWidth: 1,
        borderTopColor: Colors.border,
    },
    extrasLabel: {
        fontSize: FontSizes.xs,
        fontWeight: 'bold',
        color: Colors.textSecondary,
        marginBottom: 4,
    },
    extrasText: {
        fontSize: FontSizes.xs,
        color: Colors.textSecondary,
        lineHeight: 18,
    },
});
