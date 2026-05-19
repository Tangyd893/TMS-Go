import { defineStore } from 'pinia'

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    permissions: [] as string[],
    roles: [] as string[],
  }),
  getters: {
    hasPermission: (state) => (code: string) => state.permissions.includes(code),
  },
  actions: {
    setPermissions(permissions: string[]) {
      this.permissions = permissions
    },
    setRoles(roles: string[]) {
      this.roles = roles
    },
    clearAll() {
      this.permissions = []
      this.roles = []
    },
  },
})
