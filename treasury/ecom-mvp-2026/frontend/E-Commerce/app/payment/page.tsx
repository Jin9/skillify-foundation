/**
 * /payment — Mock payment screen.
 * 3 outcome buttons: สำเร็จ / ไม่สำเร็จ / หมดเวลา
 * Client Component. After a 1.4s "processing" delay, redirects to /orders/[id].
 */

export const metadata = { title: 'ชำระเงิน' };

import PaymentSimulator from './PaymentSimulator';

export default function PaymentPage() {
  return (
    <div className="px-4 py-6 max-w-[390px] mx-auto">
      <h1 className="text-[17px] font-semibold text-[var(--text-primary)] mb-1">ชำระเงิน</h1>
      <p className="text-[12px] text-[var(--text-tertiary)] mb-6">เลือกผลการชำระเงิน (โหมดทดสอบ)</p>
      <PaymentSimulator />
    </div>
  );
}
