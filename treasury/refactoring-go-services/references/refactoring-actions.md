# Refactoring actions (executable operations)

Loaded on demand by `refactoring-go-services` step 2 ("Choose ONE refactor
action"). Address each detected smell with ONE precise action per iteration.

- **Extract Function / Extract Method**
  - **Smell**: Long function / mixed abstractions.
  - **Action**: Isolate a distinct chunk into a private, intention-revealing method. Pass required parameters explicitly.

- **Rename for Intention**
  - **Smell**: Vague or overly technical names (`Process()`, `CalcData()`).
  - **Action**: Rename to strictly reflect business intent (`CalculateMonthlyInterest()`).

- **Move Function to Domain Layer (Tell, Don't Ask)**
  - **Smell**: Orchestrator forcibly mutating an object's fields directly.
  - **Action**: Shift the exact mutation logic to a method on the target domain model.

- **Introduce Value Object**
  - **Smell**: Signatures requiring bare metrics (`amount float64, currency string`).
  - **Action**: Encapsulate these parameters into a domain struct (`Money{}`) to guard against invalid states.

- **Replace Conditional with Strategy / Polymorphism**
  - **Smell**: A giant `switch` dictating variations in logic.
  - **Action**: Define a common interface and delegate to bounded implementations; use a mapping/factory for the lookup.

- **Separate Calculation from Presentation (Introduce Intermediate DTO)**
  - **Smell**: JSON marshalling tags or HTTP codes written directly on a rich DB-coupled domain model.
  - **Action**: Create pure `ResponseDTO` / `RequestDTO` structs in the handler layer plus a simple `ToDTO()` map function.

- **Encapsulate State**
  - **Smell**: Exported core entity fields operated on across multiple handler files.
  - **Action**: Make the struct fields private; provide controlled mutators that name domain events (`account.ActivateProfile()`).
