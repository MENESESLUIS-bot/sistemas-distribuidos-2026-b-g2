// Exposed to lms-front via Module Federation (vite.config.ts's
// federation({ exposes })). A Vue remote can't export a React
// RouteObject[] the way the other three portals do, so the canonical
// interop pattern for a non-React remote in this system (ADR-011) is a
// mount/unmount function instead of a component tree.
import { createApp, type App } from 'vue'
import LoginView from './views/LoginView.vue'

export function mountAccessLogin(el: HTMLElement, onAuthenticated: () => void): () => void {
  const app: App = createApp(LoginView, { onAuthenticated })
  app.mount(el)
  return () => app.unmount()
}
