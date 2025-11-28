<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { getUserById, deleteUser, updateUser } from "../../utils/api";
    import { toast } from "svelte-sonner";
    import Cookies from "js-cookie";

    interface User {
        id: string;
        name: string;
        email: string;
        language: string;
        level: string;
    }

    let user: User | null = null;
    let loading = true;
    let editing = false;
    let showDeleteModal = false;
    let deleteConfirmation = "";

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
        } catch (error) {
            console.error("Erro ao carregar dados do usuário:", error);
            toast.error("Erro ao carregar perfil");
        } finally {
            loading = false;
        }
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

    onMount(() => {
        loadUserData();
    });
</script>

<div class="container">
    <div class="header">
        <button class="back-btn" on:click={() => goto("/chat")}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <path
                    d="M19 12H5M5 12l7 7M5 12l7-7"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
            </svg>
            Voltar
        </button>
        <h1>Meu Perfil</h1>
    </div>

    {#if loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Carregando perfil...</p>
        </div>
    {:else if user}
        <div class="profile-card">
            <div class="profile-header">
                <div class="avatar">
                    <svg width="60" height="60" viewBox="0 0 24 24" fill="none">
                        <circle
                            cx="12"
                            cy="8"
                            r="4"
                            stroke="currentColor"
                            stroke-width="2"
                        />
                        <path
                            d="M5 20c0-4 3-7 7-7s7 3 7 7"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                        />
                    </svg>
                </div>
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
                                ? "Inglês"
                                : "Japonês"}</span
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
                            <option value="english">Inglês</option>
                            <option value="japanese">Japonês</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="level">Nível</label>
                        <select id="level" bind:value={editForm.level} required>
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
        background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
        padding: 2rem;
        color: white;
    }

    .header {
        max-width: 800px;
        margin: 0 auto 2rem;
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .back-btn {
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

    .back-btn:hover {
        background: rgba(255, 255, 255, 0.1);
        transform: translateX(-4px);
    }

    h1 {
        font-size: 2rem;
        font-weight: 700;
        background: linear-gradient(135deg, #ffffff, #a0a0a0);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
    }

    .loading {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 400px;
        gap: 1rem;
    }

    .spinner {
        width: 50px;
        height: 50px;
        border: 3px solid rgba(255, 255, 255, 0.1);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .profile-card {
        max-width: 800px;
        margin: 0 auto;
        background: linear-gradient(145deg, #1e1e1e, #1a1a1a);
        border-radius: 20px;
        padding: 2.5rem;
        border: 1px solid rgba(255, 255, 255, 0.1);
        box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
    }

    .profile-header {
        text-align: center;
        padding-bottom: 2rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        margin-bottom: 2rem;
    }

    .avatar {
        width: 100px;
        height: 100px;
        margin: 0 auto 1rem;
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.15),
            rgba(255, 255, 255, 0.05)
        );
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 2px solid rgba(255, 255, 255, 0.2);
    }

    .avatar svg {
        color: rgba(255, 255, 255, 0.8);
    }

    h2 {
        font-size: 1.8rem;
        margin: 0 0 0.5rem 0;
        font-weight: 600;
    }

    .email {
        color: rgba(255, 255, 255, 0.6);
        font-size: 1rem;
    }

    .profile-info {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        margin-bottom: 2rem;
    }

    .info-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        background: rgba(255, 255, 255, 0.05);
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.1);
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

    .actions {
        display: flex;
        gap: 1rem;
    }

    .btn {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 1rem;
        border: none;
        border-radius: 12px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.3s ease;
    }

    .btn-primary {
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.15),
            rgba(255, 255, 255, 0.05)
        );
        color: white;
        border: 1px solid rgba(255, 255, 255, 0.2);
    }

    .btn-primary:hover {
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.25),
            rgba(255, 255, 255, 0.15)
        );
        transform: translateY(-2px);
        box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
    }

    .btn-danger {
        background: linear-gradient(
            135deg,
            rgba(255, 59, 48, 0.2),
            rgba(255, 59, 48, 0.1)
        );
        color: #ff3b30;
        border: 1px solid rgba(255, 59, 48, 0.3);
    }

    .btn-danger:hover {
        background: linear-gradient(
            135deg,
            rgba(255, 59, 48, 0.3),
            rgba(255, 59, 48, 0.2)
        );
        transform: translateY(-2px);
        box-shadow: 0 8px 20px rgba(255, 59, 48, 0.2);
    }

    .btn-secondary {
        background: rgba(255, 255, 255, 0.05);
        color: rgba(255, 255, 255, 0.8);
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .btn-secondary:hover {
        background: rgba(255, 255, 255, 0.1);
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

    .form-group input:focus,
    .form-group select:focus {
        outline: none;
        border-color: rgba(255, 255, 255, 0.3);
        background: rgba(255, 255, 255, 0.08);
    }

    .form-actions {
        display: flex;
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
        animation: fadeIn 0.2s ease;
    }

    .modal {
        background: linear-gradient(145deg, #1e1e1e, #1a1a1a);
        border-radius: 16px;
        max-width: 500px;
        width: 90%;
        border: 1px solid rgba(255, 255, 255, 0.1);
        box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
        animation: slideUp 0.3s ease;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    @keyframes slideUp {
        from {
            transform: translateY(20px);
            opacity: 0;
        }
        to {
            transform: translateY(0);
            opacity: 1;
        }
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
