import { Link, Outlet, useChildMatches } from '@tanstack/react-router';

import { loginImage } from '@/features/auth/login';
import { registerImage } from '@/features/auth/register';

import { APP_NAME } from '@/shared/config';

export function AuthLayout() {
  const isLogin = useChildMatches({
    select: (matches) => matches.some((m) => m.staticData.authVariant === 'login'),
  });

  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex justify-center gap-2 md:justify-start">
          <Link to="/" className="flex items-center gap-1 text-lg font-bold">
            {APP_NAME}
          </Link>
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-xs">
            <Outlet />
          </div>
        </div>
      </div>

      <div className="bg-muted hidden items-center justify-center p-10 lg:flex">
        <img
          src={isLogin ? loginImage : registerImage}
          alt={isLogin ? 'Login illustration' : 'Register illustration'}
          className="w-full max-w-md object-contain"
        />
      </div>
    </div>
  );
}
