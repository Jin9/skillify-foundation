'use client';

/**
 * /login — Login page (Client Component)
 *
 * react-hook-form + zod for validation.
 * POSTs to /api/proxy/identity/auth/login.
 * On success, route handler sets cookies; client redirects to `next` param or /.
 */

import { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import Link from 'next/link';

const loginSchema = z.object({
  email: z.string().email('Enter a valid email address').max(254),
  password: z.string().min(8, 'Password must be at least 8 characters').max(72),
});

type LoginFields = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const nextParam = searchParams.get('next') ?? '/';
  const [apiError, setApiError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFields>({ resolver: zodResolver(loginSchema) });

  const onSubmit = async (data: LoginFields) => {
    setApiError(null);
    try {
      const res = await fetch('/api/proxy/identity/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: data.email, password: data.password }),
      });

      const envelope = await res.json() as { code: string; message: string };

      if (res.ok && envelope.code === 'SUCCESS') {
        router.push(nextParam);
        router.refresh();
        return;
      }

      setApiError(envelope.message ?? 'Login failed. Please try again.');
    } catch {
      setApiError('Network error. Please check your connection and try again.');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--bg)] px-4">
      <div className="w-full max-w-[390px] bg-[var(--card)] rounded-2xl shadow-sm border border-[var(--divider)] p-6">
        <h1 className="text-[22px] font-bold text-[var(--text-primary)] mb-1">เข้าสู่ระบบ</h1>
        <p className="text-[var(--text-secondary)] text-[13px] mb-6">
          ยังไม่มีบัญชี?{' '}
          <Link href="/register" className="text-[var(--brand)] hover:underline font-medium">
            สมัครสมาชิก
          </Link>
        </p>

        {apiError && (
          <div className="mb-4 rounded-xl bg-red-50 border border-red-200 text-red-700 px-4 py-3 text-[12px]">
            {apiError}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label htmlFor="email" className="block text-[12px] font-medium text-[var(--text-secondary)] mb-1">
              อีเมล
            </label>
            <input
              {...register('email')}
              id="email"
              type="email"
              autoComplete="email"
              className="w-full border border-[var(--divider)] rounded-xl px-3 py-2.5 text-[14px] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
              placeholder="your@email.com"
            />
            {errors.email && (
              <p className="mt-1 text-[11px] text-red-600">{errors.email.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="password" className="block text-[12px] font-medium text-[var(--text-secondary)] mb-1">
              รหัสผ่าน
            </label>
            <input
              {...register('password')}
              id="password"
              type="password"
              autoComplete="current-password"
              className="w-full border border-[var(--divider)] rounded-xl px-3 py-2.5 text-[14px] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
            />
            {errors.password && (
              <p className="mt-1 text-[11px] text-red-600">{errors.password.message}</p>
            )}
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full bg-[var(--brand)] hover:opacity-90 disabled:opacity-50 text-white font-semibold py-3 rounded-xl transition text-[14px]"
          >
            {isSubmitting ? 'กำลังเข้าสู่ระบบ…' : 'เข้าสู่ระบบ'}
          </button>
        </form>
      </div>
    </div>
  );
}
