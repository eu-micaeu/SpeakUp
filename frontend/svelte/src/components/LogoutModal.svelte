<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { goto } from "$app/navigation";
  import Cookies from "js-cookie";

  export let show = false;

  const dispatch = createEventDispatcher();

  function close() {
    show = false;
    dispatch("close");
  }

  function confirmLogout() {
    Cookies.remove("authToken");
    show = false;
    goto("/");
  }
</script>

{#if show}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="logout-modal-overlay" on:click={close}>
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="logout-modal-box" on:click|stopPropagation>
      <div class="logout-icon-wrapper">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M16 17l5-5-5-5M21 12H9" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
      <h3>Deseja realmente sair?</h3>
      <p>Você precisará fazer login novamente para acessar sua conta.</p>
      <div class="logout-modal-actions">
        <button type="button" class="btn-cancel" on:click={close}>
          Cancelar
        </button>
        <button type="button" class="btn-confirm" on:click={confirmLogout}>
          Sim, sair
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .logout-modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.75);
    backdrop-filter: blur(6px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    padding: 1.5rem;
    animation: fadeIn 0.15s ease-out;
  }

  .logout-modal-box {
    background: #141414;
    border: 1px solid #262626;
    border-radius: 12px;
    padding: 2rem;
    width: 100%;
    max-width: 380px;
    text-align: center;
    animation: scaleUp 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .logout-icon-wrapper {
    width: 48px;
    height: 48px;
    margin: 0 auto 1rem auto;
    border-radius: 10px;
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.25);
    color: #ef4444;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .logout-modal-box h3 {
    color: #ffffff;
    font-size: 1.2rem;
    font-weight: 700;
    margin: 0 0 0.5rem 0;
  }

  .logout-modal-box p {
    color: #888888;
    font-size: 0.9rem;
    line-height: 1.5;
    margin: 0 0 1.75rem 0;
  }

  .logout-modal-actions {
    display: flex;
    gap: 0.75rem;
  }

  .btn-cancel,
  .btn-confirm {
    flex: 1;
    padding: 0.7rem 1rem;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-cancel {
    background: transparent;
    border: 1px solid #2e2e2e;
    color: #cccccc;
  }

  .btn-cancel:hover {
    background: #1f1f1f;
    color: #ffffff;
    border-color: #444444;
  }

  .btn-confirm {
    background: #ef4444;
    border: 1px solid #ef4444;
    color: #ffffff;
  }

  .btn-confirm:hover {
    background: #dc2626;
    border-color: #dc2626;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes scaleUp {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }
</style>
