// Ambient types for what lms-front (the container) exposes, so this portal's
// views can import them with full type-checking even though the actual
// modules only exist at runtime, resolved by @module-federation/vite against
// the container's dev server (vite.config.ts's federation({ remotes }) — set
// on the container's side, not this repo's, since a remote never declares
// its own container). Same shape as every React portal's shell.d.ts
// (rules/2-anexos/H-front.md) — framework-agnostic, plain TS ambient module.
declare module 'shell/apiClient' {
  // Minimal shape this portal actually calls — not the full AxiosInstance
  // type, so this repo doesn't need axios as its own dependency just for a
  // type import; the real instance (and its real type) lives in lms-front.
  export interface ShellApiClient {
    post<T>(url: string, body?: unknown, config?: { headers?: Record<string, string> }): Promise<{ data: T }>
  }
  export const apiClient: ShellApiClient

  export interface ShellError {
    error: string
    message: string
    details?: { field: string; message: string }[]
    traceId?: string
  }
}

declare module 'shell/session' {
  export function getToken(): string | null
  export function setToken(token: string): void
  export function isAuthenticated(): boolean
  export function clearSession(): void
}
