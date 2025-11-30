import React, { useState, useRef } from 'react';
import { View, TouchableOpacity, StyleSheet, Text, Platform } from 'react-native';
import { Audio } from 'expo-av';
import { Ionicons } from '@expo/vector-icons';
import { Colors, BorderRadius, Spacing, FontSizes } from '../constants/theme';
import * as FileSystem from 'expo-file-system';

interface AudioRecorderProps {
    onTranscription: (text: string) => void;
    disabled?: boolean;
}

export const AudioRecorder: React.FC<AudioRecorderProps> = ({ onTranscription, disabled = false }) => {
    const [isRecording, setIsRecording] = useState(false);
    const [isTranscribing, setIsTranscribing] = useState(false);
    const [recordingTime, setRecordingTime] = useState(0);
    const recordingRef = useRef<Audio.Recording | null>(null);
    const timerRef = useRef<NodeJS.Timeout | null>(null);

    const startRecording = async () => {
        try {
            // Request permissions
            const permission = await Audio.requestPermissionsAsync();
            if (!permission.granted) {
                alert('Permissão de microfone é necessária para gravar áudio');
                return;
            }

            // Configure audio mode
            await Audio.setAudioModeAsync({
                allowsRecordingIOS: true,
                playsInSilentModeIOS: true,
            });

            // Start recording
            const { recording } = await Audio.Recording.createAsync(
                Audio.RecordingOptionsPresets.HIGH_QUALITY
            );

            recordingRef.current = recording;
            setIsRecording(true);
            setRecordingTime(0);

            // Start timer
            timerRef.current = setInterval(() => {
                setRecordingTime((prev) => prev + 1);
            }, 1000);
        } catch (error) {
            console.error('Failed to start recording:', error);
            alert('Erro ao iniciar gravação');
        }
    };

    const stopRecording = async () => {
        if (!recordingRef.current) return;

        try {
            setIsRecording(false);

            // Stop timer
            if (timerRef.current) {
                clearInterval(timerRef.current);
                timerRef.current = null;
            }

            // Stop recording
            await recordingRef.current.stopAndUnloadAsync();
            const uri = recordingRef.current.getURI();
            recordingRef.current = null;

            if (uri) {
                setIsTranscribing(true);

                // Read file as base64
                const audioBase64 = await FileSystem.readAsStringAsync(uri, {
                    encoding: FileSystem.EncodingType.Base64,
                });

                // Convert to blob for transcription
                const audioBlob = await fetch(`data:audio/m4a;base64,${audioBase64}`).then(r => r.blob());

                // TODO: Call transcription service
                // const text = await audioService.transcribeAudio(audioBlob);
                // onTranscription(text);

                // Temporary placeholder
                onTranscription('Transcrição em desenvolvimento...');
                setIsTranscribing(false);
            }
        } catch (error) {
            console.error('Failed to stop recording:', error);
            setIsTranscribing(false);
        }
    };

    const toggleRecording = () => {
        if (isRecording) {
            stopRecording();
        } else {
            startRecording();
        }
    };

    const formatTime = (seconds: number): string => {
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    };

    return (
        <View style={styles.container}>
            <TouchableOpacity
                style={[
                    styles.button,
                    isRecording && styles.recordingButton,
                    disabled && styles.disabledButton,
                ]}
                onPress={toggleRecording}
                disabled={disabled || isTranscribing}
            >
                {isTranscribing ? (
                    <Text style={styles.buttonText}>Transcrevendo...</Text>
                ) : (
                    <>
                        <Ionicons
                            name={isRecording ? 'stop-circle' : 'mic'}
                            size={24}
                            color={Colors.background}
                        />
                        {isRecording && (
                            <Text style={styles.timeText}>{formatTime(recordingTime)}</Text>
                        )}
                    </>
                )}
            </TouchableOpacity>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        alignItems: 'center',
    },
    button: {
        backgroundColor: Colors.primary,
        borderRadius: BorderRadius.full,
        width: 56,
        height: 56,
        alignItems: 'center',
        justifyContent: 'center',
        flexDirection: 'row',
        gap: Spacing.xs,
        paddingHorizontal: Spacing.md,
    },
    recordingButton: {
        backgroundColor: Colors.error,
        width: 'auto',
        minWidth: 120,
    },
    disabledButton: {
        opacity: 0.5,
    },
    buttonText: {
        color: Colors.background,
        fontSize: FontSizes.sm,
        fontWeight: '600',
    },
    timeText: {
        color: Colors.background,
        fontSize: FontSizes.md,
        fontWeight: '600',
    },
});
