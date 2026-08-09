<script lang="ts">
  import { onMount } from "svelte";
  import Cookies from "js-cookie";
  import { goto } from "$app/navigation";

  interface Flashcard {
    id: string;
    front: string;
    back: string;
    context_sentence: string;
    explanation: string;
    ease_factor: number;
    interval: number;
    repetitions: number;
    next_review: string;
    created_at: string;
    last_reviewed_at?: string;
  }

  let flashcards: Flashcard[] = [];
  let dueCards: Flashcard[] = [];
  let loading = true;

  // Study mode state
  let currentCardIndex = 0;
  let isFlipped = false;

  // Add/Generate Modal state
  let showModal = false;
  let newTerm = "";
  let newContext = "";
  let newBack = "";
  let newExplanation = "";
  let isGenerating = false;
  let isSaving = false;

  const API_URL = "http://localhost:8082/api";

  function formatSingleWord(str: string): string {
    if (!str) return "";
    let cleaned = str.trim().replace(/^["'\s.,;:!?()]+|["'\s.,;:!?()]+$/g, "");
    if (!cleaned) return "";
    const parts = cleaned.split(/[\s\/,;]+/);
    let target = parts[0] || "";
    const lowerFirst = parts[0]?.toLowerCase() || "";
    if (parts.length > 1 && ["a", "o", "um", "uma", "to"].includes(lowerFirst)) {
      target = parts[1] || target;
    }
    target = target.replace(/^["'\s.,;:!?()]+|["'\s.,;:!?()]+$/g, "");
    if (!target) return "";
    return target.charAt(0).toUpperCase() + target.slice(1).toLowerCase();
  }

  onMount(async () => {
    await fetchFlashcards();
  });

  async function fetchFlashcards() {
    loading = true;
    try {
      const res = await fetch(`${API_URL}/flashcards`, {
        headers: { "X-User-ID": "default_user" }
      });
      if (res.ok) {
        const raw: Flashcard[] = await res.json();
        flashcards = raw.map(c => ({ ...c, back: formatSingleWord(c.back) }));
        const now = new Date();
        dueCards = flashcards.filter(c => new Date(c.next_review) <= now);
      }
    } catch (e) {
      console.error("Erro ao carregar flashcards:", e);
    } finally {
      loading = false;
    }
  }

  function flipCard() {
    isFlipped = !isFlipped;
  }

  function speakText(text: string) {
    if (!("speechSynthesis" in window)) return;
    window.speechSynthesis.cancel();
    const utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = "en-US";
    utterance.rate = 0.9;
    window.speechSynthesis.speak(utterance);
  }

  async function handleReview(rating: number) {
    if (dueCards.length === 0) return;
    const card = dueCards[currentCardIndex];
    try {
      const res = await fetch(`${API_URL}/flashcards/${card.id}/review`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ rating })
      });

      if (res.ok) {
        isFlipped = false;
        dueCards = dueCards.filter((_, idx) => idx !== currentCardIndex);
        if (currentCardIndex >= dueCards.length) {
          currentCardIndex = 0;
        }
        await fetchFlashcards();
      }
    } catch (e) {
      console.error("Erro ao registrar revisão:", e);
    }
  }

  async function handleGenerateAI() {
    if (!newTerm.trim()) {
      alert("Digite primeiro uma palavra ou expressão em Inglês.");
      return;
    }
    isGenerating = true;
    try {
      const res = await fetch(`${API_URL}/flashcards/generate`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          term: newTerm,
          context_sentence: newContext
        })
      });
      const data = await res.json();
      if (res.ok && data) {
        newBack = formatSingleWord(data.back || newTerm);
        if (data.context_sentence && !newContext) newContext = data.context_sentence;
        if (data.explanation) newExplanation = data.explanation;
      } else {
        console.error("Erro na resposta da IA:", data);
        alert("Não foi possível gerar detalhes com IA no momento. Por favor, insira a tradução manualmente.");
      }
    } catch (e) {
      console.error("Erro ao gerar com IA:", e);
      alert("Erro ao se comunicar com o serviço de IA.");
    } finally {
      isGenerating = false;
    }
  }

  async function handleSaveFlashcard() {
    if (!newTerm.trim() || !newBack.trim()) return;

    const termClean = newTerm.trim().toLowerCase();
    if (flashcards.some(c => c.front.trim().toLowerCase() === termClean)) {
      alert(`Já existe um flashcard cadastrado para a palavra "${newTerm.trim()}".`);
      return;
    }

    const formattedBack = formatSingleWord(newBack);
    isSaving = true;
    try {
      const res = await fetch(`${API_URL}/flashcards`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-User-ID": "default_user"
        },
        body: JSON.stringify({
          front: newTerm.trim(),
          back: formattedBack,
          context_sentence: newContext,
          explanation: newExplanation
        })
      });

      if (res.ok) {
        showModal = false;
        newTerm = "";
        newContext = "";
        newBack = "";
        newExplanation = "";
        await fetchFlashcards();
      } else {
        const errData = await res.json();
        alert(errData.error || "Erro ao salvar flashcard.");
      }
    } catch (e) {
      console.error("Erro ao salvar cartão:", e);
    } finally {
      isSaving = false;
    }
  }

  function handleLogout() {
    Cookies.remove("authToken");
    goto("/");
  }
</script>

<div class="flashcards-wrapper">
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

  <main class="flashcards-content">
    <!-- Minimalist Header -->
    <header class="page-header">
      <div class="header-title">
        <h1>Flashcards & Memória</h1>
        <p class="subtitle">Repetição Espaçada (SuperMemo-2)</p>
      </div>

      <button class="btn-create" on:click={() => (showModal = true)}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M12 5v14m7-7H5" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <span>Criar Flashcard</span>
      </button>
    </header>

    <!-- Minimalist Stats Bar -->
    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">Pendentes hoje</span>
        <span class="stat-value highlight">{dueCards.length}</span>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <span class="stat-label">Total salvos</span>
        <span class="stat-value">{flashcards.length}</span>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <span class="stat-label">Dominados</span>
        <span class="stat-value">{flashcards.filter(c => c.repetitions >= 3).length}</span>
      </div>
    </div>

    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
      </div>
    {:else}
      <!-- Minimalist Study View -->
      <div class="study-area">
        {#if dueCards.length === 0}
          <div class="empty-card">
            <h2>Revisões concluídas!</h2>
            <p>Você revisou todos os flashcards de hoje. Volte mais tarde para novas surpresas!</p>
          </div>
        {:else}
          {@const card = dueCards[currentCardIndex]}
          <div class="progress-indicator">
            {currentCardIndex + 1} de {dueCards.length}
          </div>

          <!-- Minimalist 3D Card -->
          <div
            class="card-wrapper"
            class:flipped={isFlipped}
            on:click={flipCard}
            role="button"
            tabindex="0"
            on:keydown={(e) => e.key === "Enter" && flipCard()}
          >
            <div class="card-inner">
              <!-- FRONT -->
              <div class="card-face card-front">
                <div class="face-header">
                  <span class="tag">ENGLISH</span>
                  <button
                    class="btn-audio"
                    on:click|stopPropagation={() => speakText(card.front)}
                    title="Ouvir pronúncia"
                  >
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/>
                      <path d="M15.54 8.46a5 5 0 0 1 0 7.07"/>
                    </svg>
                  </button>
                </div>

                <div class="face-content">
                  <h2 class="main-word">{card.front}</h2>
                  {#if card.context_sentence}
                    <p class="sentence">"{card.context_sentence}"</p>
                  {/if}
                </div>

                <div class="face-footer">
                  <span>Clique para revelar a tradução</span>
                </div>
              </div>

              <!-- BACK -->
              <div class="card-face card-back">
                <div class="face-header">
                  <span class="tag tag-accent">PORTUGUÊS</span>
                  <button
                    class="btn-audio"
                    on:click|stopPropagation={() => speakText(card.front)}
                    title="Ouvir pronúncia"
                  >
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/>
                      <path d="M15.54 8.46a5 5 0 0 1 0 7.07"/>
                    </svg>
                  </button>
                </div>

                <div class="face-content">
                  <h2 class="translated-word">{card.back}</h2>
                  {#if card.explanation}
                    <p class="explanation">{card.explanation}</p>
                  {/if}
                </div>

                <div class="face-footer">
                  <span>Selecione sua facilidade de resposta abaixo</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Minimalist Rating Bar -->
          {#if isFlipped}
            <div class="rating-bar">
              <button class="rate-btn rate-again" on:click={() => handleReview(1)}>
                <span>Errei</span>
                <small>Hoje</small>
              </button>

              <button class="rate-btn rate-hard" on:click={() => handleReview(2)}>
                <span>Difícil</span>
                <small>1d</small>
              </button>

              <button class="rate-btn rate-good" on:click={() => handleReview(3)}>
                <span>Bom</span>
                <small>{Math.max(card.interval * 2, 3)}d</small>
              </button>

              <button class="rate-btn rate-easy" on:click={() => handleReview(4)}>
                <span>Fácil</span>
                <small>{Math.max(card.interval * 3, 6)}d</small>
              </button>
            </div>
          {/if}
        {/if}
      </div>
    {/if}
  </main>

  <!-- Minimalist SpeakUp Modal -->
  {#if showModal}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="modal-backdrop" on:click={() => (showModal = false)}>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="modal-box" on:click={(e) => e.stopPropagation()}>
        <div class="modal-header">
          <h3>Novo Flashcard</h3>
          <button class="close-btn" on:click={() => (showModal = false)}>✕</button>
        </div>

        <div class="modal-body">
          <div class="field">
            <label for="new-term">Palavra ou Expressão (Inglês)</label>
            <div class="input-ai-row">
              <input
                id="new-term"
                type="text"
                placeholder="Ex: Ubiquitous"
                bind:value={newTerm}
              />
              <button
                type="button"
                class="btn-ai"
                disabled={isGenerating || !newTerm.trim()}
                on:click={handleGenerateAI}
              >
                {isGenerating ? "Gerando..." : "🤖 IA"}
              </button>
            </div>
          </div>

          <div class="field">
            <label for="new-context">Frase de Exemplo (Opcional)</label>
            <input
              id="new-context"
              type="text"
              placeholder="Ex: Smartphones are ubiquitous in modern life."
              bind:value={newContext}
            />
          </div>

          <div class="field">
            <label for="new-back">Tradução (Português)</label>
            <input
              id="new-back"
              type="text"
              placeholder="Ex: Onipresente / Presente em toda parte"
              bind:value={newBack}
            />
          </div>

          <div class="field">
            <label for="new-exp">Explicação / Nota Gramatical (Opcional)</label>
            <textarea
              id="new-exp"
              rows="2"
              placeholder="Dica didática sobre o uso do termo..."
              bind:value={newExplanation}
            ></textarea>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-secondary" on:click={() => (showModal = false)}>
            Cancelar
          </button>
          <button
            class="btn-primary"
            disabled={isSaving || !newTerm.trim() || !newBack.trim()}
            on:click={handleSaveFlashcard}
          >
            {isSaving ? "Salvando..." : "Salvar Flashcard"}
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  :global(body) {
    background-color: #0a0a0a;
    margin: 0;
    font-family: 'Inter', system-ui, -apple-system, sans-serif;
  }

  .flashcards-wrapper {
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
    border-color: rgba(239, 68, 68, 0.4);
    color: #ef4444;
  }

  /* Flashcards Content Area */
  .flashcards-content {
    flex: 1;
    max-width: 720px;
    width: 100%;
    margin: 0 auto;
    padding: 2.5rem 2rem;
    box-sizing: border-box;
  }

  /* Header */
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .header-title h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 0.25rem 0;
    letter-spacing: -0.02em;
  }

  .subtitle {
    color: #888888;
    margin: 0;
    font-size: 0.9rem;
  }

  .btn-create {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.6rem 1.1rem;
    background-color: #1a1a1a;
    color: #ffffff;
    border: 1px solid #333333;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-create:hover {
    background-color: #262626;
    border-color: #5c6dff;
    color: #ffffff;
  }

  /* Stats Bar */
  .stats-bar {
    display: flex;
    align-items: center;
    background-color: #121212;
    border: 1px solid #222222;
    border-radius: 8px;
    padding: 1rem 1.5rem;
    margin-bottom: 2rem;
  }

  .stat-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    flex: 1;
  }

  .stat-label {
    font-size: 0.8rem;
    color: #777777;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .stat-value {
    font-size: 1.25rem;
    font-weight: 600;
    color: #ffffff;
  }

  .stat-value.highlight {
    color: #5c6dff;
  }

  .stat-divider {
    width: 1px;
    height: 30px;
    background-color: #222222;
    margin: 0 1.5rem;
  }

  /* Segmented Control Tabs */
  .segment-control {
    display: inline-flex;
    background-color: #121212;
    border: 1px solid #222222;
    padding: 3px;
    border-radius: 8px;
    margin-bottom: 2rem;
  }

  .segment-btn {
    padding: 0.5rem 1.25rem;
    background: transparent;
    border: none;
    color: #888888;
    font-size: 0.875rem;
    font-weight: 500;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .segment-btn.active {
    background-color: #1a1a1a;
    color: #ffffff;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
  }

  /* Study Area */
  .study-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
  }

  .progress-indicator {
    font-size: 0.85rem;
    color: #666666;
    margin-bottom: 1rem;
  }

  .empty-card {
    text-align: center;
    padding: 4rem 2rem;
    background-color: #121212;
    border: 1px solid #222222;
    border-radius: 12px;
    width: 100%;
    box-sizing: border-box;
  }

  .empty-card h2 {
    font-size: 1.25rem;
    color: #ffffff;
    margin: 1rem 0 0.5rem 0;
  }

  .empty-card p {
    color: #888888;
    font-size: 0.9rem;
    margin-bottom: 1.5rem;
  }

  /* Minimalist 3D Card */
  .card-wrapper {
    width: 100%;
    height: 300px;
    perspective: 1000px;
    cursor: pointer;
    margin-bottom: 1.5rem;
  }

  .card-inner {
    position: relative;
    width: 100%;
    height: 100%;
    transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    transform-style: preserve-3d;
  }

  .card-wrapper.flipped .card-inner {
    transform: rotateY(180deg);
  }

  .card-face {
    position: absolute;
    inset: 0;
    -webkit-backface-visibility: hidden;
    backface-visibility: hidden;
    background-color: #141414;
    border: 1px solid #262626;
    border-radius: 12px;
    padding: 1.75rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    box-sizing: border-box;
    transition: border-color 0.2s ease;
  }

  .card-wrapper:hover .card-face {
    border-color: #333333;
  }

  .card-back {
    transform: rotateY(180deg);
    border-color: rgba(92, 109, 255, 0.4);
  }

  .face-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .tag {
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    color: #666666;
    background-color: #1e1e1e;
    padding: 0.2rem 0.6rem;
    border-radius: 4px;
  }

  .tag-accent {
    color: #5c6dff;
    background-color: rgba(92, 109, 255, 0.1);
  }

  .btn-audio {
    background: transparent;
    border: none;
    color: #777777;
    cursor: pointer;
    padding: 0.4rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    transition: color 0.2s ease;
  }

  .btn-audio:hover {
    color: #ffffff;
  }

  .face-content {
    text-align: center;
  }

  .main-word {
    font-size: 2rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 0.5rem 0;
  }

  .translated-word {
    font-size: 1.75rem;
    font-weight: 700;
    color: #5c6dff;
    margin: 0 0 0.5rem 0;
  }

  .sentence {
    color: #aaaaaa;
    font-size: 0.95rem;
    font-style: italic;
    margin: 0;
  }

  .explanation {
    color: #888888;
    font-size: 0.85rem;
    line-height: 1.5;
    margin: 0;
  }

  .face-footer {
    text-align: center;
    font-size: 0.75rem;
    color: #555555;
  }

  /* Rating Bar */
  .rating-bar {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 0.75rem;
    width: 100%;
  }

  .rate-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0.6rem 0.5rem;
    background-color: #121212;
    border: 1px solid #222222;
    border-radius: 8px;
    color: #cccccc;
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .rate-btn small {
    font-size: 0.7rem;
    color: #666666;
    margin-top: 0.2rem;
  }

  .rate-again:hover {
    border-color: #ef4444;
    color: #ef4444;
  }

  .rate-hard:hover {
    border-color: #f97316;
    color: #f97316;
  }

  .rate-good:hover {
    border-color: #5c6dff;
    color: #5c6dff;
  }

  .rate-easy:hover {
    border-color: #22c55e;
    color: #22c55e;
  }

  /* Library Grid */
  .search-box {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background-color: #121212;
    border: 1px solid #222222;
    border-radius: 8px;
    padding: 0.6rem 1rem;
    margin-bottom: 1.5rem;
  }

  .search-box input {
    background: transparent;
    border: none;
    color: #ffffff;
    width: 100%;
    font-size: 0.9rem;
    outline: none;
  }

  .library-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 1rem;
  }

  .library-card {
    background-color: #121212;
    border: 1px solid #222222;
    border-radius: 8px;
    padding: 1.1rem;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .library-card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .front-text {
    font-weight: 700;
    font-size: 1.05rem;
    color: #ffffff;
  }

  .btn-delete {
    background: transparent;
    border: none;
    color: #555555;
    cursor: pointer;
    padding: 0.2rem;
  }

  .btn-delete:hover {
    color: #ef4444;
  }

  .back-text {
    font-weight: 600;
    font-size: 0.95rem;
    color: #5c6dff;
  }

  .context-text {
    font-size: 0.85rem;
    color: #888888;
    font-style: italic;
  }

  .card-meta {
    font-size: 0.75rem;
    color: #555555;
    margin-top: 0.5rem;
    padding-top: 0.5rem;
    border-top: 1px solid #1a1a1a;
  }

  /* Modal */
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.75);
    backdrop-filter: blur(4px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-box {
    background-color: #141414;
    border: 1px solid #262626;
    border-radius: 10px;
    width: 90%;
    max-width: 460px;
    padding: 1.5rem;
    color: #e9e9e9;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.25rem;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.1rem;
    color: #ffffff;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: #777777;
    cursor: pointer;
    font-size: 1rem;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    margin-bottom: 1rem;
  }

  .field label {
    font-size: 0.8rem;
    color: #888888;
  }

  .input-ai-row {
    display: flex;
    gap: 0.5rem;
  }

  .input-ai-row input,
  .field input,
  .field textarea {
    flex: 1;
    background-color: #0a0a0a;
    border: 1px solid #262626;
    border-radius: 6px;
    padding: 0.6rem 0.8rem;
    color: #ffffff;
    font-family: inherit;
    font-size: 0.9rem;
    outline: none;
  }

  .field input:focus,
  .field textarea:focus {
    border-color: #5c6dff;
  }

  .btn-ai {
    background-color: #1e1e1e;
    color: #5c6dff;
    border: 1px solid #333333;
    border-radius: 6px;
    padding: 0 0.9rem;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    white-space: nowrap;
    transition: all 0.2s ease;
  }

  .btn-ai:hover {
    background-color: rgba(92, 109, 255, 0.15);
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    margin-top: 1.5rem;
  }

  .btn-primary {
    padding: 0.6rem 1.2rem;
    background-color: #5c6dff;
    color: #ffffff;
    border: none;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.2s ease;
  }

  .btn-primary:hover {
    opacity: 0.9;
  }

  .btn-secondary {
    padding: 0.6rem 1.2rem;
    background-color: transparent;
    color: #888888;
    border: 1px solid #262626;
    border-radius: 6px;
    font-size: 0.9rem;
    cursor: pointer;
  }

  .btn-secondary:hover {
    color: #ffffff;
    border-color: #444444;
  }

  .loading-state {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 4rem;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 2px solid rgba(255, 255, 255, 0.1);
    border-top-color: #5c6dff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin { to { transform: rotate(360deg); } }

  @media (max-width: 768px) {
    .top-nav {
      padding: 1rem 1.25rem;
    }

    .flashcards-content {
      padding: 2rem 1.25rem;
    }
  }
</style>
