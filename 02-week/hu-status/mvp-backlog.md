# MVP 1 Backlog Slice

| ID | Story | Acceptance criteria | Priority |
|---|---|---|---|
| HU-01 | Administrator login | Valid credentials create a JWT session; invalid credentials return a generic error. | Must |
| HU-02 | Register student | Required data is stored; duplicate document IDs are rejected. | Must |
| HU-04 | Register book | Required bibliographic data is stored; initial availability equals total copies; duplicate ISBNs are rejected. | Must |
| HU-06 | Register loan | Only an eligible student and available book can be selected; due date is seven days after loan date; availability decreases by one. | Must |
| HU-07 | Register return | An active loan becomes returned; availability increases by one; late status is recorded. | Must |

## Dependency order

`HU-01` and `HU-02` enable `HU-06`; `HU-04` enables `HU-06`; `HU-06` enables `HU-07`. HU-03, HU-05, HU-08, and HU-09 remain in the next slice because they are Should Have conveniences or depend on the core lifecycle.
