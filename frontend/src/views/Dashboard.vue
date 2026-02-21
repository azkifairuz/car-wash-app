<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">📊 Dashboard</h2>
      <div style="display:flex;gap:10px;align-items:center">
        <div class="date-input-group">
          <input type="date" v-model="dateFrom" class="input input-sm" @change="fetchStats">
          <span style="color:var(--t3)">s/d</span>
          <input type="date" v-model="dateTo" class="input input-sm" @change="fetchStats">
        </div>
      </div>
    </div>

    <div class="grid-4" style="margin-bottom:24px">
      <div class="stat-card"><div class="stat-icon" style="background:rgba(52,211,153,.12)">💰</div><div><span class="stat-label">Pendapatan Periode</span><div class="stat-value" style="color:var(--green)">{{ fmt(stats.todayRevenue) }}</div></div></div>
      <div class="stat-card"><div class="stat-icon" style="background:rgba(34,211,238,.12)">🧾</div><div><span class="stat-label">Transaksi Periode</span><div class="stat-value">{{ stats.todayTransactions }}</div></div></div>
      <div class="stat-card"><div class="stat-icon" style="background:rgba(251,191,36,.12)">🔄</div><div><span class="stat-label">Sedang Dicuci</span><div class="stat-value" style="color:var(--yellow)">{{ stats.activeWashes }}</div></div></div>
      <div class="stat-card"><div class="stat-icon" style="background:rgba(251,146,60,.12)">🍽️</div><div><span class="stat-label">Pesanan Makanan</span><div class="stat-value">{{ stats.foodOrdersToday }}</div></div></div>
    </div>
    <div class="stat-card" style="margin-bottom:24px"><div class="stat-icon" style="background:rgba(167,139,250,.12)">📈</div><div><span class="stat-label">Pendapatan Bulan Ini</span><div class="stat-value" style="color:var(--purple)">{{ fmt(stats.monthlyRevenue) }}</div></div></div>
    <div class="grid-3" style="margin-bottom:24px">
      <router-link to="/new-order" class="card" style="text-decoration:none;color:inherit;cursor:pointer"><span style="font-size:28px">➕</span><div style="font-size:15px;font-weight:700;margin-top:8px">Order Baru</div><div style="font-size:12px;color:var(--t3)">Buat pesanan cuci & makanan</div></router-link>
      <router-link to="/active" class="card" style="text-decoration:none;color:inherit;cursor:pointer"><span style="font-size:28px">🔄</span><div style="font-size:15px;font-weight:700;margin-top:8px">Lihat Antrian</div><div style="font-size:12px;color:var(--t3)">Monitor proses cuci aktif</div></router-link>
      <router-link to="/history" class="card" style="text-decoration:none;color:inherit;cursor:pointer"><span style="font-size:28px">📋</span><div style="font-size:15px;font-weight:700;margin-top:8px">Riwayat</div><div style="font-size:12px;color:var(--t3)">Lihat semua transaksi</div></router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { callBackend, formatCurrency } from '../stores/auth'

const fmt = formatCurrency
const stats = ref({ todayRevenue:0, todayTransactions:0, activeWashes:0, foodOrdersToday:0, monthlyRevenue:0 })

const today = new Date().toISOString().slice(0, 10)
const dateFrom = ref(today)
const dateTo = ref(today)

async function fetchStats() {
  try {
    const s = await callBackend('GetDashboardStats', dateFrom.value, dateTo.value)
    if (s) stats.value = s
  } catch(e) {
    console.error('Fetch stats error:', e)
  }
}

onMounted(fetchStats)
</script>

<style scoped>
.date-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--surface);
  padding: 4px 8px;
  border-radius: var(--rs);
  border: 1px solid var(--border);
}
.input-sm {
  background: transparent;
  border: none;
  color: var(--t1);
  font-size: 13px;
  outline: none;
}
</style>
