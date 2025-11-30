import React, { useState } from 'react';
import {
    View,
    Text,
    StyleSheet,
    KeyboardAvoidingView,
    Platform,
    ScrollView,
    Alert,
    Image,
} from 'react-native';
import { useRouter } from 'expo-router';
import { Input } from '../../src/components/Input';
import { Button } from '../../src/components/Button';
import { authService } from '../../src/services/auth.service';
import { Colors, Spacing, FontSizes } from '../../src/constants/theme';

export default function LoginScreen() {
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);

    console.log('📱 [Login] Screen rendered');

    const handleLogin = async () => {
        console.log('🔐 [Login] Attempting login...');
        if (!email || !password) {
            Alert.alert('Erro', 'Por favor, preencha todos os campos');
            return;
        }

        setLoading(true);
        try {
            console.log('📡 [Login] Calling auth service...');
            await authService.login({ email, password });
            console.log('✅ [Login] Login successful, redirecting...');
            router.replace('/(tabs)/home');
        } catch (error: any) {
            console.error('❌ [Login] Error:', error);
            Alert.alert('Erro', error.response?.data?.error || 'Erro ao fazer login');
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
                    <Image
                        source={require('../../assets/logo.png')}
                        style={styles.logo}
                        resizeMode="contain"
                    />
                    <Text style={styles.title}>SpeakUp</Text>
                    <Text style={styles.subtitle}>Aprenda idiomas com IA</Text>
                </View>

                <View style={styles.form}>
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

                    <Button
                        title="Entrar"
                        onPress={handleLogin}
                        loading={loading}
                        style={styles.loginButton}
                    />

                    <Button
                        title="Criar conta"
                        onPress={() => router.push('/(auth)/register')}
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
        marginBottom: Spacing.xxl,
    },
    logo: {
        width: 75,
        height: 75,
        marginBottom: Spacing.md,
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
    loginButton: {
        marginBottom: Spacing.md,
    },
});
