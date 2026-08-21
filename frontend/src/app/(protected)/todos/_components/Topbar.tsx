'use client';

import { Fragment } from 'react';
import { Icon } from '@squantumengine/horizon';

interface TopbarProps {
  /** Breadcrumb segments. Last item is rendered as bold (current). */
  crumbs?: string[];
}

export function Topbar({ crumbs = ['Workspace', 'Todos'] }: TopbarProps) {
  return (
    <header className="tb">
      <div className="tb__crumbs">
        {crumbs.map((c, i) => (
          <Fragment key={`${i}-${c}`}>
            {i > 0 && <Icon name="chevron-right" />}
            {i === crumbs.length - 1 ? <b>{c}</b> : <span>{c}</span>}
          </Fragment>
        ))}
      </div>
      <div className="tb__right">
        <div className="search-field">
          <Icon name="search" />
          <input disabled placeholder="Search todos (coming soon)" />
          <span className="kbd-hint">⌘K</span>
        </div>
      </div>
    </header>
  );
}
