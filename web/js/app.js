// Logika Aplikasi Dashboard Garuda Auto Service

checkAuth();

document.getElementById("currentUsername").textContent = getCurrentUser();

let currentCustomers = [];
let activeViewingCustomer = null;

// Elemen DOM
const searchInput = document.getElementById("searchInput");
const brandFilter = document.getElementById("brandFilter");
const customerList = document.getElementById("customerList");
const emptyState = document.getElementById("emptyState");
const alertFeedback = document.getElementById("alertFeedback");

// Inisialisasi awal
document.addEventListener("DOMContentLoaded", () => {
    loadCarBrands();
    fetchCustomers();
    loadMonthlyStats();

    searchInput.addEventListener("input", debounce(() => {
        fetchCustomers();
    }, 300));

    brandFilter.addEventListener("change", () => {
        fetchCustomers();
    });
});

// Menampilkan pesan notifikasi sementara
function showAlert(message, type = "success") {
    alertFeedback.className = `alert alert-${type}`;
    alertFeedback.textContent = message;
    alertFeedback.style.display = "block";
    setTimeout(() => {
        alertFeedback.style.display = "none";
    }, 3500);
}

// Mengambil daftar merk mobil unik untuk dropdown filter
async function loadCarBrands() {
    try {
        const res = await authFetch("/api/brands");
        if (!res || !res.ok) return;
        const brands = await res.json();

        // Pertahankan opsi pertama "Semua Merk Mobil"
        const currentSelected = brandFilter.value;
        brandFilter.innerHTML = '<option value="">Semua Merk Mobil</option>';

        brands.forEach(brand => {
            if (brand) {
                const opt = document.createElement("option");
                opt.value = brand;
                opt.textContent = brand;
                if (brand === currentSelected) opt.selected = true;
                brandFilter.appendChild(opt);
            }
        });
    } catch (err) {
        console.error("Gagal memuat merk mobil:", err);
    }
}

// Mengambil daftar konsumen berdasarkan pencarian & filter merk
async function fetchCustomers() {
    const search = encodeURIComponent(searchInput.value.trim());
    const brand = encodeURIComponent(brandFilter.value.trim());

    try {
        const res = await authFetch(`/api/customers?search=${search}&brand=${brand}`);
        if (!res || !res.ok) return;

        const data = await res.json();
        currentCustomers = data || [];
        renderCustomerList(currentCustomers);
    } catch (err) {
        showAlert("Gagal memuat daftar konsumen", "danger");
    }
}

// Mengambil statistik konsumen & servis per bulan untuk diagram batang
async function loadMonthlyStats() {
    const chart = document.getElementById("monthlyChart");
    if (!chart) return;

    try {
        const res = await authFetch("/api/stats/monthly?months=12");
        if (!res || !res.ok) {
            renderMonthlyChart(null);
            return;
        }
        renderMonthlyChart(await res.json());
    } catch (err) {
        console.error("Gagal memuat statistik bulanan:", err);
        renderMonthlyChart(null);
    }
}

// Menampilkan diagram batang konsumen baru dan riwayat servis per bulan
function renderMonthlyChart(data) {
    const chart = document.getElementById("monthlyChart");
    const subtitle = document.getElementById("monthlySubtitle");
    const monthlyCustomers = document.getElementById("monthlyCustomers");
    const monthlyServices = document.getElementById("monthlyServices");
    if (!chart) return;

    const months = data && Array.isArray(data.months) ? data.months : [];

    if (months.length === 0) {
        chart.style.gridTemplateColumns = "1fr";
        chart.innerHTML = '<p class="chart-empty-note">Belum ada data statistik bulanan yang bisa ditampilkan.</p>';
        if (subtitle) subtitle.textContent = "Data 12 bulan terakhir belum tersedia.";
        if (monthlyCustomers) monthlyCustomers.textContent = "0";
        if (monthlyServices) monthlyServices.textContent = "0";
        return;
    }

    // Hitung tinggi batang relatif terhadap nilai terbesar (skala 135 piksel)
    const maxValue = Math.max(
        1,
        ...months.map(m => Math.max(Number(m.new_customers) || 0, Number(m.services) || 0))
    );

    chart.style.gridTemplateColumns = `repeat(${months.length}, minmax(0, 1fr))`;
    chart.innerHTML = months.map((m, index) => {
        const newCustomers = Number(m.new_customers) || 0;
        const services = Number(m.services) || 0;
        const customerHeight = newCustomers > 0
            ? Math.max(4, Math.round((newCustomers / maxValue) * 135))
            : 0;
        const serviceHeight = services > 0
            ? Math.max(4, Math.round((services / maxValue) * 135))
            : 0;
        const isCurrent = index === months.length - 1;
        const tooltip = `${m.label} ${m.year}: ${newCustomers} konsumen baru, ${services} riwayat servis`;

        return `
            <div class="chart-col${isCurrent ? " is-current" : ""}" title="${escapeHtml(tooltip)}">
                <div class="chart-bars">
                    <div class="chart-bar bar-customer" style="height: ${customerHeight}px">
                        ${newCustomers > 0 ? `<span>${newCustomers}</span>` : ""}
                    </div>
                    <div class="chart-bar bar-service" style="height: ${serviceHeight}px">
                        ${services > 0 ? `<span>${services}</span>` : ""}
                    </div>
                </div>
                <span class="chart-month">${escapeHtml(m.label)}</span>
            </div>
        `;
    }).join("");

    const lastMonth = months[months.length - 1];
    if (monthlyCustomers) monthlyCustomers.textContent = Number(lastMonth.new_customers) || 0;
    if (monthlyServices) monthlyServices.textContent = Number(lastMonth.services) || 0;

    if (subtitle) {
        const totals = data.totals || {};
        const totalCustomers = Number(totals.new_customers) || 0;
        const totalServices = Number(totals.services) || 0;
        subtitle.textContent = `${months.length} bulan terakhir · total ${totalCustomers} konsumen baru dan ${totalServices} riwayat servis.`;
    }

    chart.setAttribute(
        "aria-label",
        `Diagram batang ${months.length} bulan terakhir. Konsumen baru bulan ini ${Number(lastMonth.new_customers) || 0}, riwayat servis bulan ini ${Number(lastMonth.services) || 0}.`
    );
}

// Render data konsumen ke antarmuka kartu
function renderCustomerList(customers) {
    customerList.innerHTML = "";

    // Perbarui statistik ringkas pada hero
    const statVehicles = document.getElementById("statVehicles");
    const statServices = document.getElementById("statServices");
    if (statVehicles) statVehicles.textContent = customers ? customers.length : 0;
    if (statServices) {
        const totalServices = (customers || []).reduce(
            (total, cust) => total + (cust.services ? cust.services.length : 0),
            0
        );
        statServices.textContent = totalServices;
    }

    if (!customers || customers.length === 0) {
        emptyState.style.display = "block";
        return;
    }

    emptyState.style.display = "none";

    customers.forEach(cust => {
        const card = document.createElement("div");
        card.className = "customer-card";

        const services = cust.services || [];
        const serviceCount = services.length;

        // Informasi servis terakhir
        let serviceSummaryHtml = "";
        if (serviceCount > 0) {
            const latest = services[services.length - 1];
            serviceSummaryHtml = `
                <div class="service-summary">
                    <div><strong>Servis Terakhir:</strong> ${escapeHtml(latest.service_date)} (${latest.mileage ? latest.mileage.toLocaleString() + ' KM' : 'KM -'})</div>
                    <div><strong>Pengerjaan:</strong> ${escapeHtml(latest.description)}</div>
                    ${latest.next_service_date || latest.next_mileage ? `
                        <div style="margin-top: 4px; color: #d9534f; font-weight: 600;">
                            📅 Jadwal Follow-Up: ${latest.next_service_date ? escapeHtml(latest.next_service_date) : ''} 
                            ${latest.next_mileage ? `| Target: ${latest.next_mileage.toLocaleString()} KM` : ''}
                        </div>
                    ` : ''}
                </div>
            `;
        } else {
            serviceSummaryHtml = `
                <div class="service-summary" style="color: #888;">
                    <em>Belum ada catatan riwayat servis untuk kendaraan ini.</em>
                </div>
            `;
        }

        card.innerHTML = `
            <div class="card-header-row">
                <div>
                    <span class="plate-badge">${escapeHtml(cust.plate_number)}</span>
                    <span class="car-title" style="margin-left: 0.5rem;">${escapeHtml(cust.car_brand)} ${escapeHtml(cust.car_type)}</span>
                </div>
            </div>
            <div class="customer-meta">
                <div>👤 <strong>Pemilik:</strong> ${escapeHtml(cust.name)}</div>
                <div>📞 <strong>Kontak:</strong> ${cust.phone ? escapeHtml(cust.phone) : '-'}</div>
                <div>📍 <strong>Alamat:</strong> ${cust.address ? escapeHtml(cust.address) : '-'}</div>
            </div>
            ${serviceSummaryHtml}
            <div class="card-actions">
                <button onclick="openServiceModal(${cust.id}, '${escapeHtml(cust.name)}', '${escapeHtml(cust.plate_number)}')" class="btn btn-success btn-sm">+ Tambah Servis</button>
                <button onclick="viewHistory(${cust.id})" class="btn btn-outline btn-sm">Riwayat (${serviceCount})</button>
                <button onclick="editCustomer(${cust.id})" class="btn btn-secondary btn-sm">Edit</button>
                <button onclick="deleteCustomer(${cust.id}, '${escapeHtml(cust.plate_number)}')" class="btn btn-danger btn-sm">Hapus</button>
            </div>
        `;

        customerList.appendChild(card);
    });
}

// Modal Konsumen (Tambah & Edit)
function openCustomerModal(isEdit = false) {
    document.getElementById("customerModal").classList.add("active");
    if (!isEdit) {
        document.getElementById("customerModalTitle").textContent = "Tambah Konsumen Baru";
        document.getElementById("customerForm").reset();
        document.getElementById("customerId").value = "";
    }
}

function closeCustomerModal() {
    document.getElementById("customerModal").classList.remove("active");
}

function editCustomer(id) {
    const cust = currentCustomers.find(c => c.id === id);
    if (!cust) return;

    document.getElementById("customerModalTitle").textContent = "Edit Data Konsumen & Kendaraan";
    document.getElementById("customerId").value = cust.id;
    document.getElementById("custName").value = cust.name;
    document.getElementById("custPhone").value = cust.phone || "";
    document.getElementById("custAddress").value = cust.address || "";
    document.getElementById("custPlate").value = cust.plate_number;
    document.getElementById("custBrand").value = cust.car_brand;
    document.getElementById("custType").value = cust.car_type;

    openCustomerModal(true);
}

document.getElementById("customerForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const id = document.getElementById("customerId").value;
    const payload = {
        name: document.getElementById("custName").value.trim(),
        phone: document.getElementById("custPhone").value.trim(),
        address: document.getElementById("custAddress").value.trim(),
        plate_number: document.getElementById("custPlate").value.trim(),
        car_brand: document.getElementById("custBrand").value.trim(),
        car_type: document.getElementById("custType").value.trim(),
    };

    const isEdit = !!id;
    const url = isEdit ? `/api/customers/${id}` : "/api/customers";
    const method = isEdit ? "PUT" : "POST";

    try {
        const res = await authFetch(url, {
            method: method,
            body: JSON.stringify(payload)
        });

        if (!res || !res.ok) {
            const errData = await res.json();
            alert(errData.error || "Gagal menyimpan data.");
            return;
        }

        closeCustomerModal();
        showAlert(isEdit ? "Data konsumen berhasil diperbarui!" : "Konsumen baru berhasil ditambahkan!");
        loadCarBrands();
        fetchCustomers();
        loadMonthlyStats();
    } catch (err) {
        alert("Terjadi kesalahan server saat menyimpan data konsumen.");
    }
});

// Hapus Konsumen
async function deleteCustomer(id, plateNumber) {
    if (!confirm(`Apakah Anda yakin ingin menghapus data kendaraan [${plateNumber}] beserta seluruh riwayat servisnya?`)) {
        return;
    }

    try {
        const res = await authFetch(`/api/customers/${id}`, { method: "DELETE" });
        if (!res || !res.ok) {
            showAlert("Gagal menghapus data konsumen", "danger");
            return;
        }

        showAlert("Data konsumen dan riwayat servis berhasil dihapus.");
        loadCarBrands();
        fetchCustomers();
        loadMonthlyStats();
    } catch (err) {
        showAlert("Terjadi kesalahan server saat menghapus data", "danger");
    }
}

// Modal Tambah Servis
function openServiceModal(customerId, name, plateNumber) {
    document.getElementById("serviceCustomerId").value = customerId;
    document.getElementById("serviceModalTitle").textContent = `Tambah Servis: ${plateNumber} (${name})`;
    document.getElementById("serviceForm").reset();

    // Default tanggal hari ini (YYYY-MM-DD)
    const today = new Date().toISOString().split("T")[0];
    document.getElementById("serviceDate").value = today;

    document.getElementById("serviceModal").classList.add("active");
}

function closeServiceModal() {
    document.getElementById("serviceModal").classList.remove("active");
}

document.getElementById("serviceForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const customerId = document.getElementById("serviceCustomerId").value;
    const mileageVal = document.getElementById("serviceMileage").value;
    const nextMileageVal = document.getElementById("nextMileage").value;

    const payload = {
        service_date: document.getElementById("serviceDate").value,
        mileage: mileageVal ? parseInt(mileageVal, 10) : 0,
        description: document.getElementById("serviceDesc").value.trim(),
        next_service_date: document.getElementById("nextServiceDate").value,
        next_mileage: nextMileageVal ? parseInt(nextMileageVal, 10) : 0,
        notes: document.getElementById("serviceNotes").value.trim()
    };

    try {
        const res = await authFetch(`/api/customers/${customerId}/services`, {
            method: "POST",
            body: JSON.stringify(payload)
        });

        if (!res || !res.ok) {
            const errData = await res.json();
            alert(errData.error || "Gagal menyimpan riwayat servis.");
            return;
        }

        closeServiceModal();
        showAlert("Riwayat servis baru berhasil dicatat!");
        fetchCustomers();
        loadMonthlyStats();
    } catch (err) {
        alert("Terjadi kesalahan server saat mencatat servis.");
    }
});

// Modal Lihat Riwayat Servis Lengkap
async function viewHistory(customerId) {
    try {
        const res = await authFetch(`/api/customers/${customerId}`);
        if (!res || !res.ok) return;

        const cust = await res.json();
        activeViewingCustomer = cust;

        document.getElementById("viewHistoryTitle").textContent = `Riwayat Servis [${cust.plate_number}]`;
        document.getElementById("viewHistoryCustomerInfo").innerHTML = `
            <strong>Konsumen:</strong> ${escapeHtml(cust.name)} | 
            <strong>Mobil:</strong> ${escapeHtml(cust.car_brand)} ${escapeHtml(cust.car_type)} | 
            <strong>Kontak:</strong> ${cust.phone ? escapeHtml(cust.phone) : '-'}
        `;

        renderHistoryTable(cust.services || []);
        document.getElementById("viewHistoryModal").classList.add("active");
    } catch (err) {
        showAlert("Gagal memuat detail riwayat servis", "danger");
    }
}

function closeViewHistoryModal() {
    document.getElementById("viewHistoryModal").classList.remove("active");
    activeViewingCustomer = null;
}

function renderHistoryTable(services) {
    const container = document.getElementById("viewHistoryTableContainer");
    if (!services || services.length === 0) {
        container.innerHTML = '<p style="color: #777; padding: 1rem 0;">Belum ada riwayat servis yang tercatat untuk kendaraan ini.</p>';
        return;
    }

    let rowsHtml = services.map(s => `
        <tr>
            <td>${escapeHtml(s.service_date)}</td>
            <td>${s.mileage ? s.mileage.toLocaleString() + ' KM' : '-'}</td>
            <td>${escapeHtml(s.description)}</td>
            <td>
                ${s.next_service_date ? `📅 ${escapeHtml(s.next_service_date)}` : ''}
                ${s.next_mileage ? `<br>🎯 ${s.next_mileage.toLocaleString()} KM` : ''}
                ${!s.next_service_date && !s.next_mileage ? '-' : ''}
            </td>
            <td>${s.notes ? escapeHtml(s.notes) : '-'}</td>
            <td>
                <button onclick="deleteServiceRecord(${s.id})" class="btn btn-danger btn-sm">Hapus</button>
            </td>
        </tr>
    `).join("");

    container.innerHTML = `
        <table class="service-table">
            <thead>
                <tr>
                    <th>Tanggal</th>
                    <th>KM</th>
                    <th>Pengerjaan</th>
                    <th>Jadwal Berikutnya</th>
                    <th>Catatan</th>
                    <th>Aksi</th>
                </tr>
            </thead>
            <tbody>
                ${rowsHtml}
            </tbody>
        </table>
    `;
}

// Hapus Satu Riwayat Servis
async function deleteServiceRecord(serviceId) {
    if (!confirm("Hapus catatan servis ini?")) return;

    try {
        const res = await authFetch(`/api/services/${serviceId}`, { method: "DELETE" });
        if (!res || !res.ok) {
            alert("Gagal menghapus riwayat servis.");
            return;
        }

        showAlert("Catatan servis berhasil dihapus.");
        if (activeViewingCustomer) {
            viewHistory(activeViewingCustomer.id);
        }
        fetchCustomers();
        loadMonthlyStats();
    } catch (err) {
        alert("Terjadi kesalahan server saat menghapus catatan.");
    }
}

// Utility: Debounce untuk input pencarian cepat
function debounce(func, wait) {
    let timeout;
    return function (...args) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    };
}

// Utility: Mencegah serangan XSS pada tampilan data
function escapeHtml(str) {
    if (!str) return "";
    return String(str)
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}
