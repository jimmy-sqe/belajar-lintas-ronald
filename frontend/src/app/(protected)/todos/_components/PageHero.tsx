'use client';

import { Button, Icon } from '@squantumengine/horizon';
import type { TodoStats } from '../_hooks/useTodoStats';
import { formatTodoStatsSub } from '../_hooks/useTodoStats';

interface PageHeroProps {
  stats: TodoStats;
  onCreate: () => void;
}

export function PageHero({ stats, onCreate }: PageHeroProps) {
  return (
    <div className="ph">
      <div>
        <div className="ph__eyebrow">Your list</div>
        <h1>Todos</h1>
        <div className="ph__sub">{formatTodoStatsSub(stats)}</div>
      </div>
      <div className="ph__actions">
        <Button variant="secondary" size="sm" disabled title="Sorting coming soon">
          <Icon name="arrow-down-wide-short" />
          Sort: Due date
        </Button>
        <Button variant="primary" size="sm" onClick={onCreate}>
          <Icon name="plus" />
          New todo
        </Button>
      </div>
    </div>
  );
}
