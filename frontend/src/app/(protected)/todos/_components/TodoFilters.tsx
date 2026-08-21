'use client';

import { Chip } from '@squantumengine/horizon';

interface TodoFiltersProps {
  totalCount: number;
}

/**
 * Filter chip row per spec §7.4. Only "All" is functional (active);
 * Today/Overdue/Done are visual stubs disabled until BE gains a
 * `completed` field or due-date filter support.
 *
 * Horizon Chip API requires a string `label`, so the count is embedded
 * directly in the label text. Disabled chips wrap in a span so we can
 * apply opacity/cursor styling without relying on a `style` prop the
 * Chip API doesn't expose.
 */
export function TodoFilters({ totalCount }: TodoFiltersProps) {
  return (
    <div className="filters">
      <Chip label={`All ${totalCount}`} size="sm" isActive />
      <span className="filters__disabled-wrap">
        <Chip label="Today" size="sm" isDisable />
      </span>
      <span className="filters__disabled-wrap">
        <Chip label="Overdue" size="sm" isDisable />
      </span>
      <span className="filters__disabled-wrap">
        <Chip label="Done" size="sm" isDisable />
      </span>
    </div>
  );
}
