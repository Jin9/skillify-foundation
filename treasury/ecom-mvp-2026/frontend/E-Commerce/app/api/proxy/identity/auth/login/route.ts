/**
 * POST /api/proxy/identity/auth/login
 *
 * Forwards the login request to identity.login.
 * On success: sets HttpOnly access + refresh cookies; strips raw tokens from the
 * client response (returns user object only per BA).
 * Public endpoint — no auth required.
 */

import { type NextRequest } from 'next/server';
import { buildAccessCookie, buildRefreshCookie } from '@/lib/auth';

const IDENTITY_URL = process.env.BACKEND_IDENTITY_URL ?? '';

export const dynamic = 'force-dynamic';

export async function POST(req: NextRequest): Promise<Response> {
  let body: unknown;
  try {
    body = await req.json();
  } catch {
    return Response.json(
      { code: 'VALIDATION_ERROR', message: 'Invalid JSON body', data: null, traceId: '' },
      { status: 400 },
    );
  }

  let upstream: Response;
  try {
    upstream = await fetch(`${IDENTITY_URL}/api/v1/identity/user/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
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
    data: { accessToken?: string; refreshToken?: string; user?: unknown } | null;
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

  if (upstream.status !== 200 || envelope.code !== 'SUCCESS' || !envelope.data) {
    return Response.json(envelope, { status: upstream.status });
  }

  const { accessToken, refreshToken, ...rest } = envelope.data as {
    accessToken: string;
    refreshToken: string;
    [key: string]: unknown;
  };

  // Build response — strip raw tokens, return user object only
  const clientBody = { code: 'SUCCESS', message: 'Login successful', data: rest.user ?? rest, traceId: envelope.traceId };

  const res = Response.json(clientBody, { status: 200 });
  res.headers.append('Set-Cookie', buildAccessCookie(accessToken));
  res.headers.append('Set-Cookie', buildRefreshCookie(refreshToken));
  return res;
}
