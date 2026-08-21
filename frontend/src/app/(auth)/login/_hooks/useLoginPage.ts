import { useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useRouter } from 'next/navigation';
import { z } from 'zod';
import auth, { AuthError } from '@/services/auth';
import { defaultBaseUrlPage } from '@/static/route';

const loginSchema = z.object({
  email: z.string().min(1, 'Email wajib diisi').email('Format email tidak valid'),
  password: z.string().min(8, 'Password minimal 8 karakter')
});

type LoginFormValues = z.infer<typeof loginSchema>;

export default function useLoginPage() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [isPending, setIsPending] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);

  const {
    control,
    handleSubmit,
    formState: { errors, isDirty, isValid }
  } = useForm<LoginFormValues>({
    mode: 'onChange',
    resolver: zodResolver(loginSchema)
  });

  const submit = handleSubmit(async (loginData: LoginFormValues) => {
    setIsPending(true);
    setIsError(false);
    setError(null);
    try {
      await auth.signIn(loginData);
      const callbackUrl = new URLSearchParams(window.location.search).get('callbackUrl');
      router.push(callbackUrl || defaultBaseUrlPage);
    } catch (err) {
      setIsError(true);
      setError(err instanceof Error ? err : new Error('Login failed'));
    } finally {
      setIsPending(false);
    }
  });

  const errorMessage = useMemo(() => {
    if (!error) return null;

    if (error instanceof AuthError) {
      return error.message;
    }

    if (error.cause && typeof error.cause === 'object' && 'data' in error.cause) {
      const apiError = error.cause as {
        data?: {
          msg?: string;
          error?: string;
        };
      };

      if (apiError.data?.msg) {
        return apiError.data.msg;
      } else if (apiError.data?.error) {
        return apiError.data.error;
      }
    }

    return error.message || 'Login failed. Please check your credentials.';
  }, [error]);

  return {
    control,
    errors,
    isDirty,
    isValid,
    showPassword,
    setShowPassword,
    submit,
    isPending,
    isError,
    router,
    errorMessage
  };
}
