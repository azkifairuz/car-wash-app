<template>
  <div class="fade-in">
    <h2 class="page-title">➕ Order Baru</h2>
    <div class="order-layout">
      <!-- Left: Menu -->
      <div class="menu-panel">
        <div class="tab-bar">
          <button :class="['tab-btn', { active: tab === 'wash' }]" @click="tab='wash'">🚗 Cuci Mobil</button>
          <button :class="['tab-btn', { active: tab === 'food' }]" @click="tab='food'">🍽️ Makanan & Minuman</button>
        </div>

        <div v-if="tab==='wash'" class="list-container">
          <table class="data-table">
            <thead>
              <tr>
                <th width="50" class="center-text">Icon</th>
                <th>Paket Layanan</th>
                <th width="100">Waktu</th>
                <th width="120">Harga</th>
                <th width="80" class="center-text">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in packages" :key="p.id" @click="addToCart(p,'package')">
                <td class="center-text"><span class="td-icon">{{ p.icon }}</span></td>
                <td>
                  <div class="td-name">{{ p.name }}</div>
                  <div class="td-desc">{{ p.description }}</div>
                </td>
                <td><span class="td-meta">⏱ {{ p.duration }} mnt</span></td>
                <td class="td-price">{{ fmt(p.price) }}</td>
                <td class="center-text"><button class="td-add-btn">+ Tambahkan</button></td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="tab==='food'" class="list-container-wrap">
          <div class="sub-tabs">
            <button v-for="c in foodCats" :key="c.key" :class="['sub-tab',{active:foodCat===c.key}]" @click="foodCat=c.key">{{ c.icon }} {{ c.label }}</button>
          </div>
          <div class="list-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th width="50" class="center-text">Icon</th>
                  <th>Nama Menu</th>
                  <th width="100">Stok</th>
                  <th width="120">Harga</th>
                  <th width="80" class="center-text">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="m in filteredMenu" :key="m.id" @click="addToCart(m,'menu')">
                  <td class="center-text"><span class="td-icon">{{ m.icon }}</span></td>
                  <td><div class="td-name">{{ m.name }}</div></td>
                  <td><span class="td-meta">Stok: {{ m.stock }}</span></td>
                  <td class="td-price">{{ fmt(m.price) }}</td>
                  <td class="center-text"><button class="td-add-btn">+ Tambahkan</button></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Right: Cart -->
      <div class="cart-panel">
        <div class="cart-box">
          <h3 class="cart-title">🛒 Keranjang</h3>
          <div v-if="cart.length > 0" class="cart-items">
            <div v-for="(item,i) in cart" :key="i" class="cart-item">
              <div class="ci-top"><span class="ci-name">{{ item.itemName }}</span><span class="ci-price">{{ fmt(item.price) }}</span></div>
              <div class="ci-bottom">
                <div class="ci-qty"><button class="qb" @click="changeQty(i,-1)">−</button><span class="qv">{{ item.quantity }}</span><button class="qb" @click="changeQty(i,1)">+</button></div>
                <span class="ci-sub">{{ fmt(item.price * item.quantity) }}</span>
                <button class="rm-btn" @click="cart.splice(i,1)">✕</button>
              </div>
            </div>
          </div>
          <div v-else class="empty-state" style="padding:20px"><span class="empty-icon" style="font-size:32px">🛒</span><p>Belum ada item</p></div>

          <!-- Discount -->
          <div v-if="cart.length > 0" style="margin-top:12px">
            <label class="label">Diskon</label>
            <select class="select" v-model="selectedDiscount">
              <option :value="null">— Tanpa Diskon —</option>
              <option v-for="d in discounts" :key="d.id" :value="d.id" :disabled="subtotal < d.minOrder">
                {{ d.name }} {{ d.type === 'percentage' ? `(${d.value}%)` : `(${fmt(d.value)})` }}
                {{ subtotal < d.minOrder ? `(min ${fmt(d.minOrder)})` : '' }}
              </option>
            </select>
          </div>

          <div v-if="cart.length > 0" class="cart-summary">
            <div class="sum-row"><span>Subtotal</span><span>{{ fmt(subtotal) }}</span></div>
            <div v-if="discountAmount > 0" class="sum-row discount"><span>Diskon</span><span>-{{ fmt(discountAmount) }}</span></div>
            <div class="sum-row total"><span>Total</span><span>{{ fmt(total) }}</span></div>
          </div>
        </div>

        <div class="cart-box">
          <h3 class="cart-title">👤 Info Pelanggan</h3>
          <div class="form-group"><label class="label">Nama</label><input class="input" v-model="customer.name" placeholder="Nama pelanggan..." /></div>
          <div class="form-row">
            <div class="form-group"><label class="label">No. Plat</label><input class="input" v-model="customer.plate" placeholder="B 1234 XX" /></div>
            <div class="form-group"><label class="label">Kendaraan</label><input class="input" v-model="customer.car" placeholder="Toyota Avanza" /></div>
          </div>
        </div>

        <button class="btn btn-primary submit-btn" :disabled="cart.length===0||submitting" @click="submitOrder">
          {{ submitting ? '⏳ Memproses...' : '✅ Buat Order' }}
        </button>

        <!-- Success + Print buttons -->
        <div v-if="lastOrderId" class="print-section">
          <div class="toast toast-success">✅ Order berhasil dibuat!</div>
          <div class="print-actions">
            <button class="btn print-btn" @click="handlePrint" :disabled="printing">
              {{ printing ? '⏳ Mencetak...' : '🖨️ Cetak Struk' }}
            </button>
            <button class="btn print-btn-secondary" @click="lastOrderId = null">✕ Tutup</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { callBackend, formatCurrency } from '../stores/auth'
const fmt = formatCurrency

const tab = ref('wash')
const foodCat = ref('all')
const packages = ref([])
const menuItems = ref([])
const discounts = ref([])
const cart = ref([])
const selectedDiscount = ref(null)
const submitting = ref(false)
const printing = ref(false)
const customer = ref({ name: '', plate: '', car: '' })

// Last order tracking for print
const lastOrderId = ref(null)
const lastOrderData = ref(null)

const foodCats = [{ key:'all', icon:'📋', label:'Semua' },{ key:'makanan', icon:'🍛', label:'Makanan' },{ key:'minuman', icon:'🥤', label:'Minuman' },{ key:'snack', icon:'🍟', label:'Snack' }]

const filteredMenu = computed(() => foodCat.value === 'all' ? menuItems.value : menuItems.value.filter(m => m.category === foodCat.value))
const subtotal = computed(() => cart.value.reduce((s, i) => s + i.price * i.quantity, 0))
const discountAmount = computed(() => {
  if (!selectedDiscount.value) return 0
  const d = discounts.value.find(x => x.id === selectedDiscount.value)
  if (!d || subtotal.value < d.minOrder) return 0
  return d.type === 'percentage' ? subtotal.value * d.value / 100 : d.value
})
const total = computed(() => Math.max(0, subtotal.value - discountAmount.value))

function addToCart(item, type) {
  const existing = cart.value.find(c => c.itemType === type && c.itemId === item.id)
  if (existing) { existing.quantity++; return }
  cart.value.push({ itemType: type, itemId: item.id, itemName: item.name, price: item.price, quantity: 1 })
}
function changeQty(i, d) { cart.value[i].quantity += d; if (cart.value[i].quantity <= 0) cart.value.splice(i, 1) }

async function submitOrder() {
  submitting.value = true
  try {
    const discId = selectedDiscount.value || null
    const tx = await callBackend('CreateTransaction', customer.value.name || 'Walk-in', customer.value.plate || '-', customer.value.car || '-', 'pending', cart.value, discId)

    // Store order ID for printing
    if (tx && tx.id) {
      lastOrderId.value = tx.id
    }

    // Also build direct print data as fallback
    lastOrderData.value = JSON.stringify({
      invoiceNo: tx?.invoiceNo || ('INV-' + Date.now().toString().slice(-6)),
      date: new Date().toLocaleString('id-ID'),
      customerName: customer.value.name || 'Walk-in',
      plateNumber: customer.value.plate || '-',
      carType: customer.value.car || '-',
      items: cart.value.map(item => ({
        name: item.itemName,
        qty: item.quantity,
        price: item.price,
        subtotal: item.price * item.quantity,
      })),
      subtotal: subtotal.value,
      discountAmount: discountAmount.value,
      total: total.value,
      paymentMethod: 'pending',
      cashierName: '',
    })

    // Reset form
    cart.value = []
    selectedDiscount.value = null
    customer.value = { name:'', plate:'', car:'' }

    await loadData()

    // Auto print if enabled
    try {
      const cfg = await callBackend('GetPrinterConfig')
      if (cfg && cfg.autoPrint) {
        await handlePrint()
      }
    } catch (e) {
      // Auto print not critical, ignore
    }
  } catch (e) {
    alert(e.message || 'Gagal membuat order')
  } finally {
    submitting.value = false
  }
}

async function handlePrint() {
  printing.value = true
  try {
    // Method 1: Print by transaction ID (preferred - data from DB)
    if (lastOrderId.value) {
      await callBackend('PrintReceipt', lastOrderId.value)
      return
    }
    // Method 2: Print from direct data (fallback)
    if (lastOrderData.value) {
      await callBackend('PrintReceiptDirect', lastOrderData.value)
      return
    }
    alert('Tidak ada data order untuk dicetak')
  } catch (e) {
    alert('Gagal cetak: ' + (e.message || e))
  } finally {
    printing.value = false
  }
}

async function loadData() {
  try {
    const [p, m, d] = await Promise.all([
      callBackend('GetActivePackages'),
      callBackend('GetActiveMenuItems'),
      callBackend('GetActiveDiscounts'),
    ])
    if (p) packages.value = p
    if (m) menuItems.value = m
    if (d) discounts.value = d
  } catch (e) {}
}

onMounted(loadData)
</script>

<style scoped>
.order-layout{display:grid;grid-template-columns:1fr 380px;gap:20px;height:calc(100vh - 120px)}
.menu-panel{display:flex;flex-direction:column;overflow:hidden}
.tab-bar{display:flex;gap:8px;margin-bottom:14px}
.tab-btn{padding:9px 18px;border-radius:var(--rs);border:1px solid var(--border);background:var(--card);color:var(--t2);font-family:inherit;font-size:13px;font-weight:600;cursor:pointer;transition:.2s}
.tab-btn.active{background:var(--cyan-g);border-color:var(--cyan);color:var(--cyan)}
.sub-tabs{display:flex;gap:6px;margin-bottom:12px}
.sub-tab{padding:5px 12px;border-radius:16px;border:1px solid var(--border);background:transparent;color:var(--t3);font-family:inherit;font-size:11px;font-weight:600;cursor:pointer;transition:.2s}
.sub-tab.active{background:var(--surface);border-color:var(--cyan);color:var(--cyan)}
.list-container { overflow-y: auto; flex: 1; padding-right: 6px; }
.list-container-wrap { display: flex; flex-direction: column; overflow-y: hidden; flex: 1; }
.data-table { width: 100%; border-collapse: separate; border-spacing: 0; background: var(--card); border-radius: var(--r); border: 1px solid var(--border); }
.data-table thead th { background: var(--surface); padding: 10px 14px; text-align: left; font-size: 11px; color: var(--t3); font-weight: 600; text-transform: uppercase; border-bottom: 1px solid var(--border); }
.data-table thead th:first-child { border-top-left-radius: var(--r); }
.data-table thead th:last-child { border-top-right-radius: var(--r); }
.data-table tbody tr { cursor: pointer; transition: 0.2s; }
.data-table tbody tr:hover { background: var(--card-h); }
.data-table tbody td { padding: 12px 14px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.data-table tbody tr:last-child td { border-bottom: none; }
.data-table tbody tr:last-child td:first-child { border-bottom-left-radius: var(--r); }
.data-table tbody tr:last-child td:last-child { border-bottom-right-radius: var(--r); }
.center-text { text-align: center !important; }
.td-icon { font-size: 20px; }
.td-name { font-weight: 700; font-size: 13px; color: var(--t1); margin-bottom: 2px; }
.td-desc { font-size: 11px; color: var(--t3); line-height: 1.3; }
.td-meta { font-size: 10px; color: var(--t3); background: var(--surface); padding: 3px 8px; border-radius: 12px; }
.td-price { font-weight: 700; font-size: 13px; color: var(--cyan); }
.td-add-btn { padding: 6px 12px; border-radius: var(--rs); border: 1px solid var(--border); background: var(--surface); color: var(--t2); font-size: 11px; font-weight: 600; cursor: pointer; transition: 0.2s; white-space: nowrap; }
.data-table tbody tr:hover .td-add-btn { border-color: var(--cyan); background: var(--cyan-g); color: var(--cyan); }
.cart-panel{display:flex;flex-direction:column;gap:14px;overflow-y:auto}
.cart-box{background:var(--card);border:1px solid var(--border);border-radius:var(--r);padding:14px}
.cart-title{font-size:14px;font-weight:700;margin-bottom:10px}
.cart-items{display:flex;flex-direction:column;gap:6px;max-height:200px;overflow-y:auto}
.cart-item{padding:8px 10px;background:var(--surface);border-radius:var(--rs)}
.ci-top{display:flex;justify-content:space-between}.ci-name{font-weight:600;font-size:12px}.ci-price{font-size:11px;color:var(--t3)}
.ci-bottom{display:flex;align-items:center;gap:8px;margin-top:4px}
.ci-qty{display:flex;align-items:center;gap:6px}
.qb{width:24px;height:24px;border-radius:6px;border:1px solid var(--border);background:var(--card);color:var(--t1);font-size:13px;cursor:pointer;display:flex;align-items:center;justify-content:center}
.qb:hover{border-color:var(--cyan);color:var(--cyan)}
.qv{font-weight:700;font-size:13px;min-width:20px;text-align:center}
.ci-sub{margin-left:auto;font-weight:700;font-size:12px;color:var(--cyan)}
.rm-btn{width:24px;height:24px;border-radius:6px;border:none;background:rgba(248,113,113,.15);color:var(--red);cursor:pointer;font-size:11px;display:flex;align-items:center;justify-content:center;margin-left:6px}
.cart-summary{margin-top:12px;padding-top:10px;border-top:1px solid var(--border);display:flex;flex-direction:column;gap:6px}
.sum-row{display:flex;justify-content:space-between;font-size:13px;color:var(--t2)}
.sum-row.discount{color:var(--green)}
.sum-row.total span:last-child{color:var(--cyan)}
.submit-btn{width:100%;padding:12px;justify-content:center;font-size:15px}

/* Print section */
.print-section{margin-top:8px;display:flex;flex-direction:column;gap:8px}
.print-actions{display:flex;gap:8px}
.print-btn{flex:1;padding:10px;border-radius:var(--rs);border:1px solid var(--cyan);background:var(--cyan-g);color:var(--cyan);font-family:inherit;font-size:13px;font-weight:700;cursor:pointer;transition:.2s;text-align:center}
.print-btn:hover:not(:disabled){background:var(--cyan);color:#fff}
.print-btn:disabled{opacity:.6;cursor:not-allowed}
.print-btn-secondary{padding:10px 16px;border-radius:var(--rs);border:1px solid var(--border);background:var(--surface);color:var(--t3);font-family:inherit;font-size:13px;font-weight:600;cursor:pointer;transition:.2s}
.print-btn-secondary:hover{border-color:var(--red);color:var(--red)}
</style>