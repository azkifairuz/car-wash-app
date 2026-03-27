<template>
  <div class="fade-in">
    <h2 class="page-title">🖨️ Pengaturan Printer</h2>

    <div class="settings-layout">
      <!-- Left: Config -->
      <div class="config-panel">
        <div class="card">
          <h3 class="card-title">Mode Koneksi Printer</h3>
          <p class="card-desc">Pilih cara printer thermal terhubung ke komputer</p>

          <div class="mode-grid">
            <div :class="['mode-card', { active: cfg.mode === 'browser' }]" @click="cfg.mode='browser'">
              <span class="mode-icon">🌐</span>
              <div class="mode-info">
                <span class="mode-name">Browser Print</span>
                <span class="mode-desc">Cetak via browser OS. Universal, bisa pakai printer apapun yang sudah terinstall. Muncul dialog print.</span>
              </div>
              <span v-if="cfg.mode==='browser'" class="mode-check">✓</span>
            </div>

            <div :class="['mode-card', { active: cfg.mode === 'network' }]" @click="cfg.mode='network'">
              <span class="mode-icon">📡</span>
              <div class="mode-info">
                <span class="mode-name">Network / LAN</span>
                <span class="mode-desc">Direct print via IP. Cepat, tanpa dialog. Cocok untuk printer thermal yang konek via ethernet/wifi.</span>
              </div>
              <span v-if="cfg.mode==='network'" class="mode-check">✓</span>
            </div>

            <div :class="['mode-card', { active: cfg.mode === 'usb' }]" @click="cfg.mode='usb'">
              <span class="mode-icon">🔌</span>
              <div class="mode-info">
                <span class="mode-name">USB / Shared Printer</span>
                <span class="mode-desc">Print via nama printer yang di-share di Windows. Cepat, tanpa dialog. Perlu setup sharing di Windows dulu.</span>
              </div>
              <span v-if="cfg.mode==='usb'" class="mode-check">✓</span>
            </div>
          </div>
        </div>

        <!-- Network Settings -->
        <div v-if="cfg.mode === 'network'" class="card">
          <h3 class="card-title">📡 Pengaturan Network</h3>
          <div class="form-row">
            <div class="form-group" style="flex:2">
              <label class="label">IP Address Printer</label>
              <input class="input" v-model="cfg.networkIp" placeholder="192.168.1.100" />
            </div>
            <div class="form-group" style="flex:1">
              <label class="label">Port</label>
              <input class="input" v-model="cfg.networkPort" placeholder="9100" />
            </div>
          </div>
          <div class="hint">Biasanya port default printer thermal adalah 9100. Cek IP di setting printer atau label di belakang printer.</div>
        </div>

        <!-- USB Settings -->
        <div v-if="cfg.mode === 'usb'" class="card">
          <h3 class="card-title">🔌 Pengaturan USB / Shared Printer</h3>
          <div class="form-group">
            <label class="label">Nama Printer (Windows Share Name)</label>
            <input class="input" v-model="cfg.usbName" placeholder="ThermalPrinter" />
          </div>
          <div class="hint">
            Cara setup di Windows: Buka Control Panel → Devices and Printers → Klik kanan printer → Printer Properties → tab Sharing → centang "Share this printer" → masukkan nama share.
          </div>
        </div>

        <!-- Paper Size -->
        <div class="card">
          <h3 class="card-title">📐 Ukuran Kertas</h3>
          <div class="paper-grid">
            <div :class="['paper-opt', { active: cfg.paperWidth === '58' }]" @click="cfg.paperWidth='58'">
              <span class="paper-size">58mm</span>
              <span class="paper-desc">32 karakter/baris</span>
            </div>
            <div :class="['paper-opt', { active: cfg.paperWidth === '80' }]" @click="cfg.paperWidth='80'">
              <span class="paper-size">80mm</span>
              <span class="paper-desc">48 karakter/baris</span>
            </div>
          </div>
        </div>

        <!-- Store Info -->
        <div class="card">
          <h3 class="card-title">🏪 Info Toko (Header Struk)</h3>
          <div class="form-group">
            <label class="label">Nama Toko</label>
            <input class="input" v-model="cfg.storeName" placeholder="SparkleWash" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="label">Alamat</label>
              <input class="input" v-model="cfg.storeAddr" placeholder="Jl. Raya No. 123" />
            </div>
            <div class="form-group">
              <label class="label">Telepon</label>
              <input class="input" v-model="cfg.storePhone" placeholder="0812-3456-7890" />
            </div>
          </div>
          <div class="form-group">
            <label class="label">Teks Footer Struk</label>
            <input class="input" v-model="cfg.footerText" placeholder="Terima Kasih!" />
          </div>
        </div>

        <!-- Auto Print -->
        <div class="card">
          <h3 class="card-title">⚡ Otomatis</h3>
          <label class="toggle-row">
            <input type="checkbox" v-model="cfg.autoPrint" class="toggle-input" />
            <span class="toggle-slider"></span>
            <span class="toggle-label">Auto cetak struk setelah order dibuat</span>
          </label>
        </div>

        <!-- Actions -->
        <div class="actions-row">
          <button class="btn btn-primary" @click="saveConfig" :disabled="saving">
            {{ saving ? '⏳ Menyimpan...' : '💾 Simpan Pengaturan' }}
          </button>
          <button class="btn btn-secondary" @click="testPrint" :disabled="testing">
            {{ testing ? '⏳ Testing...' : '🖨️ Test Print' }}
          </button>
        </div>

        <div v-if="saveMsg" :class="['toast', saveMsg.startsWith('✅') ? 'toast-success' : 'toast-error']">{{ saveMsg }}</div>
      </div>

      <!-- Right: Preview & Help -->
      <div class="help-panel">
        <div class="card">
          <h3 class="card-title">📋 Panduan Setup</h3>

          <div v-if="cfg.mode === 'browser'" class="help-content">
            <p><strong>Browser Print</strong> adalah mode paling mudah. Tidak perlu konfigurasi tambahan.</p>
            <ol>
              <li>Pastikan printer thermal sudah terinstall di komputer (lewat driver USB atau network)</li>
              <li>Set printer sebagai default printer di Windows</li>
              <li>Saat cetak, browser akan terbuka dan dialog print muncul</li>
              <li>Pilih printer yang benar, lalu klik Print</li>
            </ol>
          </div>

          <div v-if="cfg.mode === 'network'" class="help-content">
            <p><strong>Network Print</strong> mengirim perintah ESC/POS langsung ke printer via TCP.</p>
            <ol>
              <li>Pastikan printer dan komputer di jaringan yang sama (LAN/WiFi)</li>
              <li>Cek IP printer di: panel printer, label belakang, atau print self-test</li>
              <li>Default port biasanya <code>9100</code></li>
              <li>Test koneksi: buka CMD → ketik <code>telnet [IP] 9100</code></li>
            </ol>
          </div>

          <div v-if="cfg.mode === 'usb'" class="help-content">
            <p><strong>USB / Shared Printer</strong> mengirim perintah ESC/POS lewat file sharing Windows.</p>
            <ol>
              <li>Colok printer USB, install driver</li>
              <li>Buka <strong>Control Panel → Devices and Printers</strong></li>
              <li>Klik kanan printer → <strong>Printer Properties</strong></li>
              <li>Tab <strong>Sharing</strong> → centang <strong>Share this printer</strong></li>
              <li>Masukkan nama share (contoh: <code>ThermalPrinter</code>)</li>
              <li>Masukkan nama yang sama di kolom "Nama Printer" di atas</li>
            </ol>
          </div>
        </div>

        <div class="card">
          <h3 class="card-title">🖨️ Printer yang Didukung</h3>
          <div class="help-content">
            <p>Semua printer thermal ESC/POS compatible, termasuk:</p>
            <ul>
              <li>Epson TM series (TM-T82, TM-T88, dll)</li>
              <li>Xprinter (XP-58, XP-80, dll)</li>
              <li>BIXOLON</li>
              <li>Star Micronics</li>
              <li>Zjiang</li>
              <li>GOOJPRT</li>
              <li>Dan printer POS lainnya</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { callBackend } from '../stores/auth'

const cfg = ref({
  mode: 'browser',
  networkIp: '',
  networkPort: '9100',
  usbName: '',
  paperWidth: '80',
  storeName: 'SparkleWash',
  storeAddr: '',
  storePhone: '',
  footerText: 'Terima Kasih!',
  autoPrint: false,
})

const saving = ref(false)
const testing = ref(false)
const saveMsg = ref('')

async function loadConfig() {
  try {
    const c = await callBackend('GetPrinterConfig')
    if (c) {
      cfg.value = { ...cfg.value, ...c }
    }
  } catch (e) {}
}

async function saveConfig() {
  saving.value = true
  saveMsg.value = ''
  try {
    await callBackend('SavePrinterConfig',
      cfg.value.mode,
      cfg.value.networkIp,
      cfg.value.networkPort,
      cfg.value.usbName,
      cfg.value.paperWidth,
      cfg.value.storeName,
      cfg.value.storeAddr,
      cfg.value.storePhone,
      cfg.value.footerText,
      cfg.value.autoPrint,
    )
    saveMsg.value = '✅ Pengaturan printer berhasil disimpan!'
    setTimeout(() => saveMsg.value = '', 3000)
  } catch (e) {
    saveMsg.value = '❌ Gagal menyimpan: ' + (e.message || e)
  } finally {
    saving.value = false
  }
}

async function testPrint() {
  testing.value = true
  saveMsg.value = ''
  try {
    // Save first, then test
    await saveConfig()
    await callBackend('TestPrinter')
    saveMsg.value = '✅ Test print berhasil dikirim!'
    setTimeout(() => saveMsg.value = '', 3000)
  } catch (e) {
    saveMsg.value = '❌ Test print gagal: ' + (e.message || e)
  } finally {
    testing.value = false
  }
}

onMounted(loadConfig)
</script>

<style scoped>
.settings-layout{display:grid;grid-template-columns:1fr 380px;gap:20px;height:calc(100vh - 120px);overflow-y:auto}
.config-panel{display:flex;flex-direction:column;gap:14px;overflow-y:auto;padding-right:6px}
.help-panel{display:flex;flex-direction:column;gap:14px;overflow-y:auto}

.card{background:var(--card);border:1px solid var(--border);border-radius:var(--r);padding:16px}
.card-title{font-size:14px;font-weight:700;margin-bottom:6px}
.card-desc{font-size:12px;color:var(--t3);margin-bottom:12px}

/* Mode selection */
.mode-grid{display:flex;flex-direction:column;gap:8px}
.mode-card{display:flex;align-items:flex-start;gap:12px;padding:12px;border:2px solid var(--border);border-radius:var(--r);cursor:pointer;transition:.2s;position:relative}
.mode-card:hover{border-color:var(--t3)}
.mode-card.active{border-color:var(--cyan);background:var(--cyan-g)}
.mode-icon{font-size:24px;flex-shrink:0;margin-top:2px}
.mode-info{display:flex;flex-direction:column;gap:2px;flex:1}
.mode-name{font-weight:700;font-size:13px}
.mode-desc{font-size:11px;color:var(--t3);line-height:1.4}
.mode-check{position:absolute;top:10px;right:12px;width:22px;height:22px;border-radius:50%;background:var(--cyan);color:#fff;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:700}

/* Paper size */
.paper-grid{display:flex;gap:10px}
.paper-opt{flex:1;padding:14px;border:2px solid var(--border);border-radius:var(--r);cursor:pointer;text-align:center;transition:.2s}
.paper-opt:hover{border-color:var(--t3)}
.paper-opt.active{border-color:var(--cyan);background:var(--cyan-g)}
.paper-size{display:block;font-size:18px;font-weight:800}
.paper-desc{display:block;font-size:11px;color:var(--t3);margin-top:2px}

/* Toggle */
.toggle-row{display:flex;align-items:center;gap:12px;cursor:pointer}
.toggle-input{display:none}
.toggle-slider{width:44px;height:24px;border-radius:12px;background:var(--border);position:relative;transition:.3s;flex-shrink:0}
.toggle-slider::after{content:'';position:absolute;top:3px;left:3px;width:18px;height:18px;border-radius:50%;background:#fff;transition:.3s}
.toggle-input:checked+.toggle-slider{background:var(--cyan)}
.toggle-input:checked+.toggle-slider::after{transform:translateX(20px)}
.toggle-label{font-size:13px;font-weight:600}

/* Help */
.help-content{font-size:12px;color:var(--t2);line-height:1.6}
.help-content ol,.help-content ul{margin:8px 0;padding-left:20px}
.help-content li{margin:4px 0}
.help-content code{background:var(--surface);padding:1px 5px;border-radius:4px;font-size:11px}
.help-note{margin-top:10px;padding:8px 10px;border-radius:var(--rs);background:var(--surface);font-size:11px;font-weight:600}

.hint{font-size:11px;color:var(--t3);margin-top:6px;line-height:1.4}

.actions-row{display:flex;gap:10px}
.actions-row .btn{flex:1;padding:12px;justify-content:center;font-size:14px}
.btn-secondary{background:var(--surface);border:1px solid var(--border);color:var(--t2);border-radius:var(--rs);font-family:inherit;font-weight:600;cursor:pointer;transition:.2s}
.btn-secondary:hover{border-color:var(--cyan);color:var(--cyan)}
</style>