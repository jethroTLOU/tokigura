import { createFileRoute } from '@tanstack/react-router';

import { LoginForm } from '@/features/auth/login';

export const Route = createFileRoute('/_auth/login')({
  staticData: {
    authVariant: 'login',
  },
  pendingComponent: () => 'loading..',
  component: () => <LoginForm />,
});
