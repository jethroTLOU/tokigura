import { useMutation } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

import { login, sessionQueries } from '@/entities/session';

import { queryClient } from '@/shared/api';

export function useLoginMutation() {
  const navigate = useNavigate();
  const { t } = useTranslation('auth');

  return useMutation({
    mutationFn: login,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: sessionQueries.all() });
      navigate({ to: '/' });
      toast.success(t('login.success'));
    },
    onError: (error) => {
      console.error(error);
      toast.error(t('login.error'));
    },
  });
}
