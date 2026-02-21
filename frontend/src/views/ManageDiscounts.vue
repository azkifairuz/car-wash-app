<template>
  <div class="fade-in">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 class="page-title" style="margin-bottom:0">🏷️ Kelola Diskon</h2>
      <button class="btn btn-primary" @click="openForm()">➕ Tambah Diskon</button>
    </div>
    <div class="table-wrap">
      <table>
        <thead><tr><th>Nama</th><th>Tipe</th><th>Nilai</th><th>Min. Order</th><th>Berlaku</th><th>Status</th><th>Aksi</th></tr></thead>
        <tbody>
          <tr v-for="d in discounts" :key="d.id">
            <td style="font-weight:700">{{ d.name }}</td>
            <td><span style="padding:2px 8px;background:var(--surface);border-radius:4px;font-size:11px;font-weight:600;text-transform:uppercase">{{ d.type === 'percentage' ? 'Persen' : 'Nominal' }}</span></td>
            <td style="font-weight:700;color:var(--green)">{{ d.type === 'percentage' ? d.value + '%' : fmt(d.value) }}</td>
            <td>{{ fmt(d.minOrder) }}</td>
            <td style="font-size:11px;color:var(--t3)">{{ d.validFrom }} s/d {{ d.validUntil }}</td>
            <td><button :class="['toggle',{on:d.isActive}]" @click="toggleActive(d)"></button></td>
            <td><div style="display:flex;gap:6px"><button class="btn btn-ghost btn-sm" @click="openForm(d)">✏️</button><button class="btn btn-danger btn-sm" @click="del(d.id)">🗑️</button></div></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showForm" class="modal-overlay" @click.self="showForm=false">
      <div class="modal">
        <div class="modal-header"><h3>{{ form.id ? 'Edit' : 'Tambah' }} Diskon</h3><button class="close-btn" @click="showForm=false">✕</button></div>
        <div class="modal-body">
          <div class="form-group"><label class="label">Nama Diskon</label><input class="input" v-model="form.name" /></div>
          <div class="form-row">
            <div class="form-group"><label class="label">Tipe</label><select class="select" v-model="form.type"><option value="percentage">Persentase (%)</option><option value="fixed">Nominal (Rp)</option></select></div>
            <div class="form-group"><label class="label">Nilai</label><input class="input" type="number" v-model.number="form.value" /></div>
          </div>
          <div class="form-group"><label class="label">Min. Order (Rp)</label><input class="input" type="number" v-model.number="form.minOrder" /></div>
          <div class="form-row">
            <div class="form-group"><label class="label">Berlaku Dari</label><input class="input" type="date" v-model="form.validFrom" /></div>
            <div class="form-group"><label class="label">Berlaku Sampai</label><input class="input" type="date" v-model="form.validUntil" /></div>
          </div>
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
const discounts = ref([])
const showForm = ref(false)
const form = ref({})

async function load() { try { const r = await callBackend('GetDiscounts'); if (r) discounts.value = r } catch(e){} }
function openForm(d) { form.value = d ? { ...d } : { name:'', type:'percentage', value:10, minOrder:0, validFrom:'2026-01-01', validUntil:'2026-12-31' }; showForm.value = true }
async function save() {
  try {
    const f = form.value
    if (f.id) await callBackend('UpdateDiscount', f.id, f.name, f.type, f.value, f.minOrder, f.validFrom, f.validUntil, f.isActive !== false)
    else await callBackend('CreateDiscount', f.name, f.type, f.value, f.minOrder, f.validFrom, f.validUntil)
    showForm.value = false; await load()
  } catch (e) { alert(e.message) }
}
async function toggleActive(d) { try { await callBackend('UpdateDiscount', d.id, d.name, d.type, d.value, d.minOrder, d.validFrom, d.validUntil, !d.isActive); await load() } catch(e){} }
async function del(id) { if (!confirm('Hapus diskon ini?')) return; try { await callBackend('DeleteDiscount', id); await load() } catch(e) { alert(e.message) } }
onMounted(load)
</script>
