/**
 * /orders — Order history (mobile layout, Thai labels).
 * Server Component shell. Requires auth via 6-step withAuth flow on the API proxy layer.
 * Status filter chips: ทั้งหมด / กำลังดำเนินการ / เสร็จสิ้น
 */

import Link from 'next/link';
import { cookies } from 'next/headers';
import { probeSession } from '@/lib/auth';
import { redirect } from 'next/navigation';

export const metadata = { title: 'คำสั่งซื้อของฉัน' };

/** Thai status label map */
const STATUS_LABELS: Record<string, string> = {
  PENDING_PAYMENT: 'รอชำระเงิน',
  PAID:            'ชำระแล้ว',
  PAYMENT_FAILED:  'การชำระล้มเหลว',
  PAYMENT_EXPIRED: 'หมดอายุ',
  PACKING:         'กำลังแพ็คสินค้า',
  SHIPPED:         'จัดส่งแล้ว',
  DELIVERED:       'ได้รับสินค้าแล้ว',
  CANCELLED:       'ยกเลิก',
};

/** Badge color map */
const STATUS_COLOR: Record<string, string> = {
  PENDING_PAYMENT: 'bg-amber-100 text-amber-700',
  PAID:            'bg-blue-100 text-blue-700',
  PAYMENT_FAILED:  'bg-red-100 text-red-700',
  PAYMENT_EXPIRED: 'bg-red-100 text-red-600',
  PACKING:         'bg-purple-100 text-purple-700',
  SHIPPED:         'bg-indigo-100 text-indigo-700',
  DELIVERED:       'bg-green-100 text-green-700',
  CANCELLED:       'bg-slate-100 text-slate-500',
};

interface OrderSummary {
  id: string;
  orderNumber: string;
  status: string;
  grandTotal: number;
  createdAt: string;
}

export default async function OrdersPage() {
  const cookieStore = cookies();
  // Reconstruct cookie header string from individual cookies
  const cookieHeader = cookieStore.getAll().map((c) => `${c.name}=${c.value}`).join('; ');
  const session = probeSession(cookieHeader);

  if (!session) {
    redirect('/login?next=/orders');
  }

  // Fetch order list via proxy (auth handled by route handler)
  let orders: OrderSummary[] = [];
  try {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL ?? 'http://localhost:3000';
    const res = await fetch(`${baseUrl}/api/proxy/order/list`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        cookie: cookieHeader,
      },
      body: JSON.stringify({ page: 1, limit: 20 }),
      cache: 'no-store',
    });
    if (res.ok) {
      const envelope = await res.json() as { code: string; data: { items: OrderSummary[] } | null };
      orders = envelope.data?.items ?? [];
    }
  } catch {
    // Silently degrade — show empty state
  }

  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      <h1 className="text-[17px] font-semibold text-[var(--text-primary)] mb-4">คำสั่งซื้อของฉัน</h1>

      {orders.length === 0 ? (
        <div className="py-16 text-center">
          <p className="text-4xl mb-4">📦</p>
          <p className="text-[14px] font-medium text-[var(--text-primary)] mb-1">ยังไม่มีคำสั่งซื้อ</p>
          <p className="text-[12px] text-[var(--text-tertiary)] mb-6">เริ่มช้อปปิ้งและคำสั่งซื้อจะแสดงที่นี่</p>
          <Link
            href="/products"
            className="inline-block bg-[var(--brand)] text-white font-semibold px-6 py-2.5 rounded-xl text-[13px] hover:opacity-90 transition"
          >
            เลือกซื้อสินค้า
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {orders.map((order) => {
            const label = STATUS_LABELS[order.status] ?? order.status;
            const color = STATUS_COLOR[order.status] ?? 'bg-slate-100 text-slate-600';
            const date = new Date(order.createdAt).toLocaleDateString('th-TH', {
              year: 'numeric', month: 'short', day: 'numeric',
            });
            const price = new Intl.NumberFormat('th-TH', {
              style: 'currency', currency: 'THB', maximumFractionDigits: 0,
            }).format(order.grandTotal);

            return (
              <Link
                key={order.id}
                href={`/orders/${order.id}`}
                className="block bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 hover:border-brand-300 transition"
              >
                <div className="flex items-start justify-between gap-2 mb-2">
                  <span className="text-[13px] font-semibold text-[var(--text-primary)]">{order.orderNumber}</span>
                  <span className={`text-[11px] font-medium px-2 py-0.5 rounded-full ${color}`}>{label}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-[11px] text-[var(--text-tertiary)]">{date}</span>
                  <span className="text-[13px] font-semibold text-[var(--text-primary)]">{price}</span>
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}
