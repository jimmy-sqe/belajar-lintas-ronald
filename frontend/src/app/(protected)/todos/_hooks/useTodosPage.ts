'use client';

import { useState } from 'react';
import {
  useCreateTodo,
  useDeleteTodo,
  useTodos,
  useUpdateTodo,
  type TodoResponse
} from '@/services/todos';

const PAGE_SIZE = 10;

export function useTodosPage() {
  const [page, setPage] = useState(1);
  const [createOpen, setCreateOpen] = useState(false);
  const [editingTodo, setEditingTodo] = useState<TodoResponse | null>(null);
  const [deletingTodo, setDeletingTodo] = useState<TodoResponse | null>(null);

  const query = useTodos({ page, page_size: PAGE_SIZE });
  const createMutation = useCreateTodo();
  const updateMutation = useUpdateTodo();
  const deleteMutation = useDeleteTodo();

  const todos = (query.data?.data ?? []) as TodoResponse[];
  const pagination = query.data?.pagination;

  return {
    todos,
    pagination,
    isLoading: query.isLoading,
    error: query.error,
    page,
    setPage,
    createOpen,
    setCreateOpen,
    editingTodo,
    setEditingTodo,
    deletingTodo,
    setDeletingTodo,
    createMutation,
    updateMutation,
    deleteMutation,
    refetch: query.refetch
  };
}
