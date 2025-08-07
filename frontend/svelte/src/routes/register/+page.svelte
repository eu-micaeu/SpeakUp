<script lang="ts">
    import { goto } from "$app/navigation";
    import toast from "svelte-french-toast";
    import { register } from "../../utils/api";

    let isLoading = false;
    let error = "";
    let selectedLanguage = "";
    let levels: string[] = [];

    function handleLanguageChange(e: Event) {
        const target = e.target as HTMLSelectElement;
        const language = target.value;
        selectedLanguage = language;

        if (language === "english") {
            levels = ["A1", "A2", "B1", "B2", "C1", "C2"];
        } else if (language === "japanese") {
            levels = ["N5", "N4", "N3", "N2", "N1"];
        } else {
            levels = [];
        }
    }

    function isValidPassword(password: string): boolean {
        const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
        return regex.test(password);
    }

    async function handleRegister(event: SubmitEvent) {
        event.preventDefault();
        isLoading = true;
        error = "";

        try {
            const form = event.target as HTMLFormElement;
            const formData = new FormData(form);

            const userData = {
                name: formData.get("name") as string,
                email: formData.get("email") as string,
                password: formData.get("password") as string,
                language: formData.get("language") as string,
                level: formData.get("level") as string,
            };

            const confirmPassword = formData.get("confirm-password") as string;

            if (userData.password !== confirmPassword) {
                error = "As senhas não coincidem";
                isLoading = false;
                return;
            }

            if (!isValidPassword(userData.password)) {
                error =
                    "A senha deve ter pelo menos 8 caracteres, uma letra maiúscula, uma minúscula, um número e um caractere especial.";
                isLoading = false;
                return;
            }

            const response = await register(userData);

            if (response.message === "User created successfully") {
                toast.success("Usuário criado com sucesso!");
                setTimeout(() => {
                    goto("/login");
                }, 2000);
            } else {
                error = "Erro ao criar usuário.";
            }
        } catch (err: unknown) {
            if (err && typeof err === "object" && "response" in err) {
                const axiosError = err as { response: { status: number } };
                if (axiosError.response && axiosError.response.status === 409) {
                    error = "Email já cadastrado.";
                } else {
                    error = "Erro ao criar usuário.";
                }
            } else {
                error = "Erro ao criar usuário.";
            }
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="pageRegister">
    <!-- Toast container is handled automatically by svelte-french-toast -->
    <form on:submit={handleRegister}>
        <div class="card">
            <div class="textCenter">
                <a href="/">
                    <img src="./logo.png" alt="SpeakUp Logo" class="logo" />
                </a>
                <h2 class="title">Crie sua conta</h2>
                <p class="subtitle">Entre para continuar a conversa</p>
            </div>

            {#if error}
                <div class="errorMessage">
                    {error}
                </div>
            {/if}

            <div class="inputContainerFull">
                <label for="name" class="inputLabel"> Nome completo </label>
                <input
                    type="text"
                    id="name"
                    name="name"
                    placeholder="Digite seu nome"
                    required
                    class="input"
                />
            </div>

            <div class="inputRow">
                <div class="inputContainer">
                    <label for="language" class="inputLabel">
                        Idioma que deseja aprender
                    </label>
                    <select
                        id="language"
                        name="language"
                        on:change={handleLanguageChange}
                        required
                        class="select"
                    >
                        <option value="" disabled selected>Selecione</option>
                        <option value="english">Inglês</option>
                        <option value="japanese">Japonês</option>
                    </select>
                </div>

                <div class="inputContainer">
                    <label for="level" class="inputLabel">
                        Nível atual no idioma
                    </label>
                    <select
                        id="level"
                        name="level"
                        required
                        disabled={!selectedLanguage}
                        class="select"
                    >
                        <option value="" disabled selected>Selecione</option>
                        {#each levels as level}
                            <option value={level}>{level}</option>
                        {/each}
                    </select>
                </div>
            </div>

            <div class="inputContainerFull">
                <label for="email" class="inputLabel"> Email </label>
                <input
                    type="email"
                    id="email"
                    name="email"
                    placeholder="Digite seu email"
                    required
                    class="input"
                />
            </div>

            <div class="inputRow">
                <div class="inputContainer">
                    <label for="password" class="inputLabel"> Senha </label>
                    <input
                        type="password"
                        id="password"
                        name="password"
                        placeholder="Digite sua senha"
                        required
                        class="input"
                    />
                </div>

                <div class="inputContainer">
                    <label for="confirm-password" class="inputLabel">
                        Confirmar Senha
                    </label>
                    <input
                        type="password"
                        id="confirm-password"
                        name="confirm-password"
                        placeholder="Confirme sua senha"
                        required
                        class="input"
                    />
                </div>
            </div>

            <button type="submit" disabled={isLoading} class="button">
                {isLoading ? "Carregando..." : "Cadastrar"}
            </button>

            <a href="/login" class="styledLink">
                Já tem uma conta? <span>Entre agora!</span>
            </a>
        </div>
    </form>
</div>

<style>
    .pageRegister {
        display: flex;
        justify-content: center;
        align-items: center;
        background-color: rgb(0, 0, 0);
        color: white;
        height: 100vh;
    }

    form {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    .card {
        width: 100%;
        max-width: 28rem;
        padding: 2rem;
        background-color: rgb(31, 31, 31);
        border-radius: 0px 20px 20px 20px;
    }

    .textCenter {
        text-align: center;
    }

    .logo {
        height: 3rem;
        margin: 0 auto 1rem;
        cursor: pointer;
    }

    .title {
        font-size: 1.875rem;
        font-weight: 700;
        margin-bottom: 0.5rem;
    }

    .subtitle {
        color: rgb(156, 163, 175);
        margin-bottom: 2rem;
    }

    .errorMessage {
        background-color: rgb(239, 68, 68);
        color: white;
        padding: 1rem;
        border-radius: 0.5rem;
        margin-bottom: 1.5rem;
    }

    .inputRow {
        display: flex;
        justify-content: space-between;
        align-items: center;
        flex-wrap: wrap;
    }

    .inputContainer {
        margin-bottom: 1.5rem;
        width: 48%;
    }

    .inputContainerFull {
        margin-bottom: 1.5rem;
        width: 100%;
    }

    .inputLabel {
        display: block;
        font-size: 0.875rem;
        color: rgb(209, 213, 219);
        margin-bottom: 0.25rem;
    }

    .input {
        width: 100%;
        padding: 0.75rem;
        border-radius: 0.375rem;
        background-color: rgb(55, 65, 81);
        border: 1px solid rgb(75, 85, 99);
        color: white;
        box-sizing: border-box;
    }

    .input::placeholder {
        color: rgb(156, 163, 175);
    }

    .button {
        width: 100%;
        padding: 0.75rem;
        margin: 0 0 1.5rem 0;
        border-radius: 0.375rem;
        border: none;
        background-color: rgb(53, 53, 53);
        color: white;
        cursor: pointer;
        transition: background-color 0.2s;
        box-sizing: border-box;
    }

    .button:hover:not(:disabled) {
        background-color: rgb(37, 99, 235);
    }

    .button:disabled {
        opacity: 0.75;
        cursor: not-allowed;
    }

    .select {
        width: 100%;
        padding: 0.75rem;
        border-radius: 0.375rem;
        background-color: rgb(55, 65, 81);
        border: 1px solid rgb(75, 85, 99);
        color: white;
        box-sizing: border-box;
    }

    .styledLink {
        display: block;
        margin-top: 10px;
        text-align: center;
        color: rgb(156, 163, 175);
        text-decoration: none;
        font-size: 0.875rem;
    }

    .styledLink span {
        color: rgb(59, 130, 246);
        font-weight: 500;
    }

    .styledLink span:hover {
        color: rgb(37, 99, 235);
    }

    @media screen and (max-width: 768px) {
        .card {
            padding: 1.5rem;
        }

        .logo {
            height: 2rem;
            margin-bottom: 1rem;
        }

        .title {
            font-size: 1.5rem;
            text-align: center;
        }

        .subtitle {
            font-size: 0.95rem;
            text-align: center;
        }

        .inputRow {
            flex-direction: column;
        }

        .inputContainer {
            width: 100%;
        }

        .inputContainerFull {
            width: 100%;
        }

        .input,
        .select,
        .button {
            padding: 0.65rem;
            font-size: 0.95rem;
        }

        .errorMessage {
            font-size: 0.875rem;
            padding: 0.75rem;
        }

        .styledLink {
            font-size: 0.8rem;
        }
    }
</style>
