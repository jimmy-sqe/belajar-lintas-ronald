'use client';

import { Dialog, Spinner } from '@squantumengine/horizon';
import { useLoadingStore } from '@/store/loadingStore';

const LoadingDialog = () => {
  const isLoading = useLoadingStore((state) => state.isLoading);
  return (
    <Dialog open={isLoading} onClose={() => {}} hideCloseBtn className="rounded-lg p-6">
      <Spinner size="lg" />
    </Dialog>
  );
};

export default LoadingDialog;
