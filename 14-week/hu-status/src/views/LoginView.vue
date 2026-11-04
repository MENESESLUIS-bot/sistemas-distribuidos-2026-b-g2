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
