import { createFileRoute, isRedirect, redirect } from '@tanstack/react-router';

import { AuthLayout } from '@/widgets/auth';

import { sessionQueries } from '@/entities/session';

import { loadNamespace } from '@/shared/config/i18n';

export const Route = createFileRoute('/_auth')({
  beforeLoad: async ({ context }) => {
    try {
      await context.queryClient.query(sessionQueries.me());
      throw redirect({ to: '/', replace: true });
    } catch (error) {
      if (isRedirect(error)) throw error;
    }
  },
  component: AuthLayout,
  loader: async () => {
    await loadNamespace(['auth']);
  },
});
