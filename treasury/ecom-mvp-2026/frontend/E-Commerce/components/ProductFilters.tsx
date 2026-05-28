'use client';

/**
 * ProductFilters — Client Component
 *
 * Filter pane for /products page.
 * Updates URL search params via router.replace so links are shareable.
 * Filters: search, categoryId, minPrice, maxPrice, inStock.
 * Uses a 300ms debounce on text inputs.
 */

import { useCallback, useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

interface CurrentParams {
  categoryId?: string;
  minPrice?: string;
  maxPrice?: string;
  inStock?: string;
  search?: string;
  sort?: string;
  page?: string;
}

interface Props {
  currentParams: CurrentParams;
}

export default function ProductFilters({ currentParams }: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const debounceRef = useRef<ReturnType<typeof setTimeout>>();

  const [search, setSearch] = useState(currentParams.search ?? '');
  const [minPrice, setMinPrice] = useState(currentParams.minPrice ?? '');
  const [maxPrice, setMaxPrice] = useState(currentParams.maxPrice ?? '');
  const [inStock, setInStock] = useState(currentParams.inStock === 'true');
  const [sort, setSort] = useState(currentParams.sort ?? 'created_desc');

  const applyFilters = useCallback(
    (overrides: Partial<CurrentParams> = {}) => {
      const params = new URLSearchParams(searchParams.toString());
      const merged: Record<string, string | undefined> = {
        search,
        minPrice,
        maxPrice,
        inStock: inStock ? 'true' : undefined,
        sort,
        page: '1', // reset to first page on filter change
        ...overrides,
      };
      Object.entries(merged).forEach(([key, value]) => {
        if (value) {
          params.set(key, value);
        } else {
          params.delete(key);
        }
      });
      router.replace(`/products?${params.toString()}`);
    },
    [router, searchParams, search, minPrice, maxPrice, inStock, sort],
  );

  // Debounce text inputs
  const debouncedApply = useCallback(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => applyFilters(), 300);
  }, [applyFilters]);

  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  return (
    <div className="bg-[var(--card)] border border-[var(--divider)] rounded-2xl p-4 space-y-4">
      <h3 className="text-[13px] font-semibold text-[var(--text-primary)]">ตัวกรอง</h3>

      {/* Search */}
      <div>
        <label className="block text-[11px] font-medium text-[var(--text-tertiary)] mb-1">ค้นหาสินค้า</label>
        <input
          type="text"
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            debouncedApply();
          }}
          placeholder="ชื่อสินค้า…"
          className="w-full border border-[var(--divider)] rounded-xl px-3 py-2 text-[13px] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
        />
      </div>

      {/* Price range */}
      <div>
        <label className="block text-[11px] font-medium text-[var(--text-tertiary)] mb-1">ราคา (บาท)</label>
        <div className="flex gap-2">
          <input
            type="number"
            value={minPrice}
            onChange={(e) => {
              setMinPrice(e.target.value);
              debouncedApply();
            }}
            placeholder="ต่ำสุด"
            className="w-1/2 border border-[var(--divider)] rounded-xl px-2 py-2 text-[13px] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
          />
          <input
            type="number"
            value={maxPrice}
            onChange={(e) => {
              setMaxPrice(e.target.value);
              debouncedApply();
            }}
            placeholder="สูงสุด"
            className="w-1/2 border border-[var(--divider)] rounded-xl px-2 py-2 text-[13px] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
          />
        </div>
      </div>

      {/* In Stock */}
      <label className="flex items-center gap-2 cursor-pointer">
        <input
          type="checkbox"
          checked={inStock}
          onChange={(e) => {
            setInStock(e.target.checked);
            applyFilters({ inStock: e.target.checked ? 'true' : undefined });
          }}
          className="rounded text-[var(--brand)]"
        />
        <span className="text-[12px] text-[var(--text-secondary)]">มีสินค้าเท่านั้น</span>
      </label>

      {/* Sort */}
      <div>
        <label className="block text-[11px] font-medium text-[var(--text-tertiary)] mb-1">เรียงลำดับ</label>
        <select
          value={sort}
          onChange={(e) => {
            setSort(e.target.value);
            applyFilters({ sort: e.target.value });
          }}
          className="w-full border border-[var(--divider)] rounded-xl px-2 py-2 text-[13px] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
        >
          <option value="created_desc">ใหม่ล่าสุด</option>
          <option value="price_asc">ราคา: ต่ำ-สูง</option>
          <option value="price_desc">ราคา: สูง-ต่ำ</option>
        </select>
      </div>

      {/* Reset */}
      <button
        onClick={() => {
          setSearch('');
          setMinPrice('');
          setMaxPrice('');
          setInStock(false);
          setSort('created_desc');
          router.replace('/products');
        }}
        className="w-full text-[11px] text-[var(--text-tertiary)] hover:text-red-500 transition"
      >
        ล้างตัวกรองทั้งหมด
      </button>
    </div>
  );
}
