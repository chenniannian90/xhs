import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import api from '../utils/api';
import type { User } from '../types';

interface AuthStore {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, token: string) => void;
  clearAuth: () => void;
  updateUser: (user: Partial<User>) => void;
  login: (email: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setAuth: (user, token) => set({ user, token, isAuthenticated: true }),
      clearAuth: () => set({ user: null, token: null, isAuthenticated: false }),
      updateUser: (updates) => set((state) => ({
        user: state.user ? { ...state.user, ...updates } : null,
      })),
      login: async (email, password) => {
        try {
          const response = await api.post('/auth/login', { email, password });
          const { user, token } = response.data.data;

          // Debug logging
          console.log('✅ Login successful!');
          console.log('  User:', user.username);
          console.log('  Token:', token ? token.substring(0, 20) + '...' : 'null');

          set({ user, token, isAuthenticated: true });

          // Verify storage
          setTimeout(() => {
            const stored = localStorage.getItem('navhub-auth');
            if (stored) {
              const parsed = JSON.parse(stored);
              console.log('💾 Stored in localStorage:');
              console.log('  Has user:', !!parsed.state?.user);
              console.log('  Has token:', !!parsed.state?.token);
              console.log('  Is authenticated:', parsed.state?.isAuthenticated);
            }
          }, 100);
        } catch (error: any) {
          console.error('❌ Login failed:', error);
          throw new Error(error.response?.data?.message || '登录失败');
        }
      },
      register: async (username, email, password) => {
        try {
          const response = await api.post('/auth/register', { username, email, password });
          const { user, token } = response.data.data;
          set({ user, token, isAuthenticated: true });
        } catch (error: any) {
          const errorData = error.response?.data;
          // Preserve full error structure for better UX
          if (errorData) {
            const enhancedError: any = new Error(errorData.error || errorData.message || '注册失败');
            enhancedError.details = errorData.details;
            enhancedError.suggestions = errorData.suggestions;
            throw enhancedError;
          }
          throw new Error('注册失败');
        }
      },
    }),
    {
      name: 'navhub-auth',
    }
  )
);
