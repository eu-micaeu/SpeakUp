<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
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

    function handleBack() {
        if (typeof window !== "undefined" && window.history.length > 1) {
            window.history.back();
            return;
        }

        goto("/chat");
    }

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
        } catch (error) {
            console.error("Erro ao iniciar assinatura:", error);
            toast.error("Erro ao iniciar assinatura");
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
    <title>Planos | SpeakUp</title>
</svelte:head>

<section class="plans">
    <header class="plans-header">
        <button class="back-button" on:click={handleBack} aria-label="Voltar">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path
                    d="M15 18l-6-6 6-6"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
            </svg>
            <span>Voltar</span>
        </button>
        <h1>Planos do SpeakUp</h1>
        <p>
            Escolha o plano ideal para continuar evoluindo no idioma com prática
            diária, correções inteligentes e acompanhamento do progresso.
        </p>
    </header>

    <div class="plans-grid">
        <article class="plan-card free">
            <div class="highlight-pill neutral">Free</div>
            <h2>Plano Free</h2>
            <p class="price">R$ 0 <span>/ sempre</span></p>
            <p class="discount warning">Até 10 interações de IA/dia</p>
            <p class="subtitle">Comece sem pagar</p>
            <ul>
                <li>Chat com IA e correções básicas.</li>
                <li>Traduções e tópicos com limite diário.</li>
                <li>Histórico recente de conversas.</li>
            </ul>
            <div class="plan-footer">
                <span class="tag">Sem cartão</span>
            </div>
        </article>

        <article class="plan-card">
            <h2>Plano Mensal</h2>
            <p class="price">R$ 18 <span>/ mês</span></p>
            <p class="subtitle">Flexibilidade para começar agora</p>
            <ul>
                <li>Chat com IA ilimitado.</li>
                <li>Correções e traduções ilimitadas.</li>
                <li>Geração de tópicos sem limite.</li>
            </ul>
            <div class="plan-footer">
                <span class="tag">Cobrança mensal</span>
            </div>
            <div class="plan-actions">
                <button
                    class="action-button"
                    on:click={() => handleSubscribe("monthly")}
                    disabled={billingActionLoading}
                >
                    Assinar mensal
                </button>
            </div>
        </article>

        <article class="plan-card highlight">
            <div class="highlight-pill">Mais vantajoso</div>
            <h2>Plano Anual</h2>
            <p class="price">R$ 180 <span>/ ano</span></p>
            <p class="discount">Economize R$ 36 (2 meses grátis)</p>
            <p class="subtitle">Economia e consistência</p>
            <ul>
                <li>Chat com IA ilimitado.</li>
                <li>Correções e traduções ilimitadas.</li>
                <li>Geração de tópicos sem limite.</li>
            </ul>
            <div class="plan-footer">
                <span class="tag">Cobrança anual</span>
            </div>
            <div class="plan-actions">
                <button
                    class="action-button"
                    on:click={() => handleSubscribe("annual")}
                    disabled={billingActionLoading}
                >
                    Assinar anual
                </button>
            </div>
        </article>
    </div>

    <div class="billing-panel">
        {#if billingLoading}
            <p class="billing-note">Carregando status da assinatura...</p>
        {:else if billingStatus && isActive(billingStatus.stripe_status)}
            <p class="billing-note active">Assinatura ativa</p>
            <button
                class="action-button secondary"
                on:click={handleManageSubscription}
                disabled={billingActionLoading}
            >
                Gerenciar assinatura
            </button>
        {:else}
            <p class="billing-note inactive">
                Você ainda não possui assinatura ativa.
            </p>
            {#if usageLoading}
                <p class="usage-text">Carregando créditos...</p>
            {:else if usageStatus && !usageStatus.is_pro}
                <p class="usage-text">
                    Créditos hoje: {usageStatus.used_today} / {usageStatus.daily_limit}
                    · Restam {usageStatus.remaining}
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
            {/if}
        {/if}
    </div>

    <div class="helper">
        <p>
            Para gerenciar sua assinatura, abra o menu de configurações e siga
            para a área de cobrança quando estiver disponível.
        </p>
    </div>
</section>

<style>
    .plans {
        height: calc(100vh - 80px);
        padding: 2.5rem clamp(1.5rem, 4vw, 4rem) 4rem;
        color: #e9e9e9;
        background-color: #0a0a0a;
        position: relative;
        overflow: hidden;
    }

    .plans-header {
        max-width: 720px;
        margin-bottom: 2.5rem;
        display: grid;
        gap: 0.75rem;
        animation: slideDown 0.6s ease-out;
    }

    @keyframes slideDown {
        from {
            opacity: 0;
            transform: translateY(-20px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .back-button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.25rem;
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(92, 109, 255, 0.3);
        border-radius: 12px;
        color: white;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        font-size: 0.95rem;
        font-weight: 500;
        backdrop-filter: blur(10px);
        width: fit-content;
    }

    .back-button:hover {
        background: rgba(255, 255, 255, 0.08);
        border-color: rgba(255, 255, 255, 0.3);
        transform: translateX(-4px);
    }

    .back-button:active {
        transform: translateX(-2px);
    }

    .plans-header h1 {
        font-size: 2.5rem;
        font-weight: 700;
        color: #ffffff;
        margin: 0;
        letter-spacing: -0.5px;
    }

    .plans-header p {
        font-size: 1.1rem;
        color: #b5b5b5;
        margin: 0;
        line-height: 1.6;
    }

    @media (max-width: 768px) {
        .plans-header h1 {
            font-size: 1.8rem;
        }

        .plans-header p {
            font-size: 1rem;
        }
    }

    .plan-card {
        position: relative;
        padding: 2.5rem;
        border-radius: 20px;
        background: linear-gradient(
            145deg,
            rgba(28, 28, 28, 0.9),
            rgba(24, 24, 24, 0.8)
        );
        border: 1px solid rgba(255, 255, 255, 0.08);
        backdrop-filter: blur(10px);
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        display: flex;
        flex-direction: column;
        animation: fadeInUp 0.6s ease-out both;
        height: 80%;
    }

    @keyframes fadeInUp {
        from {
            opacity: 0;
            transform: translateY(30px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .plans-grid > .plan-card:nth-child(1) {
        animation-delay: 0.1s;
    }

    .plans-grid > .plan-card:nth-child(2) {
        animation-delay: 0.2s;
    }

    .plans-grid > .plan-card:nth-child(3) {
        animation-delay: 0.3s;
    }

    .plan-card:hover {
        transform: translateY(-8px);
    }

    .plan-card.highlight {
        background: linear-gradient(
            145deg,
            rgba(64, 86, 255, 0.15),
            rgba(24, 24, 24, 0.8)
        );
    }

    .plan-card.free {
        background: linear-gradient(
            145deg,
            rgba(24, 24, 24, 0.8),
            rgba(18, 18, 18, 0.8)
        );
        border-color: rgba(255, 255, 255, 0.12);
    }

    .plans-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
        gap: 2.5rem;
        margin-bottom: 3rem;
        align-items: start;
    }

    .plan-card h2 {
        margin: 0 0 1rem 0;
        font-size: 1.5rem;
        font-weight: 700;
        letter-spacing: -0.5px;
    }

    .price {
        margin: 0.6rem 0 0.2rem;
        font-size: 2rem;
        font-weight: 800;
        color: #f3f4ff;
        letter-spacing: -1px;
    }

    .price span {
        font-size: 0.95rem;
        font-weight: 500;
        color: #a5a9c8;
    }

    .discount {
        margin: -0.4rem 0 1rem;
        font-size: 0.95rem;
        font-weight: 700;
        color: #7ee5b3;
    }

    .discount.warning {
        color: #f0c98b;
    }

    .subtitle {
        margin: 0.35rem 0 1.5rem;
        color: #a4a4a4;
        font-size: 0.95rem;
    }

    .plan-card ul {
        margin: 0;
        padding-left: 1.5rem;
        display: grid;
        gap: 0.8rem;
        color: #d1d1d1;
        flex-grow: 1;
    }

    .plan-card ul li {
        line-height: 1.5;
    }

    .plan-footer {
        margin-top: 1.75rem;
    }

    .plan-actions {
        margin-top: auto;
        padding-top: 1.5rem;
    }

    .action-button {
        width: 100%;
        padding: 0.875rem 1.5rem;
        border-radius: 14px;
        border: 1px solid rgba(255, 255, 255, 0.2);
        background: rgba(255, 255, 255, 0.1);
        color: #fff;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;
        font-size: 1rem;
    }

    .action-button::before {
        content: "";
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: rgba(255, 255, 255, 0.2);
        transition: left 0.3s ease;
        z-index: 1;
    }

    .action-button:hover {
        transform: translateY(-3px);
    }

    .action-button:hover::before {
        left: 100%;
    }

    .action-button:active {
        transform: translateY(-1px);
    }

    .action-button:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        transform: none;
    }

    .action-button.secondary {
        width: auto;
        background: rgba(255, 255, 255, 0.08);
        color: #e9e9e9;
        border: 1px solid rgba(255, 255, 255, 0.15);
    }

    .action-button.secondary:hover {
        background: rgba(255, 255, 255, 0.12);
        border-color: rgba(255, 255, 255, 0.25);
    }

    .tag {
        display: inline-flex;
        padding: 0.4rem 0.9rem;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.1);
        color: #cccccc;
        font-size: 0.8rem;
        font-weight: 600;
        letter-spacing: 0.05em;
    }

    .highlight-pill {
        position: absolute;
        top: 1.5rem;
        right: 1.5rem;
        padding: 0.35rem 0.9rem;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.15);
        border: 1px solid rgba(255, 255, 255, 0.3);
        color: #fff;
        font-size: 0.75rem;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        font-weight: 700;
    }

    .highlight-pill.neutral {
        background: rgba(255, 255, 255, 0.12);
        color: #e9e9e9;
    }

    .billing-panel {
        margin-top: 3rem;
        padding: 2rem;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 1.5rem;
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(92, 109, 255, 0.25);
        border-radius: 18px;
        backdrop-filter: blur(10px);
        animation: fadeInUp 0.7s ease-out 0.4s both;
    }

    .billing-note {
        margin: 0;
        color: #b5b5b5;
        font-weight: 500;
    }

    .billing-note.active {
        color: #7ee5b3;
        font-weight: 700;
    }

    .billing-note.inactive {
        color: #f0a8a8;
    }

    .usage-text {
        margin: 0.5rem 0 0;
        color: #f0c98b;
        font-weight: 700;
        font-size: 0.95rem;
    }

    .usage-bar {
        width: 100%;
        max-width: 320px;
        height: 10px;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.1);
    }

    .usage-bar span {
        display: block;
        height: 100%;
        background: linear-gradient(90deg, #f0c98b, #f39b6d);
        border-radius: inherit;
        animation: progressFill 1.5s ease-out;
    }

    @keyframes progressFill {
        from {
            width: 0% !important;
        }
    }

    .helper {
        max-width: 720px;
        margin-top: 3rem;
        padding: 1.5rem;
        color: #9c9c9c;
        line-height: 1.7;
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(255, 255, 255, 0.08);
        border-radius: 14px;
        font-size: 0.95rem;
        animation: fadeInUp 0.7s ease-out 0.5s both;
    }

    @media (max-width: 768px) {
        .plans {
            padding: 1.5rem 1rem 2.5rem;
        }

        .plan-card {
            padding: 2rem;
        }

        .highlight-pill {
            position: static;
            margin-bottom: 0.75rem;
            display: inline-flex;
            width: fit-content;
        }

        .plan-card h2 {
            font-size: 1.3rem;
        }

        .price {
            font-size: 1.75rem;
        }

        .billing-panel {
            flex-direction: column;
            align-items: flex-start;
            gap: 1rem;
        }

        .action-button {
            padding: 0.75rem 1.25rem;
        }

        .action-button.secondary {
            width: 100%;
        }
    }

    @media (max-width: 480px) {
        .plans {
            padding: 1rem 0.5rem 2rem;
        }

        .plans-header {
            margin-bottom: 2rem;
        }

        .back-button {
            padding: 0.6rem 1rem;
            font-size: 0.9rem;
        }

        .plans-header h1 {
            font-size: 1.5rem;
        }

        .plans-header p {
            font-size: 0.9rem;
        }

        .plan-card {
            padding: 1.5rem;
            border-radius: 16px;
        }

        .plan-card h2 {
            font-size: 1.1rem;
        }

        .price {
            font-size: 1.5rem;
        }

        .subtitle {
            font-size: 0.9rem;
        }

        .plan-card ul {
            gap: 0.6rem;
            padding-left: 1.2rem;
        }

        .action-button {
            padding: 0.7rem 1rem;
            font-size: 0.9rem;
            border-radius: 12px;
        }

        .helper {
            padding: 1.25rem;
            font-size: 0.9rem;
        }
    }
</style>
