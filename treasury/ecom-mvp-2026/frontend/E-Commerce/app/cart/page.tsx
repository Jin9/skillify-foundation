'use client';

/**
 * /cart — Cart page (Client Component)
 *
 * Fetches cart data from /api/proxy/cart/cart/read.
 * Shows line items + subtotal + Checkout button (stub link to /checkout for v2).
 */

import { useEffect, useState } from 'react';
import Link from 'next/link';
import CartLineItem from '@/components/CartLineItem';

interface CartItem {
  productId: string;
  sku: string;
  productName: string;
  productImageUrl?: string;
  price: number;
  quantity: number;
  subtotal: number;
  checkoutable: boolean;
  updatedAt: string;
}

interface Cart {
  cartId: string;
  items: CartItem[];
  subtotal: number;
}

function formatTHB(amount: number): string {
  return new Intl.NumberFormat('th-TH', { style: 'currency', currency: 'THB', maximumFractionDigits: 0 }).format(amount);
}

export default function CartPage() {
  const [cart, setCart] = useState<Cart | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchCart = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch('/api/proxy/cart/cart/read', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}),
      });

      if (res.status === 401) {
        window.location.href = '/login?next=/cart';
        return;
      }

      const envelope = await res.json() as { code: string; message: string; data: Cart | null };
      if (envelope.code === 'SUCCESS' && envelope.data) {
        setCart(envelope.data);
      } else {
        setError(envelope.message ?? 'Failed to load cart.');
      }
    } catch {
      setError('Network error. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchCart();
  }, []);

  if (loading) {
    return (
      <div className="px-4 py-8 max-w-[390px] mx-auto">
        <div className="animate-pulse space-y-3">
          {[1, 2, 3].map((n) => (
            <div key={n} className="h-20 bg-[var(--divider)] rounded-2xl" />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="px-4 py-12 max-w-[390px] mx-auto text-center">
        <p className="text-red-600 text-[13px] mb-4">{error}</p>
        <button onClick={() => void fetchCart()} className="text-[var(--brand)] text-[13px] hover:underline">
          ลองใหม่อีกครั้ง
        </button>
      </div>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="px-4 py-16 max-w-[390px] mx-auto text-center">
        <p className="text-5xl mb-4">🛒</p>
        <h2 className="text-[17px] font-semibold text-[var(--text-primary)] mb-2">ตะกร้าว่างเปล่า</h2>
        <Link href="/products" className="text-[var(--brand)] text-[13px] hover:underline">
          เลือกซื้อสินค้า
        </Link>
      </div>
    );
  }

  const hasUncheckoutableItems = cart.items.some((item) => !item.checkoutable);

  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      <h1 className="text-[17px] font-semibold text-[var(--text-primary)] mb-4">ตะกร้าสินค้า</h1>

      <div className="space-y-3 mb-4">
        {cart.items.map((item) => (
          <CartLineItem key={item.productId} item={item} onCartChange={() => void fetchCart()} />
        ))}
      </div>

      {/* Cart Summary */}
      <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4">
        <div className="flex justify-between items-center mb-2">
          <span className="text-[13px] text-[var(--text-secondary)]">ยอดรวม</span>
          <span className="text-[15px] font-semibold text-[var(--text-primary)]">{formatTHB(cart.subtotal)}</span>
        </div>
        <p className="text-[11px] text-[var(--text-tertiary)] mb-4">
          ค่าจัดส่งคำนวณที่หน้าชำระเงิน
        </p>

        {hasUncheckoutableItems && (
          <div className="mb-4 rounded-xl bg-amber-50 border border-amber-200 text-amber-700 px-4 py-3 text-[12px]">
            มีสินค้าบางรายการที่ไม่สามารถสั่งซื้อได้ กรุณาลบออกก่อน
          </div>
        )}

        <Link
          href="/checkout"
          className={`block w-full text-center font-semibold py-3 rounded-xl transition text-[14px] ${
            hasUncheckoutableItems
              ? 'bg-[var(--divider)] text-[var(--text-tertiary)] cursor-not-allowed pointer-events-none'
              : 'bg-[var(--brand)] hover:opacity-90 text-white'
          }`}
        >
          ดำเนินการชำระเงิน
        </Link>
      </div>
    </div>
  );
}
