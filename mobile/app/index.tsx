import { Redirect } from 'expo-router';
import { useEffect, useState } from 'react';
import { View, ActivityIndicator, Text } from 'react-native';
import { storage } from '../src/utils/storage';

export default function Index() {
    const [loading, setLoading] = useState(true);
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        checkAuth();
    }, []);

    const checkAuth = async () => {
        try {
            console.log('🔍 [SpeakUp] Checking authentication...');
            const token = await storage.getToken();
            console.log('🔑 [SpeakUp] Token found:', !!token);
            setIsAuthenticated(!!token);
        } catch (error) {
            console.error('❌ [SpeakUp] Error checking auth:', error);
            setIsAuthenticated(false);
        } finally {
            console.log('✅ [SpeakUp] Auth check complete');
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: '#000000' }}>
                <ActivityIndicator size="large" color="#ffffff" />
                <Text style={{ color: '#ffffff', marginTop: 16, fontSize: 16 }}>
                    Carregando SpeakUp...
                </Text>
            </View>
        );
    }

    console.log('🚀 [SpeakUp] Redirecting to:', isAuthenticated ? 'home' : 'login');
    return <Redirect href={isAuthenticated ? '/(tabs)/home' : '/(auth)/login'} />;
}
