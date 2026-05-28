'use client';

/**
 * CartLineItem — Client Component
 *
 * Renders a single cart line with quantity stepper, remove button, line subtotal,
 * and a non-checkoutable badge when checkoutable=false.
 */

import { useState } from 'react';

interface CartItem {
  productId: string;
  sku: string;
  productName: string;
  productImageUrl?: string;
  price: number;
  quantity: number;
  subtotal: number;
  checkoutable: boolean;
}

interface Props {
  item: CartItem;
  onCartChange: () => void;
}

function formatTHB(amount: number): string {
  return new Intl.NumberFormat('th-TH', { style: 'currency', currency: 'THB', maximumFractionDigits: 0 }).format(amount);
}

export default function CartLineItem({ item, onCartChange }: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const updateQty = async (newQty: number) => {
    if (newQty < 1) return;
    setLoading(true);
    setError(null);
    try {
      const res = await fetch('/api/proxy/cart/cart/update-item', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: item.productId, quantity: newQty }),
      });
      const envelope = await res.json() as { code: string; message: string };
      if (res.ok && envelope.code === 'SUCCESS') {
        onCartChange();
      } else {
        setError(envelope.message ?? 'อัปเดตไม่สำเร็จ');
      }
    } catch {
      setError('เกิดข้อผิดพลาดด้านเครือข่าย');
    } finally {
      setLoading(false);
    }
  };

  const removeItem = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch('/api/proxy/cart/cart/remove-item', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: item.productId }),
      });
      const envelope = await res.json() as { code: string; message: string };
      if (res.ok && envelope.code === 'SUCCESS') {
        onCartChange();
      } else {
        setError(envelope.message ?? 'ลบสินค้าไม่สำเร็จ');
      }
    } catch {
      setError('เกิดข้อผิดพลาดด้านเครือข่าย');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={`bg-[var(--card)] border rounded-2xl p-3 flex gap-3 ${!item.checkoutable ? 'border-amber-300 bg-amber-50' : 'border-[var(--divider)]'}`}>
      {/* Image */}
      <div className="w-16 h-16 flex-shrink-0 bg-[var(--bg)] rounded-xl overflow-hidden flex items-center justify-center">
        {item.productImageUrl ? (
          <img src={item.productImageUrl} alt={item.productName} className="w-full h-full object-cover" />
        ) : (
          <span className="text-2xl">📦</span>
        )}
      </div>

      {/* Details */}
      <div className="flex-1 min-w-0">
        <p className="text-[13px] font-medium text-[var(--text-primary)] truncate">{item.productName}</p>
        <p className="text-[10px] text-[var(--text-tertiary)] mb-0.5">SKU: {item.sku}</p>
        <p className="text-[12px] text-[var(--text-secondary)]">{formatTHB(item.price)} / ชิ้น</p>

        {!item.checkoutable && (
          <span className="inline-block mt-1 text-[10px] bg-amber-100 text-amber-700 px-2 py-0.5 rounded-full">
            สินค้าไม่พร้อมจำหน่าย
          </span>
        )}

        {error && <p className="mt-1 text-[10px] text-red-600">{error}</p>}
      </div>

      {/* Controls */}
      <div className="flex flex-col items-end justify-between">
        <p className="text-[13px] font-semibold text-[var(--text-primary)]">{formatTHB(item.subtotal)}</p>
        <div className="flex items-center gap-1.5">
          <button
            onClick={() => void updateQty(item.quantity - 1)}
            disabled={loading || item.quantity <= 1}
            className="w-7 h-7 rounded-full border border-[var(--divider)] text-[var(--text-secondary)] hover:bg-[var(--bg)] disabled:opacity-40 transition flex items-center justify-center text-[13px]"
          >
            -
          </button>
          <span className="text-[13px] w-5 text-center">{item.quantity}</span>
          <button
            onClick={() => void updateQty(item.quantity + 1)}
            disabled={loading}
            className="w-7 h-7 rounded-full border border-[var(--divider)] text-[var(--text-secondary)] hover:bg-[var(--bg)] disabled:opacity-40 transition flex items-center justify-center text-[13px]"
          >
            +
          </button>
          <button
            onClick={() => void removeItem()}
            disabled={loading}
            className="ml-1 text-[10px] text-red-500 hover:text-red-700 disabled:opacity-40 transition"
          >
            ลบ
          </button>
        </div>
      </div>
    </div>
  );
}
