import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, Alert } from 'react-native';
import { useRouter, useFocusEffect } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { chatService } from '../../src/services/chat.service';
import { Loading } from '../../src/components/Loading';
import { Button } from '../../src/components/Button';
import { Colors, Spacing, FontSizes, BorderRadius } from '../../src/constants/theme';
import type { Chat } from '../../src/types';

export default function HomeScreen() {
    const router = useRouter();
    const [chats, setChats] = useState<Chat[]>([]);
    const [loading, setLoading] = useState(true);

    // Recarregar chats sempre que a tela ganhar foco
    useFocusEffect(
        useCallback(() => {
            loadChats();
        }, [])
    );

    const loadChats = async () => {
        try {
            const data = await chatService.getChatsByUserId();
            setChats(data);
        } catch (error) {
            Alert.alert('Erro', 'Não foi possível carregar as conversas');
        } finally {
            setLoading(false);
        }
    };

    const handleNewChat = () => {
        router.push('/(tabs)/chat');
    };

    const handleChatPress = (chatId: string) => {
        router.push(`/(tabs)/chat?id=${chatId}`);
    };

    if (loading) {
        return <Loading />;
    }

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Minhas Conversas</Text>
                <Button
                    title="Nova Conversa"
                    onPress={handleNewChat}
                    style={styles.newButton}
                />
            </View>

            {chats.length === 0 ? (
                <View style={styles.emptyState}>
                    <Ionicons name="chatbubbles-outline" size={64} color={Colors.textSecondary} />
                    <Text style={styles.emptyText}>Nenhuma conversa ainda</Text>
                    <Text style={styles.emptySubtext}>
                        Comece uma nova conversa para praticar seu idioma
                    </Text>
                </View>
            ) : (
                <FlatList
                    data={chats}
                    keyExtractor={(item) => item.id}
                    renderItem={({ item }) => (
                        <TouchableOpacity
                            style={styles.chatItem}
                            onPress={() => handleChatPress(item.id)}
                        >
                            <View style={styles.chatIcon}>
                                <Ionicons name="chatbubble" size={24} color={Colors.primary} />
                            </View>
                            <View style={styles.chatContent}>
                                <Text style={styles.chatTopic}>{item.topic}</Text>
                                {item.start_time && (
                                    <Text style={styles.chatDate}>
                                        {new Date(item.start_time).toLocaleDateString('pt-BR')}
                                    </Text>
                                )}
                            </View>
                            <Ionicons name="chevron-forward" size={20} color={Colors.textSecondary} />
                        </TouchableOpacity>
                    )}
                    contentContainerStyle={styles.listContent}
                />
            )}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: Colors.background,
    },
    header: {
        padding: Spacing.lg,
        paddingTop: Spacing.xxl,
    },
    title: {
        fontSize: FontSizes.xxl,
        fontWeight: 'bold',
        color: Colors.text,
        marginBottom: Spacing.md,
    },
    newButton: {
        marginTop: Spacing.sm,
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
    listContent: {
        padding: Spacing.md,
    },
    chatItem: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: Colors.surface,
        padding: Spacing.md,
        marginBottom: Spacing.sm,
        borderRadius: BorderRadius.md,
        borderWidth: 1,
        borderColor: Colors.border,
    },
    chatIcon: {
        width: 48,
        height: 48,
        borderRadius: BorderRadius.full,
        backgroundColor: Colors.backgroundSecondary,
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: Spacing.md,
    },
    chatContent: {
        flex: 1,
    },
    chatTopic: {
        fontSize: FontSizes.md,
        fontWeight: '600',
        color: Colors.text,
        marginBottom: Spacing.xs,
    },
    chatDate: {
        fontSize: FontSizes.sm,
        color: Colors.textSecondary,
    },
});
