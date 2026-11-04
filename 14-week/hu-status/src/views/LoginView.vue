<script setup lang="ts">
// Implements HU-01 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want to log in with a secure username and password."
// Functional port of lms-front's prior LoginPage.tsx (ADR-011) — same fields,
// same guard, same generic-error handling; Vue idioms instead of React's.
import { ref } from 'vue'
import { apiClient, type ShellError } from 'shell/apiClient'
import { useSession } from '../composables/useSession'

const props = defineProps<{
  onAuthenticated: () => void
}>()

const { isAuthenticated, setToken } = useSession()

const username = ref('')
const password = ref('')
const error = ref<string | null>(null)
const isSubmitting = ref(false)

if (isAuthenticated()) {
  props.onAuthenticated()
}

async function handleSubmit() {
  error.value = null
  isSubmitting.value = true
  try {
    const { data } = await apiClient.post<{ accessToken: string }>('/auth/login', {
      username: username.value,
      password: password.value,
    })
    setToken(data.accessToken)
    props.onAuthenticated()
  } catch (err) {
    // The API deliberately returns the same generic error for a wrong username
    // or a wrong password (INV-002 on Administrator) — the UI must not try to
    // guess which one it was.
    error.value = (err as ShellError).message ?? 'Unable to reach the server. Please try again.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="access-login">
    <div class="access-login__card">
      <div class="access-login__header">
        <div class="access-login__logo">L</div>
        <h1 class="access-login__title">LMS Library</h1>
        <p class="access-login__subtitle">Sign in to manage the library</p>
      </div>

      <form class="access-login__form" @submit.prevent="handleSubmit">
        <div>
          <label for="username" class="access-login__label">Username</label>
          <input
            id="username"
            v-model="username"
            name="username"
            autocomplete="username"
            required
            class="access-login__input"
          />
        </div>

        <div>
          <label for="password" class="access-login__label">Password</label>
          <input
            id="password"
            v-model="password"
            name="password"
            type="password"
            autocomplete="current-password"
            required
            class="access-login__input"
          />
        </div>

        <p v-if="error" role="alert" class="access-login__error">{{ error }}</p>

        <button type="submit" class="access-login__submit" :disabled="isSubmitting">
          {{ isSubmitting ? 'Signing in…' : 'Sign in' }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.access-login {
  display: flex;
  min-height: 100vh;
  align-items: center;
  justify-content: center;
  background-color: #f8fafc;
  padding: 0 1rem;
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
}
.access-login__card {
  width: 100%;
  max-width: 24rem;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 0.75rem;
  padding: 2rem;
}
.access-login__header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 2rem;
}
.access-login__logo {
  display: flex;
  height: 3rem;
  width: 3rem;
  align-items: center;
  justify-content: center;
  border-radius: 0.75rem;
  background-color: #2563eb;
  color: #fff;
  font-size: 1.25rem;
  font-weight: 700;
}
.access-login__title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #0f172a;
}
.access-login__subtitle {
  font-size: 0.875rem;
  color: #64748b;
}
.access-login__form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.access-login__label {
  display: block;
  margin-bottom: 0.25rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #334155;
}
.access-login__input {
  width: 100%;
  border-radius: 0.5rem;
  border: 1px solid #cbd5e1;
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
}
.access-login__input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.15);
}
.access-login__error {
  border-radius: 0.5rem;
  background-color: #fff1f2;
  color: #e11d48;
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
}
.access-login__submit {
  width: 100%;
  border-radius: 0.5rem;
  background-color: #2563eb;
  color: #fff;
  font-size: 0.875rem;
  font-weight: 500;
  padding: 0.5rem 0.75rem;
  border: none;
  cursor: pointer;
}
.access-login__submit:disabled {
  opacity: 0.6;
  cursor: default;
}
</style>
