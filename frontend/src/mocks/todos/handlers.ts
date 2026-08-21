import { http, HttpResponse } from 'msw';

interface MockTodo {
  id: string;
  title: string;
  description?: string;
  due_date?: string;
  created_at: string;
  created_by: string;
  modified_at: string;
  modified_by: string;
}

const DEMO_USER = '11111111-1111-1111-1111-111111111111';

let todos: MockTodo[] = [
  {
    id: 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
    title: 'Buy groceries',
    description: 'Milk, bread, eggs',
    due_date: '2026-05-20T10:00:00Z',
    created_at: '2026-05-14T00:00:00Z',
    created_by: DEMO_USER,
    modified_at: '2026-05-14T00:00:00Z',
    modified_by: DEMO_USER
  },
  {
    id: 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
    title: 'Walk the dog',
    description: 'Morning walk in the park',
    created_at: '2026-05-14T00:00:00Z',
    created_by: DEMO_USER,
    modified_at: '2026-05-14T00:00:00Z',
    modified_by: DEMO_USER
  },
  {
    id: 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3',
    title: 'Finish project report',
    due_date: '2026-05-25T17:00:00Z',
    created_at: '2026-05-14T00:00:00Z',
    created_by: DEMO_USER,
    modified_at: '2026-05-14T00:00:00Z',
    modified_by: DEMO_USER
  }
];

export const todosHandlers = [
  http.get('/v1/todos', () =>
    HttpResponse.json({
      success: true,
      code: 200,
      message: 'todos fetched',
      data: todos,
      pagination: { page: 1, page_size: 10, total_data: todos.length, total_page: 1 },
      timestamp: new Date().toISOString()
    })
  ),
  http.get('/v1/todos/:id', ({ params }) => {
    const t = todos.find((x) => x.id === params.id);
    if (!t) {
      return HttpResponse.json(
        {
          success: false,
          code: 40400,
          message: 'todo not found',
          timestamp: new Date().toISOString()
        },
        { status: 404 }
      );
    }
    return HttpResponse.json({
      success: true,
      code: 200,
      message: 'todo fetched',
      data: t,
      timestamp: new Date().toISOString()
    });
  }),
  http.post('/v1/todos', async ({ request }) => {
    const body = (await request.json()) as Partial<MockTodo>;
    const now = new Date().toISOString();
    const newTodo: MockTodo = {
      id: crypto.randomUUID(),
      title: body.title ?? '',
      description: body.description,
      due_date: body.due_date,
      created_at: now,
      created_by: DEMO_USER,
      modified_at: now,
      modified_by: DEMO_USER
    };
    todos.push(newTodo);
    return HttpResponse.json(
      {
        success: true,
        code: 201,
        message: 'todo created',
        data: newTodo,
        timestamp: now
      },
      { status: 201 }
    );
  }),
  http.put('/v1/todos/:id', async ({ params, request }) => {
    const body = (await request.json()) as Partial<MockTodo>;
    const idx = todos.findIndex((x) => x.id === params.id);
    if (idx === -1) {
      return HttpResponse.json(
        {
          success: false,
          code: 40400,
          message: 'todo not found',
          timestamp: new Date().toISOString()
        },
        { status: 404 }
      );
    }
    todos[idx] = {
      ...todos[idx],
      ...body,
      modified_at: new Date().toISOString(),
      modified_by: DEMO_USER
    };
    return HttpResponse.json({
      success: true,
      code: 200,
      message: 'todo updated',
      data: todos[idx],
      timestamp: new Date().toISOString()
    });
  }),
  http.delete('/v1/todos/:id', ({ params }) => {
    const idx = todos.findIndex((x) => x.id === params.id);
    if (idx === -1) {
      return HttpResponse.json(
        {
          success: false,
          code: 40400,
          message: 'todo not found',
          timestamp: new Date().toISOString()
        },
        { status: 404 }
      );
    }
    todos.splice(idx, 1);
    return HttpResponse.json({
      success: true,
      code: 200,
      message: 'todo deleted',
      data: {},
      timestamp: new Date().toISOString()
    });
  })
];
