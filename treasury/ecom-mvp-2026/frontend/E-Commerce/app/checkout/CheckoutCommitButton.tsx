'use client';

/**
 * CheckoutCommitButton — Client Component.
 * Generates a crypto.randomUUID() idempotency key per click.
 * POSTs to /api/proxy/checkout/commit.
 * On success: redirects to /payment.
 */

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function CheckoutCommitButton() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCommit = async () => {
    setLoading(true);
    setError(null);
    const idempotencyKey = crypto.randomUUID();
    try {
      const res = await fetch('/api/proxy/checkout/commit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Idempotency-Key': idempotencyKey,
        },
        body: JSON.stringify({}),
      });

      if (res.status === 401) {
        router.push('/login?next=/checkout');
        return;
      }

      const envelope = await res.json() as { code: string; message: string; data?: { orderId?: string } };

      if (res.ok && envelope.code === 'SUCCESS') {
        const orderId = envelope.data?.orderId ?? '';
        router.push(`/payment${orderId ? `?orderId=${orderId}` : ''}`);
      } else {
        setError(envelope.message ?? 'เกิดข้อผิดพลาด กรุณาลองใหม่');
      }
    } catch {
      setError('เกิดข้อผิดพลาดด้านเครือข่าย กรุณาลองใหม่');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      {error && (
        <div className="mb-3 rounded-xl bg-red-50 border border-red-200 text-red-700 px-4 py-3 text-[12px]">
          {error}
        </div>
      )}
      <button
        onClick={() => void handleCommit()}
        disabled={loading}
        className="w-full bg-[var(--brand)] hover:opacity-90 disabled:opacity-50 text-white font-semibold py-3 rounded-xl transition text-[14px]"
      >
        {loading ? 'กำลังดำเนินการ…' : 'ยืนยันคำสั่งซื้อ'}
      </button>
    </div>
  );
}
