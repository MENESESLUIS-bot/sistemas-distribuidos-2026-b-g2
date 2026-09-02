# MVP 1 MoSCoW Commitment

## Must Have

- HU-01 Administrator authentication.
- HU-02 Student registration.
- HU-04 Book registration.
- HU-06 Loan registration with the seven-day policy, maximum two active loans, suspension validation, and availability decrement.
- HU-07 Return registration with history, late detection, and availability increment.
- `/health`, database readiness, migrations, Docker execution, and automated tests for the committed behavior.

## Should Have

- HU-03 Student search, editing, and deactivation.
- HU-05 Catalog search and real-time availability view.
- HU-08 Overdue report and automatic seven-day suspension.
- HU-09 Book editing with immutable ISBN.

## Could Have

- CSV/PDF reports.
- Notifications and reminders.
- Book reservations.

## Won't Have in MVP 1

- Multiple administrator roles.
- Student self-registration.
- Monetary fines.
- Individual copy condition tracking.
- Independent microservice deployment and a message broker.

The commitment is limited to the Must Have scope so the team can deliver a coherent loan lifecycle inside one sprint.
