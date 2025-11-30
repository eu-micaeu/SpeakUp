import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, Alert } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { authService } from '../../src/services/auth.service';
import { storage } from '../../src/utils/storage';
import { Button } from '../../src/components/Button';
import { Input } from '../../src/components/Input';
import { Loading } from '../../src/components/Loading';
import { Colors, Spacing, FontSizes, BorderRadius } from '../../src/constants/theme';
import type { User } from '../../src/types';

export default function ProfileScreen() {
    const router = useRouter();
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [editing, setEditing] = useState(false);
    const [name, setName] = useState('');
    const [language, setLanguage] = useState('');
    const [level, setLevel] = useState('');
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        loadUserData();
    }, []);

    const loadUserData = async () => {
        try {
            console.log('📱 [Profile] Loading user data...');

            // Tentar carregar do storage primeiro
            let userData = await storage.getUserData();

            // Se não tiver no storage, buscar do backend
            if (!userData) {
                console.log('📡 [Profile] No cached data, fetching from backend...');
                userData = await authService.getUserProfile();
            }

            if (userData) {
                console.log('✅ [Profile] User data loaded:', userData.name);
                setUser(userData);
                setName(userData.name || '');
                setLanguage(userData.language || '');
                setLevel(userData.level || '');
            } else {
                console.error('❌ [Profile] No user data available');
            }
        } catch (error) {
            console.error('❌ [Profile] Error loading user data:', error);
            Alert.alert('Erro', 'Não foi possível carregar o perfil. Tente fazer login novamente.');
        } finally {
            setLoading(false);
        }
    };

    const handleSaveProfile = async () => {
        if (!name.trim() || !language.trim() || !level.trim()) {
            Alert.alert('Erro', 'Preencha todos os campos');
            return;
        }

        setSaving(true);
        try {
            const updatedUser = await authService.updateUserProfile({
                name: name.trim(),
                language: language.trim(),
                level: level.trim(),
            });
            setUser(updatedUser);
            setEditing(false);
            Alert.alert('Sucesso', 'Perfil atualizado com sucesso!');
        } catch (error) {
            Alert.alert('Erro', 'Não foi possível atualizar o perfil');
        } finally {
            setSaving(false);
        }
    };

    const handleLogout = () => {
        Alert.alert('Sair', 'Tem certeza que deseja sair?', [
            { text: 'Cancelar', style: 'cancel' },
            {
                text: 'Sair',
                style: 'destructive',
                onPress: async () => {
                    await authService.logout();
                    router.replace('/(auth)/login');
                },
            },
        ]);
    };

    const handleDeleteAccount = () => {
        Alert.alert(
            'Deletar Conta',
            'Esta ação é irreversível. Todos os seus dados serão permanentemente apagados. Tem certeza?',
            [
                { text: 'Cancelar', style: 'cancel' },
                {
                    text: 'Deletar',
                    style: 'destructive',
                    onPress: async () => {
                        try {
                            if (user?.id) {
                                await authService.deleteAccount(user.id);
                                await authService.logout();
                                router.replace('/(auth)/login');
                                Alert.alert('Conta Deletada', 'Sua conta foi deletada com sucesso.');
                            }
                        } catch (error) {
                            Alert.alert('Erro', 'Não foi possível deletar a conta');
                        }
                    },
                },
            ]
        );
    };

    if (loading) {
        return <Loading />;
    }

    if (!user) {
        return (
            <View style={styles.container}>
                <Text style={styles.errorText}>Erro ao carregar perfil</Text>
            </View>
        );
    }

    return (
        <ScrollView style={styles.container}>
            <View style={styles.header}>
                <View style={styles.avatarContainer}>
                    <Ionicons name="person-circle" size={80} color={Colors.primary} />
                </View>
                <Text style={styles.headerTitle}>Meu Perfil</Text>
            </View>

            <View style={styles.content}>
                {editing ? (
                    <>
                        <Input label="Nome" value={name} onChangeText={setName} />
                        <Input label="Email" value={user.email} editable={false} />
                        <Input label="Idioma" value={language} onChangeText={setLanguage} />
                        <Input label="Nível" value={level} onChangeText={setLevel} />

                        <View style={styles.buttonGroup}>
                            <Button
                                title="Cancelar"
                                onPress={() => {
                                    setEditing(false);
                                    setName(user.name);
                                    setLanguage(user.language);
                                    setLevel(user.level);
                                }}
                                variant="outline"
                                style={styles.button}
                            />
                            <Button
                                title="Salvar"
                                onPress={handleSaveProfile}
                                loading={saving}
                                style={styles.button}
                            />
                        </View>
                    </>
                ) : (
                    <>
                        <View style={styles.infoCard}>
                            <View style={styles.infoRow}>
                                <Ionicons name="person" size={20} color={Colors.textSecondary} />
                                <View style={styles.infoContent}>
                                    <Text style={styles.infoLabel}>Nome</Text>
                                    <Text style={styles.infoValue}>{user.name}</Text>
                                </View>
                            </View>

                            <View style={styles.infoRow}>
                                <Ionicons name="mail" size={20} color={Colors.textSecondary} />
                                <View style={styles.infoContent}>
                                    <Text style={styles.infoLabel}>Email</Text>
                                    <Text style={styles.infoValue}>{user.email}</Text>
                                </View>
                            </View>

                            <View style={styles.infoRow}>
                                <Ionicons name="language" size={20} color={Colors.textSecondary} />
                                <View style={styles.infoContent}>
                                    <Text style={styles.infoLabel}>Idioma</Text>
                                    <Text style={styles.infoValue}>{user.language}</Text>
                                </View>
                            </View>

                            <View style={styles.infoRow}>
                                <Ionicons name="stats-chart" size={20} color={Colors.textSecondary} />
                                <View style={styles.infoContent}>
                                    <Text style={styles.infoLabel}>Nível</Text>
                                    <Text style={styles.infoValue}>{user.level}</Text>
                                </View>
                            </View>
                        </View>

                        <Button title="Editar Perfil" onPress={() => setEditing(true)} />
                        <Button
                            title="Sair"
                            onPress={handleLogout}
                            variant="outline"
                            style={styles.logoutButton}
                        />
                        <Button
                            title="Deletar Conta"
                            onPress={handleDeleteAccount}
                            variant="outline"
                            style={styles.deleteButton}
                        />
                    </>
                )}
            </View>
        </ScrollView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: Colors.background,
    },
    header: {
        alignItems: 'center',
        paddingTop: Spacing.xxl,
        paddingBottom: Spacing.lg,
        borderBottomWidth: 1,
        borderBottomColor: Colors.border,
    },
    avatarContainer: {
        marginBottom: Spacing.md,
    },
    headerTitle: {
        fontSize: FontSizes.xxl,
        fontWeight: 'bold',
        color: Colors.text,
    },
    content: {
        padding: Spacing.lg,
    },
    infoCard: {
        backgroundColor: Colors.surface,
        borderRadius: BorderRadius.md,
        padding: Spacing.md,
        marginBottom: Spacing.lg,
        borderWidth: 1,
        borderColor: Colors.border,
    },
    infoRow: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingVertical: Spacing.md,
        borderBottomWidth: 1,
        borderBottomColor: Colors.border,
    },
    infoContent: {
        flex: 1,
        marginLeft: Spacing.md,
    },
    infoLabel: {
        fontSize: FontSizes.sm,
        color: Colors.textSecondary,
        marginBottom: Spacing.xs,
    },
    infoValue: {
        fontSize: FontSizes.md,
        fontWeight: '600',
        color: Colors.text,
    },
    buttonGroup: {
        flexDirection: 'row',
        gap: Spacing.sm,
        marginTop: Spacing.md,
    },
    button: {
        flex: 1,
    },
    logoutButton: {
        marginTop: Spacing.md,
    },
    deleteButton: {
        marginTop: Spacing.sm,
        borderColor: Colors.error,
    },
    errorText: {
        fontSize: FontSizes.md,
        color: Colors.error,
        textAlign: 'center',
        marginTop: Spacing.xxl,
    },
});
