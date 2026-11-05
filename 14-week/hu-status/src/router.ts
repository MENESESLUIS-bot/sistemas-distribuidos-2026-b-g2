// Standalone-dev routing only (`npm run dev`, without lms-front running) —
// when composed by the shell, this portal is mounted via mount.ts instead,
// never through this router. Mirrors the other portals' "preview in
// isolation" App.tsx note: without the shell, shell/apiClient and
// shell/session aren't resolvable at all (@module-federation/vite needs the
// container's dev server running to resolve them), so the login call here
// fails until lms-front is also running on port 3000.
import { createRouter, createWebHistory } from 'vue-router'
import LoginView from './views/LoginView.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: LoginView,
      props: { onAuthenticated: () => console.log('[standalone] authenticated') },
    },
  ],
})
