/**
 * POST /api/proxy/catalog/product/list
 *
 * Public passthrough to catalog.product.list.
 * If an access cookie is present and valid, forwards Bearer (enables admin variant
 * to receive DRAFT/INACTIVE products). Unauthenticated requests omit the header.
 */

import { type NextRequest } from 'next/server';
import { readCookie } from '@/lib/auth';
import { callBackend } from '@/lib/http';

const CATALOG_URL = process.env.BACKEND_CATALOG_URL ?? '';

export const dynamic = 'force-dynamic';

export async function POST(req: NextRequest): Promise<Response> {
  let body: unknown;
  try {
    body = await req.json();
  } catch {
    body = {};
  }

  const cookieHeader = req.headers.get('cookie');
  const access = readCookie(cookieHeader, 'access');

  let upstream: Response;
  try {
    upstream = await callBackend({
      url: `${CATALOG_URL}/api/v1/catalog/product/list`,
      body,
      bearer: access,
    });
  } catch {
    return Response.json(
      { code: 'UPSTREAM_TIMEOUT', message: 'Catalog service unreachable', data: null, traceId: '' },
      { status: 502 },
    );
  }

  const data = await upstream.json();
  return Response.json(data, { status: upstream.status });
}
