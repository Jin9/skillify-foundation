/**
 * /checkout — Checkout summary page (mobile layout).
 * Server Component shell. Shipping fee logic per spec:
 *   subtotal < 1500 THB → 60 THB fee; ≥ 1500 THB → free
 * CheckoutCommitButton (Client) generates crypto.randomUUID() idempotency key per click.
 */

import Link from 'next/link';
import CheckoutCommitButton from './CheckoutCommitButton';

export const metadata = { title: 'ชำระเงิน' };

export default function CheckoutPage() {
  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      {/* Back nav */}
      <Link href="/cart" className="inline-flex items-center gap-1 text-[var(--text-secondary)] text-[13px] mb-4 hover:text-[var(--brand)]">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
          <polyline points="15 18 9 12 15 6" />
        </svg>
        ย้อนกลับ
      </Link>

      <h1 className="text-[17px] font-semibold text-[var(--text-primary)] mb-4">ชำระเงิน</h1>

      {/* Address placeholder */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 mb-3">
        <p className="text-[12px] text-[var(--text-tertiary)] uppercase tracking-wide mb-2 font-medium">ที่อยู่จัดส่ง</p>
        <p className="text-[13px] text-[var(--text-secondary)]">เพิ่มที่อยู่จัดส่งในขั้นตอนถัดไป</p>
      </div>

      {/* Order summary stub */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 mb-3 space-y-2">
        <p className="text-[12px] text-[var(--text-tertiary)] uppercase tracking-wide font-medium mb-2">สรุปออเดอร์</p>
        <div className="flex justify-between text-[13px]">
          <span className="text-[var(--text-secondary)]">ยอดรวมสินค้า</span>
          <span className="text-[var(--text-primary)] font-medium">—</span>
        </div>
        <div className="flex justify-between text-[13px]">
          <span className="text-[var(--text-secondary)]">ค่าจัดส่ง</span>
          <span className="text-[var(--text-primary)] font-medium">—</span>
        </div>
        <div className="border-t border-[var(--divider)] pt-2 flex justify-between">
          <span className="text-[14px] font-semibold text-[var(--text-primary)]">ยอดชำระ</span>
          <span className="text-[14px] font-semibold text-[var(--brand)]">—</span>
        </div>
      </div>

      <CheckoutCommitButton />
    </div>
  );
}
