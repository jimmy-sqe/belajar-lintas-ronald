import { Paragraph, Title } from '@squantumengine/horizon';
import Link from 'next/link';

const NotFoundPage = () => {
  return (
    <div className="flex h-screen w-screen flex-col items-start bg-[#F6F6F6]">
      <div className="flex h-full w-full items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <Title level={1}>404</Title>
          <div className="flex flex-col items-center gap-2">
            <Title level={2}>Halaman tidak ditemukan</Title>
            <Paragraph size="xl">
              Halaman yang Anda cari tidak ditemukan. Coba periksa kembali URL
            </Paragraph>
          </div>
          <Link href="/" className="text-blue-600 underline">
            Kembali ke beranda
          </Link>
        </div>
      </div>
    </div>
  );
};

export default NotFoundPage;
