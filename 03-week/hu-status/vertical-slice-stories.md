# First Vertical Feature — Register a Loan

## Story

**As the Administrator, I want to register a loan for an eligible student and available book, so that the physical loan is recorded and inventory remains accurate.**

## Testable acceptance criteria

```gherkin
Scenario: Register an eligible loan
  Given an active student with fewer than two active loans
  And a book with at least one available copy
  When the administrator submits the loan
  Then the loan is stored as Active
  And the due date is exactly seven days after the loan date
  And available copies decrease by one

Scenario: Reject an ineligible loan
  Given a suspended student or a student with two active loans
  When the administrator submits the loan
  Then no loan is created
  And the response identifies the eligibility reason

Scenario: Reject a book without availability
  Given a book with zero available copies
  When the administrator submits the loan
  Then no loan is created
  And the response reports insufficient availability
```

## Vertical tasks

1. Add the Circulation `RegisterLoan` use case and Loan aggregate.
2. Add Membership and Catalog eligibility/availability ports and adapters.
3. Add the PostgreSQL transaction for loan insert plus availability decrement.
4. Add `POST /api/v1/loans` and map typed errors to HTTP responses.
5. Add domain, application, persistence, and API integration tests.
