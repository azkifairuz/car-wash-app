<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">📅 Kelola Shift</h2>
      <button class="btn btn-primary" @click="openForm()">➕ Tambah Shift</button>
    </div>

    <div style="display:flex;gap:10px;margin-bottom:16px">
      <div class="form-group" style="margin:0">
        <label class="label">Dari Tanggal</label>
        <input class="input" type="date" v-model="dateFrom" @change="load" style="width:180px" />
      </div>
      <div class="form-group" style="margin:0">
        <label class="label">Sampai Tanggal</label>
        <input class="input" type="date" v-model="dateTo" @change="load" style="width:180px" />
      </div>
    </div>

    <div class="table-wrap" v-if="shifts.length > 0">
      <table>
        <thead><tr><th>Tanggal</th><th>Kasir</th><th>Jam Mulai</th><th>Jam Selesai</th><th>Catatan</th><th>Aksi</th></tr></thead>
        <tbody>
          <tr v-for="s in shifts" :key="s.id">
            <td style="font-weight:700">{{ s.date }}</td>
            <td>
              <div style="font-weight:600">{{ s.user?.fullName || '-' }}</div>
              <div style="font-size:10px;color:var(--t3)">{{ s.user?.username }}</div>
            </td>
            <td>{{ s.startTime }}</td>
            <td>{{ s.endTime }}</td>
            <td style="font-size:12px;color:var(--t3);max-width:200px">{{ s.notes || '-' }}</td>
            <td>
              <div style="display:flex;gap:6px">
                <button class="btn btn-ghost btn-sm" @click="openForm(s)">✏️</button>
                <button class="btn btn-danger btn-sm" @click="del(s.id)">🗑️</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else class="empty-state"><span class="empty-icon">📅</span><p>Belum ada shift pada rentang tanggal ini</p></div>

    <!-- Form Modal -->
    <div v-if="showForm" class="modal-overlay" @click.self="showForm=false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ form.id ? 'Edit' : 'Tambah' }} Shift</h3>
          <button class="close-btn" @click="showForm=false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="label">Kasir</label>
            <select class="select" v-model.number="form.userId">
              <option v-for="u in kasirList" :key="u.id" :value="u.id">{{ u.fullName }} ({{ u.username }})</option>
            </select>
          </div>
          <div class="form-group">
            <label class="label">Tanggal</label>
            <input class="input" type="date" v-model="form.date" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="label">Jam Mulai</label>
              <input class="input" type="time" v-model="form.startTime" />
            </div>
            <div class="form-group">
              <label class="label">Jam Selesai</label>
              <input class="input" type="time" v-model="form.endTime" />
            </div>
          </div>
          <div class="form-group">
            <label class="label">Catatan</label>
            <input class="input" v-model="form.notes" placeholder="Catatan opsional..." />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-ghost" @click="showForm=false">Batal</button>
          <button class="btn btn-primary" @click="save">💾 Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { callBackend } from '../stores/auth'

const shifts = ref([])
const kasirList = ref([])
const showForm = ref(false)
const form = ref({})

// Default: current week
const today = new Date()
const weekStart = new Date(today)
weekStart.setDate(today.getDate() - today.getDay() + 1)
const weekEnd = new Date(weekStart)
weekEnd.setDate(weekStart.getDate() + 6)

const dateFrom = ref(weekStart.toISOString().slice(0, 10))
const dateTo = ref(weekEnd.toISOString().slice(0, 10))

async function load() {
  try {
    const r = await callBackend('GetShifts', dateFrom.value, dateTo.value)
    if (r) shifts.value = r
  } catch (e) {}
}

async function loadKasir() {
  try {
    const r = await callBackend('GetUsers')
    if (r) kasirList.value = r.filter(u => u.isActive)
  } catch (e) {}
}

function openForm(s) {
  if (s && s.id) {
    form.value = { id: s.id, userId: s.userId, date: s.date, startTime: s.startTime, endTime: s.endTime, notes: s.notes || '' }
  } else {
    form.value = { userId: kasirList.value[0]?.id || 0, date: today.toISOString().slice(0, 10), startTime: '08:00', endTime: '16:00', notes: '' }
  }
  showForm.value = true
}

async function save() {
  try {
    const f = form.value
    if (f.id) {
      await callBackend('UpdateShift', f.id, f.userId, f.date, f.startTime, f.endTime, f.notes)
    } else {
      await callBackend('CreateShift', f.userId, f.date, f.startTime, f.endTime, f.notes)
    }
    showForm.value = false
    await load()
  } catch (e) { alert(e.message) }
}

async function del(id) {
  if (!confirm('Hapus shift ini?')) return
  try { await callBackend('DeleteShift', id); await load() } catch (e) { alert(e.message) }
}

onMounted(() => { load(); loadKasir() })
</script>
