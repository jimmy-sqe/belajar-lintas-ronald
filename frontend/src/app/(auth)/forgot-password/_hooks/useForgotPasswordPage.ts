import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useToasterProvider } from '@squantumengine/horizon';
import { z } from 'zod';
import auth from '@/services/auth';

const forgotPasswordSchema = z.object({
  username: z.string().min(1, 'Username wajib diisi')
});

type ForgotPasswordFormValues = z.infer<typeof forgotPasswordSchema>;

export default function useForgotPasswordPage() {
  const {
    control,
    handleSubmit,
    formState: { errors, isDirty, isValid }
  } = useForm<ForgotPasswordFormValues>({
    mode: 'onChange',
    resolver: zodResolver(forgotPasswordSchema)
  });
  const [isPending, setIsPending] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  const { open } = useToasterProvider();

  const submit = handleSubmit(async (resetPasswordData: ForgotPasswordFormValues) => {
    setIsPending(true);
    setIsError(false);
    setError(null);
    try {
      if (!auth.resetPassword) {
        throw new Error('Password reset is not supported by the active auth model.');
      }
      await auth.resetPassword({ email: resetPasswordData.username });
      open({
        id: 'forgot-password-success',
        label: 'Password sudah berhasil direset, silahkan cek email anda.'
      });
    } catch (err) {
      setIsError(true);
      setError(err instanceof Error ? err : new Error('Reset password failed'));
    } finally {
      setIsPending(false);
    }
  });

  return {
    control,
    errors,
    isDirty,
    isValid,
    submit,
    isError,
    error,
    isPending
  };
}
