<svelte:head>
    <title>Assinatura cancelada | SpeakUp</title>
</svelte:head>

<section class="status-page cancel">
    <div class="card">
        <div class="icon">⚠️</div>
        <h1>Assinatura não concluída</h1>
        <p>
            A assinatura foi cancelada ou não pôde ser finalizada. Você pode
            tentar novamente.
        </p>
        <div class="actions">
            <a class="btn primary" href="/planos">Tentar novamente</a>
            <a class="btn ghost" href="/perfil">Ir para o perfil</a>
        </div>
    </div>
</section>

<style>
    .status-page {
        min-height: 100vh;
        display: grid;
        place-items: center;
        padding: 2rem;
        background: #141414;
        color: #e9e9e9;
    }

    .card {
        max-width: 520px;
        width: 100%;
        padding: 2.5rem;
        border-radius: 20px;
        background: #1c1c1c;
        border: 1px solid rgba(255, 255, 255, 0.08);
        text-align: center;
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.25);
    }

    .icon {
        font-size: 2.5rem;
        margin-bottom: 1rem;
    }

    h1 {
        margin: 0 0 0.75rem;
        font-size: 2rem;
    }

    p {
        margin: 0 0 2rem;
        color: #bdbdbd;
        line-height: 1.6;
    }

    .actions {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
        justify-content: center;
    }

    .btn {
        padding: 0.75rem 1.25rem;
        border-radius: 12px;
        text-decoration: none;
        font-weight: 600;
    }

    .btn.primary {
        background: #5c6dff;
        color: #fff;
    }

    .btn.ghost {
        border: 1px solid rgba(255, 255, 255, 0.15);
        color: #e9e9e9;
    }
</style>
