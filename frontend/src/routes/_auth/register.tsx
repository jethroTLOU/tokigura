import { createFileRoute } from '@tanstack/react-router';

import { RegisterForm } from '@/features/auth/register';

export const Route = createFileRoute('/_auth/register')({
  staticData: {
    authVariant: 'register',
  },
  component: () => <RegisterForm />,
  pendingComponent: () => 'loading..',
});
