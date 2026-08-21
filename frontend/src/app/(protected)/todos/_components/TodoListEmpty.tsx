'use client';

import { Button, Icon } from '@squantumengine/horizon';

interface TodoListEmptyProps {
  onCreate: () => void;
}

export function TodoListEmpty({ onCreate }: TodoListEmptyProps) {
  return (
    <div className="empty">
      <div className="empty__icon">
        <Icon name="clipboard-list-check" />
      </div>
      <h3>No todos yet</h3>
      <p>
        Create your first task to get started. Title is required; description and due date are
        optional.
      </p>
      <Button variant="primary" onClick={onCreate}>
        <Icon name="plus" />
        New todo
      </Button>
    </div>
  );
}
