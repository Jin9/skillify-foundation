/**
 * /orders/[id] — Order Detail (Server Component, Thai mobile layout).
 * Fetches order detail via proxy. Item snapshots survive soft-delete per CAT-009/ORD-007.
 * Renders status timeline + address snapshot + item list.
 */

import type { Metadata } from 'next';
import Link from 'next/link';
import { cookies } from 'next/headers';
import { probeSession } from '@/lib/auth';
import { redirect } from 'next/navigation';

export const metadata: Metadata = { title: 'รายละเอียดคำสั่งซื้อ' };

// Status timeline steps per spec order lifecycle
const TIMELINE_STEPS = [
  { key: 'PENDING_PAYMENT', label: 'รอชำระ' },
  { key: 'PAID',            label: 'ชำระแล้ว' },
  { key: 'PACKING',         label: 'แพ็คสินค้า' },
  { key: 'SHIPPED',         label: 'จัดส่ง' },
  { key: 'DELIVERED',       label: 'ได้รับแล้ว' },
] as const;

const STATUS_THAI: Record<string, string> = {
  PENDING_PAYMENT: 'รอชำระเงิน',
  PAID:            'ชำระแล้ว',
  PAYMENT_FAILED:  'การชำระล้มเหลว',
  PAYMENT_EXPIRED: 'หมดอายุ',
  PACKING:         'กำลังแพ็คสินค้า',
  SHIPPED:         'จัดส่งแล้ว',
  DELIVERED:       'ได้รับสินค้าแล้ว',
  CANCELLED:       'ยกเลิกแล้ว',
};

interface OrderItem {
  productId: string;
  productNameSnapshot: string;
  priceSnapshot: number;
  productImageUrlSnapshot?: string;
  quantity: number;
  subtotal: number;
}

interface OrderDetail {
  id: string;
  orderNumber: string;
  status: string;
  createdAt: string;
  subtotal: number;
  discount: number;
  shippingFee: number;
  grandTotal: number;
  items: OrderItem[];
  address?: {
    recipientName: string;
    phone: string;
    addressLine: string;
    district: string;
    province: string;
    postalCode: string;
  };
  history?: { status: string; at: string }[];
}

function formatTHB(amount: number) {
  return new Intl.NumberFormat('th-TH', {
    style: 'currency', currency: 'THB', maximumFractionDigits: 0,
  }).format(amount);
}

export default async function OrderDetailPage({ params }: { params: { id: string } }) {
  const cookieStore = cookies();
  // Reconstruct cookie header string from individual cookies
  const cookieHeader = cookieStore.getAll().map((c) => `${c.name}=${c.value}`).join('; ');
  const session = probeSession(cookieHeader);

  if (!session) {
    redirect(`/login?next=/orders/${params.id}`);
  }

  let order: OrderDetail | null = null;
  try {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL ?? 'http://localhost:3000';
    const res = await fetch(`${baseUrl}/api/proxy/order/detail`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        cookie: cookieHeader,
      },
      body: JSON.stringify({ orderId: params.id }),
      cache: 'no-store',
    });
    if (res.ok) {
      const envelope = await res.json() as { code: string; data: OrderDetail | null };
      order = envelope.data;
    }
  } catch {
    // Silently degrade
  }

  if (!order) {
    return (
      <div className="px-4 py-12 max-w-[390px] mx-auto text-center">
        <p className="text-[14px] text-[var(--text-secondary)] mb-4">ไม่พบคำสั่งซื้อ</p>
        <Link href="/orders" className="text-[var(--brand)] text-[13px] hover:underline">
          กลับไปคำสั่งซื้อทั้งหมด
        </Link>
      </div>
    );
  }

  // Determine active timeline step index
  const timelineKeys = TIMELINE_STEPS.map((s) => s.key);
  const activeIdx = timelineKeys.indexOf(order.status as typeof TIMELINE_STEPS[number]['key']);

  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      {/* Back */}
      <Link href="/orders" className="inline-flex items-center gap-1 text-[var(--text-secondary)] text-[13px] mb-4 hover:text-[var(--brand)]">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
          <polyline points="15 18 9 12 15 6" />
        </svg>
        คำสั่งซื้อของฉัน
      </Link>

      <div className="flex items-start justify-between mb-4">
        <div>
          <h1 className="text-[17px] font-semibold text-[var(--text-primary)]">{order.orderNumber}</h1>
          <p className="text-[11px] text-[var(--text-tertiary)] mt-0.5">
            {new Date(order.createdAt).toLocaleDateString('th-TH', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>
        </div>
        <span className="text-[12px] font-medium px-3 py-1 rounded-full bg-[var(--brand)] bg-opacity-10 text-[var(--brand)]">
          {STATUS_THAI[order.status] ?? order.status}
        </span>
      </div>

      {/* Status Timeline (only for non-failed/cancelled orders) */}
      {activeIdx >= 0 && (
        <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 mb-3">
          <p className="text-[12px] text-[var(--text-tertiary)] uppercase tracking-wide font-medium mb-3">สถานะออเดอร์</p>
          <div className="flex items-center justify-between relative">
            {/* Connecting line */}
            <div className="absolute top-3 left-3 right-3 h-0.5 bg-[var(--divider)]" aria-hidden="true" />
            <div
              className="absolute top-3 left-3 h-0.5 bg-[var(--brand)] transition-all"
              style={{ width: `${activeIdx > 0 ? (activeIdx / (TIMELINE_STEPS.length - 1)) * 100 : 0}%` }}
              aria-hidden="true"
            />
            {TIMELINE_STEPS.map((step, idx) => {
              const done = idx <= activeIdx;
              return (
                <div key={step.key} className="flex flex-col items-center z-10 flex-1">
                  <div className={`w-6 h-6 rounded-full border-2 flex items-center justify-center ${
                    done ? 'bg-[var(--brand)] border-[var(--brand)]' : 'bg-white border-[var(--divider)]'
                  }`}>
                    {done && (
                      <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="20 6 9 17 4 12" />
                      </svg>
                    )}
                  </div>
                  <span className={`mt-1.5 text-[9px] text-center leading-tight ${
                    done ? 'text-[var(--brand)] font-medium' : 'text-[var(--text-tertiary)]'
                  }`}>{step.label}</span>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Address */}
      {order.address && (
        <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 mb-3">
          <p className="text-[12px] text-[var(--text-tertiary)] uppercase tracking-wide font-medium mb-2">ที่อยู่จัดส่ง</p>
          <p className="text-[13px] font-medium text-[var(--text-primary)]">{order.address.recipientName}</p>
          <p className="text-[12px] text-[var(--text-secondary)]">{order.address.phone}</p>
          <p className="text-[12px] text-[var(--text-secondary)]">
            {order.address.addressLine} {order.address.district} {order.address.province} {order.address.postalCode}
          </p>
        </div>
      )}

      {/* Items */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 mb-3">
        <p className="text-[12px] text-[var(--text-tertiary)] uppercase tracking-wide font-medium mb-3">รายการสินค้า</p>
        <div className="space-y-3">
          {order.items.map((item) => (
            <div key={item.productId} className="flex gap-3">
              <div className="w-14 h-14 flex-shrink-0 bg-[var(--bg)] rounded-xl overflow-hidden flex items-center justify-center">
                {item.productImageUrlSnapshot ? (
                  <img src={item.productImageUrlSnapshot} alt={item.productNameSnapshot} className="w-full h-full object-cover" />
                ) : (
                  <span className="text-xl">📦</span>
                )}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-[13px] font-medium text-[var(--text-primary)] line-clamp-1">{item.productNameSnapshot}</p>
                <p className="text-[11px] text-[var(--text-tertiary)]">x{item.quantity}</p>
                <p className="text-[13px] font-semibold text-[var(--text-primary)]">{formatTHB(item.subtotal)}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Price summary */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 mb-4 space-y-2">
        <div className="flex justify-between text-[13px]">
          <span className="text-[var(--text-secondary)]">ยอดรวมสินค้า</span>
          <span className="text-[var(--text-primary)]">{formatTHB(order.subtotal)}</span>
        </div>
        {order.discount > 0 && (
          <div className="flex justify-between text-[13px]">
            <span className="text-[var(--text-secondary)]">ส่วนลด</span>
            <span className="text-green-600">-{formatTHB(order.discount)}</span>
          </div>
        )}
        <div className="flex justify-between text-[13px]">
          <span className="text-[var(--text-secondary)]">ค่าจัดส่ง</span>
          <span className="text-[var(--text-primary)]">
            {order.shippingFee === 0 ? 'ฟรี' : formatTHB(order.shippingFee)}
          </span>
        </div>
        <div className="border-t border-[var(--divider)] pt-2 flex justify-between">
          <span className="text-[14px] font-semibold text-[var(--text-primary)]">ยอดชำระ</span>
          <span className="text-[14px] font-semibold text-[var(--brand)]">{formatTHB(order.grandTotal)}</span>
        </div>
      </div>
    </div>
  );
}
