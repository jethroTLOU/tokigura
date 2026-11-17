import { useMutation } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

import { register } from '@/entities/session';

export function useRegisterMutation() {
  const navigate = useNavigate();
  const { t } = useTranslation('auth');

  return useMutation({
    mutationFn: register,
    onSuccess: () => {
      navigate({ to: '/login' });
      toast.success(t('register.success'));
    },
    onError: (error) => {
      console.error(error);
      toast.error(t('register.error'));
    },
  });
}
