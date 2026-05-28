/**
 * lib/auth.ts
 *
 * 6-step token acquisition flow per cross-cutting.auth.token_acquisition_flow.
 * STATELESS per-request — no shared in-memory cache across requests.
 * No external JWT library — payload is base64url-decoded locally (no-verify).
 * The backend validates signature; this code only inspects `exp`.
 */

const IDENTITY_URL = process.env.BACKEND_IDENTITY_URL ?? '';
const isProd = process.env.NODE_ENV === 'production';

// ---------------------------------------------------------------------------
// JWT helpers (no-verify — exp only)
// ---------------------------------------------------------------------------

function base64UrlDecode(input: string): string {
  // Normalise base64url → base64, pad if needed
  const base64 = input.replace(/-/g, '+').replace(/_/g, '/');
  const padded = base64 + '='.repeat((4 - (base64.length % 4)) % 4);
  return Buffer.from(padded, 'base64').toString('utf-8');
}

function decodeJwt(token: string): { exp: number; sub?: string; role?: string } {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) throw new Error('malformed jwt');
    return JSON.parse(base64UrlDecode(parts[1])) as { exp: number; sub?: string; role?: string };
  } catch {
    // Treat unparseable tokens as already-expired
    return { exp: 0 };
  }
}

/**
 * Returns true when the token is expired (or will expire within skewSec seconds).
 */
function isExpired(exp: number, skewSec: number): boolean {
  return Date.now() / 1000 + skewSec >= exp;
}

// ---------------------------------------------------------------------------
// Cookie builders
// ---------------------------------------------------------------------------

function cookieSecureFlag(): string {
  return isProd ? ' Secure;' : '';
}

export function buildAccessCookie(token: string): string {
  return `access=${token}; Path=/; HttpOnly; SameSite=Lax;${cookieSecureFlag()} Max-Age=900`;
}

export function buildRefreshCookie(token: string): string {
  return `refresh=${token}; Path=/; HttpOnly; SameSite=Lax;${cookieSecureFlag()} Max-Age=1209600`;
}

export function clearAccessCookie(): string {
  return `access=; Path=/; HttpOnly; Max-Age=0`;
}

export function clearRefreshCookie(): string {
  return `refresh=; Path=/; HttpOnly; Max-Age=0`;
}

// ---------------------------------------------------------------------------
// Cookie readers (from raw cookie header string)
// ---------------------------------------------------------------------------

export function readCookie(cookieHeader: string | null, name: string): string | undefined {
  if (!cookieHeader) return undefined;
  const match = cookieHeader.match(new RegExp(`(?:^|;\\s*)${name}=([^;]+)`));
  return match?.[1];
}

// ---------------------------------------------------------------------------
// Standard envelope helpers
// ---------------------------------------------------------------------------

function makeEnvelope(
  code: string,
  message: string,
  data: unknown,
  traceId: string,
) {
  return JSON.stringify({ code, message, data, traceId });
}

function envelope401(code: string, message: string, traceId = ''): Response {
  return new Response(makeEnvelope(code, message, null, traceId), {
    status: 401,
    headers: { 'Content-Type': 'application/json' },
  });
}

function envelope401Cleared(code: string, message: string): Response {
  const r = new Response(makeEnvelope(code, message, null, ''), {
    status: 401,
    headers: { 'Content-Type': 'application/json' },
  });
  r.headers.append('Set-Cookie', clearAccessCookie());
  r.headers.append('Set-Cookie', clearRefreshCookie());
  return r;
}

function envelope502(code: string, message: string): Response {
  return new Response(makeEnvelope(code, message, null, ''), {
    status: 502,
    headers: { 'Content-Type': 'application/json' },
  });
}

// ---------------------------------------------------------------------------
// Role guard — call AFTER withAuth resolves the access token
// ---------------------------------------------------------------------------

export function checkRole(accessToken: string, requiredRole: string): Response | null {
  const claims = decodeJwt(accessToken);
  if ((claims as Record<string, unknown>).role !== requiredRole) {
    return new Response(
      makeEnvelope('AUTH_FORBIDDEN', 'Insufficient role', null, ''),
      { status: 403, headers: { 'Content-Type': 'application/json' } },
    );
  }
  return null;
}

// ---------------------------------------------------------------------------
// 6-step token acquisition orchestrator
// ---------------------------------------------------------------------------

/**
 * withAuth — wraps a downstream call with the full 6-step token acquisition flow.
 *
 * Steps:
 * 1. Read `access` cookie from request.
 * 2. If present and not expired (30s skew), forward as Bearer immediately.
 * 3. If missing/expired, read `refresh` cookie.  If absent → 401 AUTH_MISSING.
 * 4. Call identity.refresh.  On non-2xx revoked response → clear cookies, 401 AUTH_REVOKED.
 * 5. Rotate both cookies on the downstream response.
 * 6. Return the (possibly cookie-rotated) downstream response.
 */
export async function withAuth(
  req: Request,
  downstream: (bearer: string) => Promise<Response>,
): Promise<Response> {
  const cookieHeader = req.headers.get('cookie');

  // Step 1+2: attempt fast path with valid access token
  const access = readCookie(cookieHeader, 'access');
  if (access) {
    const { exp } = decodeJwt(access);
    if (!isExpired(exp, 30)) {
      return downstream(access);
    }
  }

  // Step 3: access is absent or expired — need refresh
  const refresh = readCookie(cookieHeader, 'refresh');
  if (!refresh) {
    return envelope401('AUTH_MISSING', 'No session found. Please log in.');
  }

  // Step 4: call identity.refresh
  let refreshRes: Response;
  try {
    refreshRes = await fetch(`${IDENTITY_URL}/api/v1/identity/user/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refreshToken: refresh }),
    });
  } catch {
    return envelope502('UPSTREAM_TIMEOUT', 'Identity service unreachable');
  }

  let refreshBody: {
    code: string;
    data: { accessToken: string; refreshToken: string } | null;
  };
  try {
    refreshBody = (await refreshRes.json()) as typeof refreshBody;
  } catch {
    return envelope502('UPSTREAM_TIMEOUT', 'Invalid response from identity service');
  }

  // Step 6: revoked / invalid refresh token
  if (
    refreshBody.code === 'AUTH_REVOKED' ||
    refreshBody.code === 'AUTH_INVALID' ||
    refreshRes.status === 401
  ) {
    return envelope401Cleared('AUTH_REVOKED', 'Session expired. Please log in again.');
  }

  if (refreshBody.code !== 'SUCCESS' || !refreshBody.data) {
    return envelope502('UPSTREAM_TIMEOUT', 'Unexpected identity service response');
  }

  // Step 5: forward with new access token, then rotate cookies on response
  const { accessToken, refreshToken } = refreshBody.data;
  let downstreamRes: Response;
  try {
    downstreamRes = await downstream(accessToken);
  } catch {
    return envelope502('UPSTREAM_TIMEOUT', 'Downstream service unreachable');
  }

  // Rotate both cookies on the downstream response
  // We must clone because headers may be immutable
  const mutableHeaders = new Headers(downstreamRes.headers);
  mutableHeaders.append('Set-Cookie', buildAccessCookie(accessToken));
  mutableHeaders.append('Set-Cookie', buildRefreshCookie(refreshToken));

  return new Response(downstreamRes.body, {
    status: downstreamRes.status,
    headers: mutableHeaders,
  });
}

// ---------------------------------------------------------------------------
// Session probe — pure read, no refresh side-effect
// ---------------------------------------------------------------------------

export interface SessionClaims {
  sub: string;
  role: string;
  exp: number;
}

/**
 * probeSession — decode the access cookie without triggering a refresh.
 * Returns null if the cookie is absent or expired.
 */
export function probeSession(cookieHeader: string | null): SessionClaims | null {
  const access = readCookie(cookieHeader, 'access');
  if (!access) return null;
  const claims = decodeJwt(access);
  if (isExpired(claims.exp, 0)) return null;
  return claims as SessionClaims;
}
