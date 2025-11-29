<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { transcribeAudio } from "../utils/api";

    export let isRecording = false;
    export let disabled = false;

    const dispatch = createEventDispatcher();

    let mediaRecorder: MediaRecorder | null = null;
    let audioChunks: Blob[] = [];
    let recordingTime = 0;
    let recordingInterval: number | null = null;
    let currentStream: MediaStream | null = null;
    let isTranscribing = false;

    async function startRecording() {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                audio: true,
            });
            currentStream = stream;
            mediaRecorder = new MediaRecorder(stream);
            audioChunks = [];
            recordingTime = 0;

            mediaRecorder.ondataavailable = (event) => {
                if (event.data.size > 0) {
                    audioChunks.push(event.data);
                }
            };

            mediaRecorder.onstop = async () => {
                const audioBlob = new Blob(audioChunks, { type: "audio/webm" });
                console.log(
                    "📹 Gravação parou. Tamanho do áudio:",
                    audioBlob.size,
                    "bytes",
                );

                // Parar o stream
                if (currentStream) {
                    currentStream.getTracks().forEach((track) => track.stop());
                }

                if (recordingInterval) {
                    clearInterval(recordingInterval);
                    recordingInterval = null;
                }

                // Transcrever áudio usando o backend
                isTranscribing = true;
                try {
                    console.log("🔄 Enviando áudio para transcrição...");
                    const text = await transcribeAudio(audioBlob);
                    console.log("🎤 Texto transcrito:", text || "(vazio)");
                    dispatch("audioRecorded", { text });
                } catch (error) {
                    console.error("❌ Erro ao transcrever áudio:", error);
                    dispatch("audioRecorded", { text: "" });
                } finally {
                    isTranscribing = false;
                }
            };

            mediaRecorder.start();
            isRecording = true;

            // Atualizar tempo de gravação
            recordingInterval = window.setInterval(() => {
                recordingTime++;
            }, 1000);
        } catch (error) {
            console.error("Erro ao acessar microfone:", error);
            alert(
                "Não foi possível acessar o microfone. Verifique as permissões.",
            );
        }
    }

    function stopRecording() {
        if (mediaRecorder && mediaRecorder.state !== "inactive") {
            mediaRecorder.stop();
            isRecording = false;
        }
    }

    function toggleRecording() {
        if (isRecording || isTranscribing) {
            stopRecording();
        } else {
            startRecording();
        }
    }

    function formatTime(seconds: number): string {
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${mins}:${secs.toString().padStart(2, "0")}`;
    }
</script>

<button
    class="audioButton"
    class:recording={isRecording}
    class:transcribing={isTranscribing}
    on:click={toggleRecording}
    disabled={disabled || isTranscribing}
    type="button"
    aria-label={isRecording
        ? "Parar gravação"
        : isTranscribing
          ? "Transcrevendo..."
          : "Iniciar gravação"}
>
    {#if isTranscribing}
        <div class="recordingIndicator">
            <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="currentColor"
                class="spinner"
            >
                <circle
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                    fill="none"
                    stroke-dasharray="31.416"
                    stroke-dashoffset="31.416"
                >
                    <animate
                        attributeName="stroke-dasharray"
                        dur="2s"
                        values="0 31.416;15.708 15.708;0 31.416"
                        repeatCount="indefinite"
                    />
                    <animate
                        attributeName="stroke-dashoffset"
                        dur="2s"
                        values="0;-15.708;-31.416"
                        repeatCount="indefinite"
                    />
                </circle>
            </svg>
            <span class="recordingTime">...</span>
        </div>
    {:else if isRecording}
        <div class="recordingIndicator">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                <rect x="6" y="6" width="12" height="12" rx="2" />
            </svg>
            <span class="recordingTime">{formatTime(recordingTime)}</span>
        </div>
    {:else}
        <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path
                d="M12 14c1.66 0 3-1.34 3-3V5c0-1.66-1.34-3-3-3S9 3.34 9 5v6c0 1.66 1.34 3 3 3z"
            />
            <path
                d="M17 11c0 2.76-2.24 5-5 5s-5-2.24-5-5H5c0 3.53 2.61 6.43 6 6.92V21h2v-3.08c3.39-.49 6-3.39 6-6.92h-2z"
            />
        </svg>
    {/if}
</button>

<style>
    .audioButton {
        background: none;
        border: none;
        color: #bbb;
        cursor: pointer;
        padding: 0.5rem;
        border-radius: 0.25rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        transition: all 0.3s ease;
    }

    .audioButton:hover:not(:disabled) {
        color: #fff;
        background-color: #333;
    }

    .audioButton:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .audioButton.recording {
        color: #ff4444;
        animation: pulse 1.5s ease-in-out infinite;
    }

    .audioButton.transcribing {
        color: #4caf50;
    }

    @keyframes pulse {
        0%,
        100% {
            opacity: 1;
        }
        50% {
            opacity: 0.6;
        }
    }

    .recordingIndicator {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .recordingTime {
        font-size: 0.75rem;
        font-weight: 600;
        font-family: "Courier New", monospace;
    }
</style>
