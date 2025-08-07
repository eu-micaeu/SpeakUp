<script lang="ts">
    import { onMount, tick } from "svelte";
    import { browser } from "$app/environment";
    import Sidebar from "../../components/Sidebar.svelte";
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
                chats = [...chats, newChat];
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

    function toggleSidebar() {
        sidebarOpen = !sidebarOpen;
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
</script>

<div class="Chat">
    <Sidebar
        chats={chats.map((chat) => ({
            id: chat.id,
            topic: chat.title || chat.topic || "Untitled Chat",
        }))}
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

    <main class="main">
        <header class="chatHeader">
            <button
                class="toggleIcon"
                class:rotated={!sidebarOpen}
                on:click={toggleSidebar}
                type="button"
                aria-label="Toggle sidebar"
            >
                <svg
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                >
                    <path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z" />
                </svg>
            </button>
        </header>

        <div class="chatBody" bind:this={chatBodyRef}>
            {#if messages.length === 0}
                <div class="welcomeBox">
                    <h2>Bem-vindo ao Chat de Prática!</h2>
                    <p>
                        Aqui você pode praticar o idioma conversando com nossa
                        IA.
                    </p>
                    <div>
                        <strong>Exemplos:</strong><br /><br />
                        <em>Entrada:</em> "I ned a car"<br />
                        <em>Saída:</em> "I need a car"<br /><br />
                        <em>Entrada:</em> "How are you doin?"<br />
                        <em>Saída:</em> "How are you doing?"<br /><br />
                        <em>Entrada:</em> "Let's go beach tomorrow?"<br />
                        <em>Saída:</em> "Let's go to the beach tomorrow?"<br
                        /><br />
                        <em>Entrada:</em> "I don't know how say this."<br />
                        <em>Saída:</em> "I don't know how to say this."
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

        <footer class="chatFooter">
            <div class="inputBox">
                <input
                    type="text"
                    placeholder="Digite sua mensagem aqui..."
                    bind:value={inputMessage}
                    on:keydown={(e) =>
                        e.key === "Enter" && !isSending && handleSend()}
                    disabled={isSending}
                />
                <button
                    on:click={handleSend}
                    disabled={isSending}
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
    }

    .chatHeader {
        background-color: #1a1a1a;
        padding: 1rem;
        border-bottom: 1px solid #333;
        display: flex;
        justify-content: space-between;
        align-items: center;
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

    .toggleIcon {
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        margin: 0.25rem 0;
        padding: 0.5rem;
        border-radius: 0.25rem;
        background: none;
        border: none;
        color: currentColor;
        transition: background-color 0.2s;
    }

    .toggleIcon:hover {
        background-color: #333;
    }

    .toggleIcon.rotated {
        transform: rotate(90deg);
    }

    .welcomeBox {
        color: #ccc;
        max-width: 600px;
        margin: auto;
        padding: 2rem;
        background-color: #1a1a1a;
        border-radius: 1rem;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
    }

    .welcomeBox h2 {
        font-size: 1.4rem;
        margin: 0 0 1rem 0;
    }

    .welcomeBox div {
        background-color: #2a2a2a;
        padding: 20px;
        border-radius: 0.5rem;
        margin-top: 1rem;
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
