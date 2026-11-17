import type { TFunction } from 'i18next';
import { z } from 'zod';

export const createLoginSchema = (t: TFunction<'auth', 'login'>) =>
  z.object({
    username: z.string().regex(/^[A-Za-z0-9@.+\-_]+$/, {
      error: t('fields.username.errors.regex'),
    }),
    password: z.string().min(1),
  });

export type Login = z.infer<ReturnType<typeof createLoginSchema>>;
