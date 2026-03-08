<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import {
        getUserById,
        deleteUser,
        updateUser,
        getBillingStatus,
        createPortalSession,
    } from "../../utils/api";
    import { toast } from "svelte-sonner";
    import Cookies from "js-cookie";

    interface User {
        id: string;
        name: string;
        email: string;
        language: string;
        level: string;
    }

    interface BillingStatus {
        stripe_customer_id: string;
        stripe_subscription_id: string;
        stripe_price_id: string;
        stripe_status: string;
        stripe_current_period_end: number;
    }

    let user: User | null = null;
    let loading = true;
    let editing = false;
    let showDeleteModal = false;
    let deleteConfirmation = "";

    let billingStatus: BillingStatus | null = null;
    let billingLoading = false;
    let billingActionLoading = false;

    let editForm = {
        name: "",
        language: "",
        level: "",
    };

    let levels: string[] = [];

    function parseJwt(token: string) {
        try {
            const base64Url = token.split(".")[1];
            const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
            const jsonPayload = decodeURIComponent(
                atob(base64)
                    .split("")
                    .map(
                        (c) =>
                            "%" +
                            ("00" + c.charCodeAt(0).toString(16)).slice(-2),
                    )
                    .join(""),
            );
            return JSON.parse(jsonPayload);
        } catch (error) {
            console.error("Erro ao decodificar token:", error);
            return null;
        }
    }

    async function loadUserData() {
        try {
            const token = Cookies.get("authToken");
            if (!token) {
                toast.error("Você precisa estar logado");
                goto("/login");
                return;
            }

            const decoded = parseJwt(token);
            console.log("Token decodificado:", decoded); // Debug
            if (!decoded || !decoded.user_id) {
                toast.error("Token inválido");
                goto("/login");
                return;
            }

            user = await getUserById(decoded.user_id);
            console.log("Dados do usuário:", user); // Debug
            editForm = {
                name: user.name,
                language: user.language,
                level: user.level,
            };

            updateLevels(user.language);

            await loadBillingStatus();
        } catch (error) {
            console.error("Erro ao carregar dados do usuário:", error);
            toast.error("Erro ao carregar perfil");
        } finally {
            loading = false;
        }
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

    function setupVisibilityListener() {
        const handleVisibilityChange = async () => {
            if (!document.hidden) {
                // Página voltou ao foco, recarrega status de billing
                console.log("Página retornou ao foco, sincronizando dados...");
                await loadBillingStatus();
            }
        };

        document.addEventListener("visibilitychange", handleVisibilityChange);

        // Cleanup function
        return () => {
            document.removeEventListener(
                "visibilitychange",
                handleVisibilityChange,
            );
        };
    }

    function updateLevels(language: string) {
        if (language === "english") {
            levels = ["A1", "A2", "B1", "B2", "C1", "C2"];
        } else if (language === "japanese") {
            levels = ["N5", "N4", "N3", "N2", "N1"];
        } else {
            levels = [];
        }
    }

    function handleLanguageChange(e: Event) {
        const target = e.target as HTMLSelectElement;
        updateLevels(target.value);
        editForm.level = levels[0] || "";
    }

    async function handleUpdate() {
        if (!user) return;

        try {
            await updateUser(user.id, editForm);
            user = { ...user, ...editForm };
            editing = false;
            toast.success("Perfil atualizado com sucesso!");
        } catch (error) {
            console.error("Erro ao atualizar perfil:", error);
            toast.error("Erro ao atualizar perfil");
        }
    }

    async function handleDelete() {
        if (!user) return;
        if (deleteConfirmation !== "DELETAR") {
            toast.error('Digite "DELETAR" para confirmar');
            return;
        }

        try {
            await deleteUser(user.id);
            Cookies.remove("authToken");
            toast.success("Conta deletada com sucesso");
            setTimeout(() => {
                goto("/");
            }, 1500);
        } catch (error) {
            console.error("Erro ao deletar conta:", error);
            toast.error("Erro ao deletar conta");
        }
    }

    async function handleCancelPlan() {
        try {
            billingActionLoading = true;
            toast.loading("Abrindo portal de gerenciamento...");
            const { url } = await createPortalSession();
            if (url) {
                toast.promise(
                    new Promise((resolve) => {
                        setTimeout(() => {
                            window.location.href = url;
                            resolve(true);
                        }, 500);
                    }),
                    {
                        loading: "Redirecionando para Stripe...",
                        success:
                            "Você será redirecionado para gerenciar sua assinatura.",
                        error: "Erro ao abrir portal",
                    },
                );
            }
        } catch (error) {
            console.error("Erro ao abrir portal:", error);
            toast.error("Erro ao abrir portal de pagamento");
        } finally {
            billingActionLoading = false;
        }
    }

    function formatPeriodEnd(timestamp: number) {
        if (!timestamp) return "-";
        const date = new Date(timestamp * 1000);
        return date.toLocaleDateString("pt-BR");
    }

    function isActive(status?: string) {
        return status === "active" || status === "trialing";
    }

    function getPlanName(priceId?: string, status?: string) {
        if (status === "canceled") {
            return "❌ Assinatura Cancelada";
        }

        if (!priceId) return "Nenhum plano ativo";

        // Mapeamento dos price IDs para nomes de planos
        const planMap: { [key: string]: string } = {
            price_1RW5MgDQYqua1knA06NgjjST: "Plano Mensal",
            price_1StFqqDQYqua1knAQAXih6ES: "Plano Anual",
        };

        return planMap[priceId] || "Plano Personalizado";
    }

    function getStatusText(status?: string) {
        const statusMap: { [key: string]: string } = {
            active: "Ativa",
            trialing: "Período de teste",
            canceled: "Cancelada",
            incomplete: "Incompleta",
            incomplete_expired: "Expirada",
            past_due: "Pagamento pendente",
            unpaid: "Não paga",
        };

        return status ? statusMap[status] || status : "Sem assinatura";
    }

    onMount(() => {
        loadUserData();
        const unsubscribeVisibility = setupVisibilityListener();

        return () => {
            unsubscribeVisibility();
        };
    });
</script>

<div class="container">
    {#if loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Carregando perfil...</p>
        </div>
    {:else if user}
        <div class="profile-card">
            <div class="profile-header">
                <h2>{user.name}</h2>
                <p class="email">{user.email}</p>
            </div>

            {#if !editing}
                <div class="profile-info">
                    <div class="info-item">
                        <span class="label">
                            <svg
                                width="18"
                                height="18"
                                viewBox="0 0 24 24"
                                fill="none"
                            >
                                <path
                                    d="M3 21h18M3 10h18M3 7l9-4 9 4M4 10v11M20 10v11M8 10v11M12 10v11M16 10v11"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                            </svg>
                            Idioma
                        </span>
                        <span class="value"
                            >{user.language === "english"
                                ? "🇺🇸 Inglês"
                                : user.language === "japanese"
                                  ? "🇯🇵 Japonês"
                                  : user.language || "Não especificado"}</span
                        >
                    </div>

                    <div class="info-item">
                        <span class="label">
                            <svg
                                width="18"
                                height="18"
                                viewBox="0 0 24 24"
                                fill="none"
                            >
                                <path
                                    d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                            </svg>
                            Nível
                        </span>
                        <span class="value">{user.level}</span>
                    </div>
                </div>

                <div class="billing-card">
                    <div class="billing-header">
                        <h3>💳 Assinatura</h3>
                        {#if billingLoading}
                            <span class="billing-status">Carregando...</span>
                        {:else if billingStatus && billingStatus.stripe_status}
                            <span
                                class="billing-status {isActive(
                                    billingStatus.stripe_status,
                                )
                                    ? 'active'
                                    : 'inactive'}"
                                >{getStatusText(
                                    billingStatus.stripe_status,
                                )}</span
                            >
                        {:else}
                            <span class="billing-status inactive"
                                >Sem assinatura</span
                            >
                        {/if}
                    </div>

                    {#if billingStatus?.stripe_status === "canceled"}
                        <div class="cancellation-alert">
                            <svg
                                width="20"
                                height="20"
                                viewBox="0 0 24 24"
                                fill="none"
                            >
                                <circle
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    stroke-width="2"
                                />
                                <path
                                    d="M8 12h8M12 8v8"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                />
                            </svg>
                            <p>
                                Sua assinatura foi cancelada e o acesso será
                                removido na data de validade acima.
                            </p>
                        </div>
                    {/if}

                    <div class="billing-info">
                        <div class="info-item">
                            <span class="label">
                                <svg
                                    width="18"
                                    height="18"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                >
                                    <path
                                        d="M20 7H4a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2zM2 10h20"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                </svg>
                                Plano atual
                            </span>
                            <span
                                class="value plan-name {billingStatus?.stripe_status ===
                                'canceled'
                                    ? 'cancelled'
                                    : ''}"
                            >
                                {getPlanName(
                                    billingStatus?.stripe_price_id,
                                    billingStatus?.stripe_status,
                                )}
                            </span>
                        </div>

                        {#if billingStatus?.stripe_current_period_end && billingStatus?.stripe_status !== "canceled"}
                            <div class="info-item">
                                <span class="label">
                                    <svg
                                        width="18"
                                        height="18"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                    >
                                        <circle
                                            cx="12"
                                            cy="12"
                                            r="10"
                                            stroke="currentColor"
                                            stroke-width="2"
                                        />
                                        <path
                                            d="M12 6v6l4 2"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                        />
                                    </svg>
                                    Próxima renovação
                                </span>
                                <span class="value">
                                    {formatPeriodEnd(
                                        billingStatus.stripe_current_period_end,
                                    )}
                                </span>
                            </div>
                        {:else if billingStatus?.stripe_current_period_end && billingStatus?.stripe_status === "canceled"}
                            <div class="info-item">
                                <span class="label">
                                    <svg
                                        width="18"
                                        height="18"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                    >
                                        <circle
                                            cx="12"
                                            cy="12"
                                            r="10"
                                            stroke="currentColor"
                                            stroke-width="2"
                                        />
                                        <path
                                            d="M12 6v6l4 2"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                        />
                                    </svg>
                                    Acesso até
                                </span>
                                <span class="value" style="color: #ff9500;">
                                    {formatPeriodEnd(
                                        billingStatus.stripe_current_period_end,
                                    )}
                                </span>
                            </div>
                        {/if}

                        {#if billingStatus?.stripe_subscription_id && billingStatus?.stripe_status !== "canceled"}
                            <div class="info-item">
                                <span class="label">
                                    <svg
                                        width="18"
                                        height="18"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                    >
                                        <path
                                            d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                        <path
                                            d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                    </svg>
                                    ID da Assinatura
                                </span>
                                <span class="value subscription-id">
                                    {billingStatus.stripe_subscription_id.substring(
                                        0,
                                        20,
                                    )}...
                                </span>
                            </div>
                        {/if}
                    </div>

                    <div class="billing-actions">
                        {#if isActive(billingStatus?.stripe_status)}
                            <button
                                class="btn btn-secondary"
                                on:click={handleCancelPlan}
                                disabled={billingActionLoading}
                            >
                                <svg
                                    width="18"
                                    height="18"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                >
                                    <path
                                        d="M9 11l3 3L22 4"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                    <path
                                        d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                </svg>
                                {billingActionLoading
                                    ? "Carregando..."
                                    : "Gerenciar assinatura"}
                            </button>
                        {:else}
                            <button
                                class="btn btn-primary"
                                on:click={() => goto("/planos")}
                            >
                                <svg
                                    width="18"
                                    height="18"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                >
                                    <path
                                        d="M12 5v14M5 12h14"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                    />
                                </svg>
                                Ver planos disponíveis
                            </button>
                        {/if}
                    </div>
                </div>

                <div class="actions">
                    <button
                        class="btn btn-primary"
                        on:click={() => (editing = true)}
                    >
                        <svg
                            width="18"
                            height="18"
                            viewBox="0 0 24 24"
                            fill="none"
                        >
                            <path
                                d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                            <path
                                d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                        </svg>
                        Editar Perfil
                    </button>
                    <button
                        class="btn btn-danger"
                        on:click={() => (showDeleteModal = true)}
                    >
                        <svg
                            width="18"
                            height="18"
                            viewBox="0 0 24 24"
                            fill="none"
                        >
                            <path
                                d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                        </svg>
                        Deletar Conta
                    </button>
                </div>
            {:else}
                <form class="edit-form" on:submit|preventDefault={handleUpdate}>
                    <div class="form-group">
                        <label for="name">Nome</label>
                        <input
                            id="name"
                            type="text"
                            bind:value={editForm.name}
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="language">Idioma</label>
                        <select
                            id="language"
                            bind:value={editForm.language}
                            on:change={handleLanguageChange}
                            required
                        >
                            <option value="" disabled>Escolha um idioma</option>
                            <option value="english">🇺🇸 Inglês</option>
                            <option value="japanese">🇯🇵 Japonês</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="level">Nível</label>
                        <select id="level" bind:value={editForm.level} required>
                            <option value="" disabled>Escolha seu nível</option>
                            {#each levels as level}
                                <option value={level}>{level}</option>
                            {/each}
                        </select>
                    </div>

                    <div class="form-actions">
                        <button type="submit" class="btn btn-primary"
                            >Salvar</button
                        >
                        <button
                            type="button"
                            class="btn btn-secondary"
                            on:click={() => (editing = false)}
                        >
                            Cancelar
                        </button>
                    </div>
                </form>
            {/if}
        </div>
    {/if}
</div>

{#if showDeleteModal}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="modal-overlay" on:click={() => (showDeleteModal = false)}>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="modal" on:click={(e) => e.stopPropagation()}>
            <div class="modal-header">
                <h3>⚠️ Deletar Conta</h3>
            </div>
            <div class="modal-body">
                <p>
                    Esta ação é <strong>irreversível</strong>. Todos os seus
                    dados serão perdidos permanentemente.
                </p>
                <p>Digite <strong>DELETAR</strong> para confirmar:</p>
                <input
                    type="text"
                    bind:value={deleteConfirmation}
                    placeholder="DELETAR"
                    class="confirmation-input"
                />
            </div>
            <div class="modal-actions">
                <button class="btn btn-danger" on:click={handleDelete}
                    >Confirmar Exclusão</button
                >
                <button
                    class="btn btn-secondary"
                    on:click={() => (showDeleteModal = false)}
                >
                    Cancelar
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    .container {
        min-height: 100vh;
        background-color: #0a0a0a;
        color: white;
        position: relative;
        overflow: hidden;
    }

    .container::before {
        content: "";
        position: fixed;
        top: -50%;
        right: -50%;
        width: 100%;
        height: 100%;
        background: radial-gradient(
            circle,
            rgba(255, 255, 255, 0.03) 0%,
            transparent 70%
        );
        pointer-events: none;
        z-index: -1;
    }

    .header {
        max-width: 1400px;
        margin: 2rem auto 2rem;
        display: flex;
        align-items: center;
    }

    h1 {
        font-size: 2.5rem;
        font-weight: 700;
        color: #ffffff;
        margin: 0;
        letter-spacing: -0.5px;
    }

    @media (max-width: 768px) {
        h1 {
            font-size: 1.8rem;
        }
    }

    .loading {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1rem;
    }

    .spinner {
        width: 50px;
        height: 50px;
        border: 3px solid rgba(255, 255, 255, 0.1);
        border-top-color: white;
        border-radius: 50%;
    }

    .profile-card {
        max-width: 1400px;
        min-height: 100vh;
        margin: 0 auto;
        background: linear-gradient(
            145deg,
            rgba(30, 30, 30, 0.8),
            rgba(26, 26, 26, 0.9)
        );
        padding: 0 2.5rem;
        backdrop-filter: blur(20px);
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    .profile-header {
        text-align: center;
        padding: 2rem 0;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        margin-bottom: 2rem;
    }

    h2 {
        font-size: 2rem;
        margin: 0 0 0.5rem 0;
        font-weight: 700;
        color: #fff;
        letter-spacing: -0.5px;
    }

    .email {
        color: rgba(255, 255, 255, 0.65);
        font-size: 1rem;
        margin: 0;
    }

    .profile-info {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 1.5rem;
        margin-bottom: 2rem;
    }

    .billing-card {
        margin-bottom: 2rem;
        padding: 1.75rem;
        border-radius: 18px;
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(255, 255, 255, 0.1);
        transition: all 0.3s ease;
    }

    .billing-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .billing-header h3 {
        margin: 0;
        font-size: 1.2rem;
    }

    .billing-status {
        padding: 0.35rem 0.75rem;
        border-radius: 999px;
        font-size: 0.85rem;
        font-weight: 600;
        text-transform: capitalize;
        background: rgba(255, 255, 255, 0.1);
        color: rgba(255, 255, 255, 0.8);
    }

    .billing-status.active {
        background: rgba(255, 255, 255, 0.15);
        color: #ffffff;
        border: 1px solid rgba(255, 255, 255, 0.3);
    }

    .billing-status.inactive {
        background: rgba(255, 255, 255, 0.1);
        color: rgba(255, 255, 255, 0.7);
        border: 1px solid rgba(255, 255, 255, 0.2);
    }

    .cancellation-alert {
        display: flex;
        gap: 1rem;
        padding: 1rem;
        margin: 0 0 1rem 0;
        border-radius: 12px;
        background: rgba(255, 59, 48, 0.1);
        border: 1px solid rgba(255, 59, 48, 0.3);
        color: #ff3b30;
    }

    .cancellation-alert svg {
        flex-shrink: 0;
        margin-top: 0.2rem;
    }

    .cancellation-alert p {
        margin: 0;
        font-size: 0.95rem;
        line-height: 1.5;
        color: #ff3b30;
    }

    .billing-info {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 1rem;
        margin-bottom: 1.25rem;
    }

    .billing-actions {
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
        justify-content: flex-start;
    }

    .billing-actions .btn {
        flex: 0 1 auto;
        min-width: 200px;
    }

    .billing-actions .btn:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .info-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.25rem;
        background: rgba(255, 255, 255, 0.03);
        border-radius: 14px;
        border: 1px solid rgba(255, 255, 255, 0.08);
        transition: all 0.3s ease;
    }

    .label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: rgba(255, 255, 255, 0.7);
        font-weight: 500;
    }

    .value {
        font-weight: 600;
        color: white;
    }

    .plan-name {
        color: #ffffff;
        font-weight: 700;
    }

    .plan-name.cancelled {
        background: none;
        -webkit-text-fill-color: unset;
        color: #ff3b30;
        font-weight: 600;
    }

    .subscription-id {
        font-family: "Courier New", monospace;
        font-size: 0.85rem;
        color: rgba(255, 255, 255, 0.7);
    }

    .actions {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 1.5rem;
    }

    .btn {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 1rem 1.5rem;
        border: none;
        border-radius: 14px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;
        white-space: nowrap;
    }

    .btn::before {
        content: "";
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: rgba(255, 255, 255, 0.1);
        transition: left 0.3s ease;
        z-index: -1;
    }

    .btn:hover::before {
        left: 100%;
    }

    .btn-primary {
        background: rgba(255, 255, 255, 0.1);
        color: white;
        border: 1px solid rgba(255, 255, 255, 0.2);
    }

    .btn-primary:hover {
        background: rgba(255, 255, 255, 0.15);
        border-color: rgba(255, 255, 255, 0.3);
        transform: translateY(-3px);
    }

    .btn-primary:active {
        transform: translateY(-1px);
    }

    .btn-danger {
        background: linear-gradient(
            135deg,
            rgba(255, 59, 48, 0.15),
            rgba(255, 59, 48, 0.08)
        );
        color: #ff6b5b;
        border: 1px solid rgba(255, 59, 48, 0.3);
    }

    .btn-danger:hover {
        background: linear-gradient(
            135deg,
            rgba(255, 59, 48, 0.25),
            rgba(255, 59, 48, 0.15)
        );
        border-color: rgba(255, 59, 48, 0.5);
        color: #ff3b30;
        transform: translateY(-3px);
    }

    .btn-danger:active {
        transform: translateY(-1px);
    }

    .btn-secondary {
        background: rgba(255, 255, 255, 0.05);
        color: rgba(255, 255, 255, 0.85);
        border: 1px solid rgba(255, 255, 255, 0.15);
    }

    .btn-secondary:hover {
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(255, 255, 255, 0.25);
        transform: translateY(-3px);
    }

    .btn-secondary:active {
        transform: translateY(-1px);
    }

    .edit-form {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .form-group label {
        color: rgba(255, 255, 255, 0.8);
        font-weight: 500;
        font-size: 0.95rem;
    }

    .form-group input,
    .form-group select {
        padding: 0.75rem;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        color: white;
        font-size: 1rem;
        transition: all 0.2s ease;
    }

    .form-group select {
        appearance: none;
        background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23ffffff' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
        background-repeat: no-repeat;
        background-position: right 0.75rem center;
        background-size: 12px;
        padding-right: 2.5rem;
        cursor: pointer;
    }

    .form-group select option {
        background: #1a1a1a;
        color: white;
        padding: 0.5rem;
    }

    .form-group select option:disabled {
        color: rgba(255, 255, 255, 0.5);
    }

    .form-group input:focus,
    .form-group select:focus {
        outline: none;
        border-color: rgba(255, 255, 255, 0.3);
        background: rgba(255, 255, 255, 0.08);
    }

    .form-actions {
        grid-column: 1 / -1;
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1rem;
        margin-top: 1rem;
    }

    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.8);
        backdrop-filter: blur(4px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }

    .modal {
        background: linear-gradient(145deg, #1e1e1e, #1a1a1a);
        border-radius: 16px;
        max-width: 500px;
        width: 90%;
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .modal-header {
        padding: 1.5rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    .modal-header h3 {
        margin: 0;
        font-size: 1.5rem;
        color: #ff3b30;
    }

    .modal-body {
        padding: 1.5rem;
    }

    .modal-body p {
        margin: 0 0 1rem 0;
        color: rgba(255, 255, 255, 0.8);
        line-height: 1.5;
    }

    .confirmation-input {
        width: 100%;
        padding: 0.75rem;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        color: white;
        font-size: 1rem;
        margin-top: 0.5rem;
    }

    .confirmation-input:focus {
        outline: none;
        border-color: rgba(255, 59, 48, 0.5);
    }

    .modal-actions {
        padding: 1.5rem;
        border-top: 1px solid rgba(255, 255, 255, 0.1);
        display: flex;
        gap: 1rem;
    }

    @media (max-width: 768px) {
        .container {
            padding: 1rem;
        }

        h1 {
            font-size: 1.5rem;
        }

        .profile-card {
            padding: 1.5rem;
        }

        .actions {
            flex-direction: column;
        }

        .form-actions {
            flex-direction: column;
        }

        .modal-actions {
            flex-direction: column-reverse;
        }
    }
</style>
