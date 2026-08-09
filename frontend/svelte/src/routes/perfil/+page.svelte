<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import {
        getUserById,
        deleteUser,
        updateUser,
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
            if (!decoded || !decoded.user_id) {
                toast.error("Token inválido");
                goto("/login");
                return;
            }

            user = await getUserById(decoded.user_id);
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
            levels = ["Iniciante", "Intermediário", "Avançado"];
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
            const response: any = await updateUser(user.id, editForm);
            if (response && response.token) {
                Cookies.set("authToken", response.token);
            }
            user = { ...user, ...editForm };
            editing = false;
            toast.success("Perfil atualizado com sucesso!");
        } catch (error) {
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
            goto("/");
        } catch (error) {
            toast.error("Erro ao deletar conta");
        }
    }

    function handleLogout() {
        Cookies.remove("authToken");
        goto("/");
    }

    onMount(() => {
        loadUserData();
    });
</script>

<svelte:head>
    <title>Meu Perfil | SpeakUp</title>
    <meta name="description" content="Gerencie suas configurações de perfil, preferências de idioma, nível de aprendizado e assinatura na SpeakUp." />
    <meta property="og:title" content="Meu Perfil | SpeakUp" />
    <meta property="og:description" content="Gerencie suas configurações de perfil, preferências de idioma, nível de aprendizado e assinatura na SpeakUp." />
    <meta name="twitter:title" content="Meu Perfil | SpeakUp" />
    <meta name="twitter:description" content="Gerencie suas configurações de perfil, preferências de idioma, nível de aprendizado e assinatura na SpeakUp." />
</svelte:head>

<div class="profile-wrapper">
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

    <div class="profile-page">

    {#if loading}
        <div class="loading-state">
            <div class="spinner"></div>
            <p>Sincronizando seus dados...</p>
        </div>
    {:else if user}
        <div class="profile-container">
            <!-- Coluna da Esquerda: Dados Básicos -->
            <section class="profile-section main-info">
                <div class="user-avatar">
                    <div class="avatar-placeholder">{user.name.charAt(0)}</div>
                    <div class="user-meta">
                        <h2>{user.name}</h2>
                        <p>{user.email}</p>
                    </div>
                </div>

                {#if !editing}
                    <div class="info-grid">
                        <div class="info-card">
                            <label>Idioma de Estudo</label>
                            <div class="value-row">
                                <span class="flag">
                                    {user.language === "english" ? "🇺🇸" : user.language === "japanese" ? "🇯🇵" : "🌐"}
                                </span>
                                <span class="text">{user.language.charAt(0).toUpperCase() + user.language.slice(1)}</span>
                            </div>
                        </div>
                        <div class="info-card">
                            <label>Nível Atual</label>
                            <div class="value-row">
                                <span class="badge">{user.level}</span>
                            </div>
                        </div>
                    </div>
                    <button class="edit-btn" on:click={() => (editing = true)}>
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>
                        Editar Perfil
                    </button>
                {:else}
                    <form class="edit-form" on:submit|preventDefault={handleUpdate}>
                        <div class="field">
                            <label for="name">Nome Completo</label>
                            <input id="name" type="text" bind:value={editForm.name} required />
                        </div>
                        <div class="field-group">
                            <div class="field">
                                <label for="language">Idioma</label>
                                <select id="language" bind:value={editForm.language} on:change={handleLanguageChange}>
                                    <option value="english">Inglês</option>
                                    <option value="japanese">Japonês</option>
                                </select>
                            </div>
                            <div class="field">
                                <label for="level">Nível</label>
                                <select id="level" bind:value={editForm.level}>
                                    {#each levels as l}
                                        <option value={l}>{l}</option>
                                    {/each}
                                </select>
                            </div>
                        </div>
                        <div class="form-actions">
                            <button type="submit" class="save-btn">Salvar Alterações</button>
                            <button type="button" class="cancel-btn" on:click={() => (editing = false)}>Cancelar</button>
                        </div>
                    </form>
                {/if}
            </section>

            <!-- Coluna da Direita: Conta -->
            <div class="secondary-columns">

                <section class="profile-section danger-zone">
                    <h3>Zona de Perigo</h3>
                    <p>A exclusão da conta é permanente e removerá todos os seus dados de progresso.</p>
                    <button class="delete-account-btn" on:click={() => (showDeleteModal = true)}>
                        Excluir minha conta
                    </button>
                </section>
            </div>
        </div>
    {/if}
    </div>
</div>

{#if showDeleteModal}
    <div class="modal-overlay">
        <div class="modal">
            <h3>Tem certeza absoluta?</h3>
            <p>Isso apagará seu histórico de mensagens e configurações. Digite <strong>DELETAR</strong> abaixo para confirmar.</p>
            <input type="text" bind:value={deleteConfirmation} placeholder="DELETAR" />
            <div class="modal-buttons">
                <button class="confirm-delete" on:click={handleDelete}>Confirmar Exclusão</button>
                <button class="cancel-modal" on:click={() => (showDeleteModal = false)}>Voltar</button>
            </div>
        </div>
    </div>
{/if}

<style>
    :global(body) {
        background-color: #0a0a0a;
        margin: 0;
        font-family: 'Inter', system-ui, sans-serif;
    }

    .profile-wrapper {
        min-height: 100vh;
        background-color: #0a0a0a;
        color: #e9e9e9;
        display: flex;
        flex-direction: column;
    }

    .profile-page {
        flex: 1;
        padding: 2rem;
        max-width: 1100px;
        width: 100%;
        margin: 0 auto;
        box-sizing: border-box;
    }

    /* Top Navbar */
    .top-nav {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.25rem 2.5rem;
        border-bottom: 1px solid #1a1a1a;
        background-color: #0f0f0f;
        box-sizing: border-box;
        width: 100%;
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

    .loading-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 50vh;
        gap: 1.5rem;
    }

    .spinner {
        width: 40px;
        height: 40px;
        border: 3px solid rgba(255, 255, 255, 0.1);
        border-top-color: #5c6dff;
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin { to { transform: rotate(360deg); } }

    .profile-container {
        display: grid;
        grid-template-columns: 1.5fr 1fr;
        gap: 2rem;
        animation: fadeIn 0.5s ease-out;
    }

    @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

    .profile-section {
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(255, 255, 255, 0.08);
        border-radius: 24px;
        padding: 2.5rem;
        backdrop-filter: blur(10px);
    }

    .user-avatar {
        display: flex;
        align-items: center;
        gap: 1.5rem;
        margin-bottom: 3rem;
    }

    .avatar-placeholder {
        width: 80px;
        height: 80px;
        background: linear-gradient(135deg, #5c6dff 0%, #a5a9c8 100%);
        border-radius: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 2rem;
        font-weight: 800;
        color: #fff;
    }

    .user-meta h2 { margin: 0; font-size: 1.8rem; font-weight: 800; }
    .user-meta p { margin: 0.3rem 0 0; color: #888; font-size: 1rem; }

    .info-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
        margin-bottom: 2rem;
    }

    .info-card {
        background: rgba(255, 255, 255, 0.02);
        padding: 1.5rem;
        border-radius: 16px;
        border: 1px solid rgba(255, 255, 255, 0.05);
    }

    .info-card label { display: block; font-size: 0.8rem; color: #666; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 0.75rem; font-weight: 700; }

    .value-row { display: flex; align-items: center; gap: 0.75rem; }
    .flag { font-size: 1.5rem; }
    .text { font-size: 1.1rem; font-weight: 600; }
    .badge { background: #5c6dff; color: #fff; padding: 0.3rem 0.8rem; border-radius: 8px; font-weight: 700; font-size: 0.9rem; }

    .edit-btn {
        width: 100%;
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.1);
        color: #fff;
        padding: 1rem;
        border-radius: 14px;
        font-weight: 600;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.6rem;
        cursor: pointer;
        transition: all 0.2s;
    }

    .edit-btn:hover { background: rgba(255, 255, 255, 0.05); border-color: rgba(255, 255, 255, 0.2); }

    /* Form Styles */
    .edit-form { display: grid; gap: 1.25rem; }
    .field { display: grid; gap: 0.5rem; }
    .field label { font-size: 0.9rem; color: #888; }
    .field-group { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    
    input, select {
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        color: #fff;
        padding: 0.8rem 1rem;
        border-radius: 12px;
        font-size: 1rem;
    }

    input:focus, select:focus { outline: none; border-color: #5c6dff; background: rgba(255, 255, 255, 0.08); }

    .form-actions { display: flex; gap: 1rem; margin-top: 1rem; }
    .save-btn { flex: 2; background: #5c6dff; color: #fff; border: none; padding: 1rem; border-radius: 12px; font-weight: 700; cursor: pointer; }
    .cancel-btn { flex: 1; background: transparent; color: #888; border: 1px solid rgba(255, 255, 255, 0.1); padding: 1rem; border-radius: 12px; font-weight: 600; cursor: pointer; }

    /* Secondary Column */
    .secondary-columns { display: grid; gap: 2rem; }

    .section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
    .section-header h3 { margin: 0; font-size: 1.2rem; font-weight: 700; }

    .status-pill { padding: 0.3rem 0.8rem; border-radius: 999px; font-size: 0.75rem; font-weight: 800; text-transform: uppercase; }
    .status-pill.active { background: #7ee5b3; color: #0a4d31; }
    .status-pill.free { background: rgba(255, 255, 255, 0.1); color: #888; }

    .sub-details { display: grid; gap: 1.25rem; }
    .sub-row { display: flex; justify-content: space-between; font-size: 0.95rem; }
    .sub-row .label { color: #888; }
    .sub-row .value { font-weight: 600; }
    .sub-row .value.active { color: #7ee5b3; }

    .manage-btn, .upgrade-btn {
        width: 100%;
        padding: 0.9rem;
        border-radius: 12px;
        font-weight: 700;
        text-align: center;
        text-decoration: none;
        cursor: pointer;
        transition: all 0.2s;
    }

    .manage-btn { background: rgba(255, 255, 255, 0.05); color: #fff; border: 1px solid rgba(255, 255, 255, 0.1); }
    .upgrade-btn { background: #5c6dff; color: #fff; border: none; }

    .danger-zone h3 { color: #ff6b6b; margin-top: 0; }
    .danger-zone p { font-size: 0.85rem; color: #888; line-height: 1.5; margin-bottom: 1.5rem; }

    .delete-account-btn {
        width: 100%;
        padding: 0.8rem;
        background: transparent;
        border: 1px solid rgba(255, 107, 107, 0.2);
        color: #ff6b6b;
        border-radius: 12px;
        font-size: 0.85rem;
        font-weight: 600;
        cursor: pointer;
    }

    .delete-account-btn:hover { background: rgba(255, 107, 107, 0.05); border-color: #ff6b6b; }

    /* Modal */
    .modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.8); backdrop-filter: blur(8px); display: flex; align-items: center; justify-content: center; z-index: 100; padding: 1.5rem; }
    .modal { background: #1a1a1a; border: 1px solid rgba(255, 255, 255, 0.1); padding: 2.5rem; border-radius: 8px; max-width: 450px; text-align: center; }
    .modal h3 { color: #fff; margin-top: 0; font-size: 1.4rem; }
    .modal p { color: #888; margin-bottom: 2rem; line-height: 1.6; }
    .modal input { width: 100%; box-sizing: border-box; margin-bottom: 1.5rem; text-align: center; border-color: #ff6b6b; border-radius: 6px; }
    .modal-buttons { display: flex; gap: 1rem; }
    .confirm-delete { flex: 1; background: #ff6b6b; color: #fff; border: none; padding: 0.9rem; border-radius: 6px; font-weight: 700; cursor: pointer; }
    .cancel-modal { flex: 1; background: transparent; color: #fff; border: 1px solid rgba(255, 255, 255, 0.1); padding: 0.9rem; border-radius: 6px; font-weight: 600; cursor: pointer; }

    @media (max-width: 900px) {
        .profile-container { grid-template-columns: 1fr; }
    }

    @media (max-width: 600px) {
        .profile-page { padding: 1rem; }
        .profile-section { padding: 1.5rem; }
        .field-group { grid-template-columns: 1fr; }
        .info-grid { grid-template-columns: 1fr; }
    }
</style>
