import { create } from 'zustand';

interface UIState {
  sidebarOpen: boolean;
  currentPath: string;
  setSidebarOpen: (open: boolean) => void;
  setCurrentPath: (path: string) => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  currentPath: '/',
  
  setSidebarOpen: (sidebarOpen) => set({ sidebarOpen }),
  setCurrentPath: (currentPath) => set({ currentPath }),
}));
