/**
 * /products/[id] — Product Detail (Server Component)
 *
 * Fetches product detail server-side. The "Add to Cart" button is a small
 * Client Component so it can POST to /api/proxy/cart/cart/add-item.
 */

import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import AddToCartButton from '@/components/AddToCartButton';

const CATALOG_URL = process.env.BACKEND_CATALOG_URL ?? '';

interface ProductDetail {
  id: string;
  sku: string;
  name: string;
  description: string;
  images: string[];
  price: number;
  stockStatus: string;
  categoryId: string;
  categoryName?: string;
}

async function fetchProduct(id: string): Promise<ProductDetail | null> {
  try {
    const res = await fetch(`${CATALOG_URL}/api/v1/catalog/product/detail`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ productId: id }),
      next: { revalidate: 30 },
    });
    if (res.status === 404) return null;
    const envelope = await res.json() as { code: string; data: ProductDetail | null };
    if (envelope.code !== 'SUCCESS') return null;
    return envelope.data;
  } catch {
    return null;
  }
}

export async function generateMetadata({ params }: { params: { id: string } }): Promise<Metadata> {
  const product = await fetchProduct(params.id);
  return { title: product?.name ?? 'Product Detail' };
}

export default async function ProductDetailPage({ params }: { params: { id: string } }) {
  const product = await fetchProduct(params.id);
  if (!product) notFound();

  // ฿ prefix per spec
  const priceFormatted = '฿' + new Intl.NumberFormat('th-TH', {
    maximumFractionDigits: 0,
  }).format(product.price);

  const inStock = product.stockStatus === 'IN_STOCK';

  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      {/* Back */}
      <a href="/products" className="inline-flex items-center gap-1 text-[var(--text-secondary)] text-[13px] mb-4 hover:text-[var(--brand)]">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
          <polyline points="15 18 9 12 15 6" />
        </svg>
        สินค้าทั้งหมด
      </a>

      {/* Product Image */}
      <div className="aspect-square bg-[var(--bg)] rounded-2xl overflow-hidden flex items-center justify-center mb-4">
        {product.images[0] ? (
          <img
            src={product.images[0]}
            alt={product.name}
            className="object-cover w-full h-full"
          />
        ) : (
          <span className="text-[var(--text-tertiary)] text-6xl">📦</span>
        )}
      </div>

      {/* Product Info */}
      <div className="flex flex-col gap-3">
        {product.categoryName && (
          <span className="text-[11px] text-[var(--brand)] font-medium uppercase tracking-wide">
            {product.categoryName}
          </span>
        )}
        <h1 className="text-[22px] font-bold text-[var(--text-primary)] leading-snug">{product.name}</h1>
        <p className="text-[11px] text-[var(--text-tertiary)]">SKU: {product.sku}</p>
        <p className="text-[22px] font-semibold text-[var(--text-primary)]">{priceFormatted}</p>

        <span
          className={`inline-flex w-fit items-center px-3 py-1 rounded-full text-[12px] font-medium ${
            inStock ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-600'
          }`}
        >
          {inStock ? 'มีสินค้า' : 'สินค้าหมด'}
        </span>

        <p className="text-[13px] text-[var(--text-secondary)] leading-relaxed">{product.description}</p>

        {/* Add to Cart — Client Component */}
        <AddToCartButton productId={product.id} disabled={!inStock} />
      </div>
    </div>
  );
}
