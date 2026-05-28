/**
 * /admin — STUB (Coming in v2)
 *
 * Planned admin pages:
 *   /admin/products       — product list with DRAFT/INACTIVE (admin variant)
 *   /admin/products/new   — product create form (sku, name, price, categoryId, images, status)
 *   /admin/products/[id]  — product edit + soft-delete
 *   /admin/categories     — category CRUD (create, update, deactivate)
 *   /admin/inventory      — stock adjustment form (sku, delta, reason)
 *   /admin/orders         — order list with status/date/search filters
 *   /admin/orders/[id]    — order detail + AdminStatusTransitionPanel
 *
 * All admin pages require role=ADMIN; CUSTOMER tokens receive 403.
 * Layout will run the session probe and redirect non-admins to /login or /403.
 */

import Link from 'next/link';

export const metadata = { title: 'Admin' };

export default function AdminPage() {
  return (
    <div className="max-w-3xl mx-auto px-4 py-12 text-center">
      <h1 className="text-2xl font-bold text-slate-800 mb-4">Admin Dashboard</h1>
      <p className="text-slate-500 mb-6">
        Admin back office is <strong>coming in v2</strong>.
      </p>
      <p className="text-sm text-slate-400 mb-4">
        Planned: product CRUD, category management, stock adjustment, order management with
        PAID→PACKING→SHIPPED→DELIVERED transitions (SHIPPED requires tracking number).
      </p>
      <Link href="/" className="text-brand-600 hover:underline text-sm">
        Back to Home
      </Link>
    </div>
  );
}
