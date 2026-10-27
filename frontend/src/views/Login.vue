<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const registration = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const body = JSON.stringify({ email: email.value, password: password.value })
    if (registration.value) await api('/auth/registration', { method: 'POST', body })
    await api('/auth/login', { method: 'POST', body })
    router.push('/records')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="login-page">
    <form class="panel login-form" @submit.prevent="submit">
      <h1>Журнал гостей</h1>
      <p class="muted">{{ registration ? 'Создайте аккаунт' : 'Войдите, чтобы работать с записями' }}</p>

      <label for="email">Почта</label>
      <input id="email" v-model="email" type="email" required autocomplete="email" />

      <label for="password">Пароль</label>
      <input id="password" v-model="password" type="password" required :autocomplete="registration ? 'new-password' : 'current-password'" />

      <p v-if="error" class="error" role="alert">{{ error }}</p>
      <button type="submit" :disabled="loading">{{ loading ? 'Подождите…' : registration ? 'Зарегистрироваться' : 'Войти' }}</button>
      <button type="button" class="text-button" @click="registration = !registration; error = ''">
        {{ registration ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться' }}
      </button>
    </form>
  </main>
</template>
