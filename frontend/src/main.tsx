import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import { RouterProvider } from '@tanstack/react-router';

import { QueryProvider, ThemeProvider } from '@/app/providers';
import { router } from '@/app/router';
import '@/app/styles/index.css';

import { queryClient } from '@/shared/api';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <QueryProvider>
        <RouterProvider router={router} context={{ queryClient }} />
      </QueryProvider>
    </ThemeProvider>
  </StrictMode>,
);
