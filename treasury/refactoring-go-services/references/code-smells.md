# Go code smells (detection catalog)

Loaded on demand by `refactoring-go-services` step 1 ("Identify ONE smell").
Before touching the code, scan for these:

- **Long Function**: a function spanning more than ~50 lines that handles multiple distinct operations.
- **God Handler / Orchestrator**: a single handler or service managing DB connections, formatting external API requests, and running core business logic simultaneously.
- **Primitive Obsession**: passing loosely typed primitives (`string`, `int`, `float64`) instead of rich domain objects (`MerchantID`, `Money`, `Email`).
- **Conditional Explosion**: deeply nested `if/else` or gigantic `switch` blocks that control application behaviour.
- **Feature Envy**: a function residing in module A that extracts and manipulates variables exclusively from module B.
- **Mixed Abstraction Levels**: combining high-level logic (validating an account) with low-level protocol details (parsing a JWT) in the same function.
- **Infrastructure Leakage**: emitting SQL-specific errors, caching logic, or Kafka `msg` objects directly into the domain layer.
