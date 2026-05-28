'use client';

/**
 * PaymentSimulator — Client Component.
 * Simulates 3 payment outcomes after 1.4s delay.
 * Uses crypto.randomUUID() for the idempotency key per submit.
 * Redirects to /orders/[id] (or /orders on missing orderId).
 */

import { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

type Outcome = 'SUCCESS' | 'FAILED' | 'TIMEOUT';

const OUTCOMES: { outcome: Outcome; label: string; variant: string }[] = [
  { outcome: 'SUCCESS', label: 'สำเร็จ',    variant: 'bg-[var(--brand)] text-white hover:opacity-90' },
  { outcome: 'FAILED',  label: 'ไม่สำเร็จ', variant: 'bg-red-500 text-white hover:bg-red-600' },
  { outcome: 'TIMEOUT', label: 'หมดเวลา',   variant: 'bg-amber-500 text-white hover:bg-amber-600' },
];

export default function PaymentSimulator() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const orderId = searchParams.get('orderId') ?? '';

  const [processing, setProcessing] = useState(false);
  const [selected, setSelected] = useState<Outcome | null>(null);
  const [error, setError] = useState<string | null>(null);

  const simulate = async (outcome: Outcome) => {
    setProcessing(true);
    setSelected(outcome);
    setError(null);
    const idempotencyKey = crypto.randomUUID();

    // 1.4s processing delay per spec
    await new Promise((r) => setTimeout(r, 1400));

    try {
      const res = await fetch('/api/proxy/payment/simulate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Idempotency-Key': idempotencyKey,
        },
        body: JSON.stringify({ orderId, outcome }),
      });

      const envelope = await res.json() as { code: string; message: string };

      if (res.ok && envelope.code === 'SUCCESS') {
        router.push(orderId ? `/orders/${orderId}` : '/orders');
      } else {
        setError(envelope.message ?? 'เกิดข้อผิดพลาด');
        setProcessing(false);
        setSelected(null);
      }
    } catch {
      setError('เกิดข้อผิดพลาดด้านเครือข่าย');
      setProcessing(false);
      setSelected(null);
    }
  };

  return (
    <div className="space-y-4">
      {/* Mock QR placeholder */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-6 flex flex-col items-center mb-2">
        <div className="w-32 h-32 bg-[var(--bg)] rounded-xl flex items-center justify-center mb-3 border border-[var(--divider)]">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="var(--brand)" strokeWidth="1.5">
            <rect x="3" y="3" width="8" height="8" rx="1" />
            <rect x="13" y="3" width="8" height="8" rx="1" />
            <rect x="3" y="13" width="8" height="8" rx="1" />
            <rect x="13" y="13" width="4" height="4" rx="0.5" />
            <rect x="19" y="13" width="2" height="2" rx="0.5" />
            <rect x="19" y="19" width="2" height="2" rx="0.5" />
            <rect x="13" y="19" width="4" height="2" rx="0.5" />
          </svg>
        </div>
        <p className="text-[12px] text-[var(--text-tertiary)] text-center">
          QR code จำลอง — สแกนด้วยแอปธนาคาร
        </p>
      </div>

      {error && (
        <div className="rounded-xl bg-red-50 border border-red-200 text-red-700 px-4 py-3 text-[12px]">
          {error}
        </div>
      )}

      <p className="text-[12px] text-[var(--text-secondary)] text-center font-medium">
        จำลองผลการชำระเงิน:
      </p>

      <div className="grid grid-cols-3 gap-3">
        {OUTCOMES.map(({ outcome, label, variant }) => (
          <button
            key={outcome}
            onClick={() => void simulate(outcome)}
            disabled={processing}
            className={`py-3 rounded-xl font-semibold text-[13px] transition disabled:opacity-50 ${variant}`}
          >
            {processing && selected === outcome ? (
              <span className="inline-block animate-pulse">กำลังดำเนินการ…</span>
            ) : (
              label
            )}
          </button>
        ))}
      </div>
    </div>
  );
}
