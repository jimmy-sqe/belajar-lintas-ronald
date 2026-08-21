'use client';

import Link from 'next/link';
import { Icon } from '@squantumengine/horizon';
import type { TodoResponse } from '@/services/todos';

interface TodoRowProps {
  todo: TodoResponse;
  onEdit: (t: TodoResponse) => void;
  onDelete: (t: TodoResponse) => void;
}

type DueState = 'empty' | 'overdue' | 'soon' | 'normal';

function classifyDue(due_date?: string): { state: DueState; label: string; icon: string } {
  if (!due_date) return { state: 'empty', label: 'No due date', icon: 'calendar-xmark' };
  const d = new Date(due_date);
  if (isNaN(d.getTime())) return { state: 'empty', label: 'No due date', icon: 'calendar-xmark' };

  const now = new Date();
  const startOfToday = new Date(now);
  startOfToday.setHours(0, 0, 0, 0);
  const startOfTomorrow = new Date(startOfToday);
  startOfTomorrow.setDate(startOfTomorrow.getDate() + 1);

  const label = d.toLocaleDateString('id-ID', { month: 'short', day: 'numeric' });
  if (d < now) return { state: 'overdue', label, icon: 'clock-rotate-left' };
  if (d >= startOfToday && d < startOfTomorrow) return { state: 'soon', label, icon: 'clock' };
  return { state: 'normal', label, icon: 'calendar' };
}

export function TodoRow({ todo, onEdit, onDelete }: TodoRowProps) {
  const due = classifyDue(todo.due_date);
  const dueClass =
    due.state === 'overdue'
      ? 'todo-row__due is-overdue'
      : due.state === 'soon'
        ? 'todo-row__due is-soon'
        : due.state === 'empty'
          ? 'todo-row__due todo-row__due--empty'
          : 'todo-row__due';

  return (
    <li className="todo-row">
      <div>
        <Link href={`/todos/${todo.id}`} className="todo-row__title">
          {todo.title}
        </Link>
        {todo.description && <p className="todo-row__desc">{todo.description}</p>}
      </div>
      <div className={dueClass}>
        <Icon name={due.icon} />
        {due.label}
      </div>
      <div className="todo-row__actions">
        <button
          type="button"
          className="icon-btn"
          onClick={() => onEdit(todo)}
          aria-label="Edit"
          title="Edit">
          <Icon name="pen" />
        </button>
        <button
          type="button"
          className="icon-btn icon-btn--danger"
          onClick={() => onDelete(todo)}
          aria-label="Delete"
          title="Delete">
          <Icon name="trash" />
        </button>
      </div>
    </li>
  );
}
