import type { TFunction } from 'i18next';
import { z } from 'zod';

export const createRegisterSchema = (t: TFunction<'auth', 'register'>) =>
  z
    .object({
      username: z.string().regex(/^[A-Za-z0-9@.+\-_]+$/, {
        message: t('fields.username.errors.regex'),
      }),
      email: z.email(),
      password: z.string().nonempty(),
      confirmPassword: z.string(),
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: t('fields.password.errors.missmatch'),
      path: ['confirmPassword'],
    });

export type Register = z.infer<ReturnType<typeof createRegisterSchema>>;
