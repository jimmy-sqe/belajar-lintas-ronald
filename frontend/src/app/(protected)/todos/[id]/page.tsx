'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { Button, Icon, useToaster } from '@squantumengine/horizon';
import { useTodoDetailPage } from './_hooks/useTodoDetailPage';
import { EditTodoModal } from '../_components/EditTodoModal';
import { DeleteTodoConfirm } from '../_components/DeleteTodoConfirm';
import { useUpdateTodo, useDeleteTodo, type TodoResponse } from '@/services/todos';

export default function TodoDetailPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { todo, isLoading, error } = useTodoDetailPage(params.id);
  const [editing, setEditing] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const updateMutation = useUpdateTodo();
  const deleteMutation = useDeleteTodo();
  const toaster = useToaster();
  const toast = (label: string, variant: 'primary' | 'secondary' = 'primary') => {
    toaster.open({ id: `todo-${Date.now()}`, label, variant } as Parameters<
      typeof toaster.open
    >[0]);
  };

  if (isLoading) {
    return (
      <div className="empty">
        <p>Loading…</p>
      </div>
    );
  }

  if (error || !todo) {
    return (
      <div className="empty">
        <div className="empty__icon">
          <Icon name="exclamation-circle" />
        </div>
        <h3>Todo not found</h3>
        <p>The todo you’re looking for may have been deleted.</p>
        <Link href="/todos">
          <Button variant="primary">Back to todos</Button>
        </Link>
      </div>
    );
  }

  const fmt = (iso: string) =>
    new Date(iso).toLocaleString('id-ID', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });

  return (
    <>
      <Link href="/todos" className="back-link">
        <Icon name="arrow-left" />
        Back to todos
      </Link>
      <article className="detail">
        <header className="detail__header">
          <h1 className="detail__title">{todo.title}</h1>
          {todo.due_date && (
            <div className="detail__chip-row">
              <span className="todo-row__due">
                <Icon name="calendar" />
                Due {fmt(todo.due_date)}
              </span>
            </div>
          )}
        </header>
        <div className="detail__body">
          <section>
            <div className="detail__section-label">Description</div>
            {todo.description ? (
              <div className="detail__desc">{todo.description}</div>
            ) : (
              <div className="detail__desc is-empty">No description</div>
            )}
          </section>
          <section>
            <div className="detail__section-label">Meta</div>
            <div className="detail__meta-grid">
              <div className="detail__meta-item">
                <div className="detail__meta-label">ID</div>
                <div className="detail__meta-value">{todo.id}</div>
              </div>
              <div className="detail__meta-item">
                <div className="detail__meta-label">Created</div>
                <div className="detail__meta-value">{fmt(todo.created_at)}</div>
              </div>
              <div className="detail__meta-item">
                <div className="detail__meta-label">Created by</div>
                <div className="detail__meta-value">{todo.created_by}</div>
              </div>
              <div className="detail__meta-item">
                <div className="detail__meta-label">Modified</div>
                <div className="detail__meta-value">{fmt(todo.modified_at)}</div>
              </div>
              <div className="detail__meta-item">
                <div className="detail__meta-label">Modified by</div>
                <div className="detail__meta-value">{todo.modified_by}</div>
              </div>
            </div>
          </section>
        </div>
        <footer className="detail__footer">
          <Button variant="secondary" onClick={() => setEditing(true)}>
            <Icon name="pen" />
            Edit
          </Button>
          <Button variant="primary" onClick={() => setDeleting(true)}>
            <Icon name="trash" />
            Delete
          </Button>
        </footer>
      </article>

      {editing && (
        <EditTodoModal
          todo={todo as TodoResponse}
          onClose={() => setEditing(false)}
          isSubmitting={updateMutation.isPending}
          onSubmit={(data) =>
            updateMutation.mutate(
              { id: todo.id, data },
              {
                onSuccess: () => {
                  setEditing(false);
                  toast('Todo updated');
                },
                onError: () => toast('Failed to update todo', 'secondary')
              }
            )
          }
        />
      )}
      {deleting && (
        <DeleteTodoConfirm
          todo={todo as TodoResponse}
          onClose={() => setDeleting(false)}
          isDeleting={deleteMutation.isPending}
          onConfirm={() =>
            deleteMutation.mutate(
              { id: todo.id },
              {
                onSuccess: () => {
                  toast('Todo deleted');
                  router.push('/todos');
                },
                onError: () => toast('Failed to delete todo', 'secondary')
              }
            )
          }
        />
      )}
    </>
  );
}
