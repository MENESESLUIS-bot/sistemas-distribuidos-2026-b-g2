# Service Contracts — MVP 1

MVP 1 uses one deployable `library-api` with internal module boundaries. These contracts are defined as ports first and can become HTTP or event contracts if a context is extracted later.

## Synchronous contracts

| Contract | Provider | Consumer | Request | Response / errors |
| --- | --- | --- | --- | --- |
| `StudentEligibility` | Membership | Circulation | `studentId` | active, suspended-until, active-loan-count; `StudentNotFound` |
| `BookAvailability` | Catalog | Circulation | `bookId`, quantity | reservation/release result; `BookNotFound`, `InsufficientAvailability` |
| `RegisterLoan` | Circulation | API | student ID, book ID, current date | loan ID, loan date, due date; eligibility or availability error |
| `RegisterReturn` | Circulation | API | loan ID, return date | returned loan, late flag; `LoanNotFound`, `LoanAlreadyReturned` |

## Domain events

```json
{
  "type": "LoanReturnedLate",
  "loanId": "uuid",
  "studentId": "uuid",
  "returnedAt": "2026-09-02T12:00:00Z",
  "occurredAt": "2026-09-02T12:00:01Z"
}
```

Membership consumes this event through an in-process event handler in MVP 1 and applies a seven-day suspension. The event payload contains identifiers and facts, not database rows.

## Contract rules

- Provider modules own validation of their invariants.
- Consumers depend on interfaces, not concrete repositories.
- Errors are typed at the application boundary and mapped to the common API error response.
- Events are idempotent by event ID and are published only after the transaction succeeds.
