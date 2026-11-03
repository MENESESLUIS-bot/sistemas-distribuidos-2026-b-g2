// Standalone-dev entry point — mounts App directly when running `npm run
// dev` outside the shell, for local iteration. Never used when composed by
// lms-front (mount.ts is used instead).
import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router'

createApp(App).use(router).mount('#app')
