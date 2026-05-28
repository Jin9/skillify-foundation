/**
 * ProductCard — Server Component
 *
 * Renders a product tile with image, name, price, and stock badge.
 */

import Link from 'next/link';

interface Product {
  id: string;
  name: string;
  price: number;
  images: string[];
  stockStatus: string;
}

interface Props {
  product: Product;
}

export default function ProductCard({ product }: Props) {
  // ฿ prefix format per spec: ฿1,290
  const priceFormatted = '฿' + new Intl.NumberFormat('th-TH', {
    maximumFractionDigits: 0,
  }).format(product.price);

  const inStock = product.stockStatus === 'IN_STOCK';

  return (
    <Link
      href={`/products/${product.id}`}
      className="group block bg-[var(--card)] border border-[var(--divider)] rounded-2xl overflow-hidden hover:shadow-sm hover:border-brand-300 transition"
    >
      {/* Image */}
      <div className="aspect-square bg-[var(--bg)] flex items-center justify-center overflow-hidden">
        {product.images[0] ? (
          <img
            src={product.images[0]}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <span className="text-4xl text-[var(--text-tertiary)]">📦</span>
        )}
      </div>

      {/* Info */}
      <div className="p-3">
        <p className="text-[12px] font-medium text-[var(--text-primary)] line-clamp-2 mb-1">{product.name}</p>
        <p className="text-[14px] font-semibold text-[var(--text-primary)] mb-1">{priceFormatted}</p>
        <span
          className={`inline-block text-[10px] px-2 py-0.5 rounded-full font-medium ${
            inStock ? 'bg-green-100 text-green-700' : 'bg-[var(--divider)] text-[var(--text-tertiary)]'
          }`}
        >
          {inStock ? 'มีสินค้า' : 'สินค้าหมด'}
        </span>
      </div>
    </Link>
  );
}
