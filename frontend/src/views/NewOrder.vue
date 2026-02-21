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

        <div v-if="tab==='wash'" class="menu-grid">
          <div v-for="p in packages" :key="p.id" class="menu-card" @click="addToCart(p,'package')">
            <span class="mc-icon">{{ p.icon }}</span>
            <div class="mc-info">
              <span class="mc-name">{{ p.name }}</span>
              <span class="mc-desc">{{ p.description }}</span>
              <div class="mc-meta"><span class="mc-price">{{ fmt(p.price) }}</span><span class="mc-dur">⏱ {{ p.duration }} mnt</span></div>
            </div>
          </div>
        </div>

        <div v-if="tab==='food'">
          <div class="sub-tabs">
            <button v-for="c in foodCats" :key="c.key" :class="['sub-tab',{active:foodCat===c.key}]" @click="foodCat=c.key">{{ c.icon }} {{ c.label }}</button>
          </div>
          <div class="menu-grid food-grid">
            <div v-for="m in filteredMenu" :key="m.id" class="menu-card" @click="addToCart(m,'menu')">
              <span class="mc-icon">{{ m.icon }}</span>
              <div class="mc-info">
                <span class="mc-name">{{ m.name }}</span>
                <div class="mc-meta"><span class="mc-price">{{ fmt(m.price) }}</span><span class="mc-stock">Stok: {{ m.stock }}</span></div>
              </div>
            </div>
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
          <div class="form-group">
            <label class="label">Pembayaran</label>
            <div class="pay-grid">
              <button v-for="pm in payMethods" :key="pm" :class="['pay-btn',{active:customer.payment===pm.v}]" @click="customer.payment=pm.v">{{ pm.i }} {{ pm.l }}</button>
            </div>
          </div>
        </div>

        <button class="btn btn-primary submit-btn" :disabled="cart.length===0||submitting" @click="submitOrder">
          {{ submitting ? '⏳ Memproses...' : '✅ Buat Order' }}
        </button>
        <div v-if="successMsg" class="toast toast-success">{{ successMsg }}</div>
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
const successMsg = ref('')
const customer = ref({ name: '', plate: '', car: '', payment: 'cash' })

const foodCats = [{ key:'all', icon:'📋', label:'Semua' },{ key:'makanan', icon:'🍛', label:'Makanan' },{ key:'minuman', icon:'🥤', label:'Minuman' },{ key:'snack', icon:'🍟', label:'Snack' }]
const payMethods = [{ v:'cash', i:'💵', l:'Cash' },{ v:'debit', i:'💳', l:'Debit' },{ v:'qris', i:'📱', l:'QRIS' },{ v:'transfer', i:'🏦', l:'Transfer' }]

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
    await callBackend('CreateTransaction', customer.value.name || 'Walk-in', customer.value.plate || '-', customer.value.car || '-', customer.value.payment, cart.value, discId)
    cart.value = []; selectedDiscount.value = null; customer.value = { name:'', plate:'', car:'', payment:'cash' }
    successMsg.value = '✅ Order berhasil dibuat!'
    setTimeout(() => successMsg.value = '', 3000)
    await loadData()
  } catch (e) { alert(e.message || 'Gagal membuat order') }
  finally { submitting.value = false }
}

async function loadData() {
  try {
    const [p, m, d] = await Promise.all([callBackend('GetActivePackages'), callBackend('GetActiveMenuItems'), callBackend('GetActiveDiscounts')])
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
.menu-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:10px;overflow-y:auto;flex:1;padding-right:6px}
.food-grid{grid-template-columns:repeat(3,1fr)}
.menu-card{background:var(--card);border:1px solid var(--border);border-radius:var(--r);padding:14px;cursor:pointer;transition:.2s;display:flex;gap:10px}
.menu-card:hover{border-color:var(--cyan);background:var(--card-h);transform:translateY(-1px)}
.menu-card:active{transform:scale(.98)}
.mc-icon{font-size:26px;flex-shrink:0}
.mc-info{display:flex;flex-direction:column;gap:3px;min-width:0}
.mc-name{font-weight:700;font-size:13px}
.mc-desc{font-size:11px;color:var(--t3);line-height:1.3}
.mc-meta{display:flex;gap:10px;margin-top:3px}
.mc-price{font-weight:700;font-size:13px;color:var(--cyan)}
.mc-dur,.mc-stock{font-size:10px;color:var(--t3)}
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
.sum-row.total{font-size:18px;font-weight:800;color:var(--t1);padding-top:6px;border-top:1px solid var(--border)}
.sum-row.total span:last-child{color:var(--cyan)}
.pay-grid{display:grid;grid-template-columns:1fr 1fr;gap:6px}
.pay-btn{padding:8px;border-radius:var(--rs);border:1px solid var(--border);background:var(--surface);color:var(--t2);font-family:inherit;font-size:12px;font-weight:600;cursor:pointer;transition:.2s}
.pay-btn.active{border-color:var(--cyan);background:var(--cyan-g);color:var(--cyan)}
.submit-btn{width:100%;padding:12px;justify-content:center;font-size:15px}
</style>
