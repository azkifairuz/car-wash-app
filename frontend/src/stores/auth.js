import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const isWails = typeof window !== 'undefined' && window['go']

async function call(method, ...args) {
  if (isWails) {
    return await window.go.backend.App[method](...args)
  }
  return null
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('currentUser') || 'null'))

  const isLoggedIn = computed(() => !!user.value)
  const isSuperAdmin = computed(() => user.value?.role === 'superadmin')
  const isKasir = computed(() => user.value?.role === 'kasir')

  async function login(username, password) {
    try {
      const res = await call('Login', username, password)
      if (res && res.success) {
        user.value = res.user
        localStorage.setItem('currentUser', JSON.stringify(res.user))
        return { success: true }
      }
      return { success: false, message: res?.message || 'Login gagal' }
    } catch (e) {
      // Mock login for dev
      if (username === 'admin' && password === 'admin123') {
        user.value = { id: 1, username: 'admin', fullName: 'Super Admin', role: 'superadmin' }
      } else if (username === 'kasir1' && password === 'kasir123') {
        user.value = { id: 2, username: 'kasir1', fullName: 'Kasir Satu', role: 'kasir' }
      } else {
        return { success: false, message: 'Username atau password salah' }
      }
      localStorage.setItem('currentUser', JSON.stringify(user.value))
      return { success: true }
    }
  }

  async function logout() {
    try { await call('Logout') } catch (e) {}
    user.value = null
    localStorage.removeItem('currentUser')
  }

  return { user, isLoggedIn, isSuperAdmin, isKasir, login, logout }
})

// Generic backend caller export
export async function callBackend(method, ...args) {
  if (isWails) {
    return await window.go.backend.App[method](...args)
  }
  return null
}

export function formatCurrency(amount) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency', currency: 'IDR', minimumFractionDigits: 0,
  }).format(amount)
}
