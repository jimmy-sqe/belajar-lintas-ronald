import {
  QueryCache,
  QueryClient,
  defaultShouldDehydrateQuery,
  isServer
} from '@tanstack/react-query';

interface QueryClientOptions {
  onError?: (error: unknown) => void;
}

function makeQueryClient(options: QueryClientOptions = {}) {
  const { onError } = options;
  return new QueryClient({
    queryCache: new QueryCache({
      onError
    }),
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000
      },
      dehydrate: {
        // include pending queries in dehydration
        shouldDehydrateQuery: (query) =>
          defaultShouldDehydrateQuery(query) || query.state.status === 'pending'
      }
    }
  });
}

let browserQueryClient: QueryClient | undefined = undefined;

export function getQueryClient(options: QueryClientOptions = {}) {
  if (isServer) {
    // Server: always make a new query client
    return makeQueryClient(options);
  } else {
    // Browser: make a new query client if we don't already have one
    // This is very important, so we don't re-make a new client if React
    // suspends during the initial render. This may not be needed if we
    // have a suspense boundary BELOW the creation of the query client
    if (!browserQueryClient) browserQueryClient = makeQueryClient(options);
    return browserQueryClient;
  }
}
