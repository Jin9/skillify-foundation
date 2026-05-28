'use client';

/**
 * AddToCartButton — Client Component
 *
 * Posts to /api/proxy/cart/cart/add-item.
 * Disabled while in-flight and when product is out of stock.
 * Shows toast on success/error.
 */

import { useState } from 'react';

interface Props {
  productId: string;
  disabled?: boolean;
}

export default function AddToCartButton({ productId, disabled = false }: Props) {
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  const showToast = (type: 'success' | 'error', message: string) => {
    setToast({ type, message });
    setTimeout(() => setToast(null), 3000);
  };

  const handleAddToCart = async () => {
    if (disabled || loading) return;
    setLoading(true);
    try {
      const res = await fetch('/api/proxy/cart/cart/add-item', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId, quantity: 1 }),
      });

      if (res.status === 401) {
        window.location.href = `/login?next=${encodeURIComponent(window.location.pathname)}`;
        return;
      }

      const envelope = await res.json() as { code: string; message: string; traceId?: string };

      if (res.ok && envelope.code === 'SUCCESS') {
        showToast('success', 'เพิ่มลงตะกร้าแล้ว');
      } else {
        showToast('error', envelope.message ?? 'ไม่สามารถเพิ่มลงตะกร้าได้');
      }
    } catch {
      showToast('error', 'เกิดข้อผิดพลาดด้านเครือข่าย กรุณาลองใหม่');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative">
      <button
        onClick={() => void handleAddToCart()}
        disabled={disabled || loading}
        className={`w-full py-3 rounded-xl font-semibold transition text-[14px] ${
          disabled
            ? 'bg-[var(--divider)] text-[var(--text-tertiary)] cursor-not-allowed'
            : loading
            ? 'bg-brand-400 text-white cursor-wait'
            : 'bg-[var(--brand)] hover:opacity-90 text-white'
        }`}
      >
        {disabled ? 'สินค้าหมด' : loading ? 'กำลังเพิ่ม…' : 'เพิ่มลงตะกร้า'}
      </button>

      {toast && (
        <div
          className={`absolute -top-12 left-0 right-0 text-center text-[12px] py-2 px-4 rounded-xl font-medium ${
            toast.type === 'success'
              ? 'bg-green-100 text-green-700'
              : 'bg-red-100 text-red-700'
          }`}
        >
          {toast.message}
        </div>
      )}
    </div>
  );
}
