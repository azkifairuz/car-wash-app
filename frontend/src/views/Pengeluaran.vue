<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">💸 Pengeluaran</h2>
      <div style="display:flex;gap:8px;align-items:center">
        <input type="date" class="input" v-model="dateFrom" style="width:150px" @change="load" />
        <span style="color:var(--t3);font-size:13px">s/d</span>
        <input type="date" class="input" v-model="dateTo" style="width:150px" @change="load" />
        <button class="btn btn-export-detail" @click="exportLaporan" :disabled="exporting">
          {{ exporting ? '⏳ Menyimpan...' : '📊 Export Laporan' }}
        </button>
      </div>
    </div>

    <!-- Add Form -->
    <div class="form-card">
      <h3 class="section-title">{{ editing ? '✏️ Edit Pengeluaran' : '➕ Tambah Pengeluaran' }}</h3>
      <div class="form-row-4">
        <div class="form-group">
          <label class="label">Tanggal</label>
          <input type="date" class="input" v-model="form.date" />
        </div>
        <div class="form-group">
          <label class="label">Kategori</label>
          <select class="select" v-model="form.category">
            <option value="upah">👷 Upah Tenaga Kerja</option>
            <option value="pengeluaran">🛒 Pengeluaran Bengkel</option>
          </select>
        </div>
        <div class="form-group">
          <label class="label">Keterangan</label>
          <input class="input" v-model="form.name" placeholder="Contoh: Upah carwash, Bensin..." />
        </div>
        <div class="form-group">
          <label class="label">Jumlah (Rp)</label>
          <input type="number" class="input" v-model.number="form.amount" placeholder="0" min="0" />
        </div>
      </div>
      <div class="form-group" style="grid-column:1/-1">
        <label class="label">Catatan (opsional)</label>
        <input class="input" v-model="form.notes" placeholder="Catatan tambahan..." />
      </div>
      <div style="display:flex;gap:8px;margin-top:8px">
        <button class="btn btn-primary" @click="submit" :disabled="submitting">
          {{ submitting ? '⏳...' : (editing ? '💾 Simpan' : '➕ Tambah') }}
        </button>
        <button v-if="editing" class="btn btn-ghost" @click="cancelEdit">✕ Batal</button>
      </div>
    </div>

    <!-- Table -->
    <div v-if="expenses.length > 0" class="table-container" style="margin-top:16px">
      <table class="data-table">
        <thead>
          <tr>
            <th>Tanggal</th>
            <th>Kategori</th>
            <th>Keterangan</th>
            <th>Jumlah</th>
            <th>Catatan</th>
            <th class="center-text">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="exp in expenses" :key="exp.id">
            <td class="td-date">{{ fmtDate(exp.date) }}</td>
            <td>
              <span :class="['cat-badge', 'cat-' + exp.category]">
                {{ exp.category === 'upah' ? '👷 Upah' : '🛒 Pengeluaran' }}
              </span>
            </td>
            <td class="td-name">{{ exp.name }}</td>
            <td class="td-amount">{{ fmt(exp.amount) }}</td>
            <td class="td-notes">{{ exp.notes || '-' }}</td>
            <td>
              <div class="td-actions">
                <button class="btn-edit" @click="startEdit(exp)">✏️</button>
                <button class="btn-del" @click="deleteExp(exp.id)">🗑️</button>
              </div>
            </td>
          </tr>
        </tbody>
        <tfoot>
          <tr class="tfoot-upah">
            <td colspan="3"><strong>Total Upah Tenaga Kerja</strong></td>
            <td colspan="3"><strong>{{ fmt(totalUpah) }}</strong></td>
          </tr>
          <tr class="tfoot-pengeluaran">
            <td colspan="3"><strong>Total Pengeluaran Bengkel</strong></td>
            <td colspan="3"><strong>{{ fmt(totalPengeluaran) }}</strong></td>
          </tr>
          <tr class="tfoot-grand">
            <td colspan="3"><strong>Grand Total</strong></td>
            <td colspan="3"><strong>{{ fmt(totalUpah + totalPengeluaran) }}</strong></td>
          </tr>
        </tfoot>
      </table>
    </div>
    <div v-else class="empty-state" style="margin-top:16px"><span class="empty-icon">📭</span><p>Belum ada pengeluaran</p></div>
  </div>

  <!-- Modal Konfirmasi Hapus -->
  <!-- Custom Confirm Modal -->
<div v-if="confirmDelete" class="modal-overlay" @click.self="confirmDelete = null">
  <div class="modal-box">
    <p style="margin-bottom:16px;font-weight:600">Hapus pengeluaran ini?</p>
    <div style="display:flex;gap:8px;justify-content:flex-end">
      <button class="btn btn-ghost" @click="confirmDelete = null">Batal</button>
      <button class="btn btn-danger" @click="doDelete">🗑️ Hapus</button>
    </div>
  </div>
</div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import * as XLSX from 'xlsx'
import { callBackend, formatCurrency } from '../stores/auth'
const fmt = formatCurrency

// ---- State ----
const expenses = ref([])
const exporting = ref(false)
const submitting = ref(false)
const editing = ref(null)

// Default range: this week Mon - Sun
function getWeekRange() {
  const now = new Date()
  const day = now.getDay() || 7
  const mon = new Date(now); mon.setDate(now.getDate() - day + 1)
  const sun = new Date(mon); sun.setDate(mon.getDate() + 6)
  return [mon.toISOString().slice(0,10), sun.toISOString().slice(0,10)]
}
const [defaultFrom, defaultTo] = getWeekRange()
const dateFrom = ref(defaultFrom)
const dateTo = ref(defaultTo)

const emptyForm = () => ({
  date: new Date().toISOString().slice(0,10),
  category: 'pengeluaran',
  name: '',
  amount: 0,
  notes: ''
})
const form = ref(emptyForm())

const totalUpah = computed(() => expenses.value.filter(e => e.category === 'upah').reduce((s, e) => s + e.amount, 0))
const totalPengeluaran = computed(() => expenses.value.filter(e => e.category === 'pengeluaran').reduce((s, e) => s + e.amount, 0))

const fmtDate = d => {
  try {
    return new Date(d + 'T00:00:00').toLocaleDateString('id-ID', { day:'2-digit', month:'short', year:'numeric' })
  } catch { return d }
}

// ---- CRUD ----
async function load() {
  try {
    const r = await callBackend('GetExpenses', dateFrom.value, dateTo.value)
    expenses.value = r || []
  } catch(e) {}
}

async function submit() {
  if (!form.value.name || !form.value.amount) return alert('Keterangan dan Jumlah wajib diisi')
  submitting.value = true
  try {
    if (editing.value) {
      await callBackend('UpdateExpense', editing.value, form.value.date, form.value.category, form.value.name, form.value.amount, form.value.notes)
      editing.value = null
    } else {
      await callBackend('CreateExpense', form.value.date, form.value.category, form.value.name, form.value.amount, form.value.notes)
    }
    form.value = emptyForm()
    await load()
  } catch(e) { alert(e.message || 'Gagal menyimpan') }
  finally { submitting.value = false }
}

function startEdit(exp) {
  editing.value = exp.id
  form.value = { date: exp.date, category: exp.category, name: exp.name, amount: exp.amount, notes: exp.notes || '' }
}

function cancelEdit() {
  editing.value = null
  form.value = emptyForm()
}

const confirmDelete = ref(null) // menyimpan id yang mau dihapus

function deleteExp(id) {
  confirmDelete.value = id // tampilkan modal
}

async function doDelete() {
  try { 
    await callBackend('DeleteExpense', confirmDelete.value)
    confirmDelete.value = null
    await load() 
  } catch(e) { alert(e.message) }
}
async function exportLaporan() {
  exporting.value = true
  console.log('Export range:', dateFrom.value, dateTo.value)
  try {
    // Gunakan fungsi baru yang sudah ada grouping per tipe kendaraan
    const exportData = await callBackend('GetTransactionsForExport', dateFrom.value, dateTo.value)
    const rows = exportData?.rows || []
    const totalCarwash = exportData?.totalPendapatan || 0
    const totalDetailing = exportData?.totalDetailing || 0
    const totalKend = exportData?.totalKend || 0
    const kendByType = exportData?.kendByType || {}

    const fmtNum = n => n > 0 ? Math.round(n).toLocaleString('id-ID') : ''
    const fmtNumAlways = n => Math.round(n).toLocaleString('id-ID')

    // Collect all car types for header columns
    const carTypes = Object.keys(kendByType).sort()

    const aoa = []
    const cellStyles = {} // track which cells need special styling (blue bg for detailing col)

    // ── TITLE ROW ──
    const dateFromFmt = new Date(dateFrom.value + 'T00:00:00').toLocaleDateString('id-ID', { day: 'numeric', month: 'long' })
    const dateToFmt = new Date(dateTo.value + 'T00:00:00').toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
    aoa.push([`${dateFromFmt} s/d ${dateToFmt}`, '', '', '', '', ...carTypes.map(() => '')])

    // blank row
    aoa.push([])

    // ── HEADER ROW ──
    // Cols: Hari | Tanggal | Pendapatan | Kend | Detailing | [carType1] | [carType2] ...
    aoa.push(['Hari', 'Tanggal', 'Pendapatan', 'Kend', 'Detailing', ...carTypes])

    // ── DATA ROWS ──
    for (const r of rows) {
      const dt = new Date(r.date + 'T00:00:00')
      const dateFmt = dt.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
      const pendapatan = r.pendapatan + r.detailing
      const typeCounts = carTypes.map(t => r.kendByType?.[t] || '')
      aoa.push([
        r.dayName,
        dateFmt,
        fmtNum(pendapatan),
        r.totalKend || '',
        r.detailing > 0 ? fmtNum(r.detailing) : '',
        ...typeCounts
      ])
    }

    // ── TOTAL ROW ──
    const totalPendapatan = totalCarwash + totalDetailing
    const totalTypeCounts = carTypes.map(t => kendByType[t] || '')
    aoa.push(['', '', fmtNumAlways(totalPendapatan), totalKend, fmtNumAlways(totalDetailing), ...totalTypeCounts])

    aoa.push([])

    // ── PENDAPATAN BENGKEL ──
    aoa.push(['PENDAPATAN BENGKEL'])
    aoa.push(['Pendapatan Carwash', '', fmtNumAlways(totalCarwash)])
    aoa.push(['Pendapatan Detailing', '', fmtNumAlways(totalDetailing)])
    aoa.push(['JUMLAH', '', '', '', fmtNumAlways(totalPendapatan)])

    aoa.push([])

    // ── UPAH TENAGA KERJA ──
    const upahRows = expenses.value.filter(e => e.category === 'upah')
    aoa.push(['UPAH TENAGA KERJA'])
    for (const e of upahRows) aoa.push([e.name, '', fmtNumAlways(e.amount)])
    const totUpah = upahRows.reduce((s, e) => s + e.amount, 0)
    aoa.push(['JUMLAH', '', '', '', fmtNumAlways(totUpah)])

    aoa.push([])

    // ── PENGELUARAN BENGKEL ──
    const pengeluaranRows = expenses.value.filter(e => e.category === 'pengeluaran')
    aoa.push(['PENGELUARAN BENGKEL'])
    for (const e of pengeluaranRows) aoa.push([e.name, '', fmtNumAlways(e.amount)])
    const totPengeluaran = pengeluaranRows.reduce((s, e) => s + e.amount, 0)
    aoa.push(['JUMLAH', '', '', '', fmtNumAlways(totPengeluaran)])

    aoa.push([])

    // ── TRANSFER ──
    const transfer = totalPendapatan - totUpah - totPengeluaran
    aoa.push(['', '', 'TRANSFER', '', fmtNumAlways(transfer)])

    // ── Build workbook ──
    const ws = XLSX.utils.aoa_to_sheet(aoa)

    // Column widths
    const baseCols = [{ wch: 12 }, { wch: 12 }, { wch: 18 }, { wch: 8 }, { wch: 16 }]
    const typeCols = carTypes.map(() => ({ wch: 12 }))
    ws['!cols'] = [...baseCols, ...typeCols]

    // ── Styling via SheetJS (cell-level) ──
    // Title: merge + bold + center
    const totalCols = 5 + carTypes.length
    if (!ws['!merges']) ws['!merges'] = []
    ws['!merges'].push({ s: { r: 0, c: 0 }, e: { r: 0, c: totalCols - 1 } }) // title merge

    // Style helper
    const styleCell = (cellAddr, style) => {
      if (!ws[cellAddr]) ws[cellAddr] = { v: '', t: 's' }
      ws[cellAddr].s = style
    }

    const headerRow = 2 // 0-indexed: row index of header (Hari, Tanggal, ...)
    const dataStartRow = 3
    const dataEndRow = dataStartRow + rows.length - 1
    const totalRow = dataEndRow + 1
    const detailingCol = 4 // E column (0-indexed)

    const blueFill = { patternType: 'solid', fgColor: { rgb: '4FC3F7' } } // light blue
    const headerStyle = {
      font: { bold: true, italic: true, color: { rgb: 'CC0000' } },
      fill: { patternType: 'solid', fgColor: { rgb: 'FFFFFF' } },
      alignment: { horizontal: 'center' },
      border: { bottom: { style: 'thin', color: { rgb: '000000' } } }
    }
    const detailingHeaderStyle = {
      ...headerStyle,
      font: { bold: true, italic: true, color: { rgb: 'CC0000' } },
      fill: blueFill
    }
    const detailingCellStyle = { fill: blueFill, alignment: { horizontal: 'right' } }
    const boldStyle = { font: { bold: true } }
    const titleStyle = { font: { bold: true, sz: 14 }, alignment: { horizontal: 'center' } }
    const jumlahRedStyle = { font: { bold: true, color: { rgb: 'CC0000' } } }
    const transferStyle = { font: { bold: true, sz: 13 }, alignment: { horizontal: 'center' } }
    const totalNumStyle = { font: { bold: true, sz: 13 }, alignment: { horizontal: 'right' } }

    const colLetter = i => XLSX.utils.encode_col(i)
    const cellAddr = (r, c) => `${colLetter(c)}${r + 1}`

    // Title row style
    styleCell(cellAddr(0, 0), titleStyle)

    // Header row styles
    for (let c = 0; c < totalCols; c++) {
      const addr = cellAddr(headerRow, c)
      if (c === detailingCol) {
        styleCell(addr, detailingHeaderStyle)
      } else {
        styleCell(addr, headerStyle)
      }
    }

    // Data + total rows: highlight detailing column
    for (let r = dataStartRow; r <= totalRow; r++) {
      const addr = cellAddr(r, detailingCol)
      const isTotal = r === totalRow
      styleCell(addr, isTotal
        ? { fill: blueFill, font: { bold: true }, alignment: { horizontal: 'right' } }
        : detailingCellStyle
      )
      // Total row bold all
      if (isTotal) {
        for (let c = 0; c < totalCols; c++) {
          styleCell(cellAddr(r, c), c === detailingCol
            ? { fill: blueFill, font: { bold: true }, alignment: { horizontal: 'right' } }
            : boldStyle
          )
        }
      }
    }

    // Section headers bold
    const sectionLabels = ['PENDAPATAN BENGKEL', 'UPAH TENAGA KERJA', 'PENGELUARAN BENGKEL']
    for (let r = 0; r < aoa.length; r++) {
      const cell = aoa[r][0]
      if (sectionLabels.includes(cell)) {
        styleCell(cellAddr(r, 0), { font: { bold: true, sz: 11 } })
      }
      if (cell === 'JUMLAH') {
        styleCell(cellAddr(r, 0), jumlahRedStyle)
        styleCell(cellAddr(r, 4), { ...jumlahRedStyle, alignment: { horizontal: 'right' } })
      }
      if (aoa[r][2] === 'TRANSFER') {
        styleCell(cellAddr(r, 2), transferStyle)
        styleCell(cellAddr(r, 4), totalNumStyle)
      }
    }

    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, 'Laporan')

    const base64 = XLSX.write(wb, { bookType: 'xlsx', type: 'base64' })
    // const dateStr = new Date().toISOString().slice(0, 10)
    const result = await callBackend('SaveExcelFile', base64, `Laporan_${dateFrom.value}_sd_${dateTo.value}.xlsx`)
    // const result = await callBackend('SaveExcelFile', base64, `Laporan_Pengeluaran_${dateStr}.xlsx`)
    if (result) alert('✅ File berhasil disimpan!')
  } catch (e) {
    console.error(e)
    alert('❌ Gagal export: ' + (e.message || e))
  } finally {
    exporting.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.form-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r); padding: 18px; }
.section-title { font-size: 14px; font-weight: 700; margin-bottom: 14px; color: var(--t1); }
.form-row-4 { display: grid; grid-template-columns: 1fr 1fr 2fr 1fr; gap: 12px; }
.table-container { overflow-x: auto; background: var(--card); border-radius: var(--r); border: 1px solid var(--border); }
.data-table { width: 100%; border-collapse: separate; border-spacing: 0; min-width: 700px; }
.data-table thead th { background: var(--surface); padding: 10px 14px; text-align: left; font-size: 11px; color: var(--t3); font-weight: 600; text-transform: uppercase; border-bottom: 1px solid var(--border); }
.data-table tbody tr { transition: 0.2s; }
.data-table tbody tr:hover { background: var(--surface); }
.data-table tbody td { padding: 12px 14px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.data-table tbody tr:last-child td { border-bottom: none; }
.data-table tfoot td { padding: 10px 14px; font-size: 13px; border-top: 1px solid var(--border); }
.center-text { text-align: center; }
.td-date { font-size: 12px; color: var(--t3); }
.td-name { font-weight: 600; font-size: 13px; }
.td-amount { font-weight: 700; font-size: 13px; color: var(--cyan); }
.td-notes { font-size: 12px; color: var(--t3); }
.td-actions { display: flex; gap: 6px; justify-content: center; }
.btn-edit { padding: 5px 10px; border-radius: var(--rs); border: 1px solid var(--border); background: var(--surface); cursor: pointer; font-size: 13px; transition: 0.2s; }
.btn-edit:hover { border-color: var(--cyan); }
.btn-del { padding: 5px 10px; border-radius: var(--rs); border: none; background: rgba(248,113,113,.15); color: var(--red); cursor: pointer; font-size: 13px; transition: 0.2s; }
.btn-del:hover { background: rgba(248,113,113,.3); }
.cat-badge { display: inline-block; padding: 3px 9px; border-radius: 12px; font-size: 11px; font-weight: 600; }
.cat-upah { background: rgba(59,130,246,.15); color: #3b82f6; }
.cat-pengeluaran { background: rgba(245,158,11,.15); color: #f59e0b; }
.tfoot-upah { background: rgba(59,130,246,.07); color: #3b82f6; }
.tfoot-pengeluaran { background: rgba(245,158,11,.07); color: #f59e0b; }
.tfoot-grand { background: rgba(34,211,238,.08); color: var(--cyan); font-size: 14px; }
.btn-export-detail { background: linear-gradient(135deg,#6366f1,#4f46e5) !important; color: #fff !important; border: none !important; font-weight: 700; letter-spacing: .3px; padding: 9px 16px; border-radius: var(--rs); font-family: inherit; font-size: 13px; cursor: pointer; transition: 0.2s; }
.btn-export-detail:hover { background: linear-gradient(135deg,#4f46e5,#4338ca) !important; transform: translateY(-1px); }
.btn-export-detail:disabled { opacity: .5; cursor: not-allowed; }
.modal-overlay { position:fixed;inset:0;background:rgba(0,0,0,.5);display:flex;align-items:center;justify-content:center;z-index:999; }
.modal-box { background:var(--card);border:1px solid var(--border);border-radius:var(--r);padding:24px;min-width:280px; }
.btn-danger { background:rgba(248,113,113,.2);color:var(--red);border:1px solid var(--red);padding:7px 14px;border-radius:var(--rs);cursor:pointer;font-weight:600; }
</style>
