/**
 * / — Home page (Server Component)
 *
 * Fetches top categories and top products (newest 12) server-side.
 * No auth required — public page.
 */

import Link from 'next/link';
import ProductCard from '@/components/ProductCard';

const CATALOG_URL = process.env.BACKEND_CATALOG_URL ?? '';

interface Category {
  id: string;
  name: string;
  slug: string;
}

interface Product {
  id: string;
  name: string;
  price: number;
  images: string[];
  stockStatus: string;
}

async function fetchCategories(): Promise<Category[]> {
  try {
    const res = await fetch(`${CATALOG_URL}/api/v1/catalog/category/list`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ activeOnly: true }),
      next: { revalidate: 60 },
    });
    if (!res.ok) return [];
    const envelope = await res.json() as { code: string; data: Category[] | null };
    return envelope.data ?? [];
  } catch {
    return [];
  }
}

async function fetchTopProducts(): Promise<Product[]> {
  try {
    const res = await fetch(`${CATALOG_URL}/api/v1/catalog/product/list`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ page: 1, limit: 12, sort: 'created_desc' }),
      next: { revalidate: 30 },
    });
    if (!res.ok) return [];
    const envelope = await res.json() as { code: string; data: { items: Product[] } | null };
    return envelope.data?.items ?? [];
  } catch {
    return [];
  }
}

export default async function HomePage() {
  const [categories, products] = await Promise.all([fetchCategories(), fetchTopProducts()]);

  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      {/* Hero Banner */}
      <section className="rounded-2xl bg-[var(--brand)] text-white px-5 py-8 mb-6 text-center">
        <h1 className="text-[22px] font-bold mb-2 leading-snug">ยินดีต้อนรับสู่ ShopPilot</h1>
        <p className="text-[13px] mb-5 opacity-90">
          สินค้าหลายพันรายการ ชำระง่าย ติดตามออเดอร์แบบเรียลไทม์
        </p>
        <Link
          href="/products"
          className="inline-block bg-white text-[var(--brand)] font-semibold px-6 py-2.5 rounded-xl text-[14px] hover:opacity-90 transition"
        >
          เลือกซื้อสินค้า
        </Link>
      </section>

      {/* Category Grid */}
      {categories.length > 0 && (
        <section className="mb-6">
          <h2 className="text-[15px] font-semibold text-[var(--text-primary)] mb-3">หมวดหมู่</h2>
          <div className="grid grid-cols-3 gap-3">
            {categories.map((cat) => (
              <Link
                key={cat.id}
                href={`/products?categoryId=${cat.id}`}
                className="flex flex-col items-center justify-center p-3 rounded-2xl bg-[var(--card)] border border-[var(--divider)] hover:border-brand-500 transition text-center"
              >
                <span className="text-2xl mb-1">🏷️</span>
                <span className="text-[11px] font-medium text-[var(--text-primary)] leading-tight">{cat.name}</span>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* New Arrivals */}
      <section>
        <h2 className="text-[15px] font-semibold text-[var(--text-primary)] mb-3">สินค้าใหม่</h2>
        {products.length === 0 ? (
          <p className="text-[var(--text-secondary)] text-[13px]">ยังไม่มีสินค้าในขณะนี้</p>
        ) : (
          <div className="grid grid-cols-2 gap-3">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
