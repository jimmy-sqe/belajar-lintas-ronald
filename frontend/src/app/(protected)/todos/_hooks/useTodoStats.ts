'use client';

import { useMemo } from 'react';
import type { TodoResponse } from '@/services/todos';

export interface TodoStats {
  total: number;
  /** Count of todos with due_date < now */
  overdue: number;
  /** Count of todos with due_date in [today 00:00, tomorrow 00:00) */
  dueToday: number;
}

/**
 * Derive total / overdue / due-today counts from a todos array.
 *
 * Pure, memoized — does not fetch. Caller passes already-fetched todos
 * (typically from useTodos in useTodosPage).
 */
export function useTodoStats(todos: TodoResponse[]): TodoStats {
  return useMemo(() => {
    const now = new Date();
    const startOfToday = new Date(now);
    startOfToday.setHours(0, 0, 0, 0);
    const startOfTomorrow = new Date(startOfToday);
    startOfTomorrow.setDate(startOfTomorrow.getDate() + 1);

    let overdue = 0;
    let dueToday = 0;
    for (const t of todos) {
      if (!t.due_date) continue;
      const d = new Date(t.due_date);
      if (isNaN(d.getTime())) continue;
      if (d < now) overdue++;
      else if (d >= startOfToday && d < startOfTomorrow) dueToday++;
    }
    return { total: todos.length, overdue, dueToday };
  }, [todos]);
}

/**
 * Format stats into the page-hero subtitle string per spec §7.3.
 *  - Always start with "<total> tasks"
 *  - Append "· <overdue> overdue" if overdue > 0
 *  - Append "· <dueToday> due today" if dueToday > 0
 */
export function formatTodoStatsSub(stats: TodoStats): string {
  const parts: string[] = [`${stats.total} task${stats.total === 1 ? '' : 's'}`];
  if (stats.overdue > 0) parts.push(`${stats.overdue} overdue`);
  if (stats.dueToday > 0) parts.push(`${stats.dueToday} due today`);
  return parts.join(' · ');
}
