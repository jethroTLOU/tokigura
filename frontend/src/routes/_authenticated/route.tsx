import { createFileRoute, redirect } from '@tanstack/react-router';

import { sessionQueries } from '@/entities/session';

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: async ({ context, location }) => {
    try {
      await context.queryClient.query({
        ...sessionQueries.me(),
        staleTime: 'static',
      });
    } catch {
      throw redirect({
        to: '/register',
        search: { redirect: location.href },
      });
    }
  },
});
