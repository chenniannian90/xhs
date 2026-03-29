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
          const data = response.data.data;
          set({
            user: data.user,
            token: data.access_token,
            isAuthenticated: true,
          });
        } catch (error: any) {
          throw new Error(error.response?.data?.error || error.response?.data?.message || '登录失败');
        }
      },
      register: async (username, email, password) => {
        try {
          const response = await api.post('/auth/register', { username, email, password });
          const data = response.data.data;
          set({
            user: data.user,
            token: data.access_token,
            isAuthenticated: true,
          });
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
