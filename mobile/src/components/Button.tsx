import React from 'react';
import { TouchableOpacity, StyleSheet, ActivityIndicator, Text, ViewStyle, TextStyle } from 'react-native';
import { Colors, BorderRadius, Spacing, FontSizes } from '../constants/theme';

interface ButtonProps {
    title: string;
    onPress: () => void;
    loading?: boolean;
    disabled?: boolean;
    variant?: 'primary' | 'secondary' | 'outline';
    style?: ViewStyle;
    textStyle?: TextStyle;
}

export const Button: React.FC<ButtonProps> = ({
    title,
    onPress,
    loading = false,
    disabled = false,
    variant = 'primary',
    style,
    textStyle,
}) => {
    const isDisabled = disabled || loading;

    return (
        <TouchableOpacity
            style={[
                styles.button,
                variant === 'primary' && styles.primary,
                variant === 'secondary' && styles.secondary,
                variant === 'outline' && styles.outline,
                isDisabled && styles.disabled,
                style,
            ]}
            onPress={onPress}
            disabled={isDisabled}
            activeOpacity={0.7}
        >
            {loading ? (
                <ActivityIndicator color={variant === 'outline' ? Colors.primary : Colors.background} />
            ) : (
                <Text
                    style={[
                        styles.text,
                        variant === 'primary' && styles.primaryText,
                        variant === 'secondary' && styles.secondaryText,
                        variant === 'outline' && styles.outlineText,
                        textStyle,
                    ]}
                >
                    {title}
                </Text>
            )}
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    button: {
        paddingVertical: Spacing.md,
        paddingHorizontal: Spacing.lg,
        borderRadius: BorderRadius.md,
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: 48,
    },
    primary: {
        backgroundColor: Colors.primary,
    },
    secondary: {
        backgroundColor: Colors.secondary,
    },
    outline: {
        backgroundColor: 'transparent',
        borderWidth: 2,
        borderColor: Colors.primary,
    },
    disabled: {
        opacity: 0.5,
    },
    text: {
        fontSize: FontSizes.md,
        fontWeight: '600',
    },
    primaryText: {
        color: Colors.background,
    },
    secondaryText: {
        color: Colors.background,
    },
    outlineText: {
        color: Colors.primary,
    },
});
