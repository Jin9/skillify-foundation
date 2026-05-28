/**
 * POST /api/proxy/[...path]
 *
 * Catch-all STUB proxy — forwards all other /api/proxy/* routes to the appropriate
 * backend with the 6-step auth flow.
 *
 * Path mapping: /api/proxy/{domain}/{aggregate}/{action}
 *   → {BACKEND_{DOMAIN}_URL}/api/v1/{domain}/{aggregate}/{action}
 *
 * STUB NOTE (v2 TODO): This handler provides auth forwarding but does NOT implement:
 *   - Per-route role guards (ADMIN check is not performed here)
 *   - Idempotency-Key forwarding (for checkout.commit, use the dedicated handler)
 *   - Request body validation
 * These are flagged for dedicated handlers in v2.
 */

import { type NextRequest } from 'next/server';
import { withAuth } from '@/lib/auth';
import { serviceUrl } from '@/lib/http';

export const dynamic = 'force-dynamic';

export async function POST(
  req: NextRequest,
  { params }: { params: { path: string[] } },
): Promise<Response> {
  const segments = params.path; // e.g. ['cart', 'cart', 'add-item']

  const domain = segments[0];
  if (!domain) {
    return Response.json(
      { code: 'NOT_FOUND', message: 'Unknown proxy route', data: null, traceId: '' },
      { status: 404 },
    );
  }

  let baseUrl: string;
  try {
    baseUrl = serviceUrl(domain);
  } catch {
    return Response.json(
      { code: 'NOT_FOUND', message: `Unknown service domain: ${domain}`, data: null, traceId: '' },
      { status: 404 },
    );
  }

  const downstreamPath = `/api/v1/${segments.join('/')}`;
  const downstreamUrl = `${baseUrl}${downstreamPath}`;

  let body: unknown;
  try {
    body = await req.json();
  } catch {
    body = {};
  }

  return withAuth(req, async (bearer: string) => {
    let upstream: Response;
    try {
      upstream = await fetch(downstreamUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${bearer}`,
        },
        body: JSON.stringify(body),
      });
    } catch {
      return Response.json(
        { code: 'UPSTREAM_TIMEOUT', message: `${domain} service unreachable`, data: null, traceId: '' },
        { status: 502 },
      );
    }

    const data = await upstream.json();
    return Response.json(data, { status: upstream.status });
  });
}
