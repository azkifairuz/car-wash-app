<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">👥 Kelola User</h2>
      <button class="btn btn-primary" @click="openForm()">➕ Tambah User</button>
    </div>
    <div class="table-wrap">
      <table>
        <thead><tr><th>Username</th><th>Nama Lengkap</th><th>Role</th><th>Status</th><th>Dibuat</th><th>Aksi</th></tr></thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td style="font-weight:700;font-family:monospace">{{ u.username }}</td>
            <td>{{ u.fullName }}</td>
            <td>
              <span :style="{padding:'2px 8px',background:u.role==='superadmin'?'rgba(167,139,250,.15)':'rgba(34,211,238,.15)',borderRadius:'4px',fontSize:'11px',fontWeight:600,color:u.role==='superadmin'?'var(--purple)':'var(--cyan)'}">
                {{ u.role }}
              </span>
            </td>
            <td><button :class="['toggle',{on:u.isActive}]" @click="toggleActive(u)"></button></td>
            <td style="font-size:11px;color:var(--t3)">{{ fmtDate(u.createdAt) }}</td>
            <td>
              <div style="display:flex;gap:6px">
                <button class="btn btn-ghost btn-sm" @click="openForm(u)">✏️</button>
                <button class="btn btn-warning btn-sm" @click="resetPw(u)">🔑</button>
                <button class="btn btn-danger btn-sm" @click="del(u.id)">🗑️</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showForm" class="modal-overlay" @click.self="showForm=false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ form.id ? 'Edit' : 'Tambah' }} User</h3>
          <button class="close-btn" @click="showForm=false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="label">Username</label>
            <input class="input" v-model="form.username" :disabled="!!form.id" placeholder="username" />
          </div>
          <div v-if="!form.id" class="form-group">
            <label class="label">Password</label>
            <input class="input" type="password" v-model="form.password" placeholder="Password untuk akun baru" />
          </div>
          <div class="form-group">
            <label class="label">Nama Lengkap</label>
            <input class="input" v-model="form.fullName" placeholder="Nama lengkap" />
          </div>
          <div class="form-group">
            <label class="label">Role</label>
            <select class="select" v-model="form.role">
              <option value="kasir">Kasir</option>
              <option value="superadmin">Super Admin</option>
            </select>
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

const users = ref([])
const showForm = ref(false)
const form = ref({})

const fmtDate = t => { try { return new Date(t).toLocaleDateString('id-ID',{day:'2-digit',month:'short',year:'numeric'}) } catch(e) { return t } }

async function load() {
  try { const r = await callBackend('GetUsers'); if (r) users.value = r } catch(e){}
}

function openForm(u) {
  form.value = u ? { id: u.id, username: u.username, fullName: u.fullName, role: u.role } : { username:'', password:'', fullName:'', role:'kasir' }
  showForm.value = true
}

async function save() {
  try {
    const f = form.value
    if (f.id) {
      await callBackend('UpdateUser', f.id, f.fullName, f.role, true)
    } else {
      if (!f.password) { alert('Password harus diisi'); return }
      await callBackend('CreateUser', f.username, f.password, f.fullName, f.role)
    }
    showForm.value = false; await load()
  } catch (e) { alert(e.message) }
}

async function toggleActive(u) {
  try { await callBackend('UpdateUser', u.id, u.fullName, u.role, !u.isActive); await load() } catch(e){}
}

async function resetPw(u) {
  const pw = prompt(`Reset password untuk ${u.username}:`)
  if (!pw) return
  try { await callBackend('ResetPassword', u.id, pw); alert('Password berhasil direset!') } catch(e) { alert(e.message) }
}

async function del(id) {
  if (!confirm('Hapus user ini?')) return
  try { await callBackend('DeleteUser', id); await load() } catch(e) { alert(e.message) }
}

onMounted(load)
</script>
