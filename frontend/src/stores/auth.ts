import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('access_token') || '',
  }),
  actions: {
    setToken(token: string) {
      this.token = token
      localStorage.setItem('access_token', token)
    },
    clearToken() {
      this.token = ''
      localStorage.removeItem('access_token')
    },
  },
})
