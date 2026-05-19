import { defineStore } from 'pinia'
import type { UserInfo } from '@/api/modules/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('access_token') || '',
    user: null as UserInfo | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.user?.username || '',
  },
  actions: {
    setToken(token: string) {
      this.token = token
      localStorage.setItem('access_token', token)
    },
    setUser(user: UserInfo) {
      this.user = user
    },
    clearAll() {
      this.token = ''
      this.user = null
      localStorage.removeItem('access_token')
    },
  },
})
