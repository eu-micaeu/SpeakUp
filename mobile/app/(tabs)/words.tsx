import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { wordService } from '../../src/services/word.service';
import { Loading } from '../../src/components/Loading';
import { Button } from '../../src/components/Button';
import { Input } from '../../src/components/Input';
import { Colors, Spacing, FontSizes, BorderRadius } from '../../src/constants/theme';
import type { Word } from '../../src/types';

export default function WordsScreen() {
    const [words, setWords] = useState<Word[]>([]);
    const [loading, setLoading] = useState(true);
    const [showAddModal, setShowAddModal] = useState(false);
    const [newWord, setNewWord] = useState('');
    const [newTranslation, setNewTranslation] = useState('');
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        loadWords();
    }, []);

    const loadWords = async () => {
        try {
            const data = await wordService.getWordsByUserId();
            setWords(data);
        } catch (error) {
            Alert.alert('Erro', 'Não foi possível carregar as palavras');
        } finally {
            setLoading(false);
        }
    };

    const handleAddWord = async () => {
        if (!newWord.trim() || !newTranslation.trim()) {
            Alert.alert('Erro', 'Preencha todos os campos');
            return;
        }

        setSaving(true);
        try {
            const word = await wordService.addWord(newWord.trim(), newTranslation.trim());
            setWords([word, ...words]);
            setNewWord('');
            setNewTranslation('');
            setShowAddModal(false);
        } catch (error) {
            Alert.alert('Erro', 'Não foi possível adicionar a palavra');
        } finally {
            setSaving(false);
        }
    };

    const handleDeleteWord = async (wordId: string) => {
        Alert.alert('Excluir palavra', 'Tem certeza que deseja excluir esta palavra?', [
            { text: 'Cancelar', style: 'cancel' },
            {
                text: 'Excluir',
                style: 'destructive',
                onPress: async () => {
                    try {
                        await wordService.deleteWord(wordId);
                        setWords(words.filter((w) => w.id !== wordId));
                    } catch (error) {
                        Alert.alert('Erro', 'Não foi possível excluir a palavra');
                    }
                },
            },
        ]);
    };

    if (loading) {
        return <Loading />;
    }

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Meu Vocabulário</Text>
                {!showAddModal && (
                    <Button
                        title="Adicionar Palavra"
                        onPress={() => setShowAddModal(true)}
                        style={styles.addButton}
                    />
                )}
            </View>

            {showAddModal && (
                <View style={styles.addModal}>
                    <Input
                        label="Palavra"
                        value={newWord}
                        onChangeText={setNewWord}
                        placeholder="Digite a palavra"
                    />
                    <Input
                        label="Tradução"
                        value={newTranslation}
                        onChangeText={setNewTranslation}
                        placeholder="Digite a tradução"
                    />
                    <View style={styles.modalButtons}>
                        <Button
                            title="Cancelar"
                            onPress={() => setShowAddModal(false)}
                            variant="outline"
                            style={styles.modalButton}
                        />
                        <Button
                            title="Salvar"
                            onPress={handleAddWord}
                            loading={saving}
                            style={styles.modalButton}
                        />
                    </View>
                </View>
            )}

            {words.length === 0 ? (
                <View style={styles.emptyState}>
                    <Ionicons name="book-outline" size={64} color={Colors.textSecondary} />
                    <Text style={styles.emptyText}>Nenhuma palavra ainda</Text>
                    <Text style={styles.emptySubtext}>Adicione palavras para criar seu vocabulário</Text>
                </View>
            ) : (
                <FlatList
                    data={words}
                    keyExtractor={(item) => item.id}
                    renderItem={({ item }) => (
                        <View style={styles.wordItem}>
                            <View style={styles.wordContent}>
                                <Text style={styles.word}>{item.word}</Text>
                                <Text style={styles.translation}>{item.translation}</Text>
                            </View>
                            <TouchableOpacity
                                style={styles.deleteButton}
                                onPress={() => handleDeleteWord(item.id)}
                            >
                                <Ionicons name="trash-outline" size={20} color={Colors.error} />
                            </TouchableOpacity>
                        </View>
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
    addButton: {
        marginTop: Spacing.sm,
    },
    addModal: {
        backgroundColor: Colors.backgroundSecondary,
        padding: Spacing.lg,
        marginHorizontal: Spacing.md,
        borderRadius: BorderRadius.md,
        marginBottom: Spacing.md,
    },
    modalButtons: {
        flexDirection: 'row',
        gap: Spacing.sm,
        marginTop: Spacing.sm,
    },
    modalButton: {
        flex: 1,
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
    wordItem: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: Colors.surface,
        padding: Spacing.md,
        marginBottom: Spacing.sm,
        borderRadius: BorderRadius.md,
        borderWidth: 1,
        borderColor: Colors.border,
    },
    wordContent: {
        flex: 1,
    },
    word: {
        fontSize: FontSizes.lg,
        fontWeight: '600',
        color: Colors.text,
        marginBottom: Spacing.xs,
    },
    translation: {
        fontSize: FontSizes.md,
        color: Colors.textSecondary,
    },
    deleteButton: {
        padding: Spacing.sm,
    },
});
