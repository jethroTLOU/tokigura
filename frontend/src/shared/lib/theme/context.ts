import { createContext, use } from 'react';

import type { ThemeProviderState } from './types';

export const ThemeProviderContext = createContext<ThemeProviderState | undefined>(undefined);

export const useTheme = () => {
  const context = use(ThemeProviderContext);

  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }

  return context;
};
