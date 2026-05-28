/**
 * lib/idempotency.ts
 *
 * Client-side idempotency key generation and in-flight de-duplication.
 * Uses crypto.randomUUID() (Node 20+ / modern browsers) exclusively.
 * Keys are held in React state, NOT localStorage.
 */

/**
 * generateIdempotencyKey — produces a UUID v4 for a single logical button click.
 * A new call on re-click generates a NEW key (fresh logical attempt).
 */
export function generateIdempotencyKey(): string {
  return crypto.randomUUID();
}

/**
 * InFlightMap — a Map that holds one in-flight Promise per named operation.
 * While a promise is pending, subsequent calls with the same key are dropped.
 *
 * Usage:
 *   const map = new InFlightMap();
 *   async function submit() {
 *     if (map.has('checkout-commit')) return; // already in flight
 *     map.set('checkout-commit', doCommit().finally(() => map.delete('checkout-commit')));
 *   }
 */
export class InFlightMap {
  private readonly _map: Map<string, Promise<unknown>> = new Map();

  has(key: string): boolean {
    return this._map.has(key);
  }

  set<T>(key: string, promise: Promise<T>): void {
    this._map.set(key, promise);
  }

  delete(key: string): void {
    this._map.delete(key);
  }

  get<T = unknown>(key: string): Promise<T> | undefined {
    return this._map.get(key) as Promise<T> | undefined;
  }
}
