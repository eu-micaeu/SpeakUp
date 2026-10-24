import Cookies from 'js-cookie';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';

let lastRedirectToastAt = 0;

function decodeBase64Url(base64Url) {
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), '=');
    return decodeURIComponent(
        atob(padded)
            .split('')
            .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
            .join('')
    );
}

// Recupera o valor de um cookie pelo nome
export function getAuthTokenFromCookies() {
    const token = Cookies.get('authToken');
    return token || null;
}

// Remove o valor de um cookie pelo nome
export function removeAuthTokenFromCookies() {
    Cookies.remove('authToken');
}

// Verifica se o token de autenticação existe e é válido
export function isAuthTokenValid() {
    const authToken = getAuthTokenFromCookies();

    if (authToken) {
        try {
            const parts = authToken.split('.');
            if (parts.length !== 3) {
                return false;
            }

            const [, payload] = parts;
            const data = JSON.parse(decodeBase64Url(payload));
            const exp = new Date(data.exp * 1000);

            if (exp <= new Date()) {
                removeAuthTokenFromCookies();
                return false;
            }

            return true;
        } catch (error) {
            console.error('Erro ao verificar o token:', error);
            removeAuthTokenFromCookies();
            return false;
        }
    }

    return false;
}

// Verifica se o token existe/é válido; caso contrário, avisa o usuário e redireciona para o login
export function requireAuthOrRedirect(message = 'Você precisa fazer login para acessar esta página.') {
    if (isAuthTokenValid()) {
        return true;
    }

    removeAuthTokenFromCookies();

    const now = Date.now();
    if (now - lastRedirectToastAt > 2000) {
        lastRedirectToastAt = now;
        try {
            toast.error(message);
        } catch (_) {
            // Ignora caso o Toaster ainda não esteja montado
        }
    }

    if (typeof window !== 'undefined' && window.location.pathname !== '/') {
        goto('/?login=true');
    }

    return false;
}

// Retorna os dados criptografados do token de autenticação
export function getDecodedToken() {
    const authToken = getAuthTokenFromCookies();

    if (authToken) {
        try {
            const parts = authToken.split('.');
            if (parts.length !== 3) {
                return null;
            }

            const [, payload] = parts;
            return JSON.parse(decodeBase64Url(payload));
        } catch (error) {
            console.error('Erro ao decodificar o token:', error);
        }
    }

    return null;
}