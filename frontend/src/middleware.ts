import { NextResponse, type NextRequest } from 'next/server';

export async function middleware(req: NextRequest) {
  const response = NextResponse.next();
  response.headers.set('x-pathname', req.nextUrl.pathname);

  // boilerplate:axis=auth option=jwt-refresh START
  if (process.env.NEXT_PUBLIC_ENABLE_REFRESH_TOKEN === 'true') {
    const { tryRefreshIfExpired } = await import(
      '@/services/auth/jwt-refresh/fetcher-middleware'
    );
    await tryRefreshIfExpired(req);
  }
  // boilerplate:axis=auth option=jwt-refresh END

  return response;
}

export const config = {
  matcher: [
    /*
     * Run on all paths except:
     * - _next/static  (bundled assets)
     * - _next/image   (image optimisation)
     * - favicon.ico, sitemap.xml, robots.txt
     * - /api          (route handlers handle their own auth)
     */
    '/((?!_next/static|_next/image|favicon\\.ico|sitemap\\.xml|robots\\.txt|api).*)'
  ]
};
