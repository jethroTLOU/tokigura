export type Theme = 'dark' | 'light' | 'system';
export type ResolvedTheme = 'dark' | 'light';

export type ThemeProviderState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
};
