'use client';

import { Button, FormTextField, Icon, Token } from '@squantumengine/horizon';
import useLoginPage from './_hooks/useLoginPage';
import './_styles/login.css';

export default function LoginPage() {
  const {
    control,
    errors,
    isDirty,
    isValid,
    showPassword,
    setShowPassword,
    submit,
    isPending,
    isError,
    errorMessage
  } = useLoginPage();

  return (
    <div className="screen">
      <div className="login">
        <aside className="login__aside">
          <div className="login__brand">
            <div className="login__brand-mark">
              <Icon name="check" />
            </div>
            TodoApp
          </div>
          <div className="login__hero">
            <h2 style={{ fontSize: 32, lineHeight: 1.15, margin: '0 0 12px' }}>
              Personal tasks, in one place.
            </h2>
            <p style={{ fontSize: 15, opacity: 0.85 }}>
              Internal preview environment for the lintas boilerplate.
              Sign in with a demo account to explore.
            </p>
          </div>
          <div className="login__meta" style={{ fontSize: 12, opacity: 0.7 }}>
            internal · v0.1
          </div>
        </aside>

        <main className="login__main">
          <div className="login__card">
            {isError && (
              <div className="inline-alert" role="alert">
                <Icon name="exclamation-circle" />
                <div>{errorMessage ?? 'Login failed.'}</div>
              </div>
            )}
            <h1 style={{ fontSize: 24, fontWeight: 700, margin: '0 0 8px' }}>
              Sign in
            </h1>
            <p style={{ fontSize: 14, color: 'var(--hz-text-tertiary)', margin: '0 0 24px' }}>
              Use your email and password.
            </p>
            <form onSubmit={submit} className="flex w-full flex-col gap-4">
              <div className="flex w-full flex-col gap-2">
                <div className="text-base font-bold">Email</div>
                <FormTextField
                  full
                  controller={{ control, name: 'email' }}
                  type="email"
                  placeholder="Masukkan email Anda"
                  isInvalid={!!errors.email?.message}
                  errorMessage={errors.email?.message as string}
                />
              </div>
              <div className="flex w-full flex-col gap-2">
                <div className="text-base font-bold">Password</div>
                <FormTextField
                  full
                  controller={{ control, name: 'password' }}
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Masukkan password Anda"
                  isInvalid={!!errors.password?.message}
                  errorMessage={errors.password?.message as string}
                  suffix={
                    <div
                      className="cursor-pointer"
                      onClick={() => setShowPassword(!showPassword)}>
                      <Icon
                        name={showPassword ? 'eye-slash' : 'eye'}
                        color={Token.COLORS.neutral[950]}
                      />
                    </div>
                  }
                />
              </div>
              <Button
                size="lg"
                full
                disabled={!isDirty || !isValid}
                variant="primary"
                type="submit"
                loading={isPending}>
                Sign in
              </Button>
            </form>
          </div>
        </main>
      </div>
    </div>
  );
}
