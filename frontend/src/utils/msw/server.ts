import { setupServer } from 'msw/node';
import { handlers as mockHandlers } from '@/mocks/handlers';
import { handlersPromise } from '@/utils/msw/handlers';

export const server = setupServer();

export async function initializeServer() {
  try {
    const openApiHandlers = await handlersPromise;
    server.use(...openApiHandlers, ...mockHandlers);
    console.log('MSW server initialized with OpenAPI and mock handlers');
  } catch (error) {
    console.error('Failed to initialize MSW server with OpenAPI handlers:', error);
  }
}
