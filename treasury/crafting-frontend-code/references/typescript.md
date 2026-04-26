# TypeScript

## Contents
- Compiler defaults
- Type-first patterns
- Discriminated unions
- Branded types
- Inference & narrowing
- Generics
- API contracts (codegen)
- Risky patterns

## Low-risk use

Respect the repo's current TypeScript settings. Tighten compiler options only for new packages or with a staged migration plan, because flipping strict flags in a legacy app can create a large unrelated diff.

## Compiler defaults

Target `tsconfig.json` options for new or owned code:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noImplicitOverride": true,
    "verbatimModuleSyntax": true,
    "moduleResolution": "bundler"
  }
}
```

- `strict: true` is the target for maintained code, not a reason to rewrite unrelated legacy areas.
- `noUncheckedIndexedAccess` catches the "array[i] is T but actually undefined" class of bugs.

## Type-first patterns

- `type` for data shapes, unions, mapped / conditional types.
- `interface` for extensible contracts (component props, plugin shapes).
- Explicit return types on exported functions and React components (helps refactors and error messages).

## Discriminated unions

Model state machines as unions, not optional fields.

```ts
type FetchState<T> =
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'success'; data: T }
  | { status: 'error'; error: Error }
```

Compiler enforces handling each case via `switch` exhaustiveness:

```ts
function assertNever(x: never): never { throw new Error('unreachable') }

switch (state.status) {
  case 'idle':    return null
  case 'loading': return <Spinner/>
  case 'success': return <View data={state.data}/>
  case 'error':   return <Err error={state.error}/>
  default: assertNever(state)
}
```

## Branded types

Stop the "userId vs orderId both `string`" class of bug.

```ts
type Brand<T, B> = T & { readonly __brand: B }
type UserId  = Brand<string, 'UserId'>
type OrderId = Brand<string, 'OrderId'>

const asUserId = (s: string) => s as UserId
```

Brand at the boundary (parsers / API layer); the body of the code stays clean.

## Inference & narrowing

- Let inference work — no annotations on local `const`s.
- Use type predicates (`x is Foo`) and assertion functions for non-trivial narrowing.
- Prefer `satisfies` over `as` for literal config:
  ```ts
  const routes = { home: '/', settings: '/settings' } satisfies Record<string, string>
  ```

## Generics

- Add a generic only when the consumer benefits. "Just in case" generics hurt readability.
- Constrain: `<T extends { id: string }>`, not bare `<T>`.
- Default type params for ergonomic call sites: `<T = unknown>`.

## API contracts — prefer codegen

- OpenAPI schema → `openapi-typescript` or `orval` → typed client.
- GraphQL schema → `graphql-codegen` → typed operations + hooks.
- Hand-written `Response` types drift easily; keep them local to tests or temporary adapter boundaries when codegen is unavailable.

## Risky Patterns

- Unchecked `any`. Use `unknown` and narrow.
- `as` casts outside parsers / boundaries.
- `// @ts-ignore` without a comment explaining the gap and a follow-up TODO.
- Re-declaring types that already exist in the codegen output.
- Optional chaining masking a real type bug (`obj?.field` when `obj` is non-nullable per types).
- `Function` or `object` as types — use specific signatures.
