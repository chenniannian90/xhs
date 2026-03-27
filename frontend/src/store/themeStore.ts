import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

type Theme = 'light' | 'dark';

interface ThemeStore {
  theme: Theme;
  toggleTheme: () => void;
  setTheme: (theme: Theme) => void;
}

export const useThemeStore = create<ThemeStore>()(
  persist(
    (set) => ({
      theme: 'light',
      toggleTheme: () => set((state) => {
        const newTheme = state.theme === 'light' ? 'dark' : 'light';
        // Apply theme to document
        document.documentElement.classList.remove('light', 'dark');
        document.documentElement.classList.add(newTheme);
        return { theme: newTheme };
      }),
      setTheme: (theme) => set(() => {
        // Apply theme to document
        document.documentElement.classList.remove('light', 'dark');
        document.documentElement.classList.add(theme);
        return { theme };
      }),
    }),
    {
      name: 'navhub-theme',
    }
  )
);

// Initialize theme on app load
export const initializeTheme = () => {
  const theme = useThemeStore.getState().theme;
  document.documentElement.classList.remove('light', 'dark');
  document.documentElement.classList.add(theme);
};
