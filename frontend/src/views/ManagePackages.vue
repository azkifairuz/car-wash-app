<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">📦 Kelola Paket Cuci</h2>
      <button class="btn btn-primary" @click="openForm()">➕ Tambah Paket</button>
    </div>
    <div class="table-wrap">
      <table>
        <thead><tr><th>Icon</th><th>Nama</th><th>Deskripsi</th><th>Harga</th><th>Durasi</th><th>Kategori</th><th>Status</th><th>Aksi</th></tr></thead>
        <tbody>
          <tr v-for="p in packages" :key="p.id">
            <td style="font-size:24px">{{ p.icon }}</td>
            <td style="font-weight:700">{{ p.name }}</td>
            <td style="font-size:12px;color:var(--t3);max-width:200px">{{ p.description }}</td>
            <td style="font-weight:700;color:var(--cyan)">{{ fmt(p.price) }}</td>
            <td>{{ p.duration }} mnt</td>
            <td><span style="padding:2px 8px;background:var(--surface);border-radius:4px;font-size:11px;font-weight:600;text-transform:uppercase">{{ p.category }}</span></td>
            <td><button :class="['toggle',{on:p.isActive}]" @click="toggleActive(p)"></button></td>
            <td><div style="display:flex;gap:6px"><button class="btn btn-ghost btn-sm" @click="openForm(p)">✏️</button><button class="btn btn-danger btn-sm" @click="del(p.id)">🗑️</button></div></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showForm" class="modal-overlay" @click.self="showForm=false">
      <div class="modal">
        <div class="modal-header"><h3>{{ form.id ? 'Edit' : 'Tambah' }} Paket</h3><button class="close-btn" @click="showForm=false">✕</button></div>
        <div class="modal-body">
          <div class="form-row"><div class="form-group"><label class="label">Nama</label><input class="input" v-model="form.name" /></div><div class="form-group"><label class="label">Icon</label><input class="input" v-model="form.icon" placeholder="⚡" /></div></div>
          <div class="form-group"><label class="label">Deskripsi</label><input class="input" v-model="form.description" /></div>
          <div class="form-row"><div class="form-group"><label class="label">Harga (Rp)</label><input class="input" type="number" v-model.number="form.price" /></div><div class="form-group"><label class="label">Durasi (menit)</label><input class="input" type="number" v-model.number="form.duration" /></div></div>
          <div class="form-group"><label class="label">Kategori</label><select class="select" v-model="form.category"><option value="cuci">Cuci</option><option value="detailing">Detailing</option><option value="coating">Coating</option></select></div>
        </div>
        <div class="modal-footer"><button class="btn btn-ghost" @click="showForm=false">Batal</button><button class="btn btn-primary" @click="save">💾 Simpan</button></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { callBackend, formatCurrency } from '../stores/auth'
const fmt = formatCurrency
const packages = ref([])
const showForm = ref(false)
const form = ref({})

async function load() { try { const r = await callBackend('GetPackages'); if (r) packages.value = r } catch(e){} }
function openForm(p) { form.value = p ? { ...p } : { name:'', description:'', price:0, duration:30, category:'cuci', icon:'🚗' }; showForm.value = true }
async function save() {
  try {
    const f = form.value
    if (f.id) await callBackend('UpdatePackage', f.id, f.name, f.description, f.price, f.duration, f.category, f.icon, f.isActive !== false)
    else await callBackend('CreatePackage', f.name, f.description, f.price, f.duration, f.category, f.icon)
    showForm.value = false; await load()
  } catch (e) { alert(e.message) }
}
async function toggleActive(p) { try { await callBackend('UpdatePackage', p.id, p.name, p.description, p.price, p.duration, p.category, p.icon, !p.isActive); await load() } catch(e){} }
async function del(id) { if (!confirm('Hapus paket ini?')) return; try { await callBackend('DeletePackage', id); await load() } catch(e) { alert(e.message) } }
onMounted(load)
</script>
