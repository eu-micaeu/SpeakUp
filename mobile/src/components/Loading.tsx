import React from 'react';
import { View, ActivityIndicator, StyleSheet, Text } from 'react-native';
import { Colors, FontSizes } from '../constants/theme';

interface LoadingProps {
    text?: string;
}

export const Loading: React.FC<LoadingProps> = ({ text = 'Carregando...' }) => {
    return (
        <View style={styles.container}>
            <ActivityIndicator size="large" color={Colors.primary} />
            {text && <Text style={styles.text}>{text}</Text>}
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: Colors.background,
    },
    text: {
        marginTop: 16,
        fontSize: FontSizes.md,
        color: Colors.textSecondary,
    },
});
