import { render, screen } from '@testing-library/react';

import ObservabilityProvider from './provider';

describe('noop ObservabilityProvider — observability noop', () => {
  it('renders children unchanged (passthrough)', () => {
    render(
      <ObservabilityProvider>
        <span>child-content</span>
      </ObservabilityProvider>
    );
    expect(screen.getByText('child-content')).toBeInTheDocument();
  });
});
