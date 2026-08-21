import { z } from 'zod';

export const todoFormSchema = z.object({
  title: z.string().min(1, 'Title is required').max(200, 'Title too long (max 200)'),
  description: z
    .string()
    .max(2000, 'Description too long (max 2000)')
    .optional()
    .or(z.literal('')),
  due_date: z.string().optional().or(z.literal(''))
});

export type TodoFormValues = z.infer<typeof todoFormSchema>;

/**
 * Normalize form values for API submission:
 *  - empty strings become undefined (so optional fields aren't sent as "")
 *  - due_date (datetime-local "YYYY-MM-DDTHH:mm") becomes ISO UTC string
 */
export function normalizeTodoFormValues(values: TodoFormValues): {
  title: string;
  description?: string;
  due_date?: string;
} {
  const description = values.description?.trim() ? values.description.trim() : undefined;
  let due_date: string | undefined;
  if (values.due_date && values.due_date.trim()) {
    const parsed = new Date(values.due_date);
    due_date = isNaN(parsed.getTime()) ? values.due_date : parsed.toISOString();
  }
  return { title: values.title.trim(), description, due_date };
}

/**
 * Convert an ISO due_date string back to the datetime-local format
 * "YYYY-MM-DDTHH:mm" so it can pre-fill the edit form.
 */
export function toDateTimeLocal(iso?: string): string {
  if (!iso) return '';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}
