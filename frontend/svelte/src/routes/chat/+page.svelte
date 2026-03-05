<script lang="ts">
    import { onMount, tick } from "svelte";
    import { browser } from "$app/environment";
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

            const fullUserMsg = `${userInput}\n\nCorreção: ${correction.response}`;

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
</script>

<div class="Chat">
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
                    <div class="welcomeHeader">
                        <h2>Vamos iniciar uma conversa?</h2>
                    </div>

                    <div class="welcomeContent">
                        <div class="featureSection">
                            <div class="featureIcon">
                                <svg
                                    width="24"
                                    height="24"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                >
                                    <path
                                        d="M22 11.08V12a10 10 0 1 1-5.93-9.14"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                    <path
                                        d="M22 4L12 14.01l-3-3"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                </svg>
                            </div>
                            <div class="featureText">
                                <h3>Como funciona?</h3>
                                <p>
                                    Digite suas frases e nossa IA irá corrigir
                                    erros gramaticais, fornecer traduções e
                                    continuar a conversa!
                                </p>
                            </div>
                        </div>

                        <div class="examplesSection">
                            <div class="examplesHeader">
                                <svg
                                    width="20"
                                    height="20"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                >
                                    <circle
                                        cx="12"
                                        cy="12"
                                        r="10"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    />
                                    <path
                                        d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                    <path
                                        d="M12 17h.01"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                </svg>
                                <h3>Exemplos de uso:</h3>
                            </div>

                            <div class="exampleCard">
                                <div class="exampleInput">
                                    <span class="exampleLabel"
                                        >Você escreve:</span
                                    >
                                    <span class="exampleText"
                                        >"I ned a car"</span
                                    >
                                </div>
                                <div class="exampleArrow">→</div>
                                <div class="exampleOutput">
                                    <span class="exampleLabel">IA corrige:</span
                                    >
                                    <span class="exampleText"
                                        >"I need a car"</span
                                    >
                                </div>
                            </div>

                            <div class="exampleCard">
                                <div class="exampleInput">
                                    <span class="exampleLabel"
                                        >Você escreve:</span
                                    >
                                    <span class="exampleText"
                                        >"Let's go beach tomorrow?"</span
                                    >
                                </div>
                                <div class="exampleArrow">→</div>
                                <div class="exampleOutput">
                                    <span class="exampleLabel">IA corrige:</span
                                    >
                                    <span class="exampleText"
                                        >"Let's go to the beach tomorrow?"</span
                                    >
                                </div>
                            </div>

                            <!-- <div class="exampleCard">
                                <div class="exampleInput">
                                    <span class="exampleLabel"
                                        >Você escreve:</span
                                    >
                                    <span class="exampleText"
                                        >"I don't know how say this"</span
                                    >
                                </div>
                                <div class="exampleArrow">→</div>
                                <div class="exampleOutput">
                                    <span class="exampleLabel">IA corrige:</span
                                    >
                                    <span class="exampleText"
                                        >"I don't know how to say this"</span
                                    >
                                </div>
                            </div> -->
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
                        <p>{msg.text.split("\n\n")[0]}</p>
                        {#if msg.type === "request"}
                            <div class="responseExtras">
                                <p>
                                    <strong>Correção:</strong>
                                    {msg.text
                                        .split("\n\n")[1]
                                        ?.replace("Correção: ", "") || ""}
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

<style>
    .Chat {
        display: flex;
        height: 100vh;
        background-color: #111;
        color: #f0f0f0;
        overflow: hidden;
    }

    .main {
        flex: 1;
        display: flex;
        flex-direction: column;
        background-color: #111;
        margin-left: 360px;
        transition: margin-left 0.3s ease;
    }

    .main.sidebarClosed {
        margin-left: 60px;
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
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .botMessage div {
        background-color: #1f1f1f;
        color: #e0e0e0;
        border-radius: 0 20px 20px 20px;
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
        border-radius: 20px;
        border: 1px solid rgba(255, 255, 255, 0.1);
        overflow: hidden;
        animation: fadeInUp 0.6s ease;
    }

    @keyframes fadeInUp {
        from {
            opacity: 0;
            transform: translateY(20px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .welcomeHeader {
        text-align: center;
        padding: 2.5rem 2rem 2rem;
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.05),
            rgba(255, 255, 255, 0.02)
        );
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    .welcomeHeader h2 {
        font-size: 2rem;
        font-weight: 700;
        margin: 0 0 0.5rem 0;
        color: #ffffff;
    }

    .welcomeContent {
        padding: 2rem;
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }

    .featureSection {
        display: flex;
        gap: 1.5rem;
        align-items: flex-start;
        padding: 1.5rem;
        background: rgba(255, 255, 255, 0.05);
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .featureIcon {
        flex-shrink: 0;
        width: 50px;
        height: 50px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid rgba(255, 255, 255, 0.2);
    }

    .featureIcon svg {
        color: #ffffff;
    }

    .featureText h3 {
        margin: 0 0 0.5rem 0;
        font-size: 1.2rem;
        color: white;
    }

    .featureText p {
        margin: 0;
        color: rgba(255, 255, 255, 0.8);
        line-height: 1.6;
    }

    .examplesSection {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .examplesHeader {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: rgba(255, 255, 255, 0.9);
        margin-bottom: 0.5rem;
    }

    .examplesHeader h3 {
        margin: 0;
        font-size: 1.1rem;
        font-weight: 600;
    }

    .exampleCard {
        display: grid;
        grid-template-columns: 1fr auto 1fr;
        gap: 1rem;
        align-items: center;
        padding: 1rem;
        background: rgba(255, 255, 255, 0.03);
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.08);
        transition: all 0.3s ease;
    }

    .exampleCard:hover {
        background: rgba(255, 255, 255, 0.06);
        border-color: rgba(255, 255, 255, 0.15);
        transform: translateX(4px);
    }

    .exampleInput,
    .exampleOutput {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .exampleLabel {
        font-size: 0.75rem;
        color: rgba(255, 255, 255, 0.5);
        text-transform: uppercase;
        font-weight: 600;
        letter-spacing: 0.5px;
    }

    .exampleInput .exampleText {
        color: rgba(255, 255, 255, 0.7);
        font-family: "Courier New", monospace;
        font-size: 0.95rem;
    }

    .exampleOutput .exampleText {
        color: rgba(255, 255, 255, 0.9);
        font-family: "Courier New", monospace;
        font-size: 0.95rem;
    }

    .exampleArrow {
        color: rgba(255, 255, 255, 0.4);
        font-size: 1.5rem;
        font-weight: bold;
    }

    @media (max-width: 768px) {
        .welcomeBox {
            margin: 1rem;
        }

        .welcomeHeader {
            padding: 2rem 1.5rem 1.5rem;
        }

        .welcomeHeader h2 {
            font-size: 1.5rem;
        }

        .welcomeContent {
            padding: 1.5rem;
        }

        .exampleCard {
            grid-template-columns: 1fr;
            gap: 0.75rem;
        }

        .exampleArrow {
            transform: rotate(90deg);
            text-align: center;
        }

        .featureSection {
            flex-direction: column;
            text-align: center;
        }

        .featureIcon {
            margin: 0 auto;
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
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
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
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.8);
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
</style>
