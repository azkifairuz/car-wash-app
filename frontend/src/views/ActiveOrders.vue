<template>
  <div class="fade-in">
    <h2 class="page-title">🔄 Sedang Proses</h2>
    <div v-if="orders.length > 0" class="orders-grid">
      <div v-for="tx in orders" :key="tx.id" class="order-card">
        <div class="oc-head"><span class="oc-plate">{{ tx.plateNumber }}</span><span :class="['status-badge','status-'+tx.status]">{{ statusLabel(tx.status) }}</span></div>
        <div class="oc-cust"><span class="oc-name">{{ tx.customerName }}</span><span class="oc-car">{{ tx.carType }}</span></div>
        <div class="oc-items">
          <div v-for="it in tx.items" :key="it.id" class="oc-item"><span>{{ it.itemName }} × {{ it.quantity }}</span><span class="oc-sub">{{ fmt(it.subtotal) }}</span></div>
        </div>
        <div v-if="tx.discountAmount > 0" class="oc-disc">Diskon: -{{ fmt(tx.discountAmount) }}</div>
        <div class="oc-foot">
          <div><span class="oc-total-label">Total</span><span class="oc-total">{{ fmt(tx.total) }}</span></div>
          <div class="oc-meta"><span>🕐 {{ fmtTime(tx.createdAt) }}</span><span class="oc-pay">{{ tx.paymentMethod?.toUpperCase() }}</span></div>
        </div>
        <div class="oc-actions">
          <button v-if="tx.status==='washing'" class="btn btn-success btn-sm" style="flex:1" @click="updateStatus(tx.id,'done')">✅ Selesai Cuci</button>
          <button v-if="tx.status==='done'" class="btn btn-primary btn-sm" style="flex:1" @click="updateStatus(tx.id,'paid')">💳 Tandai Lunas</button>
          <button v-if="tx.status==='washing'" class="btn btn-danger btn-sm" @click="cancelTx(tx.id)">✕</button>
        </div>
      </div>
    </div>
    <div v-else class="empty-state"><span class="empty-icon">🎉</span><p>Tidak ada antrian saat ini</p><router-link to="/new-order" class="btn btn-primary" style="margin-top:12px">➕ Buat Order Baru</router-link></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { callBackend, formatCurrency } from '../stores/auth'
const fmt = formatCurrency
const orders = ref([])

const statusLabel = s => ({ washing:'🔄 Dicuci', done:'✅ Selesai', paid:'💳 Lunas' }[s] || s)
const fmtTime = t => { try { return new Date(t).toLocaleTimeString('id-ID',{hour:'2-digit',minute:'2-digit'}) } catch(e) { return t } }

async function load() {
  try { const r = await callBackend('GetTransactions','washing','',''); const r2 = await callBackend('GetTransactions','done','',''); orders.value = [...(r||[]),...(r2||[])] } catch(e){}
}
async function updateStatus(id, status) { try { await callBackend('UpdateTransactionStatus', id, status); await load() } catch(e) { alert(e.message) } }
async function cancelTx(id) { if (!confirm('Batalkan transaksi ini?')) return; try { await callBackend('CancelTransaction', id); await load() } catch(e) { alert(e.message) } }

onMounted(load)
</script>

<style scoped>
.orders-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(340px,1fr));gap:14px}
.order-card{background:var(--card);border:1px solid var(--border);border-radius:var(--r);padding:18px;display:flex;flex-direction:column;gap:12px;transition:.2s}
.order-card:hover{border-color:rgba(34,211,238,.3);box-shadow:0 4px 24px rgba(0,0,0,.3)}
.oc-head{display:flex;justify-content:space-between;align-items:center}
.oc-plate{font-size:17px;font-weight:800;color:var(--cyan);letter-spacing:1px}
.oc-cust{display:flex;justify-content:space-between}.oc-name{font-weight:600;font-size:13px}.oc-car{font-size:12px;color:var(--t3)}
.oc-items{background:var(--surface);border-radius:var(--rs);padding:10px;display:flex;flex-direction:column;gap:5px}
.oc-item{display:flex;justify-content:space-between;font-size:12px;color:var(--t2)}
.oc-sub{font-weight:600;color:var(--t1)}
.oc-disc{font-size:12px;color:var(--green);text-align:right}
.oc-foot{display:flex;justify-content:space-between;align-items:flex-end}
.oc-total-label{font-size:11px;color:var(--t3);display:block}
.oc-total{font-size:18px;font-weight:800;color:var(--green)}
.oc-meta{display:flex;flex-direction:column;align-items:flex-end;gap:4px;font-size:11px;color:var(--t3)}
.oc-pay{padding:2px 7px;background:var(--surface);border-radius:4px;font-size:10px;font-weight:700}
.oc-actions{display:flex;gap:8px;padding-top:10px;border-top:1px solid var(--border)}
</style>
