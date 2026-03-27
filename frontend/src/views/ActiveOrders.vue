<template>
  <div class="fade-in">
    <h2 class="page-title">🔄 Sedang Proses</h2>
    <div v-if="orders.length > 0" class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>Waktu & Plat</th>
            <th>Pelanggan & Kendaraan</th>
            <th>Pesanan</th>
            <th>Total Bayar</th>
            <th>Status</th>
            <th class="center-text">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tx in orders" :key="tx.id">
            <td>
              <div class="td-plate">{{ tx.plateNumber }}</div>
              <div class="td-time">🕐 {{ fmtTime(tx.createdAt) }}</div>
            </td>
            <td>
              <div class="td-name">{{ tx.customerName }}</div>
              <div class="td-car">{{ tx.carType }}</div>
            </td>
            <td>
              <div class="td-items">
                <div v-for="it in tx.items" :key="it.id" class="td-item">
                  <span>• {{ it.itemName }} (×{{ it.quantity }})</span>
                </div>
              </div>
            </td>
            <td>
              <div class="td-total">{{ fmt(tx.total) }}</div>
              <div class="td-pay">{{ tx.paymentMethod?.toUpperCase() }}</div>
              <div v-if="tx.discountAmount > 0" class="td-disc">Diskon: -{{ fmt(tx.discountAmount) }}</div>
            </td>
            <td>
              <span :class="['status-badge','status-'+tx.status]">{{ statusLabel(tx.status) }}</span>
            </td>
            <td>
              <div class="td-actions">
                <button v-if="tx.status==='washing'" class="btn btn-success btn-sm" @click="updateStatus(tx.id,'done')">✅ Selesai</button>
                <button v-if="tx.status==='done'" class="btn btn-primary btn-sm" @click="openPaymentModal(tx)">💳 Lunas</button>
                <button v-if="tx.status==='washing'" class="btn btn-danger btn-sm" @click="cancelTx(tx.id)" title="Batalkan">✕</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else class="empty-state"><span class="empty-icon">🎉</span><p>Tidak ada antrian saat ini</p><router-link to="/new-order" class="btn btn-primary" style="margin-top:12px">➕ Buat Order Baru</router-link></div>

    <!-- Payment Modal -->
    <div v-if="showPaymentModal" class="modal-overlay">
      <div class="modal-content">
        <h3 class="modal-title">Pilih Pembayaran</h3>
        <p class="modal-subtitle">Total: {{ fmt(selectedOrder.total) }}</p>
        
        <div class="pay-grid">
          <button v-for="pm in payMethods" :key="pm.v" :class="['pay-btn',{active:selectedPayment===pm.v}]" @click="selectedPayment=pm.v">
            {{ pm.i }} {{ pm.l }}
          </button>
        </div>

        <div class="modal-actions">
          <button class="btn btn-secondary" @click="closePaymentModal">Batal</button>
          <button class="btn btn-primary" :disabled="!selectedPayment || processingPayment" @click="confirmPayment">
            {{ processingPayment ? '⏳ Memproses...' : '✅ Pemasukan Lunas' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { callBackend, formatCurrency } from '../stores/auth'
const fmt = formatCurrency
const orders = ref([])

const showPaymentModal = ref(false)
const selectedOrder = ref(null)
const selectedPayment = ref('cash')
const processingPayment = ref(false)
const payMethods = [
  { v: 'cash', i: '💵', l: 'Cash' },
  { v: 'debit', i: '💳', l: 'Debit' },
  { v: 'qris', i: '📱', l: 'QRIS' },
  { v: 'transfer', i: '🏦', l: 'Transfer' }
]

const statusLabel = s => ({ washing:'🔄 Dicuci', done:'✅ Selesai', paid:'💳 Lunas' }[s] || s)
const fmtTime = t => { try { return new Date(t).toLocaleTimeString('id-ID',{hour:'2-digit',minute:'2-digit'}) } catch(e) { return t } }

async function load() {
  try { const r = await callBackend('GetTransactions','washing','',''); const r2 = await callBackend('GetTransactions','done','',''); orders.value = [...(r||[]),...(r2||[])] } catch(e){}
}
async function updateStatus(id, status) { try { await callBackend('UpdateTransactionStatus', id, status); await load() } catch(e) { alert(e.message) } }
async function cancelTx(id) { if (!confirm('Batalkan transaksi ini?')) return; try { await callBackend('CancelTransaction', id); await load() } catch(e) { alert(e.message) } }

function openPaymentModal(tx) {
  selectedOrder.value = tx
  selectedPayment.value = 'cash'
  showPaymentModal.value = true
}

function closePaymentModal() {
  showPaymentModal.value = false
  selectedOrder.value = null
}

async function confirmPayment() {
  if (!selectedOrder.value || !selectedPayment.value) return
  processingPayment.value = true
  try {
    await callBackend('ProcessPayment', selectedOrder.value.id, selectedPayment.value)
    showPaymentModal.value = false
    await load()
  } catch (e) {
    alert(e.message)
  } finally {
    processingPayment.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.table-container { overflow-x: auto; background: var(--card); border-radius: var(--r); border: 1px solid var(--border); }
.data-table { width: 100%; border-collapse: separate; border-spacing: 0; min-width: 800px; }
.data-table thead th { background: var(--surface); padding: 12px 14px; text-align: left; font-size: 11px; color: var(--t3); font-weight: 600; text-transform: uppercase; border-bottom: 1px solid var(--border); }
.data-table tbody tr { transition: 0.2s; }
.data-table tbody tr:hover { background: var(--surface); }
.data-table tbody td { padding: 14px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.data-table tbody tr:last-child td { border-bottom: none; }
.center-text { text-align: center; }
.td-plate { font-size: 15px; font-weight: 800; color: var(--cyan); letter-spacing: 1px; margin-bottom: 4px; }
.td-time { font-size: 11px; color: var(--t3); }
.td-name { font-weight: 600; font-size: 13px; color: var(--t1); margin-bottom: 4px; }
.td-car { font-size: 12px; color: var(--t3); }
.td-items { display: flex; flex-direction: column; gap: 4px; }
.td-item { font-size: 12px; color: var(--t2); padding-bottom: 2px; }
.td-total { font-size: 15px; font-weight: 800; color: var(--green); margin-bottom: 4px; }
.td-pay { display: inline-block; padding: 2px 6px; background: var(--surface); border-radius: 4px; font-size: 10px; font-weight: 700; color: var(--t2); }
.td-disc { font-size: 11px; color: var(--green); margin-top: 4px; }
.td-actions { display: flex; gap: 6px; justify-content: center; align-items: center; }
.status-badge { display: inline-block; padding: 4px 8px; border-radius: 12px; font-size: 11px; font-weight: 700; }
.status-washing { background: rgba(59,130,246,0.15); color: #3b82f6; }
.status-done { background: rgba(34,197,94,0.15); color: #22c55e; }
.btn-sm { padding: 6px 10px; font-size: 11px; }

/* Modal Styles */
.modal-overlay { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: var(--bg); border-radius: var(--r); padding: 24px; width: 360px; max-width: 90vw; border: 1px solid var(--border); box-shadow: 0 10px 40px rgba(0,0,0,0.3); }
.modal-title { font-size: 16px; font-weight: 700; margin-bottom: 4px; color: var(--t1); }
.modal-subtitle { font-size: 14px; color: var(--green); font-weight: 800; margin-bottom: 20px; }
.pay-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-bottom: 24px; }
.pay-btn { padding: 12px; border-radius: var(--rs); border: 1px solid var(--border); background: var(--surface); color: var(--t2); font-size: 13px; font-weight: 600; cursor: pointer; transition: 0.2s; display: flex; align-items: center; justify-content: center; gap: 6px; }
.pay-btn.active { border-color: var(--cyan); background: var(--cyan-g); color: var(--cyan); }
.modal-actions { display: flex; gap: 10px; justify-content: flex-end; }
.btn-secondary { background: var(--surface); border: 1px solid var(--border); color: var(--t2); padding: 10px 16px; border-radius: var(--rs); font-family: inherit; font-size: 13px; font-weight: 600; cursor: pointer; }
.btn-secondary:hover { background: var(--border); }
</style>
