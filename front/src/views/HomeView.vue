<template>
  <div class="container">
    <h1>GitHub Login Demo</h1>
    <button @click="loginWithGitHub">Login with GitHub</button>
  </div>
</template>

<script setup>
function loginWithGitHub() {
  window.location.href = 'http://localhost:8080/github/login'
}
</script>

<style scoped>
.container {
  max-width: 600px;
  margin: 4rem auto;
  text-align: center;
}
button {
  padding: 0.6rem 1.2rem;
  font-size: 1rem;
  cursor: pointer;
}
</style>
