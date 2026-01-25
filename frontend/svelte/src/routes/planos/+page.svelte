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
        min-height: calc(100vh - 80px);
        padding: 2.5rem clamp(1.5rem, 4vw, 4rem) 4rem;
        color: #e9e9e9;
        background: #141414;
    }

    .plans-header {
        max-width: 720px;
        margin-bottom: 2.5rem;
        display: grid;
        gap: 0.75rem;
    }

    .back-button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        color: white;
        cursor: pointer;
        transition: all 0.2s ease;
        font-size: 0.95rem;
    }

    .back-button:hover {
        background: rgba(255, 255, 255, 0.1);
        transform: translateX(-4px);
    }

    .plan-card {
        position: relative;
        padding: 2rem;
        border-radius: 20px;
        background: #1c1c1c;
        border: 1px solid rgba(255, 255, 255, 0.08);
        margin:60px 0;
    }

    .plan-card h2 {
        margin: 0;
        font-size: 1.5rem;
    }

    .price {
        margin: 0.6rem 0 0.2rem;
        font-size: 1.9rem;
        font-weight: 700;
        color: #f3f4ff;
    }

    .price span {
        font-size: 0.95rem;
        font-weight: 500;
        color: #a5a9c8;
    }

    .discount {
        margin: -0.4rem 0 1rem;
        font-size: 0.95rem;
        font-weight: 600;
        color: #a9e6c3;
    }

    .discount.warning {
        color: #f0c98b;
    }

    .subtitle {
        margin: 0.35rem 0 1.5rem;
        color: #a4a4a4;
    }

    .plan-card ul {
        margin: 0;
        padding-left: 1.2rem;
        display: grid;
        gap: 0.65rem;
        color: #d1d1d1;
    }

    .plan-footer {
        margin-top: 1.75rem;
    }

    .plan-actions {
        margin-top: 1.5rem;
    }

    .action-button {
        width: 100%;
        padding: 0.75rem 1rem;
        border-radius: 12px;
        border: none;
        background: #5c6dff;
        color: #fff;
        font-weight: 600;
        cursor: pointer;
        transition:
            transform 0.2s ease,
            box-shadow 0.2s ease;
    }

    .action-button:hover {
        transform: translateY(-1px);
        box-shadow: 0 10px 20px rgba(92, 109, 255, 0.25);
    }

    .action-button:disabled {
        opacity: 0.7;
        cursor: not-allowed;
        box-shadow: none;
        transform: none;
    }

    .action-button.secondary {
        width: auto;
        background: rgba(255, 255, 255, 0.08);
        color: #e9e9e9;
        border: 1px solid rgba(255, 255, 255, 0.12);
    }

    .tag {
        display: inline-flex;
        padding: 0.35rem 0.8rem;
        border-radius: 999px;
        background: rgba(111, 145, 255, 0.15);
        color: #9fb1ff;
        font-size: 0.8rem;
    }

    .highlight {
        border-color: rgba(111, 145, 255, 0.5);
        background: linear-gradient(
            150deg,
            rgba(64, 86, 255, 0.18),
            #1c1c1c 60%
        );
    }

    .highlight-pill {
        position: absolute;
        top: 1.5rem;
        right: 1.5rem;
        padding: 0.25rem 0.7rem;
        border-radius: 999px;
        background: #5163ff;
        color: #fff;
        font-size: 0.7rem;
        letter-spacing: 0.08em;
        text-transform: uppercase;
    }

    .highlight-pill.neutral {
        background: rgba(255, 255, 255, 0.12);
        color: #e9e9e9;
    }

    .plan-card.free {
        background: #181818;
        border-color: rgba(255, 255, 255, 0.12);
    }

    .billing-panel {
        margin-top: 2.25rem;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 1rem;
    }

    .billing-note {
        margin: 0;
        color: #b5b5b5;
    }

    .billing-note.active {
        color: #9de0b2;
        font-weight: 600;
    }

    .billing-note.inactive {
        color: #f0b7b7;
    }

    .usage-text {
        margin: 0.25rem 0 0;
        color: #f0c98b;
        font-weight: 600;
    }

    .usage-bar {
        width: 100%;
        max-width: 320px;
        height: 8px;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.1);
        overflow: hidden;
    }

    .usage-bar span {
        display: block;
        height: 100%;
        background: linear-gradient(90deg, #f0c98b, #f39b6d);
        border-radius: inherit;
    }

    .helper {
        max-width: 720px;
        margin-top: 2.5rem;
        color: #9c9c9c;
        line-height: 1.6;
    }

    @media (max-width: 600px) {
        .plan-card {
            padding: 1.5rem;
        }

        .highlight-pill {
            position: static;
            margin-bottom: 0.75rem;
            display: inline-flex;
        }
    }
</style>
