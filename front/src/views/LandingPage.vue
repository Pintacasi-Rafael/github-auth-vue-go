<template>
  <div class="landing-page" v-if="token && user">
    <h2>✅ Logged in as: {{ user.login }}</h2>
    <p>Your JWT token:</p>
    <textarea readonly rows="5" cols="50" v-model="token"></textarea>
    <button @click="logout">Logout</button>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'

const token = ref(null)
const user = ref(null)

function logout() {
  localStorage.removeItem('jwt')
  window.location.href = '/'
}

onMounted(() => {
  const url = new URL(window.location.href)
  const jwtFromUrl = url.searchParams.get('token')

  const storedToken = localStorage.getItem('jwt')
  token.value = jwtFromUrl || storedToken

  if (token.value) {
    localStorage.setItem('jwt', token.value)

    const payload = JSON.parse(atob(token.value.split('.')[1]))
    user.value = { login: payload.username }

    if (jwtFromUrl) {
      window.history.replaceState({}, '', '/landing')
    }
  }
})
</script>

<style scoped>
.landing-page {
  max-width: 600px;
  margin: 30px auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-family: Arial, sans-serif;
  text-align: center;
}
textarea {
  width: 100%;
  margin: 10px 0;
  font-family: monospace;
}
button {
  padding: 10px 20px;
  background-color: #4267b2;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
button:hover {
  background-color: #365899;
}
</style>
