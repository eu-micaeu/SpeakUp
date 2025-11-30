import React from 'react';
import { TextInput, StyleSheet, View, Text, TextInputProps } from 'react-native';
import { Colors, BorderRadius, Spacing, FontSizes } from '../constants/theme';

interface InputProps extends TextInputProps {
    label?: string;
    error?: string;
    containerStyle?: any;
}

export const Input: React.FC<InputProps> = ({
    label,
    error,
    containerStyle,
    style,
    ...props
}) => {
    return (
        <View style={[styles.container, containerStyle]}>
            {label && <Text style={styles.label}>{label}</Text>}
            <TextInput
                style={[styles.input, error && styles.inputError, style]}
                placeholderTextColor={Colors.textSecondary}
                {...props}
            />
            {error && <Text style={styles.error}>{error}</Text>}
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        marginBottom: Spacing.md,
    },
    label: {
        fontSize: FontSizes.sm,
        fontWeight: '600',
        color: Colors.text,
        marginBottom: Spacing.xs,
    },
    input: {
        backgroundColor: Colors.background,
        borderWidth: 1,
        borderColor: Colors.border,
        borderRadius: BorderRadius.md,
        paddingHorizontal: Spacing.md,
        paddingVertical: Spacing.md,
        fontSize: FontSizes.md,
        color: Colors.text,
        minHeight: 48,
    },
    inputError: {
        borderColor: Colors.error,
    },
    error: {
        fontSize: FontSizes.xs,
        color: Colors.error,
        marginTop: Spacing.xs,
    },
});
