import { fromOpenApi } from '@mswjs/source/open-api';
import type { RequestHandler } from 'msw';
import spec from '@/openapi/openapi.json';

export async function createHandlersFromOpenApi(): Promise<RequestHandler[]> {
  try {
    const specString = JSON.stringify(spec);
    const generatedHandlers = await fromOpenApi(specString);
    return generatedHandlers;
  } catch (error) {
    console.error('Error loading OpenAPI spec:', error);
    return [];
  }
}

export const handlersPromise = createHandlersFromOpenApi();

export const handlers: RequestHandler[] = [];
