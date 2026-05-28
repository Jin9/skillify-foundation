/**
 * POST /api/proxy/catalog/product/detail
 *
 * Public passthrough to catalog.product.detail.
 */

import { type NextRequest } from 'next/server';
import { callBackend } from '@/lib/http';

const CATALOG_URL = process.env.BACKEND_CATALOG_URL ?? '';

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
    upstream = await callBackend({
      url: `${CATALOG_URL}/api/v1/catalog/product/detail`,
      body,
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
