<script lang="ts">
  import { onMount } from "svelte";
  import Cookies from "js-cookie";
  import { goto } from "$app/navigation";
  import { API_URL } from "../../utils/api";

  let userName = "Estudante";
  let userLevel = "B1";
  let dueFlashcardsCount = 0;

  onMount(async () => {
    await fetchUserData();
    await fetchFlashcardsCount();
  });

  async function fetchUserData() {
    try {
      const token = Cookies.get("authToken") || "";
      if (!token) return;
      const base64Url = token.split(".")[1];
      if (!base64Url) return;
      const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split("")
          .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
          .join("")
      );
      const decoded = JSON.parse(jsonPayload);
      if (decoded.name) userName = decoded.name.split(" ")[0];
      if (decoded.level) userLevel = decoded.level;
    } catch (e) {
      console.error("Erro ao ler token", e);
    }
  }

  async function fetchFlashcardsCount() {
    try {
      const res = await fetch(`${API_URL}/flashcards?due=true`, {
        headers: { "X-User-ID": "default_user" }
      });
      if (res.ok) {
        const cards = await res.json();
        dueFlashcardsCount = cards.length;
      }
    } catch (e) {
      console.error("Erro ao buscar pendências de flashcards", e);
    }
  }

  function handleLogout() {
    Cookies.remove("authToken");
    goto("/");
  }
</script>

<div class="dashboard-wrapper">
  <!-- Minimalist Top Navbar -->
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

  <main class="dashboard-content">
    <header class="hub-header">
      <span class="level-pill">Nível {userLevel}</span>
      <h1>Bem-vindo, {userName}!</h1>
      <p>Escolha como deseja praticar seu aprendizado hoje.</p>
    </header>

    <!-- The 2 Primary Mode Cards -->
    <div class="hub-grid">
      <!-- Option 1: AI Chat -->
      <div class="hub-card" on:click={() => goto("/chat")} role="button" tabindex="0" on:keydown={(e) => e.key === 'Enter' && goto("/chat")}>
        <div class="card-top">
          <div class="icon-box">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#5c6dff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
          </div>
          <span class="badge">CONVERSAÇÃO</span>
        </div>

        <div class="card-content">
          <h2>Tutor de IA</h2>
          <p>Conversação em tempo real com correções e sugestões gramaticais instantâneas.</p>
        </div>

        <div class="card-footer">
          <span>Iniciar Chat</span>
          <svg class="arrow-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
      </div>

      <!-- Option 2: Flashcards SRS -->
      <div class="hub-card" on:click={() => goto("/flashcards")} role="button" tabindex="0" on:keydown={(e) => e.key === 'Enter' && goto("/flashcards")}>
        <div class="card-top">
          <div class="icon-box">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#a855f7" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="5" width="18" height="14" rx="3"/>
              <path d="M3 10h18M8 15h3" stroke-linecap="round"/>
            </svg>
          </div>
          {#if dueFlashcardsCount > 0}
            <span class="badge badge-due">{dueFlashcardsCount} PENDENTES</span>
          {:else}
            <span class="badge">REVISÃO ESPAÇADA</span>
          {/if}
        </div>

        <div class="card-content">
          <h2>Flashcards & Memória</h2>
          <p>Revisão diária por repetição espaçada das suas palavras salvas.</p>
        </div>

        <div class="card-footer">
          <span>Estudar Flashcards</span>
          <svg class="arrow-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
      </div>
    </div>
  </main>
</div>

<style>
  :global(body) {
    background-color: #0a0a0a;
    margin: 0;
    font-family: 'Inter', system-ui, -apple-system, sans-serif;
  }

  .dashboard-wrapper {
    min-height: 100vh;
    background-color: #0a0a0a;
    color: #e9e9e9;
    display: flex;
    flex-direction: column;
  }

  /* Top Navbar */
  .top-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 2.5rem;
    border-bottom: 1px solid #1a1a1a;
    background-color: #0f0f0f;
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
    border-color: #ef4444;
    color: #ef4444;
  }

  /* Main Content Area */
  .dashboard-content {
    flex: 1;
    max-width: 768px;
    width: 100%;
    margin: 0 auto;
    padding: 3rem 2rem;
    box-sizing: border-box;
  }

  .hub-header {
    margin-bottom: 2.5rem;
    text-align: left;
  }

  .level-pill {
    display: inline-block;
    padding: 0.2rem 0.6rem;
    background-color: rgba(92, 109, 255, 0.1);
    color: #5c6dff;
    font-size: 0.75rem;
    font-weight: 600;
    border-radius: 4px;
    margin-bottom: 0.75rem;
    border: 1px solid rgba(92, 109, 255, 0.25);
  }

  .hub-header h1 {
    font-size: 1.8rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 0.25rem 0;
    letter-spacing: -0.02em;
  }

  .hub-header p {
    color: #777777;
    font-size: 0.95rem;
    margin: 0;
  }

  /* Hub Grid */
  .hub-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.25rem;
  }

  .hub-card {
    background-color: #121212;
    border: 1px solid #222222;
    border-radius: 12px;
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .hub-card:hover {
    border-color: #333333;
    background-color: #161616;
  }

  .card-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.25rem;
  }

  .icon-box {
    width: 40px;
    height: 40px;
    background-color: #181818;
    border: 1px solid #262626;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .badge {
    font-size: 0.68rem;
    font-weight: 600;
    letter-spacing: 0.06em;
    color: #666666;
    background-color: #181818;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
  }

  .badge-due {
    color: #a855f7;
    background-color: rgba(168, 85, 247, 0.12);
    border: 1px solid rgba(168, 85, 247, 0.25);
  }

  .card-content h2 {
    font-size: 1.2rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 0.4rem 0;
  }

  .card-content p {
    color: #888888;
    font-size: 0.88rem;
    line-height: 1.5;
    margin: 0 0 1.5rem 0;
  }

  .card-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #777777;
    font-size: 0.85rem;
    font-weight: 500;
    transition: color 0.2s ease;
    padding-top: 0.5rem;
    border-top: 1px solid #1a1a1a;
  }

  .hub-card:hover .card-footer {
    color: #ffffff;
  }

  .arrow-icon {
    transition: transform 0.2s ease, color 0.2s ease;
  }

  .hub-card:hover .arrow-icon {
    transform: translateX(4px);
    color: #5c6dff;
  }

  @media (max-width: 768px) {
    .top-nav {
      padding: 1rem 1.25rem;
    }

    .dashboard-content {
      padding: 2rem 1.25rem;
    }
  }
</style>
