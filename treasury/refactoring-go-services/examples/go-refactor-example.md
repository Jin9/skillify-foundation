# Worked example: Primitive Obsession & Orchestrator Leak

Illustrates `refactoring-go-services` actions **Introduce Value Object** +
**Move Function to Domain (Tell, Don't Ask)**. Loaded on demand.

**Smell:** business logic leaked into the service layer; balance/currency passed
as bare primitives.

```go
// BEFORE
func (s *Service) ProcessAccount(balance float64, currency string, accountType string) error {
    if accountType == "premium" {
        balance = balance - 10.0 // Business logic leaked into the service layer
    }
    // ... Save balance ...
    return nil
}
```

**Refactored — rich objects + Tell, Don't Ask:**

```go
// AFTER

// Domain Layer
type Money struct {
    Amount   float64
    Currency string
}

func (m Money) Subtract(o Money) Money {
    return Money{Amount: m.Amount - o.Amount, Currency: m.Currency}
}

type Account struct {
    Type    AccountType
    Balance Money // was: balance float64 — primitive promoted to a value object
}

func (a *Account) ApplyMonthlyFee() {
    if a.Type == AccountPremium {
        // Core business logic operates cleanly on the domain itself
        a.Balance = a.Balance.Subtract(Money{Amount: 10, Currency: "USD"})
    }
}

// Orchestrator Layer
func (s *Service) ProcessAccount(act *Account) error {
    act.ApplyMonthlyFee() // Tell, Don't Ask
    return s.repo.Save(act)
}
```
