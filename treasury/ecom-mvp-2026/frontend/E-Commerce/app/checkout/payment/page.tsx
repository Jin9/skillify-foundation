/**
 * /checkout/payment — redirects to /payment (the canonical payment route).
 * The full payment simulation screen lives at /payment.
 */

import { redirect } from 'next/navigation';

export const metadata = { title: 'ชำระเงิน' };

export default function CheckoutPaymentPage() {
  redirect('/payment');
}
