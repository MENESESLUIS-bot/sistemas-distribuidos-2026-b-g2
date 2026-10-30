# CP4 / CP8 — Security gates

- Workflow: [`security.yml`](security.yml) (proposal for the team pipeline)
- Same workflow running in this fork: [`.github/workflows/security.yml`](../../../.github/workflows/security.yml)
- Execution: <https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/actions/runs/41092565678>
- Results and how to run it locally: [`HOW-TO-RUN.md`](HOW-TO-RUN.md)

In that execution `negative authorization tests` and `govulncheck` pass and `gitleaks` fails
on a real finding (`12-week/hu-status/.env.example`), explained in `HOW-TO-RUN.md`.
