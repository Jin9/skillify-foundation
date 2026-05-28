/**
 * /profile — Customer profile + settings (mobile layout, Thai).
 * Server Component — session probe drives content.
 * Shows login CTA when unauthenticated.
 */

import Link from 'next/link';
import { cookies } from 'next/headers';
import { probeSession } from '@/lib/auth';

export const metadata = { title: 'โปรไฟล์' };

export default function ProfilePage() {
  const cookieStore = cookies();
  const cookieHeader = cookieStore.getAll().map((c) => `${c.name}=${c.value}`).join('; ');
  const session = probeSession(cookieHeader);

  if (!session) {
    return (
      <div className="px-4 py-16 max-w-[390px] mx-auto text-center">
        <p className="text-5xl mb-4">👤</p>
        <h2 className="text-[17px] font-semibold text-[var(--text-primary)] mb-2">เข้าสู่ระบบเพื่อดูโปรไฟล์</h2>
        <p className="text-[13px] text-[var(--text-secondary)] mb-6">จัดการที่อยู่ ดูประวัติออเดอร์ และตั้งค่าบัญชีของคุณ</p>
        <Link
          href="/login?next=/profile"
          className="inline-block bg-[var(--brand)] text-white font-semibold px-6 py-3 rounded-xl text-[14px] hover:opacity-90 transition"
        >
          เข้าสู่ระบบ
        </Link>
      </div>
    );
  }

  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      <h1 className="text-[17px] font-semibold text-[var(--text-primary)] mb-6">โปรไฟล์</h1>

      {/* Account info */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 mb-3">
        <p className="text-[12px] text-[var(--text-tertiary)] uppercase tracking-wide font-medium mb-3">บัญชีของฉัน</p>
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-full bg-[var(--brand)] bg-opacity-10 flex items-center justify-center flex-shrink-0">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="var(--brand)" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
              <circle cx="12" cy="7" r="4" />
            </svg>
          </div>
          <div>
            <p className="text-[14px] font-medium text-[var(--text-primary)]">สมาชิก ShopPilot</p>
            <p className="text-[11px] text-[var(--text-tertiary)]">ID: {session.sub.slice(0, 8)}…</p>
          </div>
        </div>
      </div>

      {/* Menu */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl overflow-hidden mb-3">
        {[
          { href: '/orders',   label: 'คำสั่งซื้อของฉัน' },
        ].map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className="flex items-center justify-between px-4 py-3.5 border-b border-[var(--divider)] last:border-0 hover:bg-[var(--bg)] transition"
          >
            <span className="text-[13px] text-[var(--text-primary)]">{item.label}</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--text-tertiary)" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="9 18 15 12 9 6" />
            </svg>
          </Link>
        ))}
      </div>

      {/* Logout */}
      <form action="/api/proxy/identity/auth/logout" method="post">
        <button
          type="submit"
          className="w-full text-[13px] text-red-500 hover:text-red-700 py-3 rounded-xl border border-red-100 hover:bg-red-50 transition"
        >
          ออกจากระบบ
        </button>
      </form>
    </div>
  );
}
