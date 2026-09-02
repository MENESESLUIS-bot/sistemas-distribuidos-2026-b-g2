# Anti-Corruption Layer Strategy

## Current scope

MVP 1 has no external business system integration. Therefore, no runtime ACL is required between LMS contexts. Internal module ports protect boundaries, and the API adapter translates HTTP DTOs into application commands.

## Planned ACLs

| External system | Consumer | ACL responsibility |
|---|---|---|
| University identity provider, if adopted | Access | Translate external subject, claims, and authentication errors into the internal Administrator model. |
| University student registry, if adopted | Membership | Translate external student records into the LMS Student language and reject unsupported fields. |
| Notification provider, if adopted | Circulation / Membership | Translate overdue and suspension events into provider-specific messages without leaking provider types into the domain. |

## Rules

- External SDKs are allowed only in infrastructure adapters.
- Domain entities never import external DTOs, clients, or transport types.
- Mapping failures become typed application errors and are logged with a correlation ID.
- Until an integration is approved, local LMS data remains the system of record.
