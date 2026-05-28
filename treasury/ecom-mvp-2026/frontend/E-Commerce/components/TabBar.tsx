'use client';

/**
 * TabBar — fixed bottom navigation for the 5 main tab routes.
 * Visible only on: /, /products, /cart, /orders, /profile.
 * Uses outline SVG icons + Thai labels.
 */

import Link from 'next/link';
import { usePathname } from 'next/navigation';

const TAB_ROUTES = ['/', '/products', '/cart', '/orders', '/profile'] as const;

interface Tab {
  href: string;
  label: string;
  icon: (active: boolean) => React.ReactNode;
}

function HomeIcon({ active }: { active: boolean }) {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none"
      stroke={active ? 'var(--brand)' : 'currentColor'} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 9.5L12 3l9 6.5V20a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V9.5z" />
      <path d="M9 21V12h6v9" />
    </svg>
  );
}

function SearchIcon({ active }: { active: boolean }) {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none"
      stroke={active ? 'var(--brand)' : 'currentColor'} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="11" cy="11" r="8" />
      <line x1="21" y1="21" x2="16.65" y2="16.65" />
    </svg>
  );
}

function CartIcon({ active }: { active: boolean }) {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none"
      stroke={active ? 'var(--brand)' : 'currentColor'} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
      <line x1="3" y1="6" x2="21" y2="6" />
      <path d="M16 10a4 4 0 0 1-8 0" />
    </svg>
  );
}

function OrdersIcon({ active }: { active: boolean }) {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none"
      stroke={active ? 'var(--brand)' : 'currentColor'} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="3" width="18" height="18" rx="2" />
      <path d="M3 9h18" />
      <path d="M9 21V9" />
    </svg>
  );
}

function ProfileIcon({ active }: { active: boolean }) {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none"
      stroke={active ? 'var(--brand)' : 'currentColor'} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

const TABS: Tab[] = [
  { href: '/',        label: 'หน้าแรก',  icon: (a) => <HomeIcon active={a} /> },
  { href: '/products', label: 'ค้นหา',   icon: (a) => <SearchIcon active={a} /> },
  { href: '/cart',    label: 'ตะกร้า',   icon: (a) => <CartIcon active={a} /> },
  { href: '/orders',  label: 'คำสั่งซื้อ', icon: (a) => <OrdersIcon active={a} /> },
  { href: '/profile', label: 'โปรไฟล์',  icon: (a) => <ProfileIcon active={a} /> },
];

export default function TabBar() {
  const pathname = usePathname();

  // Only show on the 5 tab routes
  const isTabRoute = (TAB_ROUTES as readonly string[]).includes(pathname);
  if (!isTabRoute) return null;

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 bg-[var(--card)] border-t border-[var(--divider)]">
      <div className="flex items-stretch justify-around max-w-[390px] mx-auto">
        {TABS.map((tab) => {
          const active = pathname === tab.href;
          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={`flex flex-col items-center justify-center gap-0.5 flex-1 py-2 text-[10px] font-medium transition-colors ${
                active ? 'text-[var(--brand)]' : 'text-[var(--text-secondary)]'
              }`}
            >
              {tab.icon(active)}
              <span>{tab.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
