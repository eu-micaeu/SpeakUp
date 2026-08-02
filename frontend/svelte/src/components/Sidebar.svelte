<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { goto } from "$app/navigation";
  import { deleteChat } from "../utils/api";
  import Cookies from "js-cookie";

  export let chats: Array<{ id: string; topic: string }>;
  export let sidebarOpen: boolean;
  export let selectedChat: { id: string } | null;

  const dispatch = createEventDispatcher();
  let openSettings = false;

  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }

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

  function goToPlans() {
    goto("/planos");
    openSettings = false;
  }

  function goToIndex() {
    Cookies.remove("authToken");
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
  <div class="divToggleButton">
    <button
      class="toggleButton"
      class:rotated={sidebarOpen}
      on:click={toggleSidebar}
      aria-label="Toggle sidebar"
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
        <path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z" />
      </svg>
    </button>
  </div>

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
      <!-- <button on:click={() => goto("/teaching-plan")}>
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
      </button> -->

      <!-- <button on:click={() => goto("/palavreco")}>
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
      </button> -->

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
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="modal" on:click={() => (openSettings = false)}>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="modal-content" on:click={(e) => e.stopPropagation()}>
        <div class="modal-header">
          <h3>Configurações</h3>
          <button
            class="close-btn"
            on:click={() => (openSettings = false)}
            aria-label="Fechar"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
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

        <div class="modal-options">
          <button class="modal-option" on:click={goToPerfil}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <circle
                cx="12"
                cy="8"
                r="4"
                stroke="currentColor"
                stroke-width="2"
              />
              <path
                d="M5 20c0-4 3-7 7-7s7 3 7 7"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
              />
            </svg>
            <span>Perfil</span>
          </button>

          <!-- <button class="modal-option" on:click={goToPlans}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path
                d="M3 7h18M5 7v10a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M8 11h8M8 15h5"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            <span>Planos</span>
          </button> -->

          <button class="modal-option" on:click={goToIndex}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path
                d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M16 17l5-5-5-5M21 12H9"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            <span>Sair</span>
          </button>
        </div>
      </div>
    </div>
  {/if}
</aside>

<style>
  .sidebar {
    width: 360px;
    background-color: #1a1a1a;
    border-right: 1px solid #333;
    display: flex;
    flex-direction: column;
    transition: width 0.1s ease;
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    z-index: 100;
  }

  .sidebarClosed {
    width: 60px;
  }

  .divToggleButton{
    display: flex;
    justify-content: flex-end;
  }

  .toggleButton {
    background: none;
    border: none;
    color: #ccc;
    cursor: pointer;
    padding: 1rem;
    align-items: right;
    justify-content: center;
    border-radius: 0.5rem;
    min-height: 60px;
    margin: 0 auto;
  }

  .toggleButton.rotated {
    transform: rotate(90deg);
    margin: 0;
  }

  .sidebarContent {
    padding: 1rem;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .sidebarClosed .sidebarContent {
    opacity: 0;
    pointer-events: none;
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
    backdrop-filter: blur(4px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    animation: fadeIn 0.1s ease;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .modal-content {
    background: linear-gradient(145deg, #1e1e1e, #1a1a1a);
    padding: 0;
    border-radius: 8px;
    color: white;
    min-width: 320px;
    max-width: 400px;
    animation: slideUp 0.2s ease;
  }

  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    background: linear-gradient(135deg, #ffffff, #e0e0e0);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .close-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: white;
    transform: rotate(90deg);
  }

  .modal-options {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .modal-option {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: rgba(255, 255, 255, 0.05);
    color: rgba(255, 255, 255, 0.9);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 1rem;
    font-weight: 500;
  }

  .modal-option:hover {
    background: rgba(255, 255, 255, 0.1);
    border-color: rgba(255, 255, 255, 0.3);
    transform: translateX(4px);
  }

  .modal-option svg {
    flex-shrink: 0;
    color: rgba(255, 255, 255, 0.8);
    transition: all 0.3s ease;
  }

  .modal-option:hover svg {
    color: white;
    transform: scale(1.1);
  }

  .modal-option span {
    flex: 1;
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
