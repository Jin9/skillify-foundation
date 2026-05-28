'use client';

/**
 * Header — Client Component
 *
 * Displays navigation with login state probed from /api/auth/session.
 * Shows cart link and admin link conditionally based on role.
 */

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface Session {
  authenticated: boolean;
  user?: { userId: string; role: string };
}

export default function Header() {
  const pathname = usePathname();
  const [session, setSession] = useState<Session | null>(null);

  useEffect(() => {
    fetch('/api/auth/session')
      .then((r) => r.json())
      .then((s: Session) => setSession(s))
      .catch(() => setSession({ authenticated: false }));
  }, [pathname]);

  const handleLogout = async () => {
    await fetch('/api/proxy/identity/auth/logout', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({}),
    });
    setSession({ authenticated: false });
    window.location.href = '/';
  };

  return (
    <header className="bg-[var(--card)] border-b border-[var(--divider)] sticky top-0 z-40">
      <div className="max-w-[390px] mx-auto px-4 h-13 flex items-center justify-between">
        {/* Logo */}
        <Link href="/" className="font-bold text-[17px] text-[var(--brand)]">
          ShopPilot
        </Link>

        {/* Nav */}
        <nav className="flex items-center gap-4 text-[13px]">
          {session === null && (
            <span className="text-[var(--text-tertiary)] text-[11px]">กำลังโหลด…</span>
          )}

          {session?.authenticated ? (
            <button
              onClick={() => void handleLogout()}
              className="text-[var(--text-secondary)] hover:text-red-600 transition"
            >
              ออกจากระบบ
            </button>
          ) : session !== null ? (
            <Link
              href={`/login?next=${encodeURIComponent(pathname)}`}
              className="bg-[var(--brand)] hover:opacity-90 text-white px-4 py-1.5 rounded-xl font-medium transition text-[13px]"
            >
              เข้าสู่ระบบ
            </Link>
          ) : null}
        </nav>
      </div>
    </header>
  );
}
