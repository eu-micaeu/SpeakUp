<script lang="ts">
    import { onMount } from "svelte";
    import {
        createCheckoutSession,
        createPortalSession,
        getAIUsageStatus,
        getBillingStatus,
    } from "../../utils/api";
    import { toast } from "svelte-sonner";

    interface BillingStatus {
        stripe_customer_id: string;
        stripe_subscription_id: string;
        stripe_price_id: string;
        stripe_status: string;
        stripe_current_period_end: number;
    }

    let billingStatus: BillingStatus | null = null;
    let billingLoading = false;
    let billingActionLoading = false;
    let usageLoading = false;
    let usageStatus: {
        is_pro: boolean;
        daily_limit: number;
        used_today: number;
        remaining: number;
    } | null = null;

    async function loadBillingStatus() {
        try {
            billingLoading = true;
            billingStatus = await getBillingStatus();
        } catch (error) {
            console.error("Erro ao carregar billing:", error);
        } finally {
            billingLoading = false;
        }
    }

    async function loadUsageStatus() {
        try {
            usageLoading = true;
            usageStatus = await getAIUsageStatus();
        } catch (error) {
            console.error("Erro ao carregar uso:", error);
        } finally {
            usageLoading = false;
        }
    }

    async function handleSubscribe(plan: "monthly" | "annual") {
        try {
            billingActionLoading = true;
            const { url } = await createCheckoutSession(plan);
            if (url) {
                window.location.href = url;
            }
        } catch (error: unknown) {
            console.error("Erro ao iniciar assinatura:", error);
            const message =
                error instanceof Error
                    ? error.message
                    : "Erro ao iniciar assinatura";
            toast.error(message);
        } finally {
            billingActionLoading = false;
        }
    }

    async function handleManageSubscription() {
        try {
            billingActionLoading = true;
            const { url } = await createPortalSession();
            if (url) {
                window.location.href = url;
            }
        } catch (error) {
            console.error("Erro ao abrir portal:", error);
            toast.error("Erro ao abrir portal");
        } finally {
            billingActionLoading = false;
        }
    }

    function isActive(status?: string) {
        return status === "active" || status === "trialing";
    }

    onMount(() => {
        loadBillingStatus();
        loadUsageStatus();
    });
</script>

<svelte:head>
    <title>Planos e Preços | SpeakUp</title>
    <meta name="description" content="Escolha o melhor plano da SpeakUp e tenha acesso ilimitado à inteligência artificial para prática de idiomas, conversação e correções em tempo real." />
    <meta property="og:title" content="Planos e Preços | SpeakUp" />
    <meta property="og:description" content="Escolha o melhor plano da SpeakUp e tenha acesso ilimitado à inteligência artificial para prática de idiomas, conversação e correções em tempo real." />
    <meta name="twitter:title" content="Planos e Preços | SpeakUp" />
    <meta name="twitter:description" content="Escolha o melhor plano da SpeakUp e tenha acesso ilimitado à inteligência artificial para prática de idiomas, conversação e correções em tempo real." />
</svelte:head>

<section class="plans">
    <div class="plans-header">
        <a href="/chat" class="back-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
            Voltar para o Chat
        </a>
        <h1>Escolha seu Plano</h1>
        <p>Acelere seu aprendizado com acesso ilimitado à nossa inteligência artificial especializada em idiomas.</p>
    </div>

    <div class="plans-grid">
        <article class="plan-card free">
            <div class="highlight-pill neutral">Gratuito</div>
            <h2>Plano Free</h2>
            <p class="price">R$ 0 <span>/ sempre</span></p>
            <p class="discount warning">Até {usageStatus?.daily_limit || 10} interações/dia</p>
            <p class="subtitle">Ideal para começar sua jornada</p>
            <ul>
                <li>Chat com IA e correções</li>
                <li>Traduções instantâneas</li>
                <li>Histórico básico de chats</li>
                <li>Suporte da comunidade</li>
            </ul>
            <div class="plan-footer">
                <span class="tag">Ativo por padrão</span>
            </div>
        </article>

        <article class="plan-card">
            <h2>Plano Mensal</h2>
            <p class="price">R$ 7 <span>/ mês</span></p>
            <p class="subtitle">Liberdade total para praticar</p>
            <ul>
                <li>Interações de IA ilimitadas</li>
                <li>Traduções sem restrições</li>
                <li>Correções gramaticais infinitas</li>
                <li>Prioridade no processamento</li>
            </ul>
            <div class="plan-footer">
                <span class="tag">Cancele quando quiser</span>
            </div>
            <div class="plan-actions">
                <button
                    class="action-button primary-btn"
                    on:click={() => handleSubscribe("monthly")}
                    disabled={billingActionLoading || (billingStatus && isActive(billingStatus.stripe_status))}
                >
                    {billingStatus && isActive(billingStatus.stripe_status) ? 'Plano Atual' : 'Assinar Mensal'}
                </button>
            </div>
        </article>

        <article class="plan-card highlight">
            <div class="highlight-pill">Melhor Valor</div>
            <h2>Plano Anual</h2>
            <p class="price">R$ 70 <span>/ ano</span></p>
            <p class="discount">Economize R$ 14 por ano</p>
            <p class="subtitle">Compromisso com a fluência</p>
            <ul>
                <li>Tudo do plano Mensal</li>
                <li>2 meses grátis inclusos</li>
                <li>Acesso antecipado a novas IAs</li>
                <li>Badges exclusivos no perfil</li>
            </ul>
            <div class="plan-footer">
                <span class="tag">Pagamento único</span>
            </div>
            <div class="plan-actions">
                <button
                    class="action-button highlight-btn"
                    on:click={() => handleSubscribe("annual")}
                    disabled={billingActionLoading || (billingStatus && isActive(billingStatus.stripe_status))}
                >
                    {billingStatus && isActive(billingStatus.stripe_status) ? 'Plano Atual' : 'Assinar Anual'}
                </button>
            </div>
        </article>
    </div>

    <div class="billing-panel">
        <div class="status-info">
            {#if billingLoading}
                <p class="billing-note">Verificando sua conta...</p>
            {:else if billingStatus && isActive(billingStatus.stripe_status)}
                <div class="active-status">
                    <span class="status-icon active"></span>
                    <p class="billing-note active">Sua assinatura está ativa! Aproveite o SpeakUp Pro.</p>
                </div>
                <button
                    class="action-button secondary"
                    on:click={handleManageSubscription}
                    disabled={billingActionLoading}
                >
                    Configurações de Pagamento
                </button>
            {:else}
                <div class="active-status">
                    <span class="status-icon inactive"></span>
                    <p class="billing-note inactive">Você está usando a versão gratuita.</p>
                </div>
                
                {#if usageLoading}
                    <p class="usage-text">Carregando uso diário...</p>
                {:else if usageStatus && !usageStatus.is_pro}
                    <div class="usage-container">
                        <p class="usage-text">
                            Interações hoje: <strong>{usageStatus.used_today} / {usageStatus.daily_limit}</strong>
                        </p>
                        <div class="usage-bar" aria-hidden="true">
                            <span
                                style={`width: ${Math.min(
                                    100,
                                    (usageStatus.used_today /
                                        (usageStatus.daily_limit || 1)) *
                                        100,
                                )}%`}
                            ></span>
                        </div>
                    </div>
                {/if}
            {/if}
        </div>
    </div>
</section>

<style>
    :global(body) {
        background-color: #0a0a0a;
        margin: 0;
        font-family: 'Inter', system-ui, -apple-system, sans-serif;
    }

    .plans {
        min-height: 100vh;
        color: #e9e9e9;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 4rem 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .plans-header {
        text-align: center;
        margin-bottom: 4rem;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1.5rem;
        animation: slideDown 0.6s ease-out;
    }

    @keyframes slideDown {
        from { opacity: 0; transform: translateY(-20px); }
        to { opacity: 1; transform: translateY(0); }
    }

    .back-button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.6rem 1.2rem;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        color: #b5b5b5;
        text-decoration: none;
        font-size: 0.9rem;
        font-weight: 500;
        transition: all 0.3s ease;
    }

    .back-button:hover {
        background: rgba(255, 255, 255, 0.1);
        color: #fff;
        transform: translateX(-5px);
    }

    .plans-header h1 {
        font-size: 3rem;
        font-weight: 800;
        color: #ffffff;
        margin: 0;
        letter-spacing: -1px;
        background: linear-gradient(135deg, #fff 0%, #a5a9c8 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
    }

    .plans-header p {
        font-size: 1.2rem;
        color: #b5b5b5;
        margin: 0;
        max-width: 600px;
        line-height: 1.6;
    }

    .plans-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
        gap: 2rem;
        width: 100%;
        margin-bottom: 4rem;
    }

    .plan-card {
        position: relative;
        padding: 3rem 2rem;
        border-radius: 24px;
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(255, 255, 255, 0.08);
        display: flex;
        flex-direction: column;
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        animation: fadeInUp 0.6s ease-out both;
    }

    @keyframes fadeInUp {
        from { opacity: 0; transform: translateY(30px); }
        to { opacity: 1; transform: translateY(0); }
    }

    .plan-card:hover {
        transform: translateY(-10px);
        border-color: rgba(92, 109, 255, 0.4);
        background: rgba(255, 255, 255, 0.05);
    }

    .plan-card.highlight {
        background: linear-gradient(180deg, rgba(92, 109, 255, 0.1) 0%, rgba(255, 255, 255, 0.03) 100%);
        border: 1px solid rgba(92, 109, 255, 0.3);
    }

    .plan-card.highlight:hover {
        border-color: rgba(92, 109, 255, 0.6);
    }

    .highlight-pill {
        position: absolute;
        top: 1.5rem;
        right: 1.5rem;
        padding: 0.4rem 1rem;
        border-radius: 999px;
        background: #5c6dff;
        color: #fff;
        font-size: 0.7rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .highlight-pill.neutral {
        background: rgba(255, 255, 255, 0.1);
        color: #b5b5b5;
    }

    .plan-card h2 {
        font-size: 1.5rem;
        margin: 0 0 1.5rem 0;
    }

    .price {
        font-size: 2.5rem;
        font-weight: 800;
        margin: 0;
        color: #fff;
    }

    .price span {
        font-size: 1rem;
        color: #888;
        font-weight: 400;
    }

    .discount {
        font-size: 0.9rem;
        color: #7ee5b3;
        font-weight: 600;
        margin-top: 0.5rem;
    }

    .discount.warning {
        color: #f0c98b;
    }

    .subtitle {
        color: #888;
        margin: 1rem 0 2rem 0;
        font-size: 0.95rem;
    }

    .plan-card ul {
        list-style: none;
        padding: 0;
        margin: 0 0 2.5rem 0;
        display: flex;
        flex-direction: column;
        gap: 1rem;
        flex: 1;
    }

    .plan-card ul li {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        color: #d1d1d1;
        font-size: 0.95rem;
    }

    .plan-card ul li::before {
        content: "✓";
        color: #5c6dff;
        font-weight: 900;
    }

    .plan-footer {
        margin-bottom: 1.5rem;
    }

    .tag {
        font-size: 0.75rem;
        color: #888;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        font-weight: 600;
    }

    .action-button {
        width: 100%;
        padding: 1rem;
        border-radius: 16px;
        font-size: 1rem;
        font-weight: 700;
        cursor: pointer;
        transition: all 0.3s ease;
        border: none;
    }

    .primary-btn {
        background: rgba(255, 255, 255, 0.1);
        color: #fff;
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .primary-btn:hover:not(:disabled) {
        background: rgba(255, 255, 255, 0.15);
        transform: translateY(-2px);
    }

    .highlight-btn {
        background: #5c6dff;
        color: #fff;
    }

    .highlight-btn:hover:not(:disabled) {
        background: #4a59e6;
        box-shadow: 0 8px 24px rgba(92, 109, 255, 0.3);
        transform: translateY(-2px);
    }

    .action-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .billing-panel {
        width: 100%;
        background: rgba(255, 255, 255, 0.02);
        border: 1px solid rgba(255, 255, 255, 0.05);
        border-radius: 24px;
        padding: 2rem;
        animation: fadeInUp 0.7s ease-out 0.4s both;
    }

    .status-info {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1.5rem;
        width: 100%;
    }

    .active-status {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .status-icon {
        width: 10px;
        height: 10px;
        border-radius: 50%;
    }

    .status-icon.active { background: #7ee5b3; box-shadow: 0 0 12px #7ee5b3; }
    .status-icon.inactive { background: #888; }

    .billing-note { margin: 0; font-weight: 600; }
    .billing-note.active { color: #7ee5b3; }

    .usage-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        width: 100%;
        max-width: 400px;
    }

    .usage-text { margin: 0; color: #b5b5b5; font-size: 0.9rem; }
    .usage-text strong { color: #fff; }

    .usage-bar {
        width: 100%;
        height: 8px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 4px;
        overflow: hidden;
    }

    .usage-bar span {
        display: block;
        height: 100%;
        background: #5c6dff;
        transition: width 1s ease-out;
    }

    .action-button.secondary {
        width: auto;
        min-width: 250px;
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.1);
        color: #fff;
        font-size: 0.9rem;
        padding: 0.8rem 1.5rem;
    }

    .action-button.secondary:hover {
        background: rgba(255, 255, 255, 0.05);
    }

    @media (max-width: 768px) {
        .plans { padding: 2rem 1rem; }
        .plans-header h1 { font-size: 2.2rem; }
        .plans-grid { grid-template-columns: 1fr; }
        .plan-card { padding: 2.5rem 1.5rem; }
    }
</style>
