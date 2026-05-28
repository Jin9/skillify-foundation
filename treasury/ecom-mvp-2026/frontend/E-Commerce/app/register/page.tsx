/**
 * /register — STUB (Coming in v2)
 *
 * Planned: Client Component form with react-hook-form + zod.
 * Fields: email (RFC5322), password (min 8, max 72), name (max 120).
 * POSTs to POST /api/auth/register → identity.register.
 * Per BA spec: registration does NOT auto-login; user is redirected to /login after success.
 * Note: On success, cookies are NOT set — the caller hits /login next (per route_handlers spec).
 */

import Link from 'next/link';

export const metadata = { title: 'Register' };

export default function RegisterPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50 px-4">
      <div className="w-full max-w-md bg-white rounded-2xl shadow-sm border border-slate-200 p-8 text-center">
        <h1 className="text-2xl font-bold text-slate-800 mb-4">Create Account</h1>
        <p className="text-slate-500 mb-6">
          Customer registration is <strong>coming in v2</strong>.
        </p>
        <p className="text-sm text-slate-400 mb-4">
          Planned: email + password + name form, with identity.register integration.
          Registration does not auto-login; the user is directed to /login after success.
        </p>
        <Link href="/login" className="text-brand-600 hover:underline text-sm">
          Sign in instead
        </Link>
      </div>
    </div>
  );
}
