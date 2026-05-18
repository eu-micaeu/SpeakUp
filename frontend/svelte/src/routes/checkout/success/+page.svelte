<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    let countdown = 3;

    onMount(() => {
        const timer = setInterval(() => {
            countdown -= 1;
            if (countdown <= 0) {
                clearInterval(timer);
                goto("/chat");
            }
        }, 1000);

        return () => clearInterval(timer);
    });
</script>

<svelte:head>
    <title>Assinatura concluída | SpeakUp</title>
</svelte:head>

<section class="success-page">
    <div class="success-card">
        <div class="icon-wrapper">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        </div>
        <h1>Assinatura Ativada!</h1>
        <p>
            Parabéns! Você agora tem acesso ilimitado ao SpeakUp Pro. 
            Prepare-se para acelerar seu aprendizado.
        </p>
        
        <div class="redirect-info">
            <div class="spinner"></div>
            <span>Redirecionando para o chat em {countdown}s...</span>
        </div>

        <a href="/chat" class="manual-btn">Ir agora</a>
    </div>
</section>

<style>
    :global(body) {
        background-color: #0a0a0a;
        margin: 0;
        font-family: 'Inter', system-ui, sans-serif;
    }

    .success-page {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 2rem;
        background: radial-gradient(circle at center, rgba(92, 109, 255, 0.1) 0%, transparent 70%);
    }

    .success-card {
        max-width: 480px;
        width: 100%;
        padding: 3.5rem 2.5rem;
        border-radius: 32px;
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(92, 109, 255, 0.3);
        backdrop-filter: blur(20px);
        text-align: center;
        animation: cardAppear 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
    }

    @keyframes cardAppear {
        from { opacity: 0; transform: scale(0.9) translateY(20px); }
        to { opacity: 1; transform: scale(1) translateY(0); }
    }

    .icon-wrapper {
        width: 80px;
        height: 80px;
        background: rgba(126, 229, 179, 0.1);
        color: #7ee5b3;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 2rem;
        box-shadow: 0 0 30px rgba(126, 229, 179, 0.2);
    }

    h1 {
        font-size: 2.2rem;
        font-weight: 800;
        color: #fff;
        margin: 0 0 1rem;
        letter-spacing: -1px;
    }

    p {
        color: #b5b5b5;
        font-size: 1.1rem;
        line-height: 1.6;
        margin: 0 0 2.5rem;
    }

    .redirect-info {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.75rem;
        color: #888;
        font-size: 0.95rem;
        margin-bottom: 2rem;
    }

    .spinner {
        width: 18px;
        height: 18px;
        border: 2px solid rgba(255, 255, 255, 0.1);
        border-top-color: #5c6dff;
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin { to { transform: rotate(360deg); } }

    .manual-btn {
        display: inline-block;
        padding: 0.8rem 2rem;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        color: #fff;
        text-decoration: none;
        font-weight: 600;
        transition: all 0.2s;
    }

    .manual-btn:hover {
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(255, 255, 255, 0.2);
    }
</style>
