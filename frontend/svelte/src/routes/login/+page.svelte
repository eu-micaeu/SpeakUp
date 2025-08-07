<script lang="ts">
    import { goto } from "$app/navigation";
    import { env } from "$env/dynamic/public";
    import {
        getAuth,
        GoogleAuthProvider,
        signInWithPopup,
    } from "firebase/auth";
    import { initializeApp } from "firebase/app";
    import {
        login as loginApi,
        register as registerApi,
    } from "../../utils/api";

    type UserData = {
        name: string;
        email: string;
        password: string;
        language: string;
        level: string;
    };

    const firebaseConfig = {
        apiKey: env.PUBLIC_FIREBASE_API_KEY,
        authDomain: env.PUBLIC_FIREBASE_AUTH_DOMAIN,
        projectId: env.PUBLIC_FIREBASE_PROJECT_ID,
        storageBucket: env.PUBLIC_FIREBASE_STORAGE_BUCKET,
        messagingSenderId: env.PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
        appId: env.PUBLIC_FIREBASE_APP_ID,
        measurementId: env.PUBLIC_FIREBASE_MEASUREMENT_ID,
    };

    const app = initializeApp(firebaseConfig);
    const auth = getAuth(app);

    let isLoading = false;
    let toastMessage = "";
    let toastType: "success" | "error" | "" = "";

    function showToast(message: string, type: "success" | "error") {
        toastMessage = message;
        toastType = type;
        setTimeout(() => {
            toastMessage = "";
            toastType = "";
        }, 3000);
    }

    async function handleLogin(event: Event) {
        event.preventDefault();
        isLoading = true;
        const formData = new FormData(event.target as HTMLFormElement);
        const email = formData.get("email") as string;
        const password = formData.get("password") as string;

        try {
            const response = await loginApi(email, password);

            if (response) {
                showToast("Login realizado com sucesso!", "success");
                setTimeout(() => {
                    goto("/chat");
                }, 2000);
            } else {
                showToast("Email ou senha inválidos", "error");
            }
        } catch (err) {
            console.error(err);
            showToast("Email ou senha inválidos", "error");
        } finally {
            isLoading = false;
        }
    }

    async function handleGoogleLogin() {
        const provider = new GoogleAuthProvider();
        try {
            const result = await signInWithPopup(auth, provider);
            const user = result.user;
            const email = user.email!;
            const password = user.uid;
            const name = user.displayName!;

            let isFirstTime = false;

            try {
                await loginApi(email, password);
            } catch (err) {
                isFirstTime = true;
                const userData: UserData = {
                    name,
                    email,
                    password,
                    language: "en",
                    level: "beginner",
                };
                await registerApi(userData);
                await loginApi(email, password);
            }

            showToast(`Bem-vindo, ${name}!`, "success");

            setTimeout(() => {
                if (isFirstTime) {
                    goto("/onboarding");
                } else {
                    goto("/chat");
                }
            }, 2000);
        } catch (error) {
            console.error(error);
            showToast("Erro ao fazer login com Google", "error");
        }
    }
</script>

{#if toastMessage}
    <div class="toast {toastType}">
        {toastMessage}
    </div>
{/if}

<div class="pageLogin">
    <form on:submit={handleLogin}>
        <div class="card">
            <div class="textCenter">
                <a href="/">
                    <img src="./logo.png" alt="SpeakUp Logo" class="logo" />
                </a>
                <h2 class="h2">Bem-vindo de volta!</h2>
                <p class="p">Entre para continuar a conversa</p>
            </div>

            <div class="inputContainer">
                <label for="email" class="inputLabel"> Email </label>
                <input
                    id="email"
                    type="email"
                    name="email"
                    placeholder="seu@email.com"
                    required
                    class="input"
                />
            </div>

            <div class="inputContainer">
                <label for="password" class="inputLabel"> Senha </label>
                <input
                    id="password"
                    type="password"
                    name="password"
                    placeholder="••••••••"
                    required
                    class="input"
                />
            </div>

            <button type="submit" disabled={isLoading} class="button">
                {isLoading ? "Entrando..." : "Entrar"}
            </button>

            <button
                type="button"
                on:click={handleGoogleLogin}
                disabled={isLoading}
                class="buttonGoogle"
            >
                <img src="./google.png" alt="Google Icon" class="googleIcon" />
            </button>

            <div class="textCenter">
                <a href="/register" class="styledLink">
                    Não tem uma conta ainda? <span>Registre-se</span>
                </a>
            </div>
        </div>
    </form>
</div>

<style>
    .pageLogin {
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
        transition: 1s ease all;
    }

    .logo:hover {
        transform: scale(1.1);
        cursor: pointer;
    }

    .h2 {
        font-size: 1.875rem;
        font-weight: 700;
        margin-bottom: 0.5rem;
    }

    .p {
        color: rgb(156, 163, 175);
        margin-bottom: 2rem;
    }

    .inputContainer {
        margin-bottom: 1.5rem;
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

    .buttonGoogle {
        background-color: rgb(53, 53, 53);
        border: none;
        margin: 0 0 1.5rem 0;
        border-radius: 10px;
        width: 100%;
        padding: 0.75rem;
    }

    .buttonGoogle:hover:not(:disabled) {
        background-color: rgb(229, 231, 235);
        cursor: pointer;
    }

    .toast {
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 1rem;
        border-radius: 0.5rem;
        color: white;
        z-index: 1000;
        animation: slideIn 0.3s ease-in;
    }

    .toast.success {
        background-color: rgb(34, 197, 94);
    }

    .toast.error {
        background-color: rgb(239, 68, 68);
    }

    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }

    .styledLink {
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

    .googleIcon {
        width: 25px;
    }

    @media screen and (max-width: 768px) {
        .card {
            padding: 1.5rem;
        }

        .h2 {
            font-size: 1.5rem;
            text-align: center;
        }

        .p {
            font-size: 0.9rem;
            text-align: center;
        }

        .logo {
            height: 2.5rem;
            margin-bottom: 0.8rem;
        }

        .input {
            padding: 0.65rem;
            font-size: 0.95rem;
        }

        .button {
            padding: 0.65rem;
            font-size: 0.95rem;
        }

        .styledLink {
            font-size: 0.8rem;
            text-align: center;
        }
    }
</style>
