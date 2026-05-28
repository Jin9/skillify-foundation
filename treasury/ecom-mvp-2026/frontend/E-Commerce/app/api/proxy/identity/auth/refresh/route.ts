/**
 * POST /api/proxy/identity/auth/refresh
 *
 * Reads the refresh cookie, calls identity.refresh, rotates both cookies on success.
 * On AUTH_REVOKED/AUTH_INVALID: clears both cookies and returns 401.
 *
 * This endpoint is the explicit browser-facing refresh path.
 * The 6-step orchestrator (lib/auth.ts withAuth) handles the implicit refresh path
 * inside other route handlers without ever hitting this URL.
 */

import { type NextRequest } from 'next/server';
import {
  readCookie,
  buildAccessCookie,
  buildRefreshCookie,
  clearAccessCookie,
  clearRefreshCookie,
} from '@/lib/auth';

const IDENTITY_URL = process.env.BACKEND_IDENTITY_URL ?? '';

export const dynamic = 'force-dynamic';

export async function POST(req: NextRequest): Promise<Response> {
  const cookieHeader = req.headers.get('cookie');
  const refresh = readCookie(cookieHeader, 'refresh');

  if (!refresh) {
    return Response.json(
      { code: 'AUTH_MISSING', message: 'No refresh token present', data: null, traceId: '' },
      { status: 401 },
    );
  }

  let upstream: Response;
  try {
    upstream = await fetch(`${IDENTITY_URL}/api/v1/identity/user/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refreshToken: refresh }),
    });
  } catch {
    return Response.json(
      { code: 'UPSTREAM_TIMEOUT', message: 'Identity service unreachable', data: null, traceId: '' },
      { status: 502 },
    );
  }

  let envelope: {
    code: string;
    message: string;
    data: { accessToken: string; refreshToken: string } | null;
    traceId: string;
  };
  try {
    envelope = (await upstream.json()) as typeof envelope;
  } catch {
    return Response.json(
      { code: 'UPSTREAM_TIMEOUT', message: 'Invalid upstream response', data: null, traceId: '' },
      { status: 502 },
    );
  }

  if (
    envelope.code === 'AUTH_REVOKED' ||
    envelope.code === 'AUTH_INVALID' ||
    upstream.status === 401
  ) {
    const r = Response.json(
      { code: 'AUTH_REVOKED', message: 'Session expired. Please log in again.', data: null, traceId: envelope.traceId ?? '' },
      { status: 401 },
    );
    r.headers.append('Set-Cookie', clearAccessCookie());
    r.headers.append('Set-Cookie', clearRefreshCookie());
    return r;
  }

  if (envelope.code !== 'SUCCESS' || !envelope.data) {
    return Response.json(envelope, { status: upstream.status });
  }

  const { accessToken, refreshToken } = envelope.data;
  const r = Response.json(
    { code: 'SUCCESS', message: 'Tokens rotated', data: null, traceId: envelope.traceId },
    { status: 200 },
  );
  r.headers.append('Set-Cookie', buildAccessCookie(accessToken));
  r.headers.append('Set-Cookie', buildRefreshCookie(refreshToken));
  return r;
}
