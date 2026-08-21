'use client';

import { useTodo, type TodoResponse } from '@/services/todos';

export function useTodoDetailPage(id: string) {
  const query = useTodo(id);
  return {
    todo: (query.data?.data ?? null) as TodoResponse | null,
    isLoading: query.isLoading,
    error: query.error
  };
}
