# State & Data

## Contents
- Server state (TanStack Query)
- Client state (Zustand)
- Context — what it's for
- URL state
- Forms (RHF + Zod)
- Caching, mutations, optimistic updates
- Risky patterns

## Low-risk use

Follow the repo's existing state and form libraries first. Use these recommendations for greenfield work, migration plans, or codebases that already use the named tools.

## Server state — TanStack Query

- Query key = stable, serializable identity. Convention: `[domain, id, params]`.
- `staleTime` per query type: read-heavy / cheap = 5–30s; expensive / static = minutes–hours.
- `gcTime` (garbage collect) higher than `staleTime`; default 5 min is usually fine.
- Use `select` for derived data — keeps cache canonical, components get the slice.
- Mutations: invalidate by key prefix; or `setQueryData` for optimistic patches.
- One client per app via `QueryClientProvider`. SSR: serialize with `dehydrate` / `hydrate`.

## Client state — Zustand (when already adopted or greenfield)

```ts
const useStore = create<State>()((set) => ({
  filters: { status: 'all' },
  setStatus: (s) => set((st) => ({ filters: { ...st.filters, status: s } })),
}))
```

- Selectors with shallow compare to avoid re-renders: `useStore(s => s.filters.status)`.
- Slice large stores by domain.
- Persist with `persist` middleware only when truly needed (filters, theme). Do not persist auth tokens in client-accessible storage.

When to use Redux Toolkit instead: deep undo / redo, time-travel debug, very large team needing strict patterns + middleware ecosystem.

## Context — what it's for

- Auth identity, theme, locale, feature flags. Read often, written rarely.
- Avoid high-frequency updates (cursor position, form draft, scroll). They can cause whole-subtree re-renders.

## URL state

- Filters, sort, pagination, tab — these belong in the URL, not in a store.
- Shareable, bookmarkable, back-button friendly.
- Use `useSearchParams` (Next/Remix) or `nuqs` for typed search params.

## Forms — React Hook Form + Zod (when already adopted or greenfield)

```ts
const schema = z.object({
  email: z.string().email(),
  amount: z.number().positive(),
})
type FormValues = z.infer<typeof schema>

const { register, handleSubmit, formState } = useForm<FormValues>({
  resolver: zodResolver(schema),
})
```

- One source of truth for shape: the Zod schema. Derive types via `z.infer`.
- Use `Controller` only for fields RHF can't `register` directly (custom inputs, libraries).
- Submit handler is async; surface server errors via `setError`.
- Re-use the schema server-side for end-to-end validation.

## Caching, mutations, optimistic updates

Optimistic update pattern:

1. `onMutate`: snapshot + write the optimistic value via `setQueryData`.
2. `onError`: roll back to snapshot.
3. `onSettled`: invalidate to reconcile with server.

For lists: invalidate the list query AND the detail query on change.

## Risky Patterns

- Storing server data in Zustand / Redux.
- `useEffect` to sync server data into local state — use TanStack Query.
- Reading the URL via `window.location` in components — use the router hook.
- Form state mirrored into both RHF and `useState` — pick one.
- Mutating a Zustand store directly — always go through the setter.
- `useState` for anything that should round-trip via the URL.
