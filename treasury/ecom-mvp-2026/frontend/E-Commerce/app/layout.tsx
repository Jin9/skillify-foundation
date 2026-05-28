/**
 * Root layout — Server Component.
 * lang="th", mobile-first viewport 390, IBM Plex Sans Thai.
 * TabBar is rendered here so all routes get the bottom nav.
 * Header is kept for non-tab deep routes (detail, checkout, etc.).
 */

import type { Metadata } from 'next';
import { IBM_Plex_Sans_Thai } from 'next/font/google';
import './globals.css';
import Header from '@/components/Header';
import TabBar from '@/components/TabBar';

const ibmPlexSansThai = IBM_Plex_Sans_Thai({
  subsets: ['thai', 'latin'],
  weight: ['400', '500', '600', '700'],
  variable: '--font-ibm',
  display: 'swap',
});

export const metadata: Metadata = {
  title: {
    default: 'ShopPilot',
    template: '%s | ShopPilot',
  },
  description: 'B2C E-Commerce Platform — ช้อป ตะกร้า ชำระเงิน ติดตามออเดอร์',
};

export const viewport = {
  width: 390,
  initialScale: 1,
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="th" className={ibmPlexSansThai.variable}>
      <body className="min-h-screen flex flex-col bg-[var(--bg)]">
        <Header />
        <main className="flex-1 pb-16">{children}</main>
        <TabBar />
      </body>
    </html>
  );
}
