import { useForm } from '@tanstack/react-form';
import { Link } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';

import { Button } from '@/shared/ui/button';
import { Field, FieldDescription, FieldError, FieldGroup, FieldLabel } from '@/shared/ui/field';
import { Input } from '@/shared/ui/input';

import { useLoginMutation } from '../model/mutation';
import { createLoginSchema, type Login } from '../model/schema';

export function LoginForm() {
  const mutation = useLoginMutation();
  const { t } = useTranslation('auth', { keyPrefix: 'login' });
  const schema = createLoginSchema(t);

  const form = useForm({
    defaultValues: {
      username: '',
      password: '',
    } satisfies Login,
    validators: { onChange: schema },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync(value);
    },
  });
  return (
    <form
      id="login"
      className="flex flex-col gap-6"
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      <FieldGroup>
        <div className="flex flex-col items-center gap-1 text-center">
          <h1 className="text-2xl font-bold">{t('title')}</h1>
          <p className="text-muted-foreground text-sm text-balance">{t('description')}</p>
        </div>

        <form.Field name="username">
          {(field) => {
            const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
            return (
              <Field data-invalid={isInvalid}>
                <FieldLabel htmlFor={field.name}>{t('fields.username.label')}</FieldLabel>
                <Input
                  id={field.name}
                  name={field.name}
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                  aria-invalid={isInvalid}
                  autoComplete="off"
                />
                <FieldError errors={field.state.meta.errors} className="text-xs" />
              </Field>
            );
          }}
        </form.Field>

        <form.Field name="password">
          {(field) => {
            const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
            return (
              <Field data-invalid={isInvalid}>
                <FieldLabel htmlFor={field.name}>{t('fields.password.label')}</FieldLabel>
                <Input
                  id={field.name}
                  name={field.name}
                  type="password"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                  aria-invalid={isInvalid}
                  autoComplete="off"
                />
                <FieldError errors={field.state.meta.errors} className="text-xs" />
              </Field>
            );
          }}
        </form.Field>

        <Field>
          <Button type="submit">{t('submit_btn')}</Button>
        </Field>
        <Field>
          <FieldDescription className="px-6 text-center">
            {t('footer.text')}{' '}
            <Link to={'/register'} search={(prev) => ({ ...prev })}>
              {t('footer.link')}
            </Link>
          </FieldDescription>
        </Field>
      </FieldGroup>
    </form>
  );
}
