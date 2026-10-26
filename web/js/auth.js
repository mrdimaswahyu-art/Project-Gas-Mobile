// Modul Autentikasi Frontend (Garuda Auto Service)

const AUTH_TOKEN_KEY = "gas_token";
const AUTH_USER_KEY = "gas_user";

// Mengambil token dari localStorage
function getToken() {
    return localStorage.getItem(AUTH_TOKEN_KEY);
}

// Mengambil nama pengguna saat ini
function getCurrentUser() {
    return localStorage.getItem(AUTH_USER_KEY) || "Staf";
}

// Memeriksa status login saat mengakses halaman utama
function checkAuth() {
    const token = getToken();
    if (!token) {
        window.location.href = "/login";
    }
}

// Mencegah user yang sudah login mengakses halaman login/register
function redirectIfLoggedIn() {
    const token = getToken();
    if (token) {
        window.location.href = "/";
    }
}

// Fungsi pembantu fetch dengan header otorisasi JWT otomatis
async function authFetch(url, options = {}) {
    const token = getToken();
    if (!token) {
        window.location.href = "/login";
        return;
    }

    const headers = {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
        ...(options.headers || {})
    };

    const response = await fetch(url, { ...options, headers });

    // Jika token kedaluwarsa atau tidak valid (401), logout dan arahkan ke login
    if (response.status === 401) {
        logout();
        return;
    }

    return response;
}

// Logout dari sistem
function logout() {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    localStorage.removeItem(AUTH_USER_KEY);
    window.location.href = "/login";
}
