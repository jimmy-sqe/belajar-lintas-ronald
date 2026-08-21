'use client';

import { useToaster } from '@squantumengine/horizon';
import { CreateTodoModal } from './_components/CreateTodoModal';
import { DeleteTodoConfirm } from './_components/DeleteTodoConfirm';
import { EditTodoModal } from './_components/EditTodoModal';
import { PageHero } from './_components/PageHero';
import { TodoFilters } from './_components/TodoFilters';
import { TodoListEmpty } from './_components/TodoListEmpty';
import { TodoRow } from './_components/TodoRow';
import { useTodosPage } from './_hooks/useTodosPage';
import { useTodoStats } from './_hooks/useTodoStats';

export default function TodosPage() {
  const p = useTodosPage();
  const stats = useTodoStats(p.todos);
  const toaster = useToaster();
  const toast = (label: string, variant: 'primary' | 'secondary' = 'primary') => {
    // Horizon OpenToaster union forces buttonLabel/onActionClick to never
    // when omitted; cast keeps the call site ergonomic.
    toaster.open({ id: `todo-${Date.now()}`, label, variant } as Parameters<
      typeof toaster.open
    >[0]);
  };

  if (p.error) {
    return (
      <div className="empty">
        <h3>Failed to load todos</h3>
        <p>{(p.error as Error)?.message ?? 'Please try again.'}</p>
      </div>
    );
  }

  return (
    <>
      <PageHero stats={stats} onCreate={() => p.setCreateOpen(true)} />
      <TodoFilters totalCount={stats.total} />

      {p.isLoading ? (
        <div className="empty">
          <p>Loading…</p>
        </div>
      ) : p.todos.length === 0 ? (
        <TodoListEmpty onCreate={() => p.setCreateOpen(true)} />
      ) : (
        <ul className="todo-list">
          {p.todos.map((todo) => (
            <TodoRow
              key={todo.id}
              todo={todo}
              onEdit={p.setEditingTodo}
              onDelete={p.setDeletingTodo}
            />
          ))}
        </ul>
      )}

      {p.createOpen && (
        <CreateTodoModal
          onClose={() => p.setCreateOpen(false)}
          isSubmitting={p.createMutation.isPending}
          onSubmit={(data) =>
            p.createMutation.mutate(
              { data },
              {
                onSuccess: () => {
                  p.setCreateOpen(false);
                  toast('Todo created');
                  p.refetch();
                },
                onError: () => toast('Failed to create todo', 'secondary')
              }
            )
          }
        />
      )}

      {p.editingTodo && (
        <EditTodoModal
          todo={p.editingTodo}
          onClose={() => p.setEditingTodo(null)}
          isSubmitting={p.updateMutation.isPending}
          onSubmit={(data) =>
            p.updateMutation.mutate(
              { id: p.editingTodo!.id, data },
              {
                onSuccess: () => {
                  p.setEditingTodo(null);
                  toast('Todo updated');
                  p.refetch();
                },
                onError: () => toast('Failed to update todo', 'secondary')
              }
            )
          }
        />
      )}

      {p.deletingTodo && (
        <DeleteTodoConfirm
          todo={p.deletingTodo}
          onClose={() => p.setDeletingTodo(null)}
          isDeleting={p.deleteMutation.isPending}
          onConfirm={() =>
            p.deleteMutation.mutate(
              { id: p.deletingTodo!.id },
              {
                onSuccess: () => {
                  p.setDeletingTodo(null);
                  toast('Todo deleted');
                  p.refetch();
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
