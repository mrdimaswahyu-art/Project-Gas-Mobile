package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"go_vibe/config"
	"go_vibe/handlers"
	"go_vibe/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// ============================================================
// ASET UI GARUDA AUTO SERVICE
// Seluruh HTML, CSS, dan JavaScript ditanam langsung di main.go
// agar aplikasi tidak bergantung pada folder web/.
// ============================================================

const indexHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#9d0208">
    <title>Garuda Auto Service — Dashboard Bengkel</title>
    <link rel="icon" type="image/svg+xml" href="/favicon.svg?v=4">
    <link rel="stylesheet" href="/css/style.css?v=4">
</head>
<body>
    <!-- Navigasi Atas -->
    <header class="navbar">
        <div class="nav-brand">
            <span class="brand-badge" aria-hidden="true">
                <svg viewBox="0 0 100 48">
                    <path d="M6 32c0-5 4-8 9-9l14-3c6-8 14-12 24-12h6c10 0 18 4 23 12l10 3c6 1 9 4 9 9v5c0 2-1 3-3 3h-8a10 10 0 0 1-20 0H33a10 10 0 0 1-20 0h-4c-2 0-3-1-3-3z" fill="currentColor"/>
                    <circle cx="23" cy="40" r="6.5" fill="#fff"/>
                    <circle cx="23" cy="40" r="2.5" fill="#7a0016"/>
                    <circle cx="72" cy="40" r="6.5" fill="#fff"/>
                    <circle cx="72" cy="40" r="2.5" fill="#7a0016"/>
                </svg>
            </span>
            <span class="brand-stack">
                <span class="brand-name">Garuda <em>Auto Service</em></span>
                <span class="brand-tagline">Bengkel Mobil · Servis Cepat &amp; Terpercaya</span>
            </span>
        </div>
        <div class="nav-user">
            <span class="nav-greeting">Halo, <strong id="currentUsername">Staf</strong></span>
            <button onclick="logout()" class="btn btn-secondary btn-sm">Logout</button>
        </div>
    </header>

    <div class="container">
        <!-- Hero Bermotif Mobil -->
        <section class="hero" aria-labelledby="heroTitle">
            <div class="hero-copy">
                <span class="hero-kicker"><span class="kicker-dot"></span> Pit Stop Servis Mobil</span>
                <h1 class="hero-title" id="heroTitle">Garuda <span>Auto Service</span></h1>
                <p class="hero-text">Catat konsumen, kendaraan, dan seluruh riwayat servis dalam satu dasbor bengkel yang cepat dan rapi.</p>
                <div class="hero-chips" aria-label="Layanan bengkel">
                    <span class="chip">Servis Berkala</span>
                    <span class="chip">Tune Up</span>
                    <span class="chip">Kaki-Kaki</span>
                    <span class="chip">Body Repair</span>
                </div>
                <div class="hero-metrics">
                    <div class="metric">
                        <span class="metric-value" id="statVehicles">0</span>
                        <span class="metric-label">Kendaraan</span>
                    </div>
                    <div class="metric">
                        <span class="metric-value" id="statServices">0</span>
                        <span class="metric-label">Riwayat Servis</span>
                    </div>
                </div>
            </div>

            <div class="hero-visual">
                <div class="hero-roundel" aria-hidden="true">
                    <span class="roundel-inner">
                        <b>01</b>
                        <small>GARUDA</small>
                    </span>
                </div>
                <div class="speed-lines" aria-hidden="true"><i></i><i></i><i></i></div>

                <svg class="hero-car" viewBox="0 0 520 200" role="img" aria-label="Ilustrasi mobil sport merah Garuda Auto Service">
                    <defs>
                        <linearGradient id="carBody" x1="0" y1="0" x2="0.2" y2="1">
                            <stop offset="0" stop-color="#ff6b78"/>
                            <stop offset="0.45" stop-color="#e11d33"/>
                            <stop offset="1" stop-color="#7a0016"/>
                        </linearGradient>
                        <linearGradient id="carGlass" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="0" stop-color="#a8bccf" stop-opacity="0.9"/>
                            <stop offset="1" stop-color="#11161d"/>
                        </linearGradient>
                        <linearGradient id="carRim" x1="0" y1="0" x2="1" y2="1">
                            <stop offset="0" stop-color="#f2f5f8"/>
                            <stop offset="1" stop-color="#8b949e"/>
                        </linearGradient>
                        <radialGradient id="headlightGlow" cx="0.5" cy="0.5" r="0.5">
                            <stop offset="0" stop-color="#fff6c9" stop-opacity="0.95"/>
                            <stop offset="1" stop-color="#ffd166" stop-opacity="0"/>
                        </radialGradient>
                    </defs>

                    <!-- Bayangan -->
                    <ellipse cx="272" cy="176" rx="202" ry="12" fill="#1a0d10" opacity="0.3"/>

                    <!-- Bodi -->
                    <path d="M50 140 C48 126 58 116 80 112 L136 104 C154 80 184 64 222 61 L306 60 C348 61 380 74 404 96 L434 104 L474 109 C497 113 508 124 509 138 C510 148 504 152 494 152 L92 152 C64 152 50 148 50 140 Z" fill="url(#carBody)"/>

                    <!-- Highlight bodi -->
                    <path d="M86 120 C170 110 380 106 482 116" fill="none" stroke="#ff9aa5" stroke-width="3" stroke-linecap="round" opacity="0.42"/>

                    <!-- Kaca samping -->
                    <path d="M154 102 C170 83 194 72 224 69 L270 68 L270 102 Z" fill="url(#carGlass)"/>
                    <path d="M282 68 L300 68 C330 69 356 80 376 98 L376 102 L282 102 Z" fill="url(#carGlass)"/>

                    <!-- Garis pintu -->
                    <path d="M276 70 L276 146" stroke="#7a0016" stroke-width="2" opacity="0.55"/>
                    <path d="M160 108 L163 146" stroke="#7a0016" stroke-width="2" opacity="0.3"/>

                    <!-- Handle pintu -->
                    <rect x="296" y="116" width="20" height="5" rx="2.5" fill="#2a0d12" opacity="0.6"/>
                    <rect x="206" y="118" width="18" height="5" rx="2.5" fill="#2a0d12" opacity="0.45"/>

                    <!-- Side skirt -->
                    <path d="M94 148 L482 148 L490 152 L92 152 Z" fill="#3d0009" opacity="0.55"/>

                    <!-- Lampu depan + sorot -->
                    <path d="M470 116 L496 122 L497 131 L474 128 Z" fill="#fff3c4"/>
                    <circle cx="508" cy="128" r="18" fill="url(#headlightGlow)"/>

                    <!-- Lampu belakang -->
                    <path d="M64 124 L88 124 L86 136 L66 134 Z" fill="#ff2d44"/>
                    <path d="M68 127 L82 127 L81 133 L69 132 Z" fill="#ffe3e6" opacity="0.7"/>

                    <!-- Grill -->
                    <rect x="482" y="136" width="20" height="6" rx="3" fill="#1a0d10" opacity="0.5"/>

                    <!-- Rumah roda -->
                    <path d="M96 152 A38 38 0 0 1 172 152 Z" fill="#12070a"/>
                    <path d="M306 152 A38 38 0 0 1 382 152 Z" fill="#12070a"/>

                    <!-- Roda belakang -->
                    <g class="wheel">
                        <circle cx="134" cy="142" r="30" fill="#0b0608"/>
                        <circle cx="134" cy="142" r="30" fill="none" stroke="#2a1b20" stroke-width="4"/>
                        <circle cx="134" cy="142" r="17" fill="url(#carRim)"/>
                        <g class="wheel-spokes" stroke="#3d444b" stroke-width="3" stroke-linecap="round">
                            <line x1="134" y1="129" x2="134" y2="155"/>
                            <line x1="121" y1="142" x2="147" y2="142"/>
                            <line x1="125" y1="133" x2="143" y2="151"/>
                            <line x1="143" y1="133" x2="125" y2="151"/>
                        </g>
                        <circle cx="134" cy="142" r="5" fill="#c1121f"/>
                    </g>

                    <!-- Roda depan -->
                    <g class="wheel">
                        <circle cx="344" cy="142" r="30" fill="#0b0608"/>
                        <circle cx="344" cy="142" r="30" fill="none" stroke="#2a1b20" stroke-width="4"/>
                        <circle cx="344" cy="142" r="17" fill="url(#carRim)"/>
                        <g class="wheel-spokes" stroke="#3d444b" stroke-width="3" stroke-linecap="round">
                            <line x1="344" y1="129" x2="344" y2="155"/>
                            <line x1="331" y1="142" x2="357" y2="142"/>
                            <line x1="335" y1="133" x2="353" y2="151"/>
                            <line x1="353" y1="133" x2="335" y2="151"/>
                        </g>
                        <circle cx="344" cy="142" r="5" fill="#c1121f"/>
                    </g>

                    <!-- Nomor balap -->
                    <circle cx="250" cy="124" r="13" fill="#fff" opacity="0.92"/>
                    <text x="250" y="129" text-anchor="middle" font-size="13" font-weight="700" fill="#9d0208" font-family="Arial, sans-serif">01</text>
                </svg>
            </div>
        </section>

        <!-- Statistik Bulanan -->
        <section class="monthly-panel" aria-labelledby="monthlyTitle">
            <div class="monthly-head">
                <div>
                    <span class="panel-kicker">Statistik Bulanan</span>
                    <h2 id="monthlyTitle">Diagram Konsumen Bulanan</h2>
                    <p class="monthly-subtitle" id="monthlySubtitle">Memuat data 12 bulan terakhir...</p>
                </div>
                <div class="monthly-summary" aria-label="Ringkasan bulan ini">
                    <div class="summary-pill">
                        <span id="monthlyCustomers">0</span>
                        <small>Konsumen Bulan Ini</small>
                    </div>
                    <div class="summary-pill">
                        <span id="monthlyServices">0</span>
                        <small>Servis Bulan Ini</small>
                    </div>
                </div>
            </div>
            <div class="chart-legend">
                <span class="legend-item"><i class="legend-dot dot-customer"></i> Konsumen Baru</span>
                <span class="legend-item"><i class="legend-dot dot-service"></i> Riwayat Servis</span>
            </div>
            <div class="chart-scroll">
                <div id="monthlyChart" class="monthly-chart" role="img" aria-label="Diagram batang jumlah konsumen baru dan riwayat servis per bulan"></div>
            </div>
        </section>

        <!-- Toolbar Pencarian & Filter -->
        <div class="toolbar">
            <div class="toolbar-controls">
                <input type="text" id="searchInput" class="form-control toolbar-search" placeholder="🔍 Cari plat nomor, tipe mobil, atau nama konsumen...">
                <select id="brandFilter" class="form-control toolbar-filter">
                    <option value="">Semua Merk Mobil</option>
                </select>
            </div>
            <button onclick="openCustomerModal()" class="btn btn-primary">+ Tambah Konsumen</button>
        </div>

        <div id="alertFeedback" class="alert"></div>

        <!-- Daftar Konsumen & Riwayat -->
        <div id="customerList" class="card-list">
            <!-- Kartu konsumen akan dirender via JavaScript -->
        </div>

        <div id="emptyState" class="empty-state" style="display: none;">
            <h3>Belum Ada Data</h3>
            <p>Tidak ditemukan data konsumen atau kendaraan yang sesuai dengan pencarian.</p>
        </div>
    </div>

    <!-- Modal Form Konsumen (Tambah / Edit) -->
    <div id="customerModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2 id="customerModalTitle">Tambah Konsumen Baru</h2>
                <button type="button" class="close-btn" onclick="closeCustomerModal()">&times;</button>
            </div>
            <form id="customerForm">
                <input type="hidden" id="customerId">
                <div class="form-group">
                    <label for="custName">Nama Konsumen *</label>
                    <input type="text" id="custName" class="form-control" placeholder="Contoh: Budi Santoso" required>
                </div>
                <div class="form-group">
                    <label for="custPhone">No. Telepon / WhatsApp</label>
                    <input type="text" id="custPhone" class="form-control" placeholder="Contoh: 081234567890">
                </div>
                <div class="form-group">
                    <label for="custAddress">Alamat Konsumen</label>
                    <textarea id="custAddress" class="form-control" rows="2" placeholder="Alamat tinggal / domisili"></textarea>
                </div>
                <div class="form-group">
                    <label for="custPlate">Plat Nomor Kendaraan *</label>
                    <input type="text" id="custPlate" class="form-control" placeholder="Contoh: B 1234 XYZ" required>
                </div>
                <div class="form-group">
                    <label for="custBrand">Merk Mobil *</label>
                    <input type="text" id="custBrand" class="form-control" placeholder="Contoh: Toyota, Honda, Suzuki" required>
                </div>
                <div class="form-group">
                    <label for="custType">Tipe Mobil *</label>
                    <input type="text" id="custType" class="form-control" placeholder="Contoh: Avanza 1.3 G, Brio Satya" required>
                </div>
                <div class="modal-footer">
                    <button type="button" onclick="closeCustomerModal()" class="btn btn-secondary">Batal</button>
                    <button type="submit" class="btn btn-primary" id="btnSaveCustomer">Simpan Data</button>
                </div>
            </form>
        </div>
    </div>

    <!-- Modal Tambah Riwayat Servis -->
    <div id="serviceModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2 id="serviceModalTitle">Tambah Riwayat Servis</h2>
                <button type="button" class="close-btn" onclick="closeServiceModal()">&times;</button>
            </div>
            <form id="serviceForm">
                <input type="hidden" id="serviceCustomerId">
                <div class="form-group">
                    <label for="serviceDate">Tanggal Servis *</label>
                    <input type="date" id="serviceDate" class="form-control" required>
                </div>
                <div class="form-group">
                    <label for="serviceMileage">Kilometer (KM) Saat Ini</label>
                    <input type="number" id="serviceMileage" class="form-control" placeholder="Contoh: 45000">
                </div>
                <div class="form-group">
                    <label for="serviceDesc">Rincian Servis / Pengerjaan *</label>
                    <textarea id="serviceDesc" class="form-control" rows="3" placeholder="Contoh: Ganti oli mesin Shell 10W-40, filter oli, kuras air radiator" required></textarea>
                </div>
                <div class="form-group">
                    <label for="nextServiceDate">Estimasi Jadwal Servis Berikutnya (Follow Up)</label>
                    <input type="date" id="nextServiceDate" class="form-control">
                </div>
                <div class="form-group">
                    <label for="nextMileage">Estimasi KM Servis Berikutnya</label>
                    <input type="number" id="nextMileage" class="form-control" placeholder="Contoh: 50000">
                </div>
                <div class="form-group">
                    <label for="serviceNotes">Catatan Tambahan / Keluhan</label>
                    <input type="text" id="serviceNotes" class="form-control" placeholder="Contoh: Bunyi decit saat rem mendadak">
                </div>
                <div class="modal-footer">
                    <button type="button" onclick="closeServiceModal()" class="btn btn-secondary">Batal</button>
                    <button type="submit" class="btn btn-success">Simpan Servis</button>
                </div>
            </form>
        </div>
    </div>

    <!-- Modal Tampil Detail Riwayat Servis -->
    <div id="viewHistoryModal" class="modal">
        <div class="modal-content" style="max-width: 650px;">
            <div class="modal-header">
                <h2 id="viewHistoryTitle">Riwayat Servis Kendaraan</h2>
                <button type="button" class="close-btn" onclick="closeViewHistoryModal()">&times;</button>
            </div>
            <div id="viewHistoryCustomerInfo" style="margin-bottom: 1rem; font-size: 0.9rem;"></div>
            <div id="viewHistoryTableContainer">
                <!-- Tabel riwayat disuntikkan melalui JS -->
            </div>
            <div class="modal-footer">
                <button type="button" onclick="closeViewHistoryModal()" class="btn btn-secondary">Tutup</button>
            </div>
        </div>
    </div>

    <script src="/js/auth.js?v=4"></script>
    <script src="/js/app.js?v=4"></script>
</body>
</html>
`

const loginHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#0b0f19">
    <title>Login - Garuda Auto Service</title>
    <link rel="icon" type="image/svg+xml" href="/favicon.svg?v=4">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Rajdhani:wght@500;600;700&family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #070a12;
            --panel: #0e1424;
            --card: rgba(18, 24, 40, 0.72);
            --card-border: rgba(255, 255, 255, 0.08);
            --primary: #e11d33;
            --primary-2: #ff5a6a;
            --accent: #ffd166;
            --text: #f2f5f8;
            --muted: #9aa7b8;
            --danger: #ff4d6d;
            --success: #2ecc71;
            --radius: 18px;
            --shadow: 0 24px 60px rgba(0, 0, 0, 0.55);
        }

        * { box-sizing: border-box; }

        html, body {
            margin: 0;
            padding: 0;
            height: 100%;
            font-family: "Inter", system-ui, -apple-system, Segoe UI, Roboto, sans-serif;
            color: var(--text);
            background: var(--bg);
            overflow: hidden;
        }

        body {
            position: relative;
            display: flex;
            align-items: stretch;
            justify-content: center;
        }

        /* ===== Animated background ===== */
        .bg {
            position: fixed;
            inset: 0;
            z-index: 0;
            overflow: hidden;
            background:
                radial-gradient(1200px 600px at 80% 20%, rgba(225, 29, 51, 0.22), transparent 60%),
                radial-gradient(900px 500px at 10% 80%, rgba(255, 209, 102, 0.12), transparent 60%),
                linear-gradient(180deg, #070a12 0%, #0b1020 100%);
        }

        .bg-grid {
            position: absolute;
            inset: -50%;
            background-image:
                linear-gradient(rgba(255,255,255,0.04) 1px, transparent 1px),
                linear-gradient(90deg, rgba(255,255,255,0.04) 1px, transparent 1px);
            background-size: 48px 48px;
            transform: perspective(500px) rotateX(60deg);
            transform-origin: center top;
            animation: gridMove 6s linear infinite;
            opacity: 0.35;
        }

        @keyframes gridMove {
            from { background-position: 0 0, 0 0; }
            to { background-position: 0 48px, 48px 0; }
        }

        /* speed lines */
        .speed-lines {
            position: absolute;
            inset: 0;
            pointer-events: none;
        }
        .speed-lines span {
            position: absolute;
            height: 2px;
            width: 120px;
            background: linear-gradient(90deg, transparent, rgba(255,255,255,0.5), transparent);
            animation: speedLine 1.8s linear infinite;
            opacity: 0.6;
        }
        @keyframes speedLine {
            from { transform: translateX(-200px); opacity: 0; }
            20% { opacity: 0.8; }
            to { transform: translateX(120vw); opacity: 0; }
        }

        /* floating particles */
        .particles {
            position: absolute;
            inset: 0;
            pointer-events: none;
        }
        .particles i {
            position: absolute;
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background: var(--accent);
            box-shadow: 0 0 12px var(--accent);
            animation: floatUp 6s ease-in-out infinite;
            opacity: 0.7;
        }
        @keyframes floatUp {
            0% { transform: translateY(20px); opacity: 0; }
            20% { opacity: 0.8; }
            80% { opacity: 0.8; }
            100% { transform: translateY(-120px); opacity: 0; }
        }

        /* ===== Layout ===== */
        .wrap {
            position: relative;
            z-index: 1;
            width: 100%;
            max-width: 1180px;
            margin: auto;
            padding: 2rem;
            display: grid;
            grid-template-columns: 1.15fr 1fr;
            gap: 2rem;
            align-items: center;
        }

        @media (max-width: 900px) {
            html, body { overflow: auto; }
            .wrap { grid-template-columns: 1fr; padding: 1.2rem; }
            .showcase { display: none; }
        }

        /* ===== Showcase (left) ===== */
        .showcase {
            position: relative;
            padding: 2rem 1rem;
            animation: slideInLeft 0.9s cubic-bezier(.2,.8,.2,1) both;
        }
        @keyframes slideInLeft {
            from { opacity: 0; transform: translateX(-40px); }
            to { opacity: 1; transform: translateX(0); }
        }

        .brand-row {
            display: flex;
            align-items: center;
            gap: 1rem;
            margin-bottom: 1.6rem;
        }
        .logo-badge {
            width: 78px;
            height: 78px;
            border-radius: 22px;
            background: linear-gradient(135deg, var(--primary), #7a0016);
            display: grid;
            place-items: center;
            box-shadow: 0 12px 30px rgba(225,29,51,0.45), inset 0 2px 4px rgba(255,255,255,0.25);
            position: relative;
            overflow: hidden;
        }
        .logo-badge::after {
            content: "";
            position: absolute;
            inset: 0;
            background: linear-gradient(120deg, transparent 30%, rgba(255,255,255,0.35) 50%, transparent 70%);
            animation: shine 3.2s ease-in-out infinite;
        }
        @keyframes shine {
            0% { transform: translateX(-120%); }
            60%, 100% { transform: translateX(120%); }
        }
        .logo-badge svg { width: 56px; height: 56px; }

        .brand-text h2 {
            font-family: "Rajdhani", sans-serif;
            font-size: 1.9rem;
            margin: 0;
            letter-spacing: 1px;
            line-height: 1;
        }
        .brand-text h2 em {
            color: var(--primary-2);
            font-style: normal;
        }
        .brand-text p {
            margin: 4px 0 0;
            color: var(--muted);
            font-size: 0.9rem;
        }

        .showcase h1 {
            font-family: "Rajdhani", sans-serif;
            font-size: clamp(2.6rem, 5vw, 4rem);
            line-height: 1.05;
            margin: 1.2rem 0 0.8rem;
            letter-spacing: 1px;
        }
        .showcase h1 span {
            background: linear-gradient(90deg, #fff, var(--primary-2));
            -webkit-background-clip: text;
            background-clip: text;
            color: transparent;
        }
        .showcase .lead {
            color: var(--muted);
            font-size: 1.05rem;
            max-width: 480px;
            line-height: 1.6;
        }

        .chips {
            display: flex;
            flex-wrap: wrap;
            gap: 0.6rem;
            margin-top: 1.4rem;
        }
        .chip {
            padding: 0.45rem 0.9rem;
            border-radius: 999px;
            background: rgba(255,255,255,0.06);
            border: 1px solid rgba(255,255,255,0.1);
            font-size: 0.82rem;
            color: var(--text);
            backdrop-filter: blur(6px);
        }

        /* animated car stage */
        .stage {
            margin-top: 2.2rem;
            position: relative;
            height: 150px;
            border-radius: var(--radius);
            background: linear-gradient(180deg, rgba(255,255,255,0.04), rgba(255,255,255,0.01));
            border: 1px solid rgba(255,255,255,0.08);
            overflow: hidden;
        }
        .road {
            position: absolute;
            bottom: 28px;
            left: 0;
            right: 0;
            height: 4px;
            background: repeating-linear-gradient(90deg, rgba(255,255,255,0.5) 0 24px, transparent 24px 48px);
            animation: roadMove 0.9s linear infinite;
        }
        @keyframes roadMove {
            from { background-position: 0 0; }
            to { background-position: 48px 0; }
        }
        .stage-car {
            position: absolute;
            bottom: 34px;
            left: -180px;
            width: 180px;
            animation: drive 5s linear infinite;
            filter: drop-shadow(0 8px 16px rgba(225,29,51,0.4));
        }
        @keyframes drive {
            0% { transform: translateX(-200px); }
            100% { transform: translateX(calc(100% + 220px)); }
        }
        .stage-car .wheel {
            transform-origin: center;
            animation: spin 0.6s linear infinite;
        }
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        .headlight {
            animation: glow 1.2s ease-in-out infinite alternate;
        }
        @keyframes glow {
            from { opacity: 0.6; }
            to { opacity: 1; }
        }

        /* ===== Card (right) ===== */
        .card {
            background: var(--card);
            border: 1px solid var(--card-border);
            border-radius: 24px;
            padding: 2.2rem 2rem;
            box-shadow: var(--shadow);
            backdrop-filter: blur(18px);
            -webkit-backdrop-filter: blur(18px);
            animation: slideInRight 0.9s cubic-bezier(.2,.8,.2,1) 0.15s both;
            position: relative;
            overflow: hidden;
        }
        @keyframes slideInRight {
            from { opacity: 0; transform: translateX(40px); }
            to { opacity: 1; transform: translateX(0); }
        }
        .card::before {
            content: "";
            position: absolute;
            top: -60%;
            right: -40%;
            width: 260px;
            height: 260px;
            background: radial-gradient(circle, rgba(225,29,51,0.25), transparent 70%);
            animation: pulse 4s ease-in-out infinite;
        }
        @keyframes pulse {
            0%, 100% { transform: scale(1); opacity: 0.6; }
            50% { transform: scale(1.2); opacity: 1; }
        }

        .card-header {
            text-align: center;
            margin-bottom: 1.6rem;
            position: relative;
        }
        .card-header .mini-logo {
            width: 52px;
            height: 52px;
            margin: 0 auto 0.8rem;
            border-radius: 14px;
            background: linear-gradient(135deg, var(--primary), #7a0016);
            display: grid;
            place-items: center;
            box-shadow: 0 8px 20px rgba(225,29,51,0.4);
        }
        .card-header .mini-logo svg { width: 36px; height: 36px; }
        .card-header h3 {
            font-family: "Rajdhani", sans-serif;
            font-size: 1.6rem;
            margin: 0;
            letter-spacing: 1px;
        }
        .card-header p {
            margin: 4px 0 0;
            color: var(--muted);
            font-size: 0.9rem;
        }

        .form-group {
            margin-bottom: 1.1rem;
            position: relative;
        }
        .form-group label {
            display: block;
            font-size: 0.82rem;
            color: var(--muted);
            margin-bottom: 0.4rem;
            letter-spacing: 0.5px;
            text-transform: uppercase;
        }
        .input-wrap {
            position: relative;
        }
        .input-wrap svg {
            position: absolute;
            left: 14px;
            top: 50%;
            transform: translateY(-50%);
            width: 18px;
            height: 18px;
            color: var(--muted);
            pointer-events: none;
            transition: color 0.2s;
        }
        .form-control {
            width: 100%;
            padding: 0.85rem 1rem 0.85rem 2.6rem;
            background: rgba(255,255,255,0.05);
            border: 1px solid rgba(255,255,255,0.1);
            border-radius: 12px;
            color: var(--text);
            font-size: 0.98rem;
            outline: none;
            transition: border-color 0.2s, background 0.2s, box-shadow 0.2s;
        }
        .form-control::placeholder { color: #5b6677; }
        .form-control:focus {
            border-color: var(--primary);
            background: rgba(225,29,51,0.08);
            box-shadow: 0 0 0 4px rgba(225,29,51,0.15);
        }
        .form-control:focus + svg,
        .input-wrap:focus-within svg { color: var(--primary-2); }

        .toggle-pass {
            position: absolute;
            right: 12px;
            top: 50%;
            transform: translateY(-50%);
            background: none;
            border: none;
            color: var(--muted);
            cursor: pointer;
            padding: 4px;
        }
        .toggle-pass:hover { color: var(--text); }

        .btn {
            width: 100%;
            padding: 0.95rem 1rem;
            border: none;
            border-radius: 12px;
            font-family: "Rajdhani", sans-serif;
            font-size: 1.15rem;
            font-weight: 700;
            letter-spacing: 1.5px;
            cursor: pointer;
            position: relative;
            overflow: hidden;
            transition: transform 0.15s, box-shadow 0.2s;
        }
        .btn-primary {
            background: linear-gradient(135deg, var(--primary), #a0001e);
            color: #fff;
            box-shadow: 0 10px 24px rgba(225,29,51,0.4);
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 14px 30px rgba(225,29,51,0.55);
        }
        .btn-primary:active { transform: translateY(0); }
        .btn-primary::after {
            content: "";
            position: absolute;
            top: 0;
            left: -100%;
            width: 60%;
            height: 100%;
            background: linear-gradient(90deg, transparent, rgba(255,255,255,0.35), transparent);
            animation: btnShine 2.6s ease-in-out infinite;
        }
        @keyframes btnShine {
            0% { left: -100%; }
            50%, 100% { left: 120%; }
        }
        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }

        .alert {
            display: none;
            padding: 0.75rem 0.9rem;
            border-radius: 10px;
            font-size: 0.88rem;
            margin-bottom: 1rem;
            border: 1px solid transparent;
            animation: shake 0.4s ease;
        }
        @keyframes shake {
            0%, 100% { transform: translateX(0); }
            25% { transform: translateX(-6px); }
            75% { transform: translateX(6px); }
        }
        .alert-danger {
            background: rgba(255, 77, 109, 0.12);
            border-color: rgba(255, 77, 109, 0.4);
            color: #ffb3c1;
        }
        .alert-success {
            background: rgba(46, 204, 113, 0.12);
            border-color: rgba(46, 204, 113, 0.4);
            color: #b7f0ce;
        }

        .card-footer {
            text-align: center;
            margin-top: 1.3rem;
            font-size: 0.88rem;
            color: var(--muted);
        }
        .card-footer a {
            color: var(--primary-2);
            text-decoration: none;
            font-weight: 600;
        }
        .card-footer a:hover { text-decoration: underline; }

        .divider {
            display: flex;
            align-items: center;
            gap: 0.8rem;
            margin: 1.2rem 0;
            color: #4b5563;
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .divider::before, .divider::after {
            content: "";
            flex: 1;
            height: 1px;
            background: rgba(255,255,255,0.08);
        }

        .status-row {
            display: flex;
            justify-content: space-between;
            margin-top: 0.6rem;
            font-size: 0.78rem;
            color: var(--muted);
        }
        .status-dot {
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: var(--success);
            margin-right: 6px;
            box-shadow: 0 0 8px var(--success);
            animation: blink 1.6s ease-in-out infinite;
        }
        @keyframes blink {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.4; }
        }

        /* loading spinner */
        .spinner {
            display: none;
            width: 20px;
            height: 20px;
            border: 2px solid rgba(255,255,255,0.3);
            border-top-color: #fff;
            border-radius: 50%;
            animation: spin 0.7s linear infinite;
            margin: 0 auto;
        }

        /* reduce motion */
        @media (prefers-reduced-motion: reduce) {
            * { animation: none !important; transition: none !important; }
        }
    </style>
</head>
<body>
    <div class="bg">
        <div class="bg-grid"></div>
        <div class="speed-lines" id="speedLines"></div>
        <div class="particles" id="particles"></div>
    </div>

    <div class="wrap">
        <!-- Showcase -->
        <section class="showcase">
            <div class="brand-row">
                <div class="logo-badge" aria-hidden="true">
                    <svg viewBox="0 0 100 48">
                        <path d="M6 32c0-5 4-8 9-9l14-3c6-8 14-12 24-12h6c10 0 18 4 23 12l10 3c6 1 9 4 9 9v5c0 2-1 3-3 3h-8a10 10 0 0 1-20 0H33a10 10 0 0 1-20 0h-4c-2 0-3-1-3-3z" fill="currentColor"/>
                        <circle cx="23" cy="40" r="6.5" fill="#fff"/>
                        <circle cx="23" cy="40" r="2.5" fill="#7a0016"/>
                        <circle cx="72" cy="40" r="6.5" fill="#fff"/>
                        <circle cx="72" cy="40" r="2.5" fill="#7a0016"/>
                    </svg>
                </div>
                <div class="brand-text">
                    <h2>Garuda <em>Auto Service</em></h2>
                    <p>Bengkel Mobil &middot; Servis Cepat &amp; Terpercaya</p>
                </div>
            </div>

            <h1>Pit Stop Digital untuk <span>Bengkel Anda</span></h1>
            <p class="lead">Kelola konsumen, kendaraan, dan seluruh riwayat servis dalam satu dasbor bengkel yang cepat, rapi, dan modern.</p>

            <div class="chips">
                <span class="chip">Servis Berkala</span>
                <span class="chip">Tune Up</span>
                <span class="chip">Kaki-Kaki</span>
                <span class="chip">Body Repair</span>
                <span class="chip">Ganti Oli</span>
            </div>

            <div class="stage" aria-hidden="true">
                <div class="road"></div>
                <svg class="stage-car" viewBox="0 0 520 200">
                    <defs>
                        <linearGradient id="carBody" x1="0" y1="0" x2="0.2" y2="1">
                            <stop offset="0" stop-color="#ff6b78"/>
                            <stop offset="0.45" stop-color="#e11d33"/>
                            <stop offset="1" stop-color="#7a0016"/>
                        </linearGradient>
                        <linearGradient id="carGlass" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="0" stop-color="#a8bccf" stop-opacity="0.9"/>
                            <stop offset="1" stop-color="#11161d"/>
                        </linearGradient>
                        <linearGradient id="carRim" x1="0" y1="0" x2="1" y2="1">
                            <stop offset="0" stop-color="#f2f5f8"/>
                            <stop offset="1" stop-color="#8b949e"/>
                        </linearGradient>
                    </defs>
                    <ellipse cx="272" cy="176" rx="202" ry="12" fill="#1a0d10" opacity="0.3"/>
                    <path d="M50 140 C48 126 58 116 80 112 L136 104 C154 80 184 64 222 61 L306 60 C348 61 380 74 404 96 L434 104 L474 109 C497 113 508 124 509 138 C510 148 504 152 494 152 L92 152 C64 152 50 148 50 140 Z" fill="url(#carBody)"/>
                    <path d="M154 102 C170 83 194 72 224 69 L270 68 L270 102 Z" fill="url(#carGlass)"/>
                    <path d="M282 68 L300 68 C330 69 356 80 376 98 L376 102 L282 102 Z" fill="url(#carGlass)"/>
                    <path d="M470 116 L496 122 L497 131 L474 128 Z" fill="#fff3c4" class="headlight"/>
                    <circle cx="508" cy="128" r="18" fill="#ffd166" opacity="0.25" class="headlight"/>
                    <path d="M64 124 L88 124 L86 136 L66 134 Z" fill="#ff2d44"/>
                    <path d="M96 152 A38 38 0 0 1 172 152 Z" fill="#12070a"/>
                    <path d="M306 152 A38 38 0 0 1 382 152 Z" fill="#12070a"/>
                    <g class="wheel">
                        <circle cx="134" cy="142" r="30" fill="#0b0608"/>
                        <circle cx="134" cy="142" r="17" fill="url(#carRim)"/>
                        <g stroke="#3d444b" stroke-width="3" stroke-linecap="round">
                            <line x1="134" y1="129" x2="134" y2="155"/>
                            <line x1="121" y1="142" x2="147" y2="142"/>
                            <line x1="125" y1="133" x2="143" y2="151"/>
                            <line x1="143" y1="133" x2="125" y2="151"/>
                        </g>
                        <circle cx="134" cy="142" r="5" fill="#c1121f"/>
                    </g>
                    <g class="wheel">
                        <circle cx="344" cy="142" r="30" fill="#0b0608"/>
                        <circle cx="344" cy="142" r="17" fill="url(#carRim)"/>
                        <g stroke="#3d444b" stroke-width="3" stroke-linecap="round">
                            <line x1="344" y1="129" x2="344" y2="155"/>
                            <line x1="331" y1="142" x2="357" y2="142"/>
                            <line x1="335" y1="133" x2="353" y2="151"/>
                            <line x1="353" y1="133" x2="335" y2="151"/>
                        </g>
                        <circle cx="344" cy="142" r="5" fill="#c1121f"/>
                    </g>
                    <circle cx="250" cy="124" r="13" fill="#fff" opacity="0.92"/>
                    <text x="250" y="129" text-anchor="middle" font-size="13" font-weight="700" fill="#9d0208" font-family="Arial, sans-serif">01</text>
                </svg>
            </div>
        </section>

        <!-- Card -->
        <section class="card">
            <div class="card-header">
                <div class="mini-logo" aria-hidden="true">
                    <svg viewBox="0 0 100 48">
                        <path d="M6 32c0-5 4-8 9-9l14-3c6-8 14-12 24-12h6c10 0 18 4 23 12l10 3c6 1 9 4 9 9v5c0 2-1 3-3 3h-8a10 10 0 0 1-20 0H33a10 10 0 0 1-20 0h-4c-2 0-3-1-3-3z" fill="currentColor"/>
                        <circle cx="23" cy="40" r="6.5" fill="#fff"/>
                        <circle cx="23" cy="40" r="2.5" fill="#7a0016"/>
                        <circle cx="72" cy="40" r="6.5" fill="#fff"/>
                        <circle cx="72" cy="40" r="2.5" fill="#7a0016"/>
                    </svg>
                </div>
                <h3>Masuk ke Sistem</h3>
                <p>Silakan login untuk mengakses data bengkel</p>
            </div>

            <div id="alertBox" class="alert"></div>

            <form id="loginForm" autocomplete="on">
                <div class="form-group">
                    <label for="username">Username</label>
                    <div class="input-wrap">
                        <input type="text" id="username" class="form-control" placeholder="Masukkan username" required autofocus>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                    </div>
                </div>
                <div class="form-group">
                    <label for="password">Password</label>
                    <div class="input-wrap">
                        <input type="password" id="password" class="form-control" placeholder="Masukkan password" required>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                        <button type="button" class="toggle-pass" id="togglePass" aria-label="Tampilkan password">
                            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                        </button>
                    </div>
                </div>

                <button type="submit" class="btn btn-primary btn-block" id="btnLogin">
                    <span id="btnText">MASUK KE SISTEM</span>
                    <span class="spinner" id="btnSpinner"></span>
                </button>

                <div class="status-row">
                    <span><span class="status-dot"></span>Sistem Online</span>
                    <span id="clock"></span>
                </div>
            </form>

            <div class="divider">atau</div>

            <div class="card-footer">
                Belum punya akun? <a href="/register">Daftar Staf Baru</a>
            </div>
        </section>
    </div>

    <script>
        // Redirect if already logged in
        (function redirectIfLoggedIn() {
            try {
                if (localStorage.getItem('gas_token')) {
                    window.location.href = '/';
                }
            } catch (e) {}
        })();

        // Clock
        function updateClock() {
            var now = new Date();
            var h = String(now.getHours()).padStart(2, '0');
            var m = String(now.getMinutes()).padStart(2, '0');
            var s = String(now.getSeconds()).padStart(2, '0');
            var el = document.getElementById('clock');
            if (el) el.textContent = h + ':' + m + ':' + s;
        }
        setInterval(updateClock, 1000);
        updateClock();

        // Particles & speed lines (auto-run)
        (function spawnDecor() {
            var sl = document.getElementById('speedLines');
            for (var i = 0; i < 8; i++) {
                var span = document.createElement('span');
                span.style.top = (10 + Math.random() * 70) + '%';
                span.style.animationDelay = (Math.random() * 1.8) + 's';
                span.style.animationDuration = (1.4 + Math.random() * 1.2) + 's';
                sl.appendChild(span);
            }
            var pt = document.getElementById('particles');
            for (var j = 0; j < 18; j++) {
                var p = document.createElement('i');
                p.style.left = (Math.random() * 100) + '%';
                p.style.bottom = (Math.random() * 40) + '%';
                p.style.animationDelay = (Math.random() * 6) + 's';
                p.style.animationDuration = (4 + Math.random() * 4) + 's';
                p.style.width = (3 + Math.random() * 4) + 'px';
                p.style.height = p.style.width;
                pt.appendChild(p);
            }
        })();

        // Toggle password
        (function () {
            var btn = document.getElementById('togglePass');
            var input = document.getElementById('password');
            if (btn && input) {
                btn.addEventListener('click', function () {
                    if (input.type === 'password') {
                        input.type = 'text';
                    } else {
                        input.type = 'password';
                    }
                });
            }
        })();

        // Login submit
        (function () {
            var form = document.getElementById('loginForm');
            var alertBox = document.getElementById('alertBox');
            var btn = document.getElementById('btnLogin');
            var btnText = document.getElementById('btnText');
            var spinner = document.getElementById('btnSpinner');

            form.addEventListener('submit', async function (e) {
                e.preventDefault();
                alertBox.style.display = 'none';

                var username = document.getElementById('username').value.trim();
                var password = document.getElementById('password').value;

                btn.disabled = true;
                btnText.style.display = 'none';
                spinner.style.display = 'block';

                try {
                    var res = await fetch('/api/auth/login', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ username: username, password: password })
                    });

                    var data = await res.json();

                    if (!res.ok) {
                        alertBox.className = 'alert alert-danger';
                        alertBox.textContent = data.error || 'Login gagal, periksa kembali kredensial Anda.';
                        alertBox.style.display = 'block';
                        return;
                    }

                    localStorage.setItem('gas_token', data.token);
                    localStorage.setItem('gas_user', data.username);

                    alertBox.className = 'alert alert-success';
                    alertBox.textContent = 'Login berhasil. Mengarahkan...';
                    alertBox.style.display = 'block';

                    setTimeout(function () {
                        window.location.href = '/';
                    }, 600);
                } catch (err) {
                    alertBox.className = 'alert alert-danger';
                    alertBox.textContent = 'Terjadi kesalahan jaringan atau server.';
                    alertBox.style.display = 'block';
                } finally {
                    btn.disabled = false;
                    btnText.style.display = '';
                    spinner.style.display = 'none';
                }
            });
        })();
    </script>
</body>
</html>
`

const registerHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#7a0016">
    <title>Registrasi Staf — Garuda Auto Service</title>
    <link rel="icon" type="image/svg+xml" href="/favicon.svg?v=3">
    <link rel="stylesheet" href="/css/style.css?v=3">
</head>
<body class="auth-wrapper">
    <!-- Dekorasi jalan & mobil -->
    <div class="auth-decor" aria-hidden="true">
        <div class="auth-road"></div>
        <svg class="auth-ghost-car" viewBox="0 0 520 200">
            <path d="M50 140 C48 126 58 116 80 112 L136 104 C154 80 184 64 222 61 L306 60 C348 61 380 74 404 96 L434 104 L474 109 C497 113 508 124 509 138 C510 148 504 152 494 152 L92 152 C64 152 50 148 50 140 Z" fill="currentColor"/>
            <circle cx="134" cy="142" r="30" fill="currentColor"/>
            <circle cx="344" cy="142" r="30" fill="currentColor"/>
        </svg>
    </div>

    <div class="auth-card">
        <div class="auth-header">
            <span class="auth-logo" aria-hidden="true">
                <svg viewBox="0 0 100 48">
                    <path d="M6 32c0-5 4-8 9-9l14-3c6-8 14-12 24-12h6c10 0 18 4 23 12l10 3c6 1 9 4 9 9v5c0 2-1 3-3 3h-8a10 10 0 0 1-20 0H33a10 10 0 0 1-20 0h-4c-2 0-3-1-3-3z" fill="currentColor"/>
                    <circle cx="23" cy="40" r="6.5" fill="#fff"/>
                    <circle cx="23" cy="40" r="2.5" fill="#7a0016"/>
                    <circle cx="72" cy="40" r="6.5" fill="#fff"/>
                    <circle cx="72" cy="40" r="2.5" fill="#7a0016"/>
                </svg>
            </span>
            <span class="auth-eyebrow">Pit Stop Staf Baru</span>
            <h1>Garuda Auto Service</h1>
            <p>Pendaftaran akun staf bengkel baru</p>
        </div>

        <div id="alertBox" class="alert"></div>

        <form id="registerForm">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" class="form-control" placeholder="Minimal 3 karakter" required minlength="3" autofocus>
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" class="form-control" placeholder="Minimal 6 karakter" required minlength="6">
            </div>
            <div class="form-group">
                <label for="confirmPassword">Konfirmasi Password</label>
                <input type="password" id="confirmPassword" class="form-control" placeholder="Ulangi password" required minlength="6">
            </div>
            <button type="submit" class="btn btn-primary btn-block">Daftar Sekarang</button>
        </form>

        <div class="auth-footer">
            Sudah punya akun? <a href="/login">Kembali ke Login</a>
        </div>
    </div>

    <script src="/js/auth.js?v=3"></script>
    <script>
        redirectIfLoggedIn();

        const form = document.getElementById('registerForm');
        const alertBox = document.getElementById('alertBox');

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            alertBox.style.display = 'none';

            const username = document.getElementById('username').value.trim();
            const password = document.getElementById('password').value;
            const confirmPassword = document.getElementById('confirmPassword').value;

            if (password !== confirmPassword) {
                alertBox.className = 'alert alert-danger';
                alertBox.textContent = 'Password dan Konfirmasi Password tidak sama.';
                alertBox.style.display = 'block';
                return;
            }

            try {
                const res = await fetch('/api/auth/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                const data = await res.json();

                if (!res.ok) {
                    alertBox.className = 'alert alert-danger';
                    alertBox.textContent = data.error || 'Pendaftaran gagal.';
                    alertBox.style.display = 'block';
                    return;
                }

                alertBox.className = 'alert alert-success';
                alertBox.textContent = data.message + ' Mengarahkan ke halaman login...';
                alertBox.style.display = 'block';

                setTimeout(() => {
                    window.location.href = '/login';
                }, 1500);
            } catch (err) {
                alertBox.className = 'alert alert-danger';
                alertBox.textContent = 'Terjadi kesalahan jaringan atau server.';
                alertBox.style.display = 'block';
            }
        });
    </script>
</body>
</html>
`

const styleCSS = `/* =========================================================
   GARUDA AUTO SERVICE — Tema Merah Bermotif Mobil
   Palet: merah balap, karbon gelap, aksen kuning speedometer
   ========================================================= */

:root {
    --red-950: #3d0009;
    --red-900: #5c0010;
    --red-800: #7a0016;
    --red-700: #9d0208;
    --red-600: #c1121f;
    --red-500: #e11d33;
    --red-400: #ff4757;
    --red-200: #ffccd2;
    --red-100: #ffe3e6;
    --red-50: #fff1f2;

    --ink: #1a0d10;
    --carbon: #231317;
    --smoke: #4d3a3e;
    --muted: #7d686c;
    --paper: #fffdfc;
    --cream: #f9efec;
    --amber: #ffd166;
    --line: #f0dcd8;

    --shadow-red: 0 18px 40px -18px rgba(157, 2, 8, 0.55);
    --shadow-soft: 0 10px 30px -18px rgba(26, 13, 16, 0.45);

    --font-display: "Bahnschrift", "Arial Narrow", "Roboto Condensed", "Segoe UI", system-ui, sans-serif;
    --font-body: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}

* {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
    font-family: var(--font-body);
}

html {
    scroll-behavior: smooth;
}

body {
    background-color: var(--cream);
    color: var(--ink);
    line-height: 1.5;
    min-height: 100vh;
    background-image:
        radial-gradient(1100px 500px at 88% -10%, rgba(225, 29, 51, 0.16), transparent 62%),
        radial-gradient(900px 420px at 0% 8%, rgba(157, 2, 8, 0.08), transparent 58%),
        repeating-linear-gradient(115deg, rgba(26, 13, 16, 0.035) 0 2px, transparent 2px 16px);
    background-attachment: fixed;
}

::selection {
    background: var(--red-500);
    color: #fff;
}

::-webkit-scrollbar {
    width: 12px;
    height: 12px;
}

::-webkit-scrollbar-track {
    background: #f3e5e2;
}

::-webkit-scrollbar-thumb {
    background: linear-gradient(180deg, var(--red-600), var(--red-800));
    border: 3px solid #f3e5e2;
    border-radius: 99px;
}

/* Container Umum */
.container {
    max-width: 1040px;
    margin: 0 auto;
    padding: 1.75rem 1.5rem 3rem;
}

/* =========================================================
   Navigasi Atas — carbon + garis balap
   ========================================================= */
.navbar {
    position: sticky;
    top: 0;
    z-index: 100;
    background:
        linear-gradient(90deg, rgba(225, 29, 51, 0.22), transparent 38%),
        linear-gradient(135deg, #2a1418 0%, #170b0e 55%, #2a0d13 100%);
    color: #fff;
    padding: 1rem clamp(1rem, 4vw, 2.5rem);
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    border-bottom: 3px solid var(--red-600);
    box-shadow: 0 12px 30px -14px rgba(90, 0, 15, 0.75);
}

/* Strip bendera start/finish di tepi atas */
.navbar::before {
    content: "";
    position: absolute;
    inset: 0 0 auto 0;
    height: 6px;
    background: conic-gradient(#f8f5f4 25%, #170b0e 0 50%, #f8f5f4 0 75%, #170b0e 0) 0 0 / 12px 12px;
    opacity: 0.9;
    pointer-events: none;
}

/* Garis glow merah-kuning di bawah navbar */
.navbar::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    bottom: -3px;
    height: 3px;
    background: linear-gradient(90deg, transparent, var(--red-500), var(--amber), var(--red-500), transparent);
    background-size: 200% 100%;
    animation: glowShift 5s linear infinite;
    pointer-events: none;
}

.navbar > * {
    position: relative;
    z-index: 1;
}

.nav-brand {
    display: flex;
    align-items: center;
    gap: 0.85rem;
    min-width: 0;
}

.brand-badge {
    flex: none;
    width: 47px;
    height: 47px;
    border-radius: 13px;
    display: grid;
    place-items: center;
    background: linear-gradient(145deg, var(--red-400), var(--red-800));
    box-shadow:
        0 10px 22px -10px rgba(225, 29, 51, 0.9),
        inset 0 1px 0 rgba(255, 255, 255, 0.45);
}

.brand-badge svg {
    width: 33px;
    height: auto;
    color: #fff;
    filter: drop-shadow(0 1px 1px rgba(0, 0, 0, 0.35));
}

.brand-stack {
    display: flex;
    flex-direction: column;
    line-height: 1.1;
    min-width: 0;
}

.brand-name {
    font-family: var(--font-display);
    font-size: 1.3rem;
    font-weight: 800;
    letter-spacing: 0.045em;
    text-transform: uppercase;
    color: #fff;
    white-space: nowrap;
}

.brand-name em {
    font-style: normal;
    background: linear-gradient(90deg, var(--amber), #ff8a95);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
}

.brand-tagline {
    font-size: 0.64rem;
    letter-spacing: 0.19em;
    text-transform: uppercase;
    color: rgba(255, 255, 255, 0.55);
    white-space: nowrap;
    margin-top: 0.2rem;
}

.nav-user {
    display: flex;
    align-items: center;
    gap: 0.85rem;
    font-size: 0.9rem;
    color: rgba(255, 255, 255, 0.82);
}

.nav-greeting strong {
    color: var(--amber);
}

.navbar .btn-secondary {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.3);
    color: #fff;
}

.navbar .btn-secondary:hover {
    background: rgba(255, 255, 255, 0.16);
    border-color: #fff;
}

/* =========================================================
   Hero — panggung mobil merah
   ========================================================= */
.hero {
    position: relative;
    overflow: hidden;
    margin-bottom: 1.6rem;
    padding: clamp(1.5rem, 2.6vw, 2.2rem);
    display: grid;
    grid-template-columns: minmax(0, 1.05fr) minmax(270px, 0.95fr);
    align-items: center;
    gap: 1rem;
    border-radius: 22px;
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.14);
    background:
        radial-gradient(700px 260px at 84% -10%, rgba(255, 209, 102, 0.22), transparent 62%),
        radial-gradient(520px 320px at 6% 120%, rgba(255, 71, 87, 0.4), transparent 62%),
        linear-gradient(118deg, #3d0009 0%, #7a0016 46%, #b3122a 100%);
    box-shadow: 0 28px 60px -30px rgba(122, 0, 22, 0.85);
}

/* Motif bendera kotak-kotak */
.hero::before {
    content: "";
    position: absolute;
    top: -46px;
    right: -34px;
    width: 230px;
    height: 230px;
    background: conic-gradient(#fff 25%, transparent 0 50%, #fff 0 75%, transparent 0) 0 0 / 26px 26px;
    opacity: 0.08;
    transform: rotate(15deg);
    pointer-events: none;
}

/* Marka jalan animasi di tepi bawah hero */
.hero::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 6px;
    background: repeating-linear-gradient(90deg, var(--amber) 0 40px, transparent 40px 96px);
    opacity: 0.65;
    animation: roadMove 1.4s linear infinite;
    pointer-events: none;
}

.hero-copy {
    position: relative;
    z-index: 2;
}

.hero-kicker {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.32rem 0.75rem 0.32rem 0.55rem;
    border-radius: 999px;
    background: rgba(0, 0, 0, 0.26);
    border: 1px solid rgba(255, 255, 255, 0.2);
    font-size: 0.68rem;
    font-weight: 700;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: #ffd9de;
}

.kicker-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--amber);
    box-shadow: 0 0 0 3px rgba(255, 209, 102, 0.25);
    animation: pulse 2s ease-in-out infinite;
}

.hero-title {
    font-family: var(--font-display);
    font-size: clamp(2rem, 4.6vw, 3.05rem);
    font-weight: 800;
    line-height: 0.96;
    letter-spacing: -0.015em;
    text-transform: uppercase;
    margin: 0.6rem 0 0.55rem;
    text-shadow: 0 6px 20px rgba(0, 0, 0, 0.35);
}

.hero-title span {
    display: inline-block;
    background: linear-gradient(92deg, var(--amber), #ffb3bb 55%, #ff8a95);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
}

.hero-text {
    font-size: 0.92rem;
    color: rgba(255, 255, 255, 0.78);
    max-width: 40ch;
}

.hero-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-top: 0.9rem;
}

.chip {
    font-size: 0.72rem;
    letter-spacing: 0.03em;
    padding: 0.24rem 0.62rem;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.09);
    border: 1px solid rgba(255, 255, 255, 0.22);
    color: rgba(255, 255, 255, 0.88);
}

.hero-metrics {
    display: flex;
    flex-wrap: wrap;
    gap: 0.65rem;
    margin-top: 1.1rem;
}

.metric {
    min-width: 96px;
    padding: 0.5rem 0.9rem;
    border-radius: 12px;
    background: rgba(0, 0, 0, 0.24);
    border: 1px solid rgba(255, 255, 255, 0.16);
    backdrop-filter: blur(4px);
}

.metric-value {
    display: block;
    font-family: var(--font-display);
    font-size: 1.45rem;
    font-weight: 800;
    line-height: 1;
    color: var(--amber);
    font-variant-numeric: tabular-nums;
}

.metric-label {
    display: block;
    margin-top: 0.2rem;
    font-size: 0.65rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: rgba(255, 255, 255, 0.7);
}

/* Visual mobil */
.hero-visual {
    position: relative;
    z-index: 1;
    min-height: 210px;
    display: flex;
    align-items: flex-end;
    justify-content: center;
}

.hero-car {
    position: relative;
    z-index: 2;
    width: 100%;
    max-width: 430px;
    height: auto;
    filter: drop-shadow(0 26px 20px rgba(0, 0, 0, 0.4));
    animation: carIdle 3.8s ease-in-out infinite;
}

.hero-car .wheel-spokes {
    transform-box: fill-box;
    transform-origin: center;
    animation: spin 2.4s linear infinite;
}

.speed-lines {
    position: absolute;
    left: 0;
    top: 37%;
    width: 54%;
    height: 60px;
    z-index: 1;
    pointer-events: none;
}

.speed-lines i {
    position: absolute;
    left: 0;
    height: 5px;
    border-radius: 99px;
    background: linear-gradient(90deg, rgba(255, 255, 255, 0), rgba(255, 255, 255, 0.85));
    animation: dash 1.9s linear infinite;
}

.speed-lines i:nth-child(1) {
    top: 0;
    width: 68%;
}

.speed-lines i:nth-child(2) {
    top: 18px;
    width: 44%;
    animation-delay: 0.35s;
}

.speed-lines i:nth-child(3) {
    top: 36px;
    width: 82%;
    animation-delay: 0.7s;
}

.hero-roundel {
    position: absolute;
    top: -6px;
    right: 2px;
    z-index: 3;
    width: 82px;
    height: 82px;
    padding: 6px;
    border-radius: 50%;
    background: repeating-conic-gradient(#170b0e 0 16deg, #fff 16deg 32deg);
    box-shadow: 0 16px 26px -12px rgba(0, 0, 0, 0.7);
    transform: rotate(9deg);
}

.roundel-inner {
    display: grid;
    place-content: center;
    width: 100%;
    height: 100%;
    border-radius: 50%;
    text-align: center;
    color: var(--red-800);
    background: radial-gradient(circle at 35% 30%, #fff, #ffe0e4);
}

.roundel-inner b {
    font-family: var(--font-display);
    font-size: 1.5rem;
    line-height: 1;
}

.roundel-inner small {
    margin-top: 2px;
    font-size: 0.46rem;
    font-weight: 800;
    letter-spacing: 0.18em;
}

/* =========================================================
   Statistik Bulanan — diagram konsumen
   ========================================================= */
.monthly-panel {
    position: relative;
    overflow: hidden;
    margin-bottom: 1.5rem;
    padding: 1.25rem 1.35rem 1rem;
    background: linear-gradient(180deg, var(--paper), #fff8f6);
    border: 1px solid var(--line);
    border-radius: 18px;
    box-shadow: var(--shadow-soft);
}

/* Aksen bendera kotak-kotak di pojok panel */
.monthly-panel::before {
    content: "";
    position: absolute;
    top: 0;
    right: 0;
    width: 130px;
    height: 52px;
    background: conic-gradient(var(--red-700) 25%, transparent 0 50%, var(--red-700) 0 75%, transparent 0) 0 0 / 16px 16px;
    opacity: 0.07;
    transform: skewX(-18deg) translateX(20px);
    pointer-events: none;
}

.monthly-head {
    position: relative;
    z-index: 1;
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
}

.panel-kicker {
    display: inline-block;
    font-size: 0.64rem;
    font-weight: 800;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: var(--red-600);
}

.monthly-head h2 {
    margin: 0.2rem 0 0;
    font-family: var(--font-display);
    font-size: 1.3rem;
    letter-spacing: 0.015em;
    text-transform: uppercase;
    color: var(--ink);
}

.monthly-subtitle {
    margin-top: 0.25rem;
    font-size: 0.82rem;
    color: var(--muted);
}

.monthly-summary {
    display: flex;
    flex-wrap: wrap;
    gap: 0.55rem;
}

.summary-pill {
    min-width: 118px;
    padding: 0.55rem 0.8rem;
    border-radius: 12px;
    background: linear-gradient(145deg, #2a1418, #120709);
    border: 1px solid rgba(225, 29, 51, 0.45);
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.08);
}

.summary-pill span {
    display: block;
    font-family: var(--font-display);
    font-size: 1.5rem;
    font-weight: 800;
    line-height: 1;
    color: var(--amber);
    font-variant-numeric: tabular-nums;
}

.summary-pill small {
    display: block;
    margin-top: 0.25rem;
    font-size: 0.6rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: rgba(255, 255, 255, 0.68);
}

.chart-legend {
    position: relative;
    z-index: 1;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-top: 0.9rem;
    padding-top: 0.7rem;
    border-top: 1px dashed #e8cfcb;
    font-size: 0.76rem;
    color: var(--smoke);
}

.legend-item {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
}

.legend-dot {
    width: 10px;
    height: 10px;
    border-radius: 3px;
}

.dot-customer {
    background: linear-gradient(180deg, var(--red-400), var(--red-700));
}

.dot-service {
    background: linear-gradient(180deg, #5b4449, #1a0d10);
}

.chart-scroll {
    overflow-x: auto;
    padding: 0.2rem 0 0.35rem;
}

.monthly-chart {
    display: grid;
    grid-auto-flow: column;
    grid-auto-columns: minmax(0, 1fr);
    gap: 0.25rem;
    align-items: end;
    height: 190px;
    min-width: 0;
    padding-top: 1.25rem;
    margin-top: 0.35rem;
}

.chart-col {
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    height: 100%;
    min-width: 0;
}

.chart-bars {
    flex: 1;
    display: flex;
    align-items: flex-end;
    justify-content: center;
    gap: 3px;
    padding: 0 2px;
    background: repeating-linear-gradient(to top, rgba(157, 2, 8, 0.07) 0 1px, transparent 1px 25%);
    border-bottom: 2px solid #e8cfcb;
}

.chart-bar {
    position: relative;
    flex: 1;
    max-width: 22px;
    min-height: 0;
    border-radius: 4px 4px 1px 1px;
    transition: height 0.45s ease;
}

.chart-bar span {
    position: absolute;
    top: -1.1rem;
    left: 50%;
    transform: translateX(-50%);
    font-size: 0.64rem;
    font-weight: 700;
    line-height: 1;
    color: var(--smoke);
    font-variant-numeric: tabular-nums;
}

.bar-customer {
    background: linear-gradient(180deg, var(--red-400), var(--red-700));
    box-shadow: 0 6px 12px -8px rgba(157, 2, 8, 0.9);
}

.bar-service {
    background: linear-gradient(180deg, #5b4449, #1a0d10);
}

.chart-month {
    display: block;
    padding-top: 0.35rem;
    text-align: center;
    font-size: 0.66rem;
    font-weight: 700;
    letter-spacing: 0.02em;
    color: var(--muted);
    white-space: nowrap;
}

.chart-col.is-current .chart-month {
    color: var(--red-700);
}

.chart-col.is-current .chart-bars {
    background-color: rgba(255, 209, 102, 0.12);
}

.chart-empty-note {
    grid-column: 1 / -1;
    align-self: center;
    text-align: center;
    padding: 1.5rem 1rem;
    color: var(--muted);
    font-size: 0.86rem;
}

/* =========================================================
   Toolbar Pencarian & Filter
   ========================================================= */
.toolbar {
    position: relative;
    overflow: hidden;
    background: linear-gradient(180deg, var(--paper), #fff7f5);
    padding: 1rem 1.1rem;
    border: 1px solid var(--line);
    border-left: 6px solid var(--red-600);
    border-radius: 14px;
    margin-bottom: 1.5rem;
    box-shadow: var(--shadow-soft);
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    align-items: center;
    justify-content: space-between;
}

.toolbar::after {
    content: "";
    position: absolute;
    right: -40px;
    top: -60px;
    width: 150px;
    height: 150px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(225, 29, 51, 0.1), transparent 70%);
    pointer-events: none;
}

.toolbar-controls {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    flex: 1;
    position: relative;
    z-index: 1;
}

.toolbar-search {
    flex: 2;
    min-width: 200px;
}

.toolbar-filter {
    flex: 1;
    min-width: 140px;
}

.toolbar .btn {
    position: relative;
    z-index: 1;
}

/* =========================================================
   Kartu Konsumen
   ========================================================= */
.card-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.customer-card {
    position: relative;
    overflow: hidden;
    background: var(--paper);
    padding: 1.2rem 1.35rem;
    border-radius: 16px;
    border: 1px solid var(--line);
    border-left: 6px solid var(--red-600);
    box-shadow: 0 8px 24px -18px rgba(26, 13, 16, 0.55);
    display: flex;
    flex-direction: column;
    gap: 0.8rem;
    transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
}

/* Aksen bendera kotak-kotak di pojok kartu */
.customer-card::before {
    content: "";
    position: absolute;
    top: 0;
    right: 0;
    width: 86px;
    height: 44px;
    background: conic-gradient(var(--red-700) 25%, transparent 0 50%, var(--red-700) 0 75%, transparent 0) 0 0 / 15px 15px;
    opacity: 0.07;
    transform: skewX(-18deg) translateX(16px);
    pointer-events: none;
}

.customer-card:hover {
    transform: translateY(-3px);
    border-left-color: var(--red-500);
    box-shadow: 0 22px 38px -24px rgba(157, 2, 8, 0.6);
}

.card-header-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
    gap: 0.5rem;
    position: relative;
    z-index: 1;
}

.plate-badge {
    display: inline-block;
    background: linear-gradient(180deg, #241014, #120709);
    color: #fff;
    font-family: "Consolas", "Courier New", monospace;
    font-size: 1.02rem;
    font-weight: 700;
    padding: 0.28rem 0.7rem;
    border-radius: 6px;
    border: 1px solid rgba(255, 255, 255, 0.16);
    letter-spacing: 0.12em;
    box-shadow: inset 0 0 0 1px rgba(225, 29, 51, 0.35);
}

.car-title {
    font-size: 1.08rem;
    font-weight: 800;
    color: var(--ink);
}

.customer-meta {
    position: relative;
    z-index: 1;
    font-size: 0.88rem;
    color: var(--muted);
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 0.4rem;
}

.customer-meta strong {
    color: var(--smoke);
}

.service-summary {
    position: relative;
    z-index: 1;
    background: var(--red-50);
    padding: 0.75rem 0.85rem;
    border-radius: 10px;
    font-size: 0.85rem;
    border: 1px dashed var(--red-200);
    color: var(--smoke);
}

.service-summary strong {
    color: var(--red-700);
}

.card-actions {
    position: relative;
    z-index: 1;
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    justify-content: flex-end;
    margin-top: 0.35rem;
    padding-top: 0.8rem;
    border-top: 1px dashed #e8cfcb;
}

/* =========================================================
   Empty State — roda berputar
   ========================================================= */
.empty-state {
    text-align: center;
    padding: 3rem 1rem;
    color: var(--muted);
    background: var(--paper);
    border: 1px dashed var(--red-200);
    border-radius: 16px;
}

.empty-state::before {
    content: "";
    display: block;
    width: 58px;
    height: 58px;
    margin: 0 auto 1.1rem;
    border-radius: 50%;
    border: 9px solid #f1d9d6;
    border-top-color: var(--red-500);
    border-right-color: var(--red-700);
    animation: spin 2.6s linear infinite;
}

.empty-state h3 {
    color: var(--ink);
    font-family: var(--font-display);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.35rem;
}

/* =========================================================
   Form & Tombol
   ========================================================= */
.form-group {
    margin-bottom: 1rem;
}

.form-group label {
    display: block;
    font-size: 0.86rem;
    font-weight: 700;
    margin-bottom: 0.35rem;
    color: var(--smoke);
}

.form-control {
    width: 100%;
    padding: 0.65rem 0.8rem;
    background: #fff;
    border: 1px solid #e6d2ce;
    border-radius: 8px;
    font-size: 0.95rem;
    color: var(--ink);
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
}

.form-control::placeholder {
    color: #b09a9d;
}

.form-control:focus {
    border-color: var(--red-500);
    box-shadow: 0 0 0 3px rgba(225, 29, 51, 0.14);
}

textarea.form-control {
    resize: vertical;
}

.btn {
    position: relative;
    overflow: hidden;
    display: inline-block;
    padding: 0.65rem 1.2rem;
    border: 1px solid transparent;
    border-radius: 9px;
    font-size: 0.92rem;
    font-weight: 700;
    letter-spacing: 0.01em;
    cursor: pointer;
    text-decoration: none;
    text-align: center;
    transition: background 0.2s, border-color 0.2s, transform 0.12s, box-shadow 0.2s;
}

.btn:hover {
    transform: translateY(-1px);
}

.btn:active {
    transform: translateY(0);
}

.btn:focus-visible {
    outline: 3px solid rgba(255, 209, 102, 0.75);
    outline-offset: 2px;
}

.btn-block {
    width: 100%;
}

.btn-primary {
    background: linear-gradient(135deg, var(--red-500), var(--red-800));
    color: #fff;
    box-shadow: 0 12px 26px -14px rgba(157, 2, 8, 0.9);
}

.btn-primary::after {
    content: "";
    position: absolute;
    top: 0;
    bottom: 0;
    left: -70%;
    width: 45%;
    background: linear-gradient(100deg, transparent, rgba(255, 255, 255, 0.35), transparent);
    transform: skewX(-18deg);
    transition: left 0.55s ease;
    pointer-events: none;
}

.btn-primary:hover {
    background: linear-gradient(135deg, var(--red-400), var(--red-700));
}

.btn-primary:hover::after {
    left: 125%;
}

.btn-success {
    background: linear-gradient(135deg, #2a1418, #120709);
    color: #fff;
    border-color: rgba(225, 29, 51, 0.55);
}

.btn-success:hover {
    background: linear-gradient(135deg, var(--red-700), var(--red-500));
    border-color: var(--red-400);
}

.btn-danger {
    background: linear-gradient(135deg, var(--red-700), var(--red-950));
    color: #fff;
}

.btn-danger:hover {
    background: linear-gradient(135deg, var(--red-600), var(--red-800));
}

.btn-secondary {
    background: #f2e9e7;
    border-color: #e2d3d0;
    color: var(--smoke);
}

.btn-secondary:hover {
    background: #eaddda;
    border-color: #d8c4c0;
}

.btn-outline {
    background: transparent;
    border: 1px solid var(--red-600);
    color: var(--red-700);
}

.btn-outline:hover {
    background: var(--red-600);
    border-color: var(--red-600);
    color: #fff;
}

.btn-sm {
    padding: 0.35rem 0.7rem;
    font-size: 0.82rem;
    border-radius: 7px;
}

/* =========================================================
   Modal
   ========================================================= */
.modal {
    position: fixed;
    inset: 0;
    background: rgba(20, 5, 8, 0.66);
    backdrop-filter: blur(5px);
    display: none;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    z-index: 1000;
}

.modal.active {
    display: flex;
}

.modal-content {
    position: relative;
    overflow: hidden auto;
    background: linear-gradient(180deg, var(--paper), #fff8f6);
    border-radius: 18px;
    border-top: 6px solid var(--red-600);
    width: 100%;
    max-width: 500px;
    max-height: 90vh;
    padding: 1.5rem;
    box-shadow: 0 40px 80px -30px rgba(0, 0, 0, 0.65);
    animation: modalPop 0.24s cubic-bezier(0.2, 0.8, 0.25, 1);
}

/* Bendera kotak-kotak di kepala modal */
.modal-content::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 6px;
    background: conic-gradient(#f8f5f4 25%, var(--red-800) 0 50%, #f8f5f4 0 75%, var(--red-800) 0) 0 0 / 12px 12px;
    opacity: 0.85;
    pointer-events: none;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.2rem;
    padding-bottom: 0.7rem;
    border-bottom: 1px dashed #e8cfcb;
}

.modal-header h2 {
    font-size: 1.2rem;
    color: var(--ink);
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.modal-header h2::before {
    content: "";
    flex: none;
    width: 10px;
    height: 10px;
    border-radius: 2px;
    background: var(--red-600);
    transform: rotate(45deg);
}

.close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    line-height: 1;
    cursor: pointer;
    color: var(--muted);
    border-radius: 8px;
    padding: 0 0.25rem;
    transition: color 0.2s, background 0.2s;
}

.close-btn:hover {
    color: var(--red-600);
    background: var(--red-50);
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1.2rem;
    padding-top: 0.9rem;
    border-top: 1px dashed #e8cfcb;
}

/* =========================================================
   Alert
   ========================================================= */
.alert {
    padding: 0.75rem 1rem;
    border-radius: 10px;
    margin-bottom: 1rem;
    font-size: 0.88rem;
    border-left: 5px solid transparent;
    display: none;
}

.alert-danger {
    background-color: var(--red-50);
    color: var(--red-800);
    border: 1px solid var(--red-200);
    border-left-color: var(--red-600);
}

.alert-success {
    background-color: #f0fbf3;
    color: #186b3a;
    border: 1px solid #c8ecd4;
    border-left-color: #28a745;
}

/* =========================================================
   Tabel Riwayat Servis
   ========================================================= */
#viewHistoryTableContainer {
    overflow-x: auto;
}

#viewHistoryCustomerInfo {
    background: var(--red-50);
    border: 1px dashed var(--red-200);
    border-radius: 10px;
    padding: 0.65rem 0.8rem;
    color: var(--muted);
}

#viewHistoryCustomerInfo strong {
    color: var(--ink);
}

.service-table {
    width: 100%;
    min-width: 560px;
    border-collapse: collapse;
    margin-top: 0.8rem;
    font-size: 0.85rem;
    border-radius: 12px;
    overflow: hidden;
}

.service-table th,
.service-table td {
    padding: 0.65rem;
    text-align: left;
    border: 1px solid var(--line);
}

.service-table th {
    background: linear-gradient(180deg, var(--red-600), var(--red-800));
    color: #fff;
    font-size: 0.72rem;
    letter-spacing: 0.09em;
    text-transform: uppercase;
    border-color: var(--red-800);
}

.service-table td {
    background: var(--paper);
    color: var(--smoke);
}

.service-table th:first-child,
.service-table td:first-child,
.service-table th:nth-child(2),
.service-table td:nth-child(2) {
    white-space: nowrap;
}

.service-table tbody tr:hover td {
    background: #fff5f3;
}

/* =========================================================
   Halaman Autentikasi (Login & Register)
   ========================================================= */
.auth-wrapper {
    position: relative;
    overflow: hidden;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2.5rem 1rem 6rem;
    background-color: var(--red-950);
    background-image:
        repeating-linear-gradient(115deg, rgba(255, 255, 255, 0.045) 0 2px, transparent 2px 18px),
        radial-gradient(900px 500px at 50% -12%, rgba(255, 71, 87, 0.5), transparent 62%),
        radial-gradient(720px 420px at 105% 105%, rgba(255, 209, 102, 0.2), transparent 62%),
        linear-gradient(160deg, #2a0007 0%, #5c0010 34%, #a80d1f 68%, #3d0009 100%);
}

/* Dekorasi jalan + mobil hantu */
.auth-decor {
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 0;
}

.auth-ghost-car {
    position: absolute;
    left: -70px;
    bottom: 92px;
    width: 470px;
    height: auto;
    color: #fff;
    opacity: 0.12;
    animation: carIdle 5.5s ease-in-out infinite;
}

.auth-road {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 76px;
    background: linear-gradient(180deg, transparent, rgba(8, 1, 3, 0.55));
}

.auth-road::before {
    content: "";
    position: absolute;
    inset: auto 0 0 0;
    height: 56px;
    background: radial-gradient(130% 150% at 50% 0%, #2a1216 0%, #100609 68%);
    border-top: 3px solid rgba(255, 255, 255, 0.12);
    box-shadow: 0 -20px 44px -22px rgba(0, 0, 0, 0.85);
}

.auth-road::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    bottom: 26px;
    height: 4px;
    background: repeating-linear-gradient(90deg, rgba(255, 209, 102, 0.85) 0 46px, transparent 46px 104px);
    animation: roadMoveAuth 1.5s linear infinite;
    opacity: 0.8;
}

.auth-card {
    position: relative;
    z-index: 1;
    width: 100%;
    max-width: 430px;
    background: linear-gradient(180deg, var(--paper) 0%, #fff6f4 100%);
    border-radius: 22px;
    padding: 0 2rem 2rem;
    overflow: hidden;
    box-shadow:
        0 40px 80px -30px rgba(0, 0, 0, 0.72),
        0 0 0 1px rgba(255, 255, 255, 0.35) inset;
    animation: cardRise 0.5s cubic-bezier(0.2, 0.8, 0.25, 1) both;
}

/* Strip kotak-kotak di atas kartu */
.auth-card::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 8px;
    background: conic-gradient(#f8f5f4 25%, var(--red-800) 0 50%, #f8f5f4 0 75%, var(--red-800) 0) 0 0 / 14px 14px;
}

/* Glow merah di sudut kartu */
.auth-card::after {
    content: "";
    position: absolute;
    right: -70px;
    bottom: -70px;
    width: 200px;
    height: 200px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(225, 29, 51, 0.14), transparent 70%);
    pointer-events: none;
}

.auth-header {
    text-align: center;
    padding-top: 1.7rem;
    margin-bottom: 1.4rem;
    position: relative;
    z-index: 1;
}

.auth-logo {
    display: grid;
    place-items: center;
    width: 66px;
    height: 66px;
    margin: 0 auto 0.9rem;
    border-radius: 18px;
    color: #fff;
    background: linear-gradient(145deg, var(--red-400), var(--red-800));
    box-shadow:
        0 18px 30px -14px rgba(157, 2, 8, 0.95),
        inset 0 1px 0 rgba(255, 255, 255, 0.45);
    transform: rotate(-4deg);
}

.auth-logo svg {
    width: 44px;
    height: auto;
    filter: drop-shadow(0 1px 1px rgba(0, 0, 0, 0.3));
}

.auth-eyebrow {
    display: block;
    font-size: 0.62rem;
    font-weight: 800;
    letter-spacing: 0.22em;
    text-transform: uppercase;
    color: var(--red-600);
    margin-bottom: 0.4rem;
}

.auth-header h1 {
    font-family: var(--font-display);
    font-size: 1.55rem;
    line-height: 1.1;
    letter-spacing: 0.015em;
    text-transform: uppercase;
    color: var(--ink);
}

.auth-header h1::after {
    content: "";
    display: block;
    width: 64px;
    height: 4px;
    margin: 0.55rem auto 0;
    border-radius: 99px;
    background: linear-gradient(90deg, var(--amber), var(--red-500));
}

.auth-header p {
    font-size: 0.85rem;
    color: var(--muted);
    margin-top: 0.55rem;
}

.auth-card form {
    position: relative;
    z-index: 1;
}

.auth-card .btn-primary {
    margin-top: 0.35rem;
}

.auth-footer {
    position: relative;
    z-index: 1;
    text-align: center;
    margin-top: 1.15rem;
    font-size: 0.88rem;
    color: var(--muted);
}

.auth-footer a {
    color: var(--red-700);
    text-decoration: none;
    font-weight: 700;
    border-bottom: 1px solid transparent;
    transition: border-color 0.2s, color 0.2s;
}

.auth-footer a:hover {
    color: var(--red-500);
    border-bottom-color: var(--red-400);
}

/* =========================================================
   Animasi
   ========================================================= */
@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

@keyframes carIdle {
    0%, 100% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(-5px);
    }
}

@keyframes dash {
    0% {
        transform: translateX(10%);
        opacity: 0;
    }
    25% {
        opacity: 0.95;
    }
    100% {
        transform: translateX(-26%);
        opacity: 0;
    }
}

@keyframes roadMove {
    to {
        background-position-x: 96px;
    }
}

@keyframes roadMoveAuth {
    to {
        background-position-x: 104px;
    }
}

@keyframes glowShift {
    to {
        background-position-x: 200%;
    }
}

@keyframes pulse {
    0%, 100% {
        box-shadow: 0 0 0 3px rgba(255, 209, 102, 0.25);
    }
    50% {
        box-shadow: 0 0 0 6px rgba(255, 209, 102, 0.08);
    }
}

@keyframes modalPop {
    from {
        opacity: 0;
        transform: translateY(16px) scale(0.97);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

@keyframes cardRise {
    from {
        opacity: 0;
        transform: translateY(22px) scale(0.98);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

/* =========================================================
   Responsif
   ========================================================= */
@media (max-width: 860px) {
    .hero {
        grid-template-columns: 1fr;
        padding-bottom: 1.6rem;
    }

    .hero-visual {
        min-height: 160px;
        margin-top: 0.35rem;
    }

    .hero-car {
        max-width: 340px;
    }

    .hero-roundel {
        width: 68px;
        height: 68px;
        top: -12px;
    }

    .roundel-inner b {
        font-size: 1.25rem;
    }

    .auth-ghost-car {
        left: -110px;
        width: 400px;
        opacity: 0.07;
    }
}

@media (max-width: 520px) {
    .card-header-row > div {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: 0.4rem;
    }

    .card-header-row .car-title {
        margin-left: 0 !important;
        line-height: 1.25;
    }
}

@media (max-width: 640px) {
    .container {
        padding: 1.15rem 0.95rem 2.5rem;
    }

    .navbar {
        padding: 0.85rem 0.95rem;
        flex-wrap: wrap;
    }

    .brand-badge {
        width: 42px;
        height: 42px;
    }

    .brand-name {
        font-size: 1.1rem;
    }

    .brand-tagline {
        display: none;
    }

    .nav-user {
        width: 100%;
        justify-content: space-between;
        gap: 0.5rem;
        font-size: 0.84rem;
    }

    .hero-title {
        font-size: 2.05rem;
    }

    .hero-text {
        font-size: 0.86rem;
    }

    .hero-visual {
        min-height: 130px;
    }

    .hero-car {
        max-width: 290px;
    }

    .hero-roundel {
        width: 60px;
        height: 60px;
        right: 0;
    }

    .roundel-inner small {
        font-size: 0.4rem;
    }

    .toolbar {
        padding: 0.9rem;
    }

    .monthly-panel {
        padding: 1rem 0.95rem 0.85rem;
    }

    .monthly-summary {
        width: 100%;
    }

    .summary-pill {
        flex: 1 1 45%;
    }

    .monthly-chart {
        height: 170px;
        min-width: 0;
    }

    .toolbar-controls {
        flex-direction: column;
    }

    .toolbar-search,
    .toolbar-filter {
        min-width: 0;
        width: 100%;
    }

    .toolbar .btn {
        width: 100%;
    }

    .card-actions {
        justify-content: stretch;
    }

    .card-actions .btn {
        flex: 1 1 45%;
    }

    .modal-content {
        padding: 1.2rem;
    }

    .auth-wrapper {
        padding: 1.5rem 0.9rem 5.5rem;
    }

    .auth-card {
        padding: 0 1.35rem 1.6rem;
        border-radius: 18px;
    }

    .auth-ghost-car {
        display: none;
    }
}

/* Hormati preferensi pengguna yang mengurangi animasi */
@media (prefers-reduced-motion: reduce) {
    *,
    *::before,
    *::after {
        animation-duration: 0.001ms !important;
        animation-iteration-count: 1 !important;
        transition-duration: 0.001ms !important;
        scroll-behavior: auto !important;
    }
}
`

const faviconSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64" role="img" aria-label="Garuda Auto Service">
    <defs>
        <linearGradient id="bg" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0" stop-color="#ff4757"/>
            <stop offset="1" stop-color="#9d0208"/>
        </linearGradient>
    </defs>
    <rect width="64" height="64" rx="14" fill="url(#bg)"/>
    <path d="M8 39c0-4 3-6 6-7l8-1c4-6 9-9 16-9h3c7 0 12 3 16 9l4 1c3 1 5 3 5 6v3c0 1-1 2-2 2h-3a7 7 0 0 1-13 0H27a7 7 0 0 1-13 0h-4c-1 0-2-1-2-2z" fill="#fff"/>
    <circle cx="21" cy="45" r="5" fill="#1a0d10"/>
    <circle cx="21" cy="45" r="1.8" fill="#fff"/>
    <circle cx="47" cy="45" r="5" fill="#1a0d10"/>
    <circle cx="47" cy="45" r="1.8" fill="#fff"/>
</svg>
`

// app.js dan auth.js memakai template literal (backtick),
// sehingga ditulis baris demi baris agar tetap valid sebagai string Go.
var appJS = strings.Join([]string{
	"// Logika Aplikasi Dashboard Garuda Auto Service",
	"",
	"checkAuth();",
	"",
	"document.getElementById(\"currentUsername\").textContent = getCurrentUser();",
	"",
	"let currentCustomers = [];",
	"let activeViewingCustomer = null;",
	"",
	"// Elemen DOM",
	"const searchInput = document.getElementById(\"searchInput\");",
	"const brandFilter = document.getElementById(\"brandFilter\");",
	"const customerList = document.getElementById(\"customerList\");",
	"const emptyState = document.getElementById(\"emptyState\");",
	"const alertFeedback = document.getElementById(\"alertFeedback\");",
	"",
	"// Inisialisasi awal",
	"document.addEventListener(\"DOMContentLoaded\", () => {",
	"    loadCarBrands();",
	"    fetchCustomers();",
	"    loadMonthlyStats();",
	"",
	"    searchInput.addEventListener(\"input\", debounce(() => {",
	"        fetchCustomers();",
	"    }, 300));",
	"",
	"    brandFilter.addEventListener(\"change\", () => {",
	"        fetchCustomers();",
	"    });",
	"});",
	"",
	"// Menampilkan pesan notifikasi sementara",
	"function showAlert(message, type = \"success\") {",
	"    alertFeedback.className = `alert alert-${type}`;",
	"    alertFeedback.textContent = message;",
	"    alertFeedback.style.display = \"block\";",
	"    setTimeout(() => {",
	"        alertFeedback.style.display = \"none\";",
	"    }, 3500);",
	"}",
	"",
	"// Mengambil daftar merk mobil unik untuk dropdown filter",
	"async function loadCarBrands() {",
	"    try {",
	"        const res = await authFetch(\"/api/brands\");",
	"        if (!res || !res.ok) return;",
	"        const brands = await res.json();",
	"",
	"        // Pertahankan opsi pertama \"Semua Merk Mobil\"",
	"        const currentSelected = brandFilter.value;",
	"        brandFilter.innerHTML = '<option value=\"\">Semua Merk Mobil</option>';",
	"",
	"        brands.forEach(brand => {",
	"            if (brand) {",
	"                const opt = document.createElement(\"option\");",
	"                opt.value = brand;",
	"                opt.textContent = brand;",
	"                if (brand === currentSelected) opt.selected = true;",
	"                brandFilter.appendChild(opt);",
	"            }",
	"        });",
	"    } catch (err) {",
	"        console.error(\"Gagal memuat merk mobil:\", err);",
	"    }",
	"}",
	"",
	"// Mengambil daftar konsumen berdasarkan pencarian & filter merk",
	"async function fetchCustomers() {",
	"    const search = encodeURIComponent(searchInput.value.trim());",
	"    const brand = encodeURIComponent(brandFilter.value.trim());",
	"",
	"    try {",
	"        const res = await authFetch(`/api/customers?search=${search}&brand=${brand}`);",
	"        if (!res || !res.ok) return;",
	"",
	"        const data = await res.json();",
	"        currentCustomers = data || [];",
	"        renderCustomerList(currentCustomers);",
	"    } catch (err) {",
	"        showAlert(\"Gagal memuat daftar konsumen\", \"danger\");",
	"    }",
	"}",
	"",
	"// Mengambil statistik konsumen & servis per bulan untuk diagram batang",
	"async function loadMonthlyStats() {",
	"    const chart = document.getElementById(\"monthlyChart\");",
	"    if (!chart) return;",
	"",
	"    try {",
	"        const res = await authFetch(\"/api/stats/monthly?months=12\");",
	"        if (!res || !res.ok) {",
	"            renderMonthlyChart(null);",
	"            return;",
	"        }",
	"        renderMonthlyChart(await res.json());",
	"    } catch (err) {",
	"        console.error(\"Gagal memuat statistik bulanan:\", err);",
	"        renderMonthlyChart(null);",
	"    }",
	"}",
	"",
	"// Menampilkan diagram batang konsumen baru dan riwayat servis per bulan",
	"function renderMonthlyChart(data) {",
	"    const chart = document.getElementById(\"monthlyChart\");",
	"    const subtitle = document.getElementById(\"monthlySubtitle\");",
	"    const monthlyCustomers = document.getElementById(\"monthlyCustomers\");",
	"    const monthlyServices = document.getElementById(\"monthlyServices\");",
	"    if (!chart) return;",
	"",
	"    const months = data && Array.isArray(data.months) ? data.months : [];",
	"",
	"    if (months.length === 0) {",
	"        chart.style.gridTemplateColumns = \"1fr\";",
	"        chart.innerHTML = '<p class=\"chart-empty-note\">Belum ada data statistik bulanan yang bisa ditampilkan.</p>';",
	"        if (subtitle) subtitle.textContent = \"Data 12 bulan terakhir belum tersedia.\";",
	"        if (monthlyCustomers) monthlyCustomers.textContent = \"0\";",
	"        if (monthlyServices) monthlyServices.textContent = \"0\";",
	"        return;",
	"    }",
	"",
	"    // Hitung tinggi batang relatif terhadap nilai terbesar (skala 135 piksel)",
	"    const maxValue = Math.max(",
	"        1,",
	"        ...months.map(m => Math.max(Number(m.new_customers) || 0, Number(m.services) || 0))",
	"    );",
	"",
	"    chart.style.gridTemplateColumns = `repeat(${months.length}, minmax(0, 1fr))`;",
	"    chart.innerHTML = months.map((m, index) => {",
	"        const newCustomers = Number(m.new_customers) || 0;",
	"        const services = Number(m.services) || 0;",
	"        const customerHeight = newCustomers > 0",
	"            ? Math.max(4, Math.round((newCustomers / maxValue) * 135))",
	"            : 0;",
	"        const serviceHeight = services > 0",
	"            ? Math.max(4, Math.round((services / maxValue) * 135))",
	"            : 0;",
	"        const isCurrent = index === months.length - 1;",
	"        const tooltip = `${m.label} ${m.year}: ${newCustomers} konsumen baru, ${services} riwayat servis`;",
	"",
	"        return `",
	"            <div class=\"chart-col${isCurrent ? \" is-current\" : \"\"}\" title=\"${escapeHtml(tooltip)}\">",
	"                <div class=\"chart-bars\">",
	"                    <div class=\"chart-bar bar-customer\" style=\"height: ${customerHeight}px\">",
	"                        ${newCustomers > 0 ? `<span>${newCustomers}</span>` : \"\"}",
	"                    </div>",
	"                    <div class=\"chart-bar bar-service\" style=\"height: ${serviceHeight}px\">",
	"                        ${services > 0 ? `<span>${services}</span>` : \"\"}",
	"                    </div>",
	"                </div>",
	"                <span class=\"chart-month\">${escapeHtml(m.label)}</span>",
	"            </div>",
	"        `;",
	"    }).join(\"\");",
	"",
	"    const lastMonth = months[months.length - 1];",
	"    if (monthlyCustomers) monthlyCustomers.textContent = Number(lastMonth.new_customers) || 0;",
	"    if (monthlyServices) monthlyServices.textContent = Number(lastMonth.services) || 0;",
	"",
	"    if (subtitle) {",
	"        const totals = data.totals || {};",
	"        const totalCustomers = Number(totals.new_customers) || 0;",
	"        const totalServices = Number(totals.services) || 0;",
	"        subtitle.textContent = `${months.length} bulan terakhir · total ${totalCustomers} konsumen baru dan ${totalServices} riwayat servis.`;",
	"    }",
	"",
	"    chart.setAttribute(",
	"        \"aria-label\",",
	"        `Diagram batang ${months.length} bulan terakhir. Konsumen baru bulan ini ${Number(lastMonth.new_customers) || 0}, riwayat servis bulan ini ${Number(lastMonth.services) || 0}.`",
	"    );",
	"}",
	"",
	"// Render data konsumen ke antarmuka kartu",
	"function renderCustomerList(customers) {",
	"    customerList.innerHTML = \"\";",
	"",
	"    // Perbarui statistik ringkas pada hero",
	"    const statVehicles = document.getElementById(\"statVehicles\");",
	"    const statServices = document.getElementById(\"statServices\");",
	"    if (statVehicles) statVehicles.textContent = customers ? customers.length : 0;",
	"    if (statServices) {",
	"        const totalServices = (customers || []).reduce(",
	"            (total, cust) => total + (cust.services ? cust.services.length : 0),",
	"            0",
	"        );",
	"        statServices.textContent = totalServices;",
	"    }",
	"",
	"    if (!customers || customers.length === 0) {",
	"        emptyState.style.display = \"block\";",
	"        return;",
	"    }",
	"",
	"    emptyState.style.display = \"none\";",
	"",
	"    customers.forEach(cust => {",
	"        const card = document.createElement(\"div\");",
	"        card.className = \"customer-card\";",
	"",
	"        const services = cust.services || [];",
	"        const serviceCount = services.length;",
	"",
	"        // Informasi servis terakhir",
	"        let serviceSummaryHtml = \"\";",
	"        if (serviceCount > 0) {",
	"            const latest = services[services.length - 1];",
	"            serviceSummaryHtml = `",
	"                <div class=\"service-summary\">",
	"                    <div><strong>Servis Terakhir:</strong> ${escapeHtml(latest.service_date)} (${latest.mileage ? latest.mileage.toLocaleString() + ' KM' : 'KM -'})</div>",
	"                    <div><strong>Pengerjaan:</strong> ${escapeHtml(latest.description)}</div>",
	"                    ${latest.next_service_date || latest.next_mileage ? `",
	"                        <div style=\"margin-top: 4px; color: #d9534f; font-weight: 600;\">",
	"                            📅 Jadwal Follow-Up: ${latest.next_service_date ? escapeHtml(latest.next_service_date) : ''} ",
	"                            ${latest.next_mileage ? `| Target: ${latest.next_mileage.toLocaleString()} KM` : ''}",
	"                        </div>",
	"                    ` : ''}",
	"                </div>",
	"            `;",
	"        } else {",
	"            serviceSummaryHtml = `",
	"                <div class=\"service-summary\" style=\"color: #888;\">",
	"                    <em>Belum ada catatan riwayat servis untuk kendaraan ini.</em>",
	"                </div>",
	"            `;",
	"        }",
	"",
	"        card.innerHTML = `",
	"            <div class=\"card-header-row\">",
	"                <div>",
	"                    <span class=\"plate-badge\">${escapeHtml(cust.plate_number)}</span>",
	"                    <span class=\"car-title\" style=\"margin-left: 0.5rem;\">${escapeHtml(cust.car_brand)} ${escapeHtml(cust.car_type)}</span>",
	"                </div>",
	"            </div>",
	"            <div class=\"customer-meta\">",
	"                <div>👤 <strong>Pemilik:</strong> ${escapeHtml(cust.name)}</div>",
	"                <div>📞 <strong>Kontak:</strong> ${cust.phone ? escapeHtml(cust.phone) : '-'}</div>",
	"                <div>📍 <strong>Alamat:</strong> ${cust.address ? escapeHtml(cust.address) : '-'}</div>",
	"            </div>",
	"            ${serviceSummaryHtml}",
	"            <div class=\"card-actions\">",
	"                <button onclick=\"openServiceModal(${cust.id}, '${escapeHtml(cust.name)}', '${escapeHtml(cust.plate_number)}')\" class=\"btn btn-success btn-sm\">+ Tambah Servis</button>",
	"                <button onclick=\"viewHistory(${cust.id})\" class=\"btn btn-outline btn-sm\">Riwayat (${serviceCount})</button>",
	"                <button onclick=\"editCustomer(${cust.id})\" class=\"btn btn-secondary btn-sm\">Edit</button>",
	"                <button onclick=\"deleteCustomer(${cust.id}, '${escapeHtml(cust.plate_number)}')\" class=\"btn btn-danger btn-sm\">Hapus</button>",
	"            </div>",
	"        `;",
	"",
	"        customerList.appendChild(card);",
	"    });",
	"}",
	"",
	"// Modal Konsumen (Tambah & Edit)",
	"function openCustomerModal(isEdit = false) {",
	"    document.getElementById(\"customerModal\").classList.add(\"active\");",
	"    if (!isEdit) {",
	"        document.getElementById(\"customerModalTitle\").textContent = \"Tambah Konsumen Baru\";",
	"        document.getElementById(\"customerForm\").reset();",
	"        document.getElementById(\"customerId\").value = \"\";",
	"    }",
	"}",
	"",
	"function closeCustomerModal() {",
	"    document.getElementById(\"customerModal\").classList.remove(\"active\");",
	"}",
	"",
	"function editCustomer(id) {",
	"    const cust = currentCustomers.find(c => c.id === id);",
	"    if (!cust) return;",
	"",
	"    document.getElementById(\"customerModalTitle\").textContent = \"Edit Data Konsumen & Kendaraan\";",
	"    document.getElementById(\"customerId\").value = cust.id;",
	"    document.getElementById(\"custName\").value = cust.name;",
	"    document.getElementById(\"custPhone\").value = cust.phone || \"\";",
	"    document.getElementById(\"custAddress\").value = cust.address || \"\";",
	"    document.getElementById(\"custPlate\").value = cust.plate_number;",
	"    document.getElementById(\"custBrand\").value = cust.car_brand;",
	"    document.getElementById(\"custType\").value = cust.car_type;",
	"",
	"    openCustomerModal(true);",
	"}",
	"",
	"document.getElementById(\"customerForm\").addEventListener(\"submit\", async (e) => {",
	"    e.preventDefault();",
	"",
	"    const id = document.getElementById(\"customerId\").value;",
	"    const payload = {",
	"        name: document.getElementById(\"custName\").value.trim(),",
	"        phone: document.getElementById(\"custPhone\").value.trim(),",
	"        address: document.getElementById(\"custAddress\").value.trim(),",
	"        plate_number: document.getElementById(\"custPlate\").value.trim(),",
	"        car_brand: document.getElementById(\"custBrand\").value.trim(),",
	"        car_type: document.getElementById(\"custType\").value.trim(),",
	"    };",
	"",
	"    const isEdit = !!id;",
	"    const url = isEdit ? `/api/customers/${id}` : \"/api/customers\";",
	"    const method = isEdit ? \"PUT\" : \"POST\";",
	"",
	"    try {",
	"        const res = await authFetch(url, {",
	"            method: method,",
	"            body: JSON.stringify(payload)",
	"        });",
	"",
	"        if (!res || !res.ok) {",
	"            const errData = await res.json();",
	"            alert(errData.error || \"Gagal menyimpan data.\");",
	"            return;",
	"        }",
	"",
	"        closeCustomerModal();",
	"        showAlert(isEdit ? \"Data konsumen berhasil diperbarui!\" : \"Konsumen baru berhasil ditambahkan!\");",
	"        loadCarBrands();",
	"        fetchCustomers();",
	"        loadMonthlyStats();",
	"    } catch (err) {",
	"        alert(\"Terjadi kesalahan server saat menyimpan data konsumen.\");",
	"    }",
	"});",
	"",
	"// Hapus Konsumen",
	"async function deleteCustomer(id, plateNumber) {",
	"    if (!confirm(`Apakah Anda yakin ingin menghapus data kendaraan [${plateNumber}] beserta seluruh riwayat servisnya?`)) {",
	"        return;",
	"    }",
	"",
	"    try {",
	"        const res = await authFetch(`/api/customers/${id}`, { method: \"DELETE\" });",
	"        if (!res || !res.ok) {",
	"            showAlert(\"Gagal menghapus data konsumen\", \"danger\");",
	"            return;",
	"        }",
	"",
	"        showAlert(\"Data konsumen dan riwayat servis berhasil dihapus.\");",
	"        loadCarBrands();",
	"        fetchCustomers();",
	"        loadMonthlyStats();",
	"    } catch (err) {",
	"        showAlert(\"Terjadi kesalahan server saat menghapus data\", \"danger\");",
	"    }",
	"}",
	"",
	"// Modal Tambah Servis",
	"function openServiceModal(customerId, name, plateNumber) {",
	"    document.getElementById(\"serviceCustomerId\").value = customerId;",
	"    document.getElementById(\"serviceModalTitle\").textContent = `Tambah Servis: ${plateNumber} (${name})`;",
	"    document.getElementById(\"serviceForm\").reset();",
	"",
	"    // Default tanggal hari ini (YYYY-MM-DD)",
	"    const today = new Date().toISOString().split(\"T\")[0];",
	"    document.getElementById(\"serviceDate\").value = today;",
	"",
	"    document.getElementById(\"serviceModal\").classList.add(\"active\");",
	"}",
	"",
	"function closeServiceModal() {",
	"    document.getElementById(\"serviceModal\").classList.remove(\"active\");",
	"}",
	"",
	"document.getElementById(\"serviceForm\").addEventListener(\"submit\", async (e) => {",
	"    e.preventDefault();",
	"",
	"    const customerId = document.getElementById(\"serviceCustomerId\").value;",
	"    const mileageVal = document.getElementById(\"serviceMileage\").value;",
	"    const nextMileageVal = document.getElementById(\"nextMileage\").value;",
	"",
	"    const payload = {",
	"        service_date: document.getElementById(\"serviceDate\").value,",
	"        mileage: mileageVal ? parseInt(mileageVal, 10) : 0,",
	"        description: document.getElementById(\"serviceDesc\").value.trim(),",
	"        next_service_date: document.getElementById(\"nextServiceDate\").value,",
	"        next_mileage: nextMileageVal ? parseInt(nextMileageVal, 10) : 0,",
	"        notes: document.getElementById(\"serviceNotes\").value.trim()",
	"    };",
	"",
	"    try {",
	"        const res = await authFetch(`/api/customers/${customerId}/services`, {",
	"            method: \"POST\",",
	"            body: JSON.stringify(payload)",
	"        });",
	"",
	"        if (!res || !res.ok) {",
	"            const errData = await res.json();",
	"            alert(errData.error || \"Gagal menyimpan riwayat servis.\");",
	"            return;",
	"        }",
	"",
	"        closeServiceModal();",
	"        showAlert(\"Riwayat servis baru berhasil dicatat!\");",
	"        fetchCustomers();",
	"        loadMonthlyStats();",
	"    } catch (err) {",
	"        alert(\"Terjadi kesalahan server saat mencatat servis.\");",
	"    }",
	"});",
	"",
	"// Modal Lihat Riwayat Servis Lengkap",
	"async function viewHistory(customerId) {",
	"    try {",
	"        const res = await authFetch(`/api/customers/${customerId}`);",
	"        if (!res || !res.ok) return;",
	"",
	"        const cust = await res.json();",
	"        activeViewingCustomer = cust;",
	"",
	"        document.getElementById(\"viewHistoryTitle\").textContent = `Riwayat Servis [${cust.plate_number}]`;",
	"        document.getElementById(\"viewHistoryCustomerInfo\").innerHTML = `",
	"            <strong>Konsumen:</strong> ${escapeHtml(cust.name)} | ",
	"            <strong>Mobil:</strong> ${escapeHtml(cust.car_brand)} ${escapeHtml(cust.car_type)} | ",
	"            <strong>Kontak:</strong> ${cust.phone ? escapeHtml(cust.phone) : '-'}",
	"        `;",
	"",
	"        renderHistoryTable(cust.services || []);",
	"        document.getElementById(\"viewHistoryModal\").classList.add(\"active\");",
	"    } catch (err) {",
	"        showAlert(\"Gagal memuat detail riwayat servis\", \"danger\");",
	"    }",
	"}",
	"",
	"function closeViewHistoryModal() {",
	"    document.getElementById(\"viewHistoryModal\").classList.remove(\"active\");",
	"    activeViewingCustomer = null;",
	"}",
	"",
	"function renderHistoryTable(services) {",
	"    const container = document.getElementById(\"viewHistoryTableContainer\");",
	"    if (!services || services.length === 0) {",
	"        container.innerHTML = '<p style=\"color: #777; padding: 1rem 0;\">Belum ada riwayat servis yang tercatat untuk kendaraan ini.</p>';",
	"        return;",
	"    }",
	"",
	"    let rowsHtml = services.map(s => `",
	"        <tr>",
	"            <td>${escapeHtml(s.service_date)}</td>",
	"            <td>${s.mileage ? s.mileage.toLocaleString() + ' KM' : '-'}</td>",
	"            <td>${escapeHtml(s.description)}</td>",
	"            <td>",
	"                ${s.next_service_date ? `📅 ${escapeHtml(s.next_service_date)}` : ''}",
	"                ${s.next_mileage ? `<br>🎯 ${s.next_mileage.toLocaleString()} KM` : ''}",
	"                ${!s.next_service_date && !s.next_mileage ? '-' : ''}",
	"            </td>",
	"            <td>${s.notes ? escapeHtml(s.notes) : '-'}</td>",
	"            <td>",
	"                <button onclick=\"deleteServiceRecord(${s.id})\" class=\"btn btn-danger btn-sm\">Hapus</button>",
	"            </td>",
	"        </tr>",
	"    `).join(\"\");",
	"",
	"    container.innerHTML = `",
	"        <table class=\"service-table\">",
	"            <thead>",
	"                <tr>",
	"                    <th>Tanggal</th>",
	"                    <th>KM</th>",
	"                    <th>Pengerjaan</th>",
	"                    <th>Jadwal Berikutnya</th>",
	"                    <th>Catatan</th>",
	"                    <th>Aksi</th>",
	"                </tr>",
	"            </thead>",
	"            <tbody>",
	"                ${rowsHtml}",
	"            </tbody>",
	"        </table>",
	"    `;",
	"}",
	"",
	"// Hapus Satu Riwayat Servis",
	"async function deleteServiceRecord(serviceId) {",
	"    if (!confirm(\"Hapus catatan servis ini?\")) return;",
	"",
	"    try {",
	"        const res = await authFetch(`/api/services/${serviceId}`, { method: \"DELETE\" });",
	"        if (!res || !res.ok) {",
	"            alert(\"Gagal menghapus riwayat servis.\");",
	"            return;",
	"        }",
	"",
	"        showAlert(\"Catatan servis berhasil dihapus.\");",
	"        if (activeViewingCustomer) {",
	"            viewHistory(activeViewingCustomer.id);",
	"        }",
	"        fetchCustomers();",
	"        loadMonthlyStats();",
	"    } catch (err) {",
	"        alert(\"Terjadi kesalahan server saat menghapus catatan.\");",
	"    }",
	"}",
	"",
	"// Utility: Debounce untuk input pencarian cepat",
	"function debounce(func, wait) {",
	"    let timeout;",
	"    return function (...args) {",
	"        clearTimeout(timeout);",
	"        timeout = setTimeout(() => func.apply(this, args), wait);",
	"    };",
	"}",
	"",
	"// Utility: Mencegah serangan XSS pada tampilan data",
	"function escapeHtml(str) {",
	"    if (!str) return \"\";",
	"    return String(str)",
	"        .replace(/&/g, \"&amp;\")",
	"        .replace(/</g, \"&lt;\")",
	"        .replace(/>/g, \"&gt;\")",
	"        .replace(/\"/g, \"&quot;\")",
	"        .replace(/'/g, \"&#039;\");",
	"}",
	"",
}, "\n")

var authJS = strings.Join([]string{
	"// Modul Autentikasi Frontend (Garuda Auto Service)",
	"",
	"const AUTH_TOKEN_KEY = \"gas_token\";",
	"const AUTH_USER_KEY = \"gas_user\";",
	"",
	"// Mengambil token dari localStorage",
	"function getToken() {",
	"    return localStorage.getItem(AUTH_TOKEN_KEY);",
	"}",
	"",
	"// Mengambil nama pengguna saat ini",
	"function getCurrentUser() {",
	"    return localStorage.getItem(AUTH_USER_KEY) || \"Staf\";",
	"}",
	"",
	"// Memeriksa status login saat mengakses halaman utama",
	"function checkAuth() {",
	"    const token = getToken();",
	"    if (!token) {",
	"        window.location.href = \"/login\";",
	"    }",
	"}",
	"",
	"// Mencegah user yang sudah login mengakses halaman login/register",
	"function redirectIfLoggedIn() {",
	"    const token = getToken();",
	"    if (token) {",
	"        window.location.href = \"/\";",
	"    }",
	"}",
	"",
	"// Fungsi pembantu fetch dengan header otorisasi JWT otomatis",
	"async function authFetch(url, options = {}) {",
	"    const token = getToken();",
	"    if (!token) {",
	"        window.location.href = \"/login\";",
	"        return;",
	"    }",
	"",
	"    const headers = {",
	"        \"Content-Type\": \"application/json\",",
	"        \"Authorization\": `Bearer ${token}`,",
	"        ...(options.headers || {})",
	"    };",
	"",
	"    const response = await fetch(url, { ...options, headers });",
	"",
	"    // Jika token kedaluwarsa atau tidak valid (401), logout dan arahkan ke login",
	"    if (response.status === 401) {",
	"        logout();",
	"        return;",
	"    }",
	"",
	"    return response;",
	"}",
	"",
	"// Logout dari sistem",
	"function logout() {",
	"    localStorage.removeItem(AUTH_TOKEN_KEY);",
	"    localStorage.removeItem(AUTH_USER_KEY);",
	"    window.location.href = \"/login\";",
	"}",
	"",
}, "\n")

// registerAsset mendaftarkan aset UI untuk metode GET dan HEAD.
func registerAsset(router *gin.Engine, urlPath, contentType, content string) {
	handler := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, contentType, []byte(content))
	}
	router.GET(urlPath, handler)
	router.HEAD(urlPath, handler)
}

func main() {
	// Muat file .env jika tersedia
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan variabel environment sistem.")
	}

	// Hubungkan ke database PostgreSQL
	if _, err := config.ConnectDatabase(); err != nil {
		log.Fatalf("[FATAL] Gagal inisialisasi database: %v", err)
	}

	router := gin.Default()

	// Sajikan aset UI yang sudah ditanam di main.go
	registerAsset(router, "/", "text/html; charset=utf-8", indexHTML)
	registerAsset(router, "/login", "text/html; charset=utf-8", loginHTML)
	registerAsset(router, "/register", "text/html; charset=utf-8", registerHTML)
	registerAsset(router, "/css/style.css", "text/css; charset=utf-8", styleCSS)
	registerAsset(router, "/js/app.js", "application/javascript; charset=utf-8", appJS)
	registerAsset(router, "/js/auth.js", "application/javascript; charset=utf-8", authJS)
	registerAsset(router, "/favicon.svg", "image/svg+xml", faviconSVG)

	// Endpoint API Publik (Autentikasi)
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}

	// Endpoint API Terlindungi (Memerlukan Token JWT)
	api := router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/auth/me", handlers.Me)

		// Data Konsumen & Riwayat Servis
		api.GET("/brands", handlers.GetCarBrands)
		api.GET("/stats/monthly", handlers.GetMonthlyStats)
		api.GET("/customers", handlers.GetCustomers)
		api.GET("/customers/:id", handlers.GetCustomerByID)
		api.POST("/customers", handlers.CreateCustomer)
		api.PUT("/customers/:id", handlers.UpdateCustomer)
		api.DELETE("/customers/:id", handlers.DeleteCustomer)

		// Riwayat Servis Khusus Konsumen
		api.POST("/customers/:id/services", handlers.AddServiceHistory)
		api.DELETE("/services/:id", handlers.DeleteServiceHistory)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[INFO] Server Garuda Auto Service berjalan di http://localhost:%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[FATAL] Gagal menjalankan server: %v", err)
	}
}
