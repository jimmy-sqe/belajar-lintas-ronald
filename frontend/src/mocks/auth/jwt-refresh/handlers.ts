import { http, HttpResponse } from 'msw';

export const authHandlers = [
  http.post('/v1/auth/login', async ({ request }) => {
    const body = (await request.json()) as { email: string; password: string };
    if (body.email === 'demo1@example.com' && body.password === 'Demo1234!') {
      return HttpResponse.json(
        {
          success: true,
          code: 201,
          message: 'logged in',
          data: {
            access_token: 'mock-access-token',
            refresh_token: 'mock-refresh-token',
            access_token_expires_in_sec: 3600,
            refresh_token_expires_in_sec: 86400,
            issued_at: Date.now(),
            user: { id: '11111111-1111-1111-1111-111111111111', name: 'Demo One', email: body.email },
            permissions: []
          },
          timestamp: new Date().toISOString()
        },
        { status: 201 }
      );
    }
    return HttpResponse.json(
      {
        success: false,
        code: 40100,
        message: 'invalid credentials',
        timestamp: new Date().toISOString()
      },
      { status: 401 }
    );
  }),
  http.post('/v1/auth/logout', () =>
    HttpResponse.json({
      success: true,
      code: 200,
      message: 'logged out',
      data: {},
      timestamp: new Date().toISOString()
    })
  ),
  http.post('/v1/auth/renew', () =>
    HttpResponse.json({
      success: true,
      code: 200,
      message: 'token renewed',
      data: {
        access_token: 'mock-refreshed-token',
        access_token_expires_in_sec: 3600,
        issued_at: Date.now()
      },
      timestamp: new Date().toISOString()
    })
  )
];
