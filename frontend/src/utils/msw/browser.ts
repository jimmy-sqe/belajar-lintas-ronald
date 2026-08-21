import { setupWorker } from 'msw/browser';
import { handlers as mockHandlers } from '@/mocks/handlers';
import { handlersPromise } from '@/utils/msw/handlers';

// Create worker with both manual and OpenAPI-generated handlers
export const worker = setupWorker(...mockHandlers);

// Initialize worker with OpenAPI generated handlers
export async function initializeWorker() {
  try {
    const openApiHandlers = await handlersPromise;
    // Use both manual handlers (priority) and OpenAPI generated handlers
    worker.use(...mockHandlers, ...openApiHandlers);
    console.log('MSW worker initialized with manual and OpenAPI handlers');
  } catch (error) {
    console.error('Failed to initialize MSW worker with OpenAPI handlers:', error);
    // Continue with just manual handlers if OpenAPI generation fails
    console.log('MSW worker running with manual handlers only');
  }
}
