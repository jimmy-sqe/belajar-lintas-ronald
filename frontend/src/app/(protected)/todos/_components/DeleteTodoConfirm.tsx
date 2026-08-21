'use client';

import { useState } from 'react';
import {
  Button,
  Dialog,
  DialogBody,
  DialogFooter,
  DialogHeader,
  Icon
} from '@squantumengine/horizon';
import type { TodoResponse } from '@/services/todos';

interface DeleteTodoConfirmProps {
  todo: TodoResponse;
  onClose: () => void;
  onConfirm: () => void;
  isDeleting: boolean;
}

export function DeleteTodoConfirm({ todo, onClose, onConfirm, isDeleting }: DeleteTodoConfirmProps) {
  const [error, setError] = useState<string | null>(null);

  const confirm = () => {
    setError(null);
    try {
      onConfirm();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to delete');
    }
  };

  return (
    <Dialog open onClose={onClose}>
      <DialogHeader>
        <div className="confirm__icon-wrap">
          <Icon name="trash" />
        </div>
        <div>
          <div className="modal__title">Delete todo?</div>
          <div className="modal__sub">This action cannot be undone.</div>
        </div>
      </DialogHeader>
      <DialogBody>
        {error && (
          <div className="inline-alert" role="alert">
            <Icon name="exclamation-circle" />
            <div>{error}</div>
          </div>
        )}
        <p>
          Are you sure you want to delete <b>{todo.title}</b>?
        </p>
      </DialogBody>
      <DialogFooter>
        <Button variant="secondary" type="button" onClick={onClose}>
          Cancel
        </Button>
        <Button variant="primary" type="button" onClick={confirm} loading={isDeleting}>
          Delete
        </Button>
      </DialogFooter>
    </Dialog>
  );
}
