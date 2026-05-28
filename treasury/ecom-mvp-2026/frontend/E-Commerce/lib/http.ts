/**
 * lib/http.ts
 *
 * Thin fetch wrapper for backend service calls from route handlers.
 * Handles JSON body forwarding and Bearer token attachment.
 */

export interface BackendCallOptions {
  url: string;
  bearer?: string;
  body?: unknown;
  method?: 'POST' | 'GET';
  extraHeaders?: Record<string, string>;
}

/**
 * callBackend — makes a JSON POST (default) to a backend service.
 * Forwards traceparent if available.
 * Returns the raw Response for the caller to handle.
 */
export async function callBackend(opts: BackendCallOptions): Promise<Response> {
  const { url, bearer, body, method = 'POST', extraHeaders = {} } = opts;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...extraHeaders,
  };

  if (bearer) {
    headers['Authorization'] = `Bearer ${bearer}`;
  }

  return fetch(url, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });
}

/**
 * serviceUrl — maps a path segment to the correct backend base URL.
 * The path segments mirror the /api/proxy/[domain]/[aggregate]/[action] structure.
 */
export function serviceUrl(domain: string): string {
  const map: Record<string, string | undefined> = {
    identity: process.env.BACKEND_IDENTITY_URL,
    catalog: process.env.BACKEND_CATALOG_URL,
    cart: process.env.BACKEND_CART_URL,
    checkout: process.env.BACKEND_CHECKOUT_URL,
    order: process.env.BACKEND_ORDER_URL,
    payment: process.env.BACKEND_PAYMENT_URL,
    inventory: process.env.BACKEND_INVENTORY_URL,
  };
  const base = map[domain];
  if (!base) throw new Error(`Unknown domain: ${domain}`);
  return base;
}
