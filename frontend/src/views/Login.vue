<template>
  <div class="login-page">
    <div class="login-card fade-in">
      <div class="login-header">
        <span class="login-logo">🚗</span>
        <h1 class="login-title">SparkleWash</h1>
        <p class="login-sub">Cuci Mobil & Cafe — Silakan Login</p>
      </div>
      <div class="login-body">
        <div class="form-group">
          <label class="label">Username</label>
          <input class="input" v-model="username" placeholder="Masukkan username" @keyup.enter="handleLogin" autofocus />
        </div>
        <div class="form-group">
          <label class="label">Password</label>
          <input class="input" type="password" v-model="password" placeholder="Masukkan password" @keyup.enter="handleLogin" />
        </div>
        <div v-if="error" class="login-error">{{ error }}</div>
        <button class="btn btn-primary login-btn" :disabled="loading" @click="handleLogin">
          {{ loading ? '⏳ Memproses...' : '🔑 Login' }}
        </button>
        <div class="login-hint">
          <p>Default akun:</p>
          <p><strong>Admin:</strong> admin / admin123</p>
          <p><strong>Kasir:</strong> kasir1 / kasir123</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  if (!username.value || !password.value) { error.value = 'Username dan password harus diisi'; return }
  loading.value = true
  error.value = ''
  const res = await auth.login(username.value, password.value)
  loading.value = false
  if (res.success) { router.push('/') }
  else { error.value = res.message }
}
</script>

<style scoped>
.login-page{height:100vh;display:flex;align-items:center;justify-content:center;background:var(--bg);background-image:radial-gradient(ellipse at 30% 20%, rgba(34,211,238,0.06) 0%,transparent 50%),radial-gradient(ellipse at 70% 80%,rgba(129,140,248,0.06) 0%,transparent 50%)}
.login-card{background:var(--card);border:1px solid var(--border);border-radius:16px;width:400px;overflow:hidden;box-shadow:0 8px 40px rgba(0,0,0,.4)}
.login-header{text-align:center;padding:32px 32px 20px;border-bottom:1px solid var(--border)}
.login-logo{font-size:48px;display:block;margin-bottom:12px;filter:drop-shadow(0 0 12px rgba(34,211,238,0.4))}
.login-title{font-size:24px;font-weight:800;background:linear-gradient(135deg,var(--cyan),#818cf8);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.login-sub{font-size:13px;color:var(--t3);margin-top:6px}
.login-body{padding:24px 32px 32px}
.login-error{padding:10px;border-radius:var(--rs);background:rgba(248,113,113,.12);border:1px solid rgba(248,113,113,.3);color:var(--red);font-size:13px;font-weight:500;margin-bottom:12px}
.login-btn{width:100%;padding:12px;justify-content:center;font-size:15px;margin-top:4px}
.login-hint{margin-top:16px;padding:12px;border-radius:var(--rs);background:var(--surface);font-size:11px;color:var(--t3);line-height:1.6}
.login-hint strong{color:var(--t2)}
</style>
