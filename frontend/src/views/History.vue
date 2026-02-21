<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">📋 Riwayat Transaksi</h2>
      <div style="display:flex;gap:6px;align-items:center">
        <button v-for="f in filters" :key="f.v" :class="['btn','btn-sm','btn-ghost',{active:filter===f.v}]" @click="filter=f.v;load()">{{ f.l }}</button>
        <button class="btn btn-sm btn-export" @click="exportExcel" :disabled="txs.length===0 || exporting" title="Export ke Excel">
          {{ exporting ? '⏳ Menyimpan...' : '📥 Export Excel' }}
        </button>
      </div>
    </div>
    <div class="table-wrap" v-if="txs.length > 0">
      <table>
        <thead><tr><th>Invoice</th><th>Pelanggan</th><th>Plat</th><th>Item</th><th>Subtotal</th><th>Diskon</th><th>Total</th><th>Bayar</th><th>Status</th><th>Kasir</th><th>Waktu</th></tr></thead>
        <tbody>
          <tr v-for="tx in txs" :key="tx.id" @click="selected=tx">
            <td style="font-family:monospace;font-size:11px;color:var(--t3)">{{ tx.invoiceNo }}</td>
            <td><div style="font-weight:600;font-size:12px">{{ tx.customerName }}</div><div style="font-size:10px;color:var(--t3)">{{ tx.carType }}</div></td>
            <td style="font-weight:700;letter-spacing:.5px;font-size:12px">{{ tx.plateNumber }}</td>
            <td>{{ tx.items?.length || 0 }} item</td>
            <td>{{ fmt(tx.subtotal) }}</td>
            <td style="color:var(--green)">{{ tx.discountAmount > 0 ? '-'+fmt(tx.discountAmount) : '-' }}</td>
            <td style="font-weight:700;color:var(--cyan)">{{ fmt(tx.total) }}</td>
            <td><span style="padding:2px 7px;background:var(--surface);border-radius:4px;font-size:10px;font-weight:700;text-transform:uppercase">{{ tx.paymentMethod }}</span></td>
            <td><span :class="['status-badge','status-'+tx.status]">{{ statusLabel(tx.status) }}</span></td>
            <td style="font-size:11px;color:var(--t3)">{{ tx.cashier?.fullName || '-' }}</td>
            <td style="font-size:11px;color:var(--t3)">{{ fmtTime(tx.createdAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else class="empty-state"><span class="empty-icon">📭</span><p>Belum ada transaksi</p></div>

    <!-- Detail Modal -->
    <div v-if="selected" class="modal-overlay" @click.self="selected=null">
      <div class="modal">
        <div class="modal-header"><h3>Detail {{ selected.invoiceNo }}</h3><button class="close-btn" @click="selected=null">✕</button></div>
        <div class="modal-body">
          <div class="d-row"><span>Pelanggan</span><span>{{ selected.customerName }}</span></div>
          <div class="d-row"><span>Plat / Kendaraan</span><span>{{ selected.plateNumber }} — {{ selected.carType }}</span></div>
          <div class="d-row"><span>Kasir</span><span>{{ selected.cashier?.fullName || '-' }}</span></div>
          <div class="d-row"><span>Waktu</span><span>{{ fmtTime(selected.createdAt) }}</span></div>
          <div class="d-items">
            <div v-for="it in selected.items" :key="it.id" class="d-item"><span>{{ it.itemName }}</span><span>{{ it.quantity }}×</span><span style="font-weight:600;color:var(--cyan)">{{ fmt(it.subtotal) }}</span></div>
          </div>
          <div class="d-row"><span>Subtotal</span><span>{{ fmt(selected.subtotal) }}</span></div>
          <div v-if="selected.discountAmount > 0" class="d-row" style="color:var(--green)"><span>Diskon</span><span>-{{ fmt(selected.discountAmount) }}</span></div>
          <div class="d-row" style="font-size:17px;font-weight:800;padding-top:8px;border-top:1px solid var(--border)"><span>Total</span><span style="color:var(--green)">{{ fmt(selected.total) }}</span></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as XLSX from 'xlsx'
import { callBackend, formatCurrency } from '../stores/auth'
const fmt = formatCurrency
const txs = ref([])
const filter = ref('all')
const selected = ref(null)
const filters = [{ v:'all', l:'Semua' },{ v:'washing', l:'🔄 Proses' },{ v:'done', l:'✅ Selesai' },{ v:'paid', l:'💳 Lunas' },{ v:'cancelled', l:'❌ Batal' }]
const statusLabel = s => ({ washing:'Proses', done:'Selesai', paid:'Lunas', cancelled:'Batal' }[s] || s)
const fmtTime = t => { try { return new Date(t).toLocaleString('id-ID',{day:'2-digit',month:'short',hour:'2-digit',minute:'2-digit'}) } catch(e) { return t } }
async function load() { try { const r = await callBackend('GetTransactions', filter.value, '', ''); if (r) txs.value = r } catch(e){} }
const exporting = ref(false)

async function exportExcel() {
  if (txs.value.length === 0) return
  exporting.value = true

  try {
    const data = txs.value.map(tx => ({
      'Invoice': tx.invoiceNo || '',
      'Pelanggan': tx.customerName || '',
      'Plat Nomor': tx.plateNumber || '',
      'Kendaraan': tx.carType || '',
      'Jumlah Item': tx.items?.length || 0,
      'Subtotal': tx.subtotal || 0,
      'Diskon': tx.discountAmount || 0,
      'Total': tx.total || 0,
      'Metode Bayar': (tx.paymentMethod || '').toUpperCase(),
      'Status': statusLabel(tx.status),
      'Kasir': tx.cashier?.fullName || '-',
      'Waktu': tx.createdAt ? new Date(tx.createdAt).toLocaleString('id-ID', { day:'2-digit', month:'short', year:'numeric', hour:'2-digit', minute:'2-digit' }) : ''
    }))

    const ws = XLSX.utils.json_to_sheet(data)

    // Auto-fit column widths
    const colWidths = Object.keys(data[0]).map(key => {
      const maxLen = Math.max(key.length, ...data.map(row => String(row[key]).length))
      return { wch: Math.min(maxLen + 2, 35) }
    })
    ws['!cols'] = colWidths

    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, 'Penjualan')

    // Generate base64 string instead of file download
    const base64 = XLSX.write(wb, { bookType: 'xlsx', type: 'base64' })

    const filterLabel = filters.find(f => f.v === filter.value)?.l || 'Semua'
    const dateStr = new Date().toISOString().slice(0, 10)
    const filename = `Laporan_Penjualan_${filterLabel}_${dateStr}.xlsx`

    // Call Go backend to save via native dialog
    const result = await callBackend('SaveExcelFile', base64, filename)
    if (result) {
      alert('✅ File berhasil disimpan!')
    }
  } catch (e) {
    console.error('Export error:', e)
    alert('❌ Gagal export: ' + (e.message || e))
  } finally {
    exporting.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.btn-ghost.active{border-color:var(--cyan)!important;color:var(--cyan)!important;background:var(--cyan-g)!important}
.btn-export{background:linear-gradient(135deg,#10b981,#059669)!important;color:#fff!important;border:none!important;font-weight:700;letter-spacing:.3px;transition:all .2s ease}
.btn-export:hover{background:linear-gradient(135deg,#059669,#047857)!important;transform:translateY(-1px);box-shadow:0 4px 12px rgba(16,185,129,.35)}
.btn-export:disabled{opacity:.4;cursor:not-allowed;transform:none;box-shadow:none}
.d-row{display:flex;justify-content:space-between;padding:7px 0;font-size:13px;color:var(--t2)}
.d-items{margin:12px 0;padding:10px;background:var(--surface);border-radius:var(--rs);display:flex;flex-direction:column;gap:6px}
.d-item{display:flex;justify-content:space-between;font-size:12px;gap:10px}.d-item span:first-child{flex:1}
</style>
