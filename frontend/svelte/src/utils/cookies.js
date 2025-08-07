import Cookies from 'js-cookie';

// Recupera o valor de um cookie pelo nome
export function getAuthTokenFromCookies() {
    const token = Cookies.get('authToken');
    if (!token) {
        console.warn('authToken não encontrado nos cookies.');
    }
    return token || null;
}

// Remove o valor de um cookie pelo nome
export function removeAuthTokenFromCookies() {
    Cookies.remove('authToken');
}

// Verifica se o token de autenticação é válido
export function isAuthTokenValid() {
    const authToken = getAuthTokenFromCookies();

    if (authToken) {
        try {
            const parts = authToken.split('.');
            if (parts.length !== 3) {
                console.error('Formato inválido do token de autenticação.');
                return false;
            }

            const [, payload] = parts;
            const data = JSON.parse(atob(payload));
            const exp = new Date(data.exp * 1000);

            if (exp <= new Date()) {
                console.warn('Token expirado.');
                return false;
            }

            return true;
        } catch (error) {
            console.error('Erro ao verificar o token:', error);
        }
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
                console.error('Formato inválido do token.');
                return null;
            }

            const [, payload] = parts;
            return JSON.parse(atob(payload));
        } catch (error) {
            console.error('Erro ao decodificar o token:', error);
        }
    }

    return null;
}