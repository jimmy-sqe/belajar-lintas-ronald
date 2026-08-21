import { redirect } from 'next/navigation';
import { defaultBaseUrlPage } from '@/static/route';

export default async function Home() {
  redirect(defaultBaseUrlPage);
}
