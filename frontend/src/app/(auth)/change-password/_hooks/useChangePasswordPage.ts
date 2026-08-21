import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useRouter } from 'next/navigation';
import { z } from 'zod';
import auth from '@/services/auth';
import type { IPasswordRule } from '@/static/types/auth';
import { getInitialPassRules, validateRules } from '@/utils/auth';

const changePasswordSchema = z.object({
  old_password: z.string().min(8, 'Password lama minimal 8 karakter'),
  new_password: z.string().min(8, 'Password baru minimal 8 karakter'),
  confirmation_password: z.string().min(8, 'Konfirmasi password minimal 8 karakter'),
  otp: z.string().optional()
});

type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>;

export default function useChangePasswordPage() {
  const router = useRouter();
  const [showOldPassword, setShowOldPassword] = useState<boolean>(false);
  const [showNewPassword, setShowNewPassword] = useState<boolean>(false);
  const [showConfirmationPassword, setShowConfirmationPassword] = useState<boolean>(false);
  const [isPending, setIsPending] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);

  const {
    control,
    handleSubmit,
    formState: { errors, isDirty, isValid },
    watch
  } = useForm<ChangePasswordFormValues>({
    mode: 'onChange',
    resolver: zodResolver(changePasswordSchema)
  });
  const newPasswordWatch = watch('new_password');

  const [passwordRules, setPasswordRules] = useState<IPasswordRule[]>(getInitialPassRules());

  const submit = handleSubmit(async (changePasswordData: ChangePasswordFormValues) => {
    setIsPending(true);
    setIsError(false);
    setError(null);
    try {
      if (!auth.changePassword) {
        throw new Error('Change password is not supported by the active auth model.');
      }
      await auth.changePassword({
        currentPassword: changePasswordData.old_password,
        newPassword: changePasswordData.new_password
      });
      router.push('/login');
    } catch (err) {
      setIsError(true);
      setError(err instanceof Error ? err : new Error('Change password failed'));
      console.error('error:', err);
    } finally {
      setIsPending(false);
    }
  });

  useEffect(() => {
    if (typeof newPasswordWatch === 'string') {
      // check each rule
      setPasswordRules((prevRules) => validateRules(prevRules, newPasswordWatch));
    }
  }, [newPasswordWatch]);

  return {
    control,
    isPending,
    errors,
    isDirty,
    isValid,
    submit,
    isError,
    error,
    showOldPassword,
    setShowOldPassword,
    showNewPassword,
    setShowNewPassword,
    showConfirmationPassword,
    setShowConfirmationPassword,
    passwordRules
  };
}
