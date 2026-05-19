import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const whiteList = ['/login']

export function setupRouterGuards(router: Router) {
  router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()

    if (whiteList.includes(to.path)) {
      if (authStore.isLoggedIn && to.path === '/login') {
        next('/')
        return
      }
      next()
      return
    }

    if (!authStore.isLoggedIn) {
      next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
      return
    }

    next()
  })
}
