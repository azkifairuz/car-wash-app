<template>
  <div v-if="!auth.isLoggedIn">
    <router-view />
  </div>
  <div v-else class="app-container">
    <aside class="sidebar">
      <div class="logo">
        <span class="logo-icon">🚗</span>
        <div><div class="logo-title">SparkleWash</div><div class="logo-sub">Cuci Mobil & Cafe</div></div>
      </div>

      <nav class="nav-menu">
        <div class="nav-section">Kasir</div>
        <router-link to="/" class="nav-item" exact-active-class="active">
          <span class="nav-icon">📊</span>Dashboard
        </router-link>
        <router-link to="/new-order" class="nav-item" active-class="active">
          <span class="nav-icon">➕</span>Order Baru
        </router-link>
        <router-link to="/active" class="nav-item" active-class="active">
          <span class="nav-icon">🔄</span>Sedang Proses
          <span v-if="activeCount > 0" class="badge">{{ activeCount }}</span>
        </router-link>
        <router-link to="/history" class="nav-item" active-class="active">
          <span class="nav-icon">📋</span>Riwayat
        </router-link>

        <template v-if="auth.isSuperAdmin">
          <div class="nav-section">Admin</div>
          <router-link to="/manage/packages" class="nav-item" active-class="active">
            <span class="nav-icon">📦</span>Kelola Paket
          </router-link>
          <router-link to="/manage/menu" class="nav-item" active-class="active">
            <span class="nav-icon">🍽️</span>Kelola Menu
          </router-link>
          <router-link to="/manage/discounts" class="nav-item" active-class="active">
            <span class="nav-icon">🏷️</span>Kelola Diskon
          </router-link>
          <router-link to="/manage/users" class="nav-item" active-class="active">
            <span class="nav-icon">👥</span>Kelola User
          </router-link>
          <router-link to="/manage/shifts" class="nav-item" active-class="active">
            <span class="nav-icon">📅</span>Kelola Shift
          </router-link>
        </template>
      </nav>

      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">{{ auth.user?.fullName?.[0] || '?' }}</div>
          <div>
            <div class="user-name">{{ auth.user?.fullName }}</div>
            <div class="user-role">{{ auth.user?.role }}</div>
          </div>
        </div>
        <div class="clock">{{ currentTime }}</div>
        <div class="clock-date">{{ currentDate }}</div>
        <button class="logout-btn" @click="handleLogout">🚪 Logout</button>
      </div>
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore, callBackend } from './stores/auth'

const auth = useAuthStore()
const router = useRouter()
const currentTime = ref('')
const currentDate = ref('')
const activeCount = ref(0)
let timer = null

function updateClock() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  currentDate.value = now.toLocaleDateString('id-ID', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' })
}

async function loadActiveCount() {
  try {
    const stats = await callBackend('GetDashboardStats', '', '')
    if (stats) activeCount.value = stats.activeWashes
  } catch (e) { activeCount.value = 0 }
}

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}

onMounted(() => {
  updateClock()
  timer = setInterval(() => { updateClock(); loadActiveCount() }, 5000)
  if (auth.isLoggedIn) loadActiveCount()
})

onUnmounted(() => clearInterval(timer))
</script>
