import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LandingPage from '../views/LandingPage.vue'

const routes = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/landing', name: 'Landing', component: LandingPage, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  let token = localStorage.getItem('jwt')

  // 🆕 Also check the URL (e.g., after GitHub redirects)
  const urlParams = new URLSearchParams(window.location.search)
  const tokenFromUrl = urlParams.get('token')

  if (tokenFromUrl) {
    localStorage.setItem('jwt', tokenFromUrl)
    token = tokenFromUrl
  }

  if (to.meta.requiresAuth && !token) {
    next({ path: '/' })
  } else {
    next()
  }
})

export default router
