import React, { useState } from 'react';
import {
    View,
    Text,
    StyleSheet,
    KeyboardAvoidingView,
    Platform,
    ScrollView,
    Alert,
} from 'react-native';
import { useRouter } from 'expo-router';
import { Input } from '../../src/components/Input';
import { Button } from '../../src/components/Button';
import { authService } from '../../src/services/auth.service';
import { Colors, Spacing, FontSizes } from '../../src/constants/theme';

export default function RegisterScreen() {
    const router = useRouter();
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [language, setLanguage] = useState('');
    const [level, setLevel] = useState('');
    const [loading, setLoading] = useState(false);

    const handleRegister = async () => {
        if (!name || !email || !password || !language || !level) {
            Alert.alert('Erro', 'Por favor, preencha todos os campos');
            return;
        }

        setLoading(true);
        try {
            await authService.register({ name, email, password, language, level });
            Alert.alert('Sucesso', 'Conta criada com sucesso!', [
                { text: 'OK', onPress: () => router.replace('/(auth)/login') },
            ]);
        } catch (error: any) {
            Alert.alert('Erro', error.response?.data?.error || 'Erro ao criar conta');
        } finally {
            setLoading(false);
        }
    };

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        >
            <ScrollView
                contentContainerStyle={styles.scrollContent}
                keyboardShouldPersistTaps="handled"
            >
                <View style={styles.header}>
                    <Text style={styles.title}>Criar Conta</Text>
                    <Text style={styles.subtitle}>Comece sua jornada de aprendizado</Text>
                </View>

                <View style={styles.form}>
                    <Input
                        label="Nome"
                        value={name}
                        onChangeText={setName}
                        placeholder="Seu nome"
                        autoComplete="name"
                    />

                    <Input
                        label="Email"
                        value={email}
                        onChangeText={setEmail}
                        placeholder="seu@email.com"
                        keyboardType="email-address"
                        autoCapitalize="none"
                        autoComplete="email"
                    />

                    <Input
                        label="Senha"
                        value={password}
                        onChangeText={setPassword}
                        placeholder="••••••••"
                        secureTextEntry
                        autoComplete="password"
                    />

                    <Input
                        label="Idioma que está aprendendo"
                        value={language}
                        onChangeText={setLanguage}
                        placeholder="Ex: Inglês, Espanhol..."
                    />

                    <Input
                        label="Nível"
                        value={level}
                        onChangeText={setLevel}
                        placeholder="Ex: Iniciante, Intermediário..."
                    />

                    <Button
                        title="Criar conta"
                        onPress={handleRegister}
                        loading={loading}
                        style={styles.registerButton}
                    />

                    <Button
                        title="Já tenho conta"
                        onPress={() => router.back()}
                        variant="outline"
                    />
                </View>
            </ScrollView>
        </KeyboardAvoidingView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: Colors.background,
    },
    scrollContent: {
        flexGrow: 1,
        justifyContent: 'center',
        padding: Spacing.lg,
    },
    header: {
        alignItems: 'center',
        marginBottom: Spacing.xl,
    },
    title: {
        fontSize: FontSizes.xxxl,
        fontWeight: 'bold',
        color: Colors.primary,
        marginBottom: Spacing.xs,
    },
    subtitle: {
        fontSize: FontSizes.md,
        color: Colors.textSecondary,
    },
    form: {
        width: '100%',
    },
    registerButton: {
        marginBottom: Spacing.md,
    },
});
