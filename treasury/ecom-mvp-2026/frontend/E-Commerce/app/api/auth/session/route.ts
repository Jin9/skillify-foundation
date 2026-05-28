/**
 * GET /api/auth/session
 *
 * Pure session probe — decodes the access cookie WITHOUT triggering a refresh.
 * Returns {authenticated: true, user: {userId, email, role}} on success,
 * or {authenticated: false} when the cookie is absent/expired.
 *
 * Used by (customer)/(admin) layout.tsx to gate protected pages.
 */

import { type NextRequest } from 'next/server';
import { probeSession } from '@/lib/auth';

export const dynamic = 'force-dynamic';

export async function GET(req: NextRequest): Promise<Response> {
  const cookieHeader = req.headers.get('cookie');
  const session = probeSession(cookieHeader);

  if (!session) {
    return Response.json({ authenticated: false });
  }

  return Response.json({
    authenticated: true,
    user: {
      userId: session.sub,
      role: session.role,
    },
  });
}
