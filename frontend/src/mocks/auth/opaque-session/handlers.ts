import { http, HttpResponse } from 'msw';

export const authHandlers = [
  http.post('/auth/sign-in', async ({ request }) => {
    const body = (await request.json()) as { username: string; password: string };
    if (body.username === 'demo' && body.password === 'demo') {
      return HttpResponse.json({
        success: true,
        code: 200,
        data: {
          session_token: 'mock-opaque-session-token',
          user: { id: 'u-1', name: 'Demo User', email: 'demo@example.com' },
          permissions: ['read', 'write'],
          requires_password_change: false
        },
        timestamp: new Date().toISOString()
      });
    }
    return HttpResponse.json(
      {
        success: false,
        code: 40100,
        message: 'Invalid credentials',
        timestamp: new Date().toISOString()
      },
      { status: 401 }
    );
  }),
  http.post('/auth/sign-out', () =>
    HttpResponse.json({ success: true, code: 200, data: {}, timestamp: new Date().toISOString() })
  ),
  http.post('/auth/change-password', () =>
    HttpResponse.json({
      success: true,
      code: 200,
      data: { changed_at: new Date().toISOString() },
      timestamp: new Date().toISOString()
    })
  )
];
