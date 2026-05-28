/**
 * /products — Product Listing (Server Component shell + Client filter pane)
 *
 * URL search params drive the SC fetch; filters are managed by the Client
 * ProductFilters component which uses router.replace to keep links shareable.
 */

import type { Metadata } from 'next';
import ProductCard from '@/components/ProductCard';
import ProductFilters from '@/components/ProductFilters';

export const metadata: Metadata = { title: 'Products' };

const CATALOG_URL = process.env.BACKEND_CATALOG_URL ?? '';

interface Product {
  id: string;
  name: string;
  price: number;
  images: string[];
  stockStatus: string;
}

interface ProductListResult {
  items: Product[];
  total: number;
  page: number;
  limit: number;
}

interface SearchParams {
  page?: string;
  limit?: string;
  sort?: string;
  categoryId?: string;
  minPrice?: string;
  maxPrice?: string;
  inStock?: string;
  search?: string;
}

async function fetchProducts(params: SearchParams): Promise<ProductListResult> {
  const body = {
    page: Number(params.page ?? 1),
    limit: Number(params.limit ?? 20),
    sort: params.sort ?? 'created_desc',
    ...(params.categoryId ? { categoryId: params.categoryId } : {}),
    ...(params.minPrice ? { minPrice: Number(params.minPrice) } : {}),
    ...(params.maxPrice ? { maxPrice: Number(params.maxPrice) } : {}),
    ...(params.inStock === 'true' ? { inStock: true } : {}),
    ...(params.search ? { search: params.search } : {}),
  };

  try {
    const res = await fetch(`${CATALOG_URL}/api/v1/catalog/product/list`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
      cache: 'no-store',
    });
    if (!res.ok) return { items: [], total: 0, page: 1, limit: 20 };
    const envelope = await res.json() as { data: ProductListResult | null };
    return envelope.data ?? { items: [], total: 0, page: 1, limit: 20 };
  } catch {
    return { items: [], total: 0, page: 1, limit: 20 };
  }
}

export default async function ProductsPage({
  searchParams,
}: {
  searchParams: SearchParams;
}) {
  const result = await fetchProducts(searchParams);
  const page = Number(searchParams.page ?? 1);
  const limit = Number(searchParams.limit ?? 20);
  const totalPages = Math.ceil(result.total / limit);

  return (
    <div className="px-4 py-4 max-w-[390px] mx-auto">
      <h1 className="text-[17px] font-semibold text-[var(--text-primary)] mb-3">ค้นหาสินค้า</h1>

      {/* Filter pane — Client Component (full width on mobile) */}
      <div className="mb-4">
        <ProductFilters currentParams={searchParams} />
      </div>

      {/* Product grid */}
      {result.items.length === 0 ? (
        <p className="text-[var(--text-secondary)] text-[13px]">ไม่พบสินค้าที่ตรงกับเงื่อนไข</p>
      ) : (
        <>
          <p className="text-[11px] text-[var(--text-tertiary)] mb-3">
            แสดง {result.items.length} จาก {result.total} รายการ
          </p>
          <div className="grid grid-cols-2 gap-3">
            {result.items.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <nav className="flex items-center justify-center gap-2 mt-6">
              {Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
                <a
                  key={p}
                  href={`/products?${new URLSearchParams({ ...searchParams, page: String(p) }).toString()}`}
                  className={`px-3 py-1.5 rounded-lg border text-[12px] font-medium transition ${
                    p === page
                      ? 'bg-[var(--brand)] text-white border-[var(--brand)]'
                      : 'border-[var(--divider)] text-[var(--text-secondary)] hover:border-brand-500'
                  }`}
                >
                  {p}
                </a>
              ))}
            </nav>
          )}
        </>
      )}
    </div>
  );
}
