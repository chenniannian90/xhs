import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Helper function to extract token from various possible storage structures
const getToken = (authData: any): string | null => {
  if (!authData) return null;

  // Zustand persist v4+ structure: { state: { token: "..." }, version: 0 }
  if (authData.state?.token) {
    return authData.state.token;
  }

  // Direct storage structure (older versions or custom): { token: "..." }
  if (authData.token) {
    return authData.token;
  }

  // Nested state structure (edge case): { state: { state: { token: "..." } } }
  if (authData.state?.state?.token) {
    return authData.state.state.token;
  }

  return null;
};

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authDataString = localStorage.getItem('navhub-auth');
    if (authDataString) {
      try {
        const authData = JSON.parse(authDataString);
        const token = getToken(authData);

        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
      } catch (error) {
        console.error('Failed to parse auth data:', error);
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle 401 Unauthorized
    if (error.response?.status === 401) {
      console.warn('Unauthorized request - clearing auth data');

      // Clear auth data from localStorage
      localStorage.removeItem('navhub-auth');

      // Clear Zustand store (if it's loaded)
      try {
        const { useAuthStore } = require('../store/authStore');
        useAuthStore.getState().clearAuth();
      } catch (e) {
        // AuthStore might not be loaded yet, that's ok
      }

      // Redirect to login (only if we're not already on login page)
      if (!window.location.pathname.includes('/login')) {
        window.location.href = '/login';
      }
    }

    // Log other errors for debugging
    if (error.response && error.response.status >= 500) {
      console.error('Server error:', error.response.status, error.response.data);
    }

    return Promise.reject(error);
  }
);

export default api;
