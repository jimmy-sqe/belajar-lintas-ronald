'use client';

import { AppShell } from './_components/AppShell';
import './_styles/todoapp.css';

export default function TodosLayout({ children }: { children: React.ReactNode }) {
  return <AppShell>{children}</AppShell>;
}
