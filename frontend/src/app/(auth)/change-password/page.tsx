'use client';

import { Button, FormTextField, Icon, Info, Token } from '@squantumengine/horizon';
import { errorResponseHandler } from '@/utils/error';
import useChangePasswordPage from './_hooks/useChangePasswordPage';

export default function ChangePasswordPage() {
  const {
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
  } = useChangePasswordPage();

  return (
    <div className="mt-4 flex flex-col items-center">
      <form className="flex flex-col items-center gap-2" onSubmit={submit}>
        {isError && (
          <Info
            className="mb-5 w-full"
            type="error"
            variant="simple"
            title=""
            description={
              <>
                {errorResponseHandler(error).message.map((msg, idx) => (
                  <p key={`error-msg-${idx}`} className="text-red-600">
                    {msg}
                  </p>
                ))}
              </>
            }
          />
        )}
        <div className="mb-10 flex w-full flex-col items-start gap-6">
          <div className="flex w-full flex-col items-start gap-2">
            <div className="text-base font-bold">Password lama</div>
            <FormTextField
              full
              controller={{ control, name: 'old_password' }}
              type={showOldPassword ? 'text' : 'password'}
              placeholder="Masukkan password lama"
              isInvalid={!!errors.old_password?.message}
              errorMessage={errors.old_password?.message as string}
              suffix={
                <div
                  className="cursor-pointer"
                  onClick={() => setShowOldPassword(!showOldPassword)}>
                  <Icon
                    name={showOldPassword ? 'eye-slash' : 'eye'}
                    color={Token.COLORS.neutral[950]}
                  />
                </div>
              }
            />
          </div>
          <div className="flex w-full flex-col items-start gap-2">
            <div className="text-base font-bold">Password baru</div>
            <FormTextField
              full
              controller={{
                control,
                name: 'new_password'
              }}
              type={showNewPassword ? 'text' : 'password'}
              placeholder="Masukkan password baru"
              isInvalid={!!errors.new_password?.message}
              errorMessage={errors.new_password?.message as string}
              suffix={
                <div
                  className="cursor-pointer"
                  onClick={() => setShowNewPassword(!showNewPassword)}>
                  <Icon
                    name={showNewPassword ? 'eye-slash' : 'eye'}
                    color={Token.COLORS.neutral[950]}
                  />
                </div>
              }
            />
          </div>
          <div className="flex w-full flex-col items-start gap-2">
            <div className="text-base font-bold">Ulangi password baru</div>
            <FormTextField
              full
              controller={{
                control,
                name: 'confirmation_password'
              }}
              type={showConfirmationPassword ? 'text' : 'password'}
              placeholder="Konfirmasi password"
              isInvalid={!!errors.confirmation_password?.message}
              errorMessage={errors.confirmation_password?.message as string}
              suffix={
                <div
                  className="cursor-pointer"
                  onClick={() => setShowConfirmationPassword(!showConfirmationPassword)}>
                  <Icon
                    name={showConfirmationPassword ? 'eye-slash' : 'eye'}
                    color={Token.COLORS.neutral[950]}
                  />
                </div>
              }
            />
          </div>
          <div className="flex flex-col items-start gap-3">
            <div className="text-hz-neutral-600">Password harus memiliki: </div>
            <div className="grid grid-cols-3 gap-x-6 gap-y-2 text-sm">
              {passwordRules.map((rule, idx) => (
                <div
                  key={`rule-${idx}`}
                  className={`${rule.passed ? 'text-primary-green' : 'text-hz-neutral-600'} flex flex-row items-center`}>
                  <Icon
                    color={rule.passed ? 'rgba(59, 147, 21, 1)' : rule.color}
                    name={rule.icon}
                  />
                  {rule.text}
                </div>
              ))}
            </div>
          </div>
        </div>
        <Button
          size="lg"
          full
          disabled={!isDirty || !isValid || passwordRules.some((rule) => !rule.passed)}
          variant="primary"
          type="submit"
          loading={isPending}>
          Ubah Password
        </Button>
      </form>
    </div>
  );
}
