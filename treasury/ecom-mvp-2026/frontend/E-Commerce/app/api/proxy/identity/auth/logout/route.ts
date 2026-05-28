/**
 * POST /api/proxy/identity/auth/logout
 *
 * Reads the refresh cookie, calls identity.logout.
 * Always clears both cookies (Max-Age=0) regardless of downstream outcome
 * to keep client state clean per spec.
 */

import { type NextRequest } from 'next/server';
import { readCookie, clearAccessCookie, clearRefreshCookie } from '@/lib/auth';

const IDENTITY_URL = process.env.BACKEND_IDENTITY_URL ?? '';

export const dynamic = 'force-dynamic';

export async function POST(req: NextRequest): Promise<Response> {
  const cookieHeader = req.headers.get('cookie');
  const refresh = readCookie(cookieHeader, 'refresh');

  // Always clear cookies, even if no refresh token (belt-and-suspenders)
  const clearCookies = (res: Response) => {
    res.headers.append('Set-Cookie', clearAccessCookie());
    res.headers.append('Set-Cookie', clearRefreshCookie());
    return res;
  };

  if (refresh) {
    try {
      await fetch(`${IDENTITY_URL}/api/v1/identity/user/logout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken: refresh }),
      });
    } catch {
      // Ignore upstream failures — we clear cookies regardless
    }
  }

  const res = Response.json(
    { code: 'SUCCESS', message: 'Logged out', data: null, traceId: '' },
    { status: 200 },
  );
  return clearCookies(res);
}
