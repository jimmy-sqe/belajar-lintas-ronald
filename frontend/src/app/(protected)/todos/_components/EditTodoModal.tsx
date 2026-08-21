'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Button,
  Dialog,
  DialogBody,
  DialogFooter,
  DialogHeader,
  FormTextField,
  Icon
} from '@squantumengine/horizon';
import {
  todoFormSchema,
  normalizeTodoFormValues,
  toDateTimeLocal,
  type TodoFormValues
} from './schema';
import type { TodoResponse, UpdateTodoRequest } from '@/services/todos';

interface EditTodoModalProps {
  todo: TodoResponse;
  onClose: () => void;
  onSubmit: (data: UpdateTodoRequest) => void;
  isSubmitting: boolean;
}

const MAX_DESC = 2000;

export function EditTodoModal({ todo, onClose, onSubmit, isSubmitting }: EditTodoModalProps) {
  const [error, setError] = useState<string | null>(null);
  const {
    control,
    handleSubmit,
    watch,
    formState: { isValid }
  } = useForm<TodoFormValues>({
    mode: 'onChange',
    resolver: zodResolver(todoFormSchema),
    defaultValues: {
      title: todo.title,
      description: todo.description ?? '',
      due_date: toDateTimeLocal(todo.due_date)
    }
  });

  const descLen = watch('description')?.length ?? 0;

  const submit = handleSubmit((values) => {
    setError(null);
    try {
      onSubmit(normalizeTodoFormValues(values));
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to submit');
    }
  });

  return (
    <Dialog open onClose={onClose}>
      <DialogHeader>
        <div>
          <div className="modal__title">Edit todo</div>
          <div className="modal__sub">Update the fields you want to change.</div>
        </div>
      </DialogHeader>
      <form onSubmit={submit}>
        <DialogBody>
          {error && (
            <div className="inline-alert" role="alert">
              <Icon name="exclamation-circle" />
              <div>{error}</div>
            </div>
          )}
          <div className="fg">
            <label className="fg__label">
              Title <span className="fg__label-req">*</span>
            </label>
            <FormTextField full controller={{ control, name: 'title' }} />
          </div>
          <div className="fg">
            <label className="fg__label">
              Description <span className="fg__optional">(optional)</span>
            </label>
            <FormTextField
              full
              multiline
              rows={4}
              controller={{ control, name: 'description' }}
            />
            <div className={'char-counter' + (descLen > 1800 ? ' is-warn' : '')}>
              {descLen}/{MAX_DESC}
            </div>
          </div>
          <div className="fg">
            <label className="fg__label">
              Due date <span className="fg__optional">(optional)</span>
            </label>
            <FormTextField full controller={{ control, name: 'due_date' }} type="datetime-local" />
          </div>
        </DialogBody>
        <DialogFooter>
          <Button variant="secondary" type="button" onClick={onClose}>
            Cancel
          </Button>
          <Button variant="primary" type="submit" disabled={!isValid} loading={isSubmitting}>
            Save changes
          </Button>
        </DialogFooter>
      </form>
    </Dialog>
  );
}
