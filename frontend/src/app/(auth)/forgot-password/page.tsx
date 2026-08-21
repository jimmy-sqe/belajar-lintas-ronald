'use client';

import { Button, FormTextField, Info } from '@squantumengine/horizon';
import { errorResponseHandler } from '@/utils/error';
import useForgotPasswordPage from './_hooks/useForgotPasswordPage';

export default function ForgotPasswordPage() {
  const { control, errors, isDirty, isValid, submit, isError, error, isPending } =
    useForgotPasswordPage();

  return (
    <div className="mt-4 flex flex-col items-center">
      <form className="mt-5 flex flex-col items-center gap-2" onSubmit={submit}>
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
            <div className="text-base font-bold">Username</div>
            <FormTextField
              full
              controller={{ control, name: 'username' }}
              type="text"
              placeholder="Masukkan username atau email Anda"
              isInvalid={!!errors.username?.message}
              errorMessage={errors.username?.message as string}
            />
          </div>
        </div>
        <Button
          size="lg"
          full
          disabled={!isDirty || !isValid}
          variant="primary"
          type="submit"
          loading={isPending}>
          Reset Password
        </Button>
      </form>
    </div>
  );
}
