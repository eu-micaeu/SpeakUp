<script lang="ts">
    import { onMount, tick } from "svelte";
    import { browser } from "$app/environment";
    import { goto } from "$app/navigation";
    import Sidebar from "../../components/Sidebar.svelte";
    import AudioRecorder from "../../components/AudioRecorder.svelte";
    import {
        getChatsByUserId,
        createChat,
        getMessagesByChatId,
        addMessageToChat,
        generateAIResponseDialog,
        generateAIResponseCorrection,
        generateAIResponseTranslation,
        generateAIResponseTopic,
    } from "../../utils/api";
    import Cookies from "js-cookie";

    function handleLogout() {
        Cookies.remove("authToken");
        goto("/");
    }

    interface Chat {
        id: string;
        title?: string;
        topic?: string;
        createdAt?: string;
    }

    interface Message {
        id: string;
        text: string;
        sender: "user" | "ai";
        type: "request" | "response";
    }

    interface ApiResponse {
        response: string;
        explanation?: string;
    }

    let chats: Chat[] = [];
    let messages: Message[] = [];
    let inputMessage = "";
    let currentChatId: string | null = null;
    let isSending = false;
    let sidebarOpen = true;
    let selectedChat: Chat | null = null;
    let showScrollDown = false;
    let messagesEndRef: HTMLDivElement;
    let chatBodyRef: HTMLDivElement;
    let isRecording = false;
    let showTranscriptionNotice = false;

    onMount(() => {
        (async () => {
            try {
                const res: any = await getChatsByUserId();
                chats = Array.isArray(res) ? res : res?.chats || [];
                selectedChat = null;
            } catch (error) {
                console.error("Error loading chats:", error);
                chats = [];
                selectedChat = null;
            }

            if (browser && window.innerWidth <= 768) {
                sidebarOpen = false;
            }
        })();

        if (chatBodyRef) {
            chatBodyRef.addEventListener("scroll", handleScroll);
        }

        return () => {
            if (chatBodyRef) {
                chatBodyRef.removeEventListener("scroll", handleScroll);
            }
        };
    });

    async function scrollToBottom() {
        await tick();
        if (!showScrollDown && messagesEndRef) {
            messagesEndRef.scrollIntoView({ behavior: "auto" });
        }
    }

    $: if (messages.length > 0) {
        scrollToBottom();
    }

    function handleScroll() {
        if (!chatBodyRef) return;
        const nearBottom =
            chatBodyRef.scrollHeight -
                chatBodyRef.scrollTop -
                chatBodyRef.clientHeight <
            100;
        showScrollDown = !nearBottom;
    }

    async function handleSend() {
        if (!inputMessage.trim() || isSending) return;
        isSending = true;
        const userInput = inputMessage.trim();
        inputMessage = "";

        try {
            let chatId = currentChatId;

            if (!chatId) {
                const topicRes: ApiResponse | undefined =
                    await generateAIResponseTopic(userInput);
                if (!topicRes?.response) {
                    throw new Error("Failed to generate topic");
                }
                const newChatRaw = await createChat(topicRes.response);
                const newChat: Chat = {
                    id: newChatRaw.id,
                    title: (newChatRaw as any).title ?? topicRes.response,
                    topic: (newChatRaw as any).topic ?? topicRes.response,
                    createdAt:
                        (newChatRaw as any).createdAt ??
                        new Date().toISOString(),
                };
                chatId = newChat.id;
                currentChatId = chatId;
                // Only add if not already in the list
                if (!chats.some((c) => c.id === newChat.id)) {
                    chats = [...chats, newChat];
                }
            }

            const correction: ApiResponse | undefined =
                await generateAIResponseCorrection(userInput);
            if (!correction?.response) {
                throw new Error("Failed to generate correction");
            }

            let fullUserMsg = `${userInput}\n\nCorreção: ${correction.response}`;

            const savedUserMsg = await addMessageToChat(
                chatId,
                fullUserMsg,
                "user",
                "request",
            );
            messages = [
                ...messages,
                {
                    id: savedUserMsg.id,
                    text: fullUserMsg,
                    sender: "user",
                    type: "request",
                },
            ];

            const dialogRes: ApiResponse | undefined =
                await generateAIResponseDialog(correction.response, chatId);
            if (!dialogRes?.response) {
                throw new Error("Failed to generate dialog response");
            }

            const translation: ApiResponse | undefined =
                await generateAIResponseTranslation(dialogRes.response);
            if (!translation?.response) {
                throw new Error("Failed to generate translation");
            }

            const fullAIResponse = `${dialogRes.response}\n\n[TRANSLATION]: ${translation.response}`;

            const savedBotMsg = await addMessageToChat(
                chatId,
                fullAIResponse,
                "ai",
                "response",
            );
            messages = [
                ...messages,
                {
                    id: savedBotMsg.id,
                    text: fullAIResponse,
                    sender: "ai",
                    type: "response",
                },
            ];
        } catch (err) {
            console.error("Error sending message:", err);
            inputMessage = userInput;
        } finally {
            isSending = false;
        }
    }

    async function handleChatSelect(event: CustomEvent<string>) {
        const chatId = event.detail;
        currentChatId = chatId;
        const chat = chats.find((c) => c.id === chatId);
        selectedChat = chat || null;
        try {
            const res = await getMessagesByChatId(chatId);
            messages = res.map((m: any) => ({
                id: m.id,
                text: m.content,
                sender: m.sender,
                type: m.type,
            }));
        } catch (error) {
            console.error("Error loading messages:", error);
            messages = [];
        }
    }

    function scrollToBottomSmooth() {
        if (messagesEndRef) {
            messagesEndRef.scrollIntoView({ behavior: "smooth" });
        }
    }

    function handleChatDeleted(event: CustomEvent<string>) {
        const deletedChatId = event.detail;
        chats = chats.filter((chat) => chat.id !== deletedChatId);

        // Se o chat deletado era o chat atual, limpa a conversa
        if (currentChatId === deletedChatId) {
            currentChatId = null;
            messages = [];
            selectedChat = null;
        }
    }

    function handleAudioRecorded(event: CustomEvent<{ text: string }>) {
        const { text } = event.detail;
        console.log("📨 Evento audioRecorded recebido no chat:", { text });
        if (text && text.trim()) {
            inputMessage = text.trim();
            console.log("✅ Enviando mensagem transcrita automaticamente");
            // Enviar automaticamente após transcrição bem-sucedida
            handleSend();
        } else {
            console.log(
                "ℹ️ Nenhum texto transcrito. Campo de texto disponível para digitação manual.",
            );
            // Mostrar notificação visual temporária
            showTranscriptionNotice = true;
            setTimeout(() => {
                showTranscriptionNotice = false;
            }, 4000);
        }
        // Se não houver transcrição, o campo de texto fica disponível para digitação manual
    }

    $: sidebarChats = chats
        .filter((chat) => chat.id != null)
        .filter(
            (chat, index, self) =>
                index === self.findIndex((c) => c.id === chat.id),
        )
        .map((chat) => ({
            id: chat.id,
            topic: chat.title || chat.topic || "Untitled Chat",
        }));
    let playingMsgId: string | null = null;
    let toastMessage = "";

    async function quickAddFlashcard(term: string, context: string) {
        if (!term) return;
        toastMessage = "🤖 Gerando Flashcard com IA...";
        try {
            const resGen = await fetch("http://localhost:8082/api/flashcards/generate", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ term, context_sentence: context })
            });
            const genData = resGen.ok ? await resGen.json() : { back: term, explanation: "" };

            const resSave = await fetch("http://localhost:8082/api/flashcards", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-User-ID": "default_user"
                },
                body: JSON.stringify({
                    front: term,
                    back: genData.back || term,
                    context_sentence: context || genData.context_sentence,
                    explanation: genData.explanation || ""
                })
            });

            if (resSave.ok) {
                toastMessage = "🎴 Flashcard adicionado com sucesso!";
                setTimeout(() => { toastMessage = ""; }, 3000);
            } else if (resSave.status === 409) {
                toastMessage = `⚠️ O flashcard para "${term.trim()}" já existe!`;
                setTimeout(() => { toastMessage = ""; }, 4000);
            } else {
                toastMessage = "Erro ao adicionar flashcard";
                setTimeout(() => { toastMessage = ""; }, 3000);
            }
        } catch (e) {
            console.error("Erro ao adicionar flashcard:", e);
            toastMessage = "Erro ao adicionar flashcard";
            setTimeout(() => { toastMessage = ""; }, 3000);
        }
    }

    let selectedText = "";
    let selectionCoords = { x: 0, y: 0 };
    let showSelectionMenu = false;

    function handleTextSelection() {
        if (typeof window === "undefined") return;
        const selection = window.getSelection();
        if (!selection || selection.isCollapsed) {
            showSelectionMenu = false;
            return;
        }
        const text = selection.toString().trim();
        if (text.length > 0 && text.length < 60) {
            selectedText = text;
            try {
                const range = selection.getRangeAt(0);
                const rect = range.getBoundingClientRect();
                selectionCoords = {
                    x: rect.left + rect.width / 2,
                    y: rect.top - 40
                };
                showSelectionMenu = true;
            } catch (e) {
                showSelectionMenu = false;
            }
        } else {
            showSelectionMenu = false;
        }
    }

    function getPracticingLanguage(): string {
        try {
            const token = Cookies.get("authToken");
            if (!token) return "en-US";
            const base64Url = token.split(".")[1];
            const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
            const jsonPayload = decodeURIComponent(
                atob(base64)
                    .split("")
                    .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
                    .join("")
            );
            const decoded = JSON.parse(jsonPayload);
            const lang = decoded.language || "english";
            
            const mapping: Record<string, string> = {
                "english": "en-US",
                "japanese": "ja-JP",
                "spanish": "es-ES",
                "french": "fr-FR",
                "german": "de-DE",
                "italian": "it-IT"
            };
            return mapping[lang.toLowerCase()] || "en-US";
        } catch (e) {
            return "en-US";
        }
    }

    function playTTS(msgId: string, text: string) {
        if (typeof window === "undefined" || !("speechSynthesis" in window)) {
            console.error("Speech synthesis not supported");
            return;
        }

        if (window.speechSynthesis.speaking && playingMsgId === msgId) {
            window.speechSynthesis.cancel();
            playingMsgId = null;
            return;
        }

        window.speechSynthesis.cancel();

        const cleanText = text
            .replace(/\[TRANSLATION\]:.*/gi, "")
            .replace(/Tradução:.*/gi, "")
            .trim();

        const utterance = new SpeechSynthesisUtterance(cleanText);
        utterance.lang = getPracticingLanguage();

        const voices = window.speechSynthesis.getVoices();
        const preferredVoices = voices.filter(v => v.lang.startsWith(utterance.lang));
        const naturalVoice = preferredVoices.find(v => v.name.toLowerCase().includes("natural") || v.name.toLowerCase().includes("google"));
        if (naturalVoice) {
            utterance.voice = naturalVoice;
        } else if (preferredVoices.length > 0) {
            utterance.voice = preferredVoices[0];
        }

        utterance.onend = () => {
            playingMsgId = null;
        };

        utterance.onerror = () => {
            playingMsgId = null;
        };

        playingMsgId = msgId;
        window.speechSynthesis.speak(utterance);
    }
</script>

<svelte:head>
    <title>SpeakUp</title>
    <meta name="description" content="Pratique conversação com Inteligência Artificial na SpeakUp. Melhore sua pronúncia, vocabulário e entonação de forma dinâmica e interativa." />
    <meta property="og:title" content="SpeakUp" />
    <meta property="og:description" content="Pratique conversação com Inteligência Artificial na SpeakUp. Melhore sua pronúncia, vocabulário e entonação de forma dinâmica e interativa." />
    <meta name="twitter:title" content="SpeakUp" />
    <meta name="twitter:description" content="Pratique conversação com Inteligência Artificial na SpeakUp. Melhore sua pronúncia, vocabulário e entonação de forma dinâmica e interativa." />
</svelte:head>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="Chat" on:mouseup={handleTextSelection}>
    <!-- Top Navbar Padronizado ocupando 100% da largura da página -->
    <nav class="top-nav">
        <div class="nav-brand" on:click={() => goto("/dashboard")} role="button" tabindex="0">
            <img src="/logo.png" alt="SpeakUp Logo" width="32" />
            <span class="brand-name">SpeakUp</span>
        </div>

        <div class="nav-user">
            <button class="nav-btn" on:click={() => goto("/dashboard")}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" stroke-linecap="round" stroke-linejoin="round"/>
                    <polyline points="9 22 9 12 15 12 15 22"/>
                </svg>
                <span>Hub</span>
            </button>

            <button class="nav-btn" on:click={() => goto("/perfil")}>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="8" r="4"/>
                    <path d="M5 20c0-4 3-7 7-7s7 3 7 7"/>
                </svg>
                <span>Perfil</span>
            </button>

            <button class="nav-btn btn-logout" on:click={handleLogout}>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                    <path d="M16 17l5-5-5-5M21 12H9"/>
                </svg>
                <span>Sair</span>
            </button>
        </div>
    </nav>

    {#if showSelectionMenu}
        <div
            class="selectionPopover"
            style="top: {selectionCoords.y}px; left: {selectionCoords.x}px;"
        >
            <button
                type="button"
                on:click={() => {
                    quickAddFlashcard(selectedText, "");
                    showSelectionMenu = false;
                }}
            >
                🎴 Salvar "{selectedText}" em Flashcards
            </button>
        </div>
    {/if}

    <div class="chatContent">
        <Sidebar
            chats={sidebarChats}
            on:selectChat={handleChatSelect}
            on:chatDeleted={handleChatDeleted}
            bind:sidebarOpen
            selectedChat={selectedChat ? { id: selectedChat.id } : null}
            on:newChat={() => {
                currentChatId = null;
                messages = [];
                inputMessage = "";
                selectedChat = null;
                sidebarOpen = false;
            }}
        />

        <main class="main" class:sidebarClosed={!sidebarOpen}>
            <div class="chatBody" bind:this={chatBodyRef}>
            {#if messages.length === 0}
                <div class="welcomeBox">
                    <div class="welcomeContent">
                        <div class="featuresGrid">
                            <div class="featureItem">
                                <div class="featureIcon magic">
                                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/><path d="M5 3v4"/><path d="M19 17v4"/><path d="M3 5h4"/><path d="M17 19h4"/></svg>
                                </div>
                                <h3>Correção</h3>
                                <p>Ajustes automáticos de gramática e ortografia.</p>
                            </div>
                            <div class="featureItem">
                                <div class="featureIcon globe">
                                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20"/><path d="M2 12h20"/></svg>
                                </div>
                                <h3>Tradução</h3>
                                <p>Veja o significado das frases instantaneamente.</p>
                            </div>
                            <div class="featureItem">
                                <div class="featureIcon chat">
                                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
                                </div>
                                <h3>IA Chat</h3>
                                <p>Diálogo natural para fluidez e confiança.</p>
                            </div>
                        </div>

                        <div class="examplesSection">
                            <div class="examplesHeader">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><path d="M12 17h.01"/></svg>
                                <h3>Exemplos de uso:</h3>
                            </div>

                            <div class="exampleCard">
                                <div class="exampleInput">
                                    <span class="exampleLabel">Você escreve:</span>
                                    <span class="exampleText">"I ned a car"</span>
                                </div>
                                <div class="exampleArrow">→</div>
                                <div class="exampleOutput">
                                    <span class="exampleLabel">IA corrige:</span>
                                    <span class="exampleText">"I need a car"</span>
                                </div>
                            </div>

                            <div class="exampleCard">
                                <div class="exampleInput">
                                    <span class="exampleLabel">Você escreve:</span>
                                    <span class="exampleText">"Let's go beach tomorrow?"</span>
                                </div>
                                <div class="exampleArrow">→</div>
                                <div class="exampleOutput">
                                    <span class="exampleLabel">IA corrige:</span>
                                    <span class="exampleText">"Let's go to the beach tomorrow?"</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}

            {#each messages as msg (msg.id)}
                <div
                    class:userMessage={msg.sender === "user"}
                    class:botMessage={msg.sender === "ai"}
                >
                    <div>
                        {#if msg.sender === "ai"}
                            <button
                                class="ttsButton botBubbleTts"
                                class:speaking={playingMsgId === msg.id}
                                on:click={() => playTTS(msg.id, msg.text.split("\n\n")[0])}
                                title="Ouvir pronúncia"
                                type="button"
                            >
                                {#if playingMsgId === msg.id}
                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                        <rect x="4" y="4" width="16" height="16" rx="2" ry="2"/>
                                    </svg>
                                {:else}
                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                        <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/>
                                        <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"/>
                                    </svg>
                                {/if}
                            </button>
                        {/if}
                        <p>{msg.text.split("\n\n")[0]}</p>

                        {#if msg.type === "request"}
                            <div class="responseExtras">
                                <p class="correctionTextLine">
                                    <span>
                                        <strong>Correção:</strong>
                                        {msg.text
                                            .split("\n\n")[1]
                                            ?.replace("Correção: ", "") || ""}
                                    </span>
                                    <button
                                        class="ttsButton inlineTts"
                                        class:speaking={playingMsgId === `${msg.id}-corr`}
                                        on:click={() => playTTS(`${msg.id}-corr`, msg.text.split("\n\n")[1]?.replace("Correção: ", "") || "")}
                                        title="Ouvir pronúncia da correção"
                                        type="button"
                                    >
                                        {#if playingMsgId === `${msg.id}-corr`}
                                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                                <rect x="4" y="4" width="16" height="16" rx="2" ry="2"/>
                                            </svg>
                                        {:else}
                                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                                <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/>
                                                <path d="M15.54 8.46a5 5 0 0 1 0 7.07"/>
                                            </svg>
                                        {/if}
                                    </button>
                                    <button
                                        class="flashcardAddBtn"
                                        on:click={() => quickAddFlashcard(msg.text.split("\n\n")[1]?.replace("Correção: ", "") || "", msg.text.split("\n\n")[0])}
                                        title="Salvar nos Flashcards"
                                        type="button"
                                    >
                                        🎴 +Flashcard
                                    </button>
                                </p>

                            </div>
                        {/if}
                        {#if msg.type === "response"}
                            <div class="responseExtras">
                                <p>
                                    <strong>Tradução:</strong>
                                    {msg.text.split("[TRANSLATION]: ")[1] || ""}
                                </p>
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}

            <div bind:this={messagesEndRef}></div>

            {#if showScrollDown}
                <button
                    class="scrollDownButton"
                    on:click={scrollToBottomSmooth}
                    type="button"
                    aria-label="Scroll to bottom"
                >
                    <svg
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="currentColor"
                    >
                        <path
                            d="M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6z"
                        />
                    </svg>
                </button>
            {/if}
        </div>

        {#if showTranscriptionNotice}
            <div class="transcriptionNotice">
                <svg
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                >
                    <path
                        d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"
                    />
                </svg>
                <span
                    >Transcrição indisponível. Digite sua mensagem abaixo.</span
                >
            </div>
        {/if}

        {#if toastMessage}
            <div class="transcriptionNotice">
                <span>{toastMessage}</span>
            </div>
        {/if}

        <footer class="chatFooter">
            <div class="inputBox">
                <input
                    type="text"
                    placeholder="Digite sua mensagem aqui..."
                    bind:value={inputMessage}
                    on:keydown={(e) =>
                        e.key === "Enter" &&
                        !isSending &&
                        !isRecording &&
                        handleSend()}
                    disabled={isSending || isRecording}
                />
                <AudioRecorder
                    bind:isRecording
                    disabled={isSending}
                    on:audioRecorded={handleAudioRecorded}
                />
                <button
                    on:click={handleSend}
                    disabled={isSending || isRecording}
                    type="button"
                >
                    {#if isSending}
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
                    {:else}
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="currentColor"
                        >
                            <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                        </svg>
                    {/if}
                </button>
            </div>
        </footer>
    </main>
    </div>
</div>

<style>
    .Chat {
        display: flex;
        flex-direction: column;
        height: 100vh;
        background-color: #111;
        color: #f0f0f0;
        overflow: hidden;
    }

    .chatContent {
        display: flex;
        flex: 1;
        position: relative;
        overflow: hidden;
    }

    .main {
        flex: 1;
        display: flex;
        flex-direction: column;
        background-color: #111;
        margin-left: 360px;
        transition: margin-left 0.3s ease;
        height: 100%;
    }

    .main.sidebarClosed {
        margin-left: 60px;
    }

    /* Top Navigation Header */
    .top-nav {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.25rem 2.5rem;
        background-color: #0f0f0f;
        border-bottom: 1px solid #1a1a1a;
        z-index: 101;
        box-sizing: border-box;
    }

    .nav-brand {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        cursor: pointer;
    }

    .brand-name {
        font-size: 1.2rem;
        font-weight: 700;
        color: #ffffff;
        letter-spacing: -0.01em;
    }

    .nav-user {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .nav-btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        background: transparent;
        border: 1px solid #262626;
        color: #cccccc;
        padding: 0.5rem 0.9rem;
        border-radius: 6px;
        font-size: 0.85rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .nav-btn:hover {
        background-color: #1a1a1a;
        color: #ffffff;
        border-color: #444444;
    }

    .btn-logout:hover {
        border-color: rgba(239, 68, 68, 0.4);
        color: #ef4444;
    }

    .chatBody {
        flex: 1;
        padding: 1rem;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        scrollbar-width: none;
    }

    .chatBody::-webkit-scrollbar {
        display: none;
    }

    .userMessage,
    .botMessage {
        display: flex;
    }

    .userMessage {
        justify-content: flex-end;
    }

    .botMessage {
        justify-content: flex-start;
    }

    .userMessage div,
    .botMessage div {
        max-width: 32rem;
        background-color: #222;
        color: white;
        padding: 1rem;
        border-radius: 0.5rem;
    }

    .botMessage div {
        position: relative;
        background-color: #1f1f1f;
        color: #e0e0e0;
        border-radius: 0 20px 20px 20px;
        padding-right: 2.5rem !important;
    }

    .userMessage div {
        color: white;
        border-radius: 20px 0 20px 20px;
    }

    .responseExtras {
        margin-top: 0.75rem;
        padding-top: 0.75rem;
        border-top: 1px solid #333;
        font-size: 0.75rem;
    }

    .responseExtras p {
        margin-bottom: 0.25rem;
        color: #aaa;
    }



    .botBubbleTts {
        position: absolute;
        top: 10px;
        right: 10px;
        z-index: 5;
    }

    .ttsButton {
        background: none;
        border: none;
        color: rgba(255, 255, 255, 0.4);
        cursor: pointer;
        padding: 6px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
        flex-shrink: 0;
        margin-top: -2px;
    }

    .ttsButton:hover {
        color: #fff;
        background-color: rgba(255, 255, 255, 0.08);
    }

    .ttsButton.speaking {
        color: #a855f7;
        background-color: rgba(168, 85, 247, 0.15);
        animation: ttsPulse 1.5s infinite;
    }

    .correctionTextLine {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 1rem;
        margin-bottom: 0.25rem;
    }

    .inlineTts {
        padding: 4px;
        margin-top: 0;
    }

    @keyframes ttsPulse {
        0% { opacity: 0.7; }
        50% { opacity: 1; }
        100% { opacity: 0.7; }
    }

    .chatFooter {
        background-color: #1a1a1a;
        padding: 1rem;
        border-top: 1px solid #333;
    }

    .inputBox {
        display: flex;
        align-items: center;
        background-color: #2a2a2a;
        border-radius: 0.5rem;
        padding: 0.25rem;
        gap: 0.5rem;
    }

    .inputBox input {
        flex: 1;
        background: transparent;
        border: none;
        color: white;
        padding: 0.75rem;
        font-size: 0.875rem;
        outline: none;
        font-family: "Karla", sans-serif;
    }

    .inputBox button {
        background: none;
        border: none;
        color: #bbb;
        cursor: pointer;
        padding: 0.5rem;
        border-radius: 0.25rem;
    }

    .inputBox button:hover:not(:disabled) {
        color: #fff;
        background-color: #333;
    }

    .inputBox button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .welcomeBox {
        max-width: 800px;
        margin: auto;
        border-radius: 12px;
        background: rgba(255, 255, 255, 0.03);
        backdrop-filter: blur(10px);
        -webkit-backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.1);
        box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
        overflow: hidden;
        animation: fadeInUp 0.8s cubic-bezier(0.16, 1, 0.3, 1);
    }

    @keyframes fadeInUp {
        from {
            opacity: 0;
            transform: translateY(30px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .welcomeContent {
        padding: 2.5rem;
        display: flex;
        flex-direction: column;
        gap: 3rem;
    }

    .featuresGrid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 1.5rem;
    }

    .featureItem {
        display: flex;
        flex-direction: column;
        align-items: center;
        text-align: center;
        gap: 1rem;
        padding: 1.5rem;
        background: rgba(255, 255, 255, 0.02);
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.05);
        transition: all 0.3s ease;
    }

    .featureItem:hover {
        background: rgba(255, 255, 255, 0.05);
        border-color: rgba(255, 255, 255, 0.1);
        transform: translateY(-5px);
    }

    .featureIcon {
        width: 48px;
        height: 48px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 0.5rem;
    }

    .featureIcon.magic { background: rgba(147, 51, 234, 0.15); color: #a855f7; }
    .featureIcon.globe { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
    .featureIcon.chat { background: rgba(16, 185, 129, 0.15); color: #10b981; }

    .featureItem h3 {
        margin: 0;
        font-size: 1.1rem;
        font-weight: 600;
        color: #fff;
    }

    .featureItem p {
        margin: 0;
        font-size: 0.9rem;
        color: rgba(255, 255, 255, 0.5);
        line-height: 1.5;
    }

    .examplesSection {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .examplesHeader {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        color: rgba(255, 255, 255, 0.7);
        margin-bottom: 0.5rem;
    }

    .examplesHeader h3 {
        margin: 0;
        font-size: 1rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 1px;
    }

    .exampleCard {
        display: grid;
        grid-template-columns: 1fr auto 1fr;
        gap: 2rem;
        align-items: center;
        padding: 1.25rem 1.5rem;
        background: rgba(255, 255, 255, 0.02);
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.05);
        transition: all 0.2s ease;
    }

    .exampleCard:hover {
        background: rgba(255, 255, 255, 0.04);
        border-color: rgba(255, 255, 255, 0.1);
    }

    .exampleInput,
    .exampleOutput {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
    }

    .exampleLabel {
        font-size: 0.7rem;
        color: rgba(255, 255, 255, 0.3);
        text-transform: uppercase;
        font-weight: 700;
        letter-spacing: 0.5px;
    }

    .exampleText {
        font-family: 'JetBrains Mono', 'Fira Code', monospace;
        font-size: 0.95rem;
    }

    .exampleInput .exampleText {
        color: rgba(255, 255, 255, 0.5);
    }

    .exampleOutput .exampleText {
        color: #fff;
        font-weight: 500;
    }

    .exampleArrow {
        color: rgba(255, 255, 255, 0.15);
        font-size: 1.25rem;
    }

    @media (max-width: 768px) {
        .welcomeBox {
            margin: 1rem;
            border-radius: 8px;
        }

        .featuresGrid {
            grid-template-columns: 1fr;
            gap: 1rem;
        }

        .featureItem {
            flex-direction: row;
            text-align: left;
            padding: 1.25rem;
        }

        .featureIcon {
            margin-bottom: 0;
            flex-shrink: 0;
        }

        .exampleCard {
            grid-template-columns: 1fr;
            gap: 1rem;
            padding: 1.25rem;
        }

        .exampleArrow {
            transform: rotate(90deg);
            text-align: center;
            margin: -0.5rem 0;
        }
    }

    .scrollDownButton {
        position: fixed;
        bottom: 100px;
        right: 20px;
        background-color: #2a2a2a;
        color: white;
        border: none;
        border-radius: 50%;
        width: 40px;
        height: 40px;
        display: flex;
        justify-content: center;
        align-items: center;
        cursor: pointer;
        z-index: 10;
        transition: background-color 0.2s;
    }

    .scrollDownButton:hover {
        background-color: #383838;
    }

    .transcriptionNotice {
        position: fixed;
        bottom: 100px;
        left: 50%;
        transform: translateX(-50%);
        background: #1a1a1a;
        color: white;
        padding: 0.75rem 1.5rem;
        border-radius: 12px;
        display: flex;
        align-items: center;
        gap: 0.75rem;
        border: 1px solid rgba(255, 255, 255, 0.2);
        z-index: 100;
        animation:
            slideInUp 0.3s ease-out,
            fadeOut 0.5s ease-in 3.5s forwards;
        font-size: 0.9rem;
        font-weight: 500;
    }

    .transcriptionNotice svg {
        flex-shrink: 0;
    }

    @keyframes slideInUp {
        from {
            opacity: 0;
            transform: translateX(-50%) translateY(20px);
        }
        to {
            opacity: 1;
            transform: translateX(-50%) translateY(0);
        }
    }

    @keyframes fadeOut {
        from {
            opacity: 1;
        }
        to {
            opacity: 0;
        }
    }

    .spinner {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    .flashcardAddBtn {
        background: rgba(168, 85, 247, 0.2);
        color: #c084fc;
        border: 1px solid rgba(168, 85, 247, 0.4);
        border-radius: 6px;
        padding: 2px 8px;
        font-size: 0.75rem;
        font-weight: 600;
        cursor: pointer;
        margin-left: 8px;
        transition: all 0.2s;
    }

    .flashcardAddBtn:hover {
        background: rgba(168, 85, 247, 0.4);
        color: white;
    }

    .selectionPopover {
        position: fixed;
        transform: translateX(-50%);
        background: linear-gradient(135deg, #1e1e2e, #2d1b4e);
        border: 1px solid rgba(168, 85, 247, 0.6);
        border-radius: 8px;
        padding: 4px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.6);
        z-index: 9999;
        animation: fadeIn 0.15s ease-out;
    }

    .selectionPopover button {
        background: #8b5cf6;
        color: white;
        border: none;
        padding: 6px 12px;
        border-radius: 6px;
        font-size: 0.8rem;
        font-weight: 600;
        cursor: pointer;
        white-space: nowrap;
    }

    .selectionPopover button:hover {
        background: #7c3aed;
    }
</style>
