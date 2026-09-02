# Domain Model — Circulation Context

## Aggregate root

**Loan** is the aggregate root because it controls the lifecycle of a borrowing transaction: creation, due-date calculation, return, and late-return classification.

## Entities and value objects

| Element | Type | Responsibility |
| --- | --- | --- |
| Loan | Entity / aggregate root | Coordinates loan state and protects lifecycle invariants. |
| LoanId | Value object | Stable identity for a loan. |
| StudentId | Value object | References the borrowing student without importing the Student entity. |
| BookId | Value object | References the borrowed book without importing the Book entity. |
| LoanStatus | Value object / enum | `ACTIVE` or `RETURNED`. |
| LoanPeriod | Value object | Encapsulates the fixed seven-day period. |
| ReturnDetails | Value object | Return date and computed late flag. |

## Invariants

1. A new loan must reference an active student and a book with at least one available copy.
2. The due date is always the loan date plus seven days; callers cannot set it independently.
3. A student cannot exceed two active loans.
4. A suspended student cannot create a new loan.
5. An active loan has no return date; a returned loan has exactly one return date.
6. A return can be registered only once.
7. A late return is one whose return date is after its due date.
8. Loan creation and book availability decrement commit atomically.

## Domain events

| Event | When emitted | Consumers |
| --- | --- | --- |
| `LoanRegistered` | A loan is successfully created | Catalog availability projection, audit log |
| `LoanReturned` | An active loan is returned | Catalog availability, history |
| `LoanReturnedLate` | A return is after the due date | Membership suspension policy, reporting |
