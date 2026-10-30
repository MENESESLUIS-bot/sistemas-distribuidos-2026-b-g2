import { federation } from '@module-federation/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// https://vite.dev/config/
//
// Exposes ./mount so lms-front's shell can lazy-load this portal's login
// screen as a remote (ADR-011 — this is the one Vue remote in an otherwise
// React system). vue is NOT in `shared` — there's no common Vue runtime
// with the React shell to share, unlike the other three (React) portals.
export default defineConfig({
  plugins: [
    vue(),
    federation({
      name: 'access_portal',
      filename: 'remoteEntry.js',
      exposes: {
        './mount': './src/mount.ts',
      },
      // 'shell/apiClient' and 'shell/session' (src/shims/shell.d.ts) resolve
      // against this — lms-front's own federation({ exposes }) is what
      // actually serves them.
      remotes: {
        shell: {
          type: 'module',
          name: 'shell',
          entry: 'http://localhost:3000/remoteEntry.js',
        },
      },
    }),
  ],
  server: {
    port: 3004,
  },
  build: {
    target: 'esnext',
  },
})
