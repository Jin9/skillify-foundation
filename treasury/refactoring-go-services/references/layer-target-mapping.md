# Layer target mapping (DDD boundaries)

Loaded on demand by `refactoring-go-services` step 3 ("Apply minimal change.
Check boundary layers."). Refactored code must strictly align with these
boundaries:

- **[Handler Layer] / Transport**
  - **Duties**: Validates raw input interfaces/JSON, builds transport responses (HTTP).
  - **Restriction**: NO core business logic, NO repository calls.

- **[Orchestrator Layer] / Application**
  - **Duties**: Standard flow control. Retrieves objects from the repository layer, delegates work to domain methods, and persists results.
  - **Restriction**: NO domain decision-making (calculating discounts, validating state transitions).

- **[Domain Layer] / Core**
  - **Duties**: Pure Go structs, strict types, specific domain errors, and value objects. All dynamic business logic lives here locally via pointer methods.
  - **Restriction**: ZERO external framework dependencies (no HTTP routers, no SQL drivers, no messaging queues).

- **[Adapter Layer] / Infrastructure**
  - **Duties**: Implementation of outward/inward protocols (Postgres repository, third-party REST client, gRPC handler).
  - **Restriction**: Exclusively maps transport/vendor data formats to pure domain interfaces.
