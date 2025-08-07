<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { goto } from "$app/navigation";
  import { deleteChat } from "../utils/api";
  // import { removeAuthTokenFromCookies } from '../../utils/cookies';

  export let chats: Array<{ id: string; topic: string }>;
  export let sidebarOpen: boolean;
  export let selectedChat: { id: string } | null;

  const dispatch = createEventDispatcher();
  let openSettings = false;

  function handleNewChat() {
    dispatch("newChat");
  }

  function handleChatClick(chat: { id: string; topic: string }) {
    dispatch("selectChat", chat.id);
  }

  function goToPerfil() {
    goto("/perfil");
    openSettings = false;
  }

  function goToIndex() {
    // removeAuthTokenFromCookies();
    goto("/");
    openSettings = false;
  }

  async function handleDeleteChat(chatId: string) {
    try {
      await deleteChat(chatId);
      dispatch("chatDeleted", chatId);
    } catch (err) {
      console.error("Error deleting chat:", err);
    }
  }
</script>

<aside class:sidebar={true} class:sidebarClosed={!sidebarOpen}>
  <div class="sidebarContent">
    <div class="header">
      <img src="/logo.png" alt="Logo" width={35} />
      <h3>SpeakUp</h3>
    </div>

    <button class="newChat" on:click={handleNewChat}>
      <!-- SVG de "Add" -->
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
        <path
          d="M12 5v14m7-7H5"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>

    {#each chats as chat (chat.id)}
      <div class="chatItem" class:selected={selectedChat?.id === chat.id}>
        <button
          type="button"
          class="navItem"
          aria-current={selectedChat?.id === chat.id ? "page" : undefined}
          on:click={() => handleChatClick(chat)}
          on:keydown={(e) => {
            if (e.key === "Enter" || e.key === " ") handleChatClick(chat);
          }}
          tabindex="0"
        >
          <span>{chat.topic}</span>
        </button>
        <button
          type="button"
          class="deleteButton"
          on:click={() => handleDeleteChat(chat.id)}
          aria-label="Delete chat"
        >
          <!-- SVG de "Delete" -->
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path
              d="M6 6l12 12M6 18L18 6"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>
      </div>
    {/each}

    <div class="footer">
      <button on:click={() => goto("/teaching-plan")}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <rect
            x="4"
            y="4"
            width="16"
            height="16"
            rx="2"
            stroke="currentColor"
            stroke-width="2"
          />
        </svg>
        <span>Plano de Estudo</span>
      </button>

      <button on:click={() => goto("/palavreco")}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <circle
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="2"
          />
        </svg>
        <span>Palavreco</span>
      </button>

      <button on:click={() => (openSettings = true)}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <circle
            cx="12"
            cy="12"
            r="3"
            stroke="currentColor"
            stroke-width="2"
          />
          <path
            d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82A1.65 1.65 0 0 0 3 12c0-.55.27-1.04.7-1.35a1.65 1.65 0 0 0 .33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33h.01A1.65 1.65 0 0 0 10 3.09V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51h.01a1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82c.43.31.7.8.7 1.35s-.27 1.04-.7 1.35z"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        <span>Configurações</span>
      </button>
    </div>
  </div>

  {#if openSettings}
    <div class="modal">
      <div class="modal-content">
        <button on:click={goToPerfil}>Perfil</button>
        <button on:click={goToIndex}>Sair</button>
        <button on:click={() => (openSettings = false)}>Fechar</button>
      </div>
    </div>
  {/if}
</aside>

<style>
  .sidebar {
    width: 280px;
    background-color: #1a1a1a;
    border-right: 1px solid #333;
    display: flex;
    flex-direction: column;
    transition: margin-left 0.3s ease;
  }

  .sidebarClosed {
    margin-left: -280px;
  }

  .sidebarContent {
    padding: 1rem;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #333;
  }

  .header h3 {
    color: white;
    margin: 0;
    font-size: 1.2rem;
  }

  .newChat {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem;
    background-color: #2a2a2a;
    color: white;
    border: none;
    border-radius: 0.5rem;
    cursor: pointer;
    margin-bottom: 1rem;
    transition: background-color 0.2s;
  }

  .newChat:hover {
    background-color: #333;
  }

  .chatItem {
    display: flex;
    align-items: center;
    margin-bottom: 0.25rem;
    border-radius: 0.25rem;
    background: transparent;
    transition: background-color 0.2s;
  }

  .chatItem:hover {
    background-color: #2a2a2a;
  }

  .chatItem.selected {
    background-color: #333;
  }

  .navItem {
    flex: 1;
    display: flex;
    align-items: center;
    padding: 0.75rem;
    background: transparent;
    color: #ccc;
    border: none;
    cursor: pointer;
    text-align: left;
    border-radius: 0.25rem;
  }

  .navItem:hover {
    color: white;
  }

  .chatItem.selected .navItem {
    color: white;
  }

  .navItem span {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .deleteButton {
    background: none;
    border: none;
    color: #666;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 0.25rem;
    transition: color 0.2s;
    opacity: 0;
    margin-right: 0.5rem;
  }

  .chatItem:hover .deleteButton {
    opacity: 1;
  }

  .deleteButton:hover {
    color: #ff6b6b;
  }

  .footer {
    border-top: 1px solid #333;
    padding-top: 1rem;
    margin-top: auto;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .footer button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    background: transparent;
    color: #ccc;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
    text-align: left;
    transition: all 0.2s;
  }

  .footer button:hover {
    background-color: #2a2a2a;
    color: white;
  }

  .modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    background: #1a1a1a;
    padding: 2rem;
    border-radius: 0.5rem;
    color: white;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .modal-content button {
    padding: 0.75rem;
    background: transparent;
    color: white;
    border: 1px solid #333;
    border-radius: 0.25rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .modal-content button:hover {
    background-color: #2a2a2a;
  }

  @media (max-width: 768px) {
    .sidebar {
      position: fixed;
      top: 0;
      left: 0;
      height: 100vh;
      z-index: 999;
    }
  }
</style>
