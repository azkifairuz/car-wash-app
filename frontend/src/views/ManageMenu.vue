<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">🍽️ Kelola Menu</h2>
      <button class="btn btn-primary" @click="openForm()">➕ Tambah Menu</button>
    </div>
    <div class="table-wrap">
      <table>
        <thead><tr><th>Icon</th><th>Nama</th><th>Harga</th><th>Kategori</th><th>Stok</th><th>Status</th><th>Aksi</th></tr></thead>
        <tbody>
          <tr v-for="m in items" :key="m.id">
            <td style="font-size:24px">{{ m.icon }}</td>
            <td style="font-weight:700">{{ m.name }}</td>
            <td style="font-weight:700;color:var(--cyan)">{{ fmt(m.price) }}</td>
            <td><span style="padding:2px 8px;background:var(--surface);border-radius:4px;font-size:11px;font-weight:600;text-transform:uppercase">{{ m.category }}</span></td>
            <td>
              <div style="display:flex;align-items:center;gap:6px">
                <span :style="{color: m.stock < 10 ? 'var(--red)' : 'var(--t1)', fontWeight:700}">{{ m.stock }}</span>
                <button class="btn btn-ghost btn-sm" @click="quickStock(m)" title="Update stok">📦</button>
              </div>
            </td>
            <td><button :class="['toggle',{on:m.isActive}]" @click="toggleActive(m)"></button></td>
            <td><div style="display:flex;gap:6px"><button class="btn btn-ghost btn-sm" @click="openForm(m)">✏️</button><button class="btn btn-danger btn-sm" @click="del(m.id)">🗑️</button></div></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showForm" class="modal-overlay" @click.self="showForm=false">
      <div class="modal">
        <div class="modal-header"><h3>{{ form.id ? 'Edit' : 'Tambah' }} Menu</h3><button class="close-btn" @click="showForm=false">✕</button></div>
        <div class="modal-body">
          <div class="form-row"><div class="form-group"><label class="label">Nama</label><input class="input" v-model="form.name" /></div><div class="form-group"><label class="label">Icon</label><input class="input" v-model="form.icon" placeholder="🍛" /></div></div>
          <div class="form-row"><div class="form-group"><label class="label">Harga (Rp)</label><input class="input" type="number" v-model.number="form.price" /></div><div class="form-group"><label class="label">Stok</label><input class="input" type="number" v-model.number="form.stock" /></div></div>
          <div class="form-group"><label class="label">Kategori</label><select class="select" v-model="form.category"><option value="makanan">Makanan</option><option value="minuman">Minuman</option><option value="snack">Snack</option></select></div>
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
const items = ref([])
const showForm = ref(false)
const form = ref({})

async function load() { try { const r = await callBackend('GetMenuItems'); if (r) items.value = r } catch(e){} }
function openForm(m) { form.value = m ? { ...m } : { name:'', price:0, category:'makanan', icon:'🍛', stock:50 }; showForm.value = true }
async function save() {
  try {
    const f = form.value
    if (f.id) await callBackend('UpdateMenuItem', f.id, f.name, f.price, f.category, f.icon, f.stock, f.isActive !== false)
    else await callBackend('CreateMenuItem', f.name, f.price, f.category, f.icon, f.stock)
    showForm.value = false; await load()
  } catch (e) { alert(e.message) }
}
async function toggleActive(m) { try { await callBackend('UpdateMenuItem', m.id, m.name, m.price, m.category, m.icon, m.stock, !m.isActive); await load() } catch(e){} }
async function quickStock(m) { const s = prompt(`Update stok ${m.name}:`, m.stock); if (s !== null) { try { await callBackend('UpdateMenuStock', m.id, parseInt(s)); await load() } catch(e){} } }
async function del(id) { if (!confirm('Hapus menu ini?')) return; try { await callBackend('DeleteMenuItem', id); await load() } catch(e) { alert(e.message) } }
onMounted(load)
</script>
