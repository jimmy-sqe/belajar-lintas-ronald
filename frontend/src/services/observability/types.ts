/**
 * Shared contract for the observability subsystem. Each option provides
 * two surfaces:
 *   - a client `ObservabilityProvider` (mounted in app/layout.tsx)
 *   - a server `register()` (invoked from src/instrumentation.ts)
 * Unsupported surfaces are explicit no-ops. Selection is structural at
 * scaffold time: the barrels (index.ts / server.ts) are repointed by the
 * pruner and unselected option folders are deleted.
 */
import type { ComponentType, ReactNode } from 'react';

export type ObservabilityProviderProps = { children: ReactNode };
export type ObservabilityProviderComponent = ComponentType<ObservabilityProviderProps>;
export type ObservabilityRegister = () => void | Promise<void>;
