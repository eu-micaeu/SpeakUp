<script lang="ts">
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { isAuthTokenValid } from "../utils/cookies";

  let isAuthenticated = false;

  onMount(() => {
    isAuthenticated = isAuthTokenValid();
  });

  function handleGoHome() {
    if (isAuthenticated) {
      goto("/dashboard");
    } else {
      goto("/");
    }
  }

  function handleGoBack() {
    if (window.history.length > 1) {
      window.history.back();
    } else {
      handleGoHome();
    }
  }
</script>

<svelte:head>
  <title>{$page.status} - Página não encontrada | SpeakUp</title>
</svelte:head>

<div class="error-wrapper">
  <nav class="top-nav">
    <div class="nav-brand" on:click={handleGoHome} role="button" tabindex="0" on:keydown={(e) => e.key === "Enter" && handleGoHome()}>
      <img src="/logo.png" alt="SpeakUp Logo" width="32" />
      <span class="brand-name">SpeakUp</span>
    </div>
  </nav>

  <main class="error-content">
    <div class="error-card">
      <div class="glow-bg"></div>

      <span class="status-pill">
        Erro {$page.status || 404}
      </span>

      <h1 class="error-code">{$page.status || 404}</h1>

      <h2 class="error-title">
        {#if $page.status === 404}
          Página não encontrada
        {:else}
          Algo deu errado
        {/if}
      </h2>

      <p class="error-description">
        {#if $page.status === 404}
          Parece que esta página se perdeu na tradução ou foi movida para outro endereço.
        {:else}
          {$page.error?.message || "Ocorreu um erro inesperado ao carregar esta página."}
        {/if}
      </p>

    </div>
  </main>
</div>

<style>
  :global(body) {
    background-color: #0a0a0a;
    margin: 0;
    font-family: 'Inter', system-ui, -apple-system, sans-serif;
  }

  .error-wrapper {
    min-height: 100vh;
    background-color: #0a0a0a;
    color: #e9e9e9;
    display: flex;
    flex-direction: column;
  }

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

  .error-content {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem 1.5rem;
  }

  .error-card {
    position: relative;
    background-color: #121212;
    border: 1px solid #222222;
    border-radius: 16px;
    padding: 3.5rem 2.5rem;
    max-width: 480px;
    width: 100%;
    text-align: center;
    overflow: hidden;
  }

  .glow-bg {
    position: absolute;
    top: -60px;
    left: 50%;
    transform: translateX(-50%);
    width: 220px;
    height: 120px;
    pointer-events: none;
  }

  .status-pill {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    background-color: rgba(92, 109, 255, 0.1);
    color: #5c6dff;
    font-size: 0.75rem;
    font-weight: 600;
    border-radius: 999px;
    border: 1px solid rgba(92, 109, 255, 0.25);
    margin-bottom: 1rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .error-code {
    font-size: 5.5rem;
    font-weight: 800;
    line-height: 1;
    margin: 0 0 0.75rem 0;
    letter-spacing: -0.04em;
    background: linear-gradient(180deg, #ffffff 0%, #666666 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .error-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 0.75rem 0;
    letter-spacing: -0.02em;
  }

  .error-description {
    color: #888888;
    font-size: 0.95rem;
    line-height: 1.6;
    margin: 0 0 2.25rem 0;
  }

  .error-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: center;
  }

  .btn-secondary,
  .btn-primary {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    flex: 1;
  }

  .btn-secondary {
    background: transparent;
    border: 1px solid #2a2a2a;
    color: #cccccc;
  }

  .btn-secondary:hover {
    background-color: #1a1a1a;
    color: #ffffff;
    border-color: #444444;
  }

  .btn-primary {
    background-color: #5c6dff;
    border: 1px solid #5c6dff;
    color: #ffffff;
  }

  .btn-primary:hover {
    background-color: #4b5ce6;
    border-color: #4b5ce6;
  }

  @media (max-width: 520px) {
    .top-nav {
      padding: 1rem 1.25rem;
    }

    .error-card {
      padding: 2.5rem 1.5rem;
    }

    .error-code {
      font-size: 4.25rem;
    }

    .error-actions {
      flex-direction: column;
    }
  }
</style>
