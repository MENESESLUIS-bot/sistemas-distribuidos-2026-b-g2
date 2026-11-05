// Isolates the federated 'shell/session' import to one file — every other
// file in this portal imports from here, never from 'shell/session' directly.
import { clearSession, getToken, isAuthenticated, setToken } from 'shell/session'

export function useSession() {
  return { getToken, setToken, isAuthenticated, clearSession }
}
