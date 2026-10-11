import {
  Stethoscope,
  Ambulance,
  Bed,
  Pill,
  Boxes,
  CreditCard,
  Database,
  Shield,
  ShieldCheck,
  BookOpen,
  Receipt,
} from "@lucide/svelte";

export interface ModuleDefinition {
  id: string;
  title: string;
  shortTitle: string;
  category: "clinical" | "support_and_diagnostic" | "admin_finance";
  categoryLabel: string;
  description: string;
  icon: any;
  color: {
    bg: string;
    text: string;
    border: string;
    hoverBg: string;
    badge: string;
  };
  defaultNav: string;
  status: "ready" | "in_progress";
  statusText: string;
  permission?: string;
  features: string[];
  quickStats?: string;
}

export const modules: ModuleDefinition[] = [
  // 1. Pelayanan Klinis & Asuhan Pasien
  {
    id: "clinical",
    title: "Rekam Medis & Status Lokalis",
    shortTitle: "Klinis EHR",
    category: "clinical",
    categoryLabel: "Pelayanan Pasien & Asuhan Medis",
    description:
      "Pemeriksaan status lokalis anatomi tubuh interaktif, catatan anamnesis, dan form SOAP pasien.",
    icon: Stethoscope,
    color: {
      bg: "bg-blue-50",
      text: "text-blue-700",
      border: "border-blue-200",
      hoverBg: "hover:border-blue-500 hover:shadow-blue-100",
      badge: "bg-blue-100 text-blue-800",
    },
    defaultNav: "physical",
    status: "ready",
    statusText: "Aktif & Terintegrasi",
    permission: "patient:read",
    features: ["Body Diagram Interaktif", "Resume SOAP", "Tanda Vital"],
    quickStats: "Sesuai SNOMED / ICD-10",
  },
  {
    id: "outpatient",
    title: "Pelayanan Rawat Jalan",
    shortTitle: "Rawat Jalan",
    category: "clinical",
    categoryLabel: "Pelayanan Pasien & Asuhan Medis",
    description:
      "Manajemen antrean poliklinik, registrasi kunjungan harian, dan konsultasi dokter spesialis.",
    icon: Stethoscope,
    color: {
      bg: "bg-teal-50",
      text: "text-teal-700",
      border: "border-teal-200",
      hoverBg: "hover:border-teal-500 hover:shadow-teal-100",
      badge: "bg-teal-100 text-teal-800",
    },
    defaultNav: "outpatient-queue",
    status: "in_progress",
    statusText: "Tahap Integrasi",
    permission: "patient:read",
    features: ["Antrean Poliklinik", "Panggilan Pasien", "Konsultasi DPJP"],
  },
  {
    id: "emergency",
    title: "Instalasi Gawat Darurat (IGD)",
    shortTitle: "Gawat Darurat",
    category: "clinical",
    categoryLabel: "Pelayanan Pasien & Asuhan Medis",
    description:
      "Triase kegawatdaruratan ATS/ESI, penanganan darurat pasien kritis, dan disposisi rawat inap.",
    icon: Ambulance,
    color: {
      bg: "bg-red-50",
      text: "text-red-700",
      border: "border-red-200",
      hoverBg: "hover:border-red-500 hover:shadow-red-100",
      badge: "bg-red-100 text-red-800",
    },
    defaultNav: "emergency-triage",
    status: "in_progress",
    statusText: "Tahap Integrasi",
    permission: "patient:read",
    features: ["Klasifikasi Triase", "Monitoring Kritis", "Disposisi Cepat"],
  },
  {
    id: "inpatient",
    title: "Pelayanan Rawat Inap",
    shortTitle: "Rawat Inap",
    category: "clinical",
    categoryLabel: "Pelayanan Pasien & Asuhan Medis",
    description:
      "Alokasi tempat tidur (bed management), admisi rawat inap, monitoring bangsal, dan proses discharge.",
    icon: Bed,
    color: {
      bg: "bg-cyan-50",
      text: "text-cyan-700",
      border: "border-cyan-200",
      hoverBg: "hover:border-cyan-500 hover:shadow-cyan-100",
      badge: "bg-cyan-100 text-cyan-800",
    },
    defaultNav: "inpatient-census",
    status: "in_progress",
    statusText: "Tahap Integrasi",
    permission: "patient:read",
    features: ["Alokasi Tempat Tidur", "Sensus Bangsal", "Discharge Summary"],
  },

  // 2. Penunjang Medis & Logistik
  {
    id: "pharmacy",
    title: "Farmasi & Apotek",
    shortTitle: "Farmasi",
    category: "support_and_diagnostic",
    categoryLabel: "Penunjang Medis & Rantai Pasok",
    description:
      "Peresepan elektronik (e-Prescription), telaah resep dokter, dispensing obat, dan kontrol stok depo.",
    icon: Pill,
    color: {
      bg: "bg-emerald-50",
      text: "text-emerald-700",
      border: "border-emerald-200",
      hoverBg: "hover:border-emerald-500 hover:shadow-emerald-100",
      badge: "bg-emerald-100 text-emerald-800",
    },
    defaultNav: "pharmacy-queue",
    status: "in_progress",
    statusText: "Tahap Integrasi",
    permission: "item:read",
    features: ["e-Resep Elektronik", "Telaah Interaksi", "Dispensing Racikan"],
  },
  {
    id: "inventory",
    title: "Logistik & Inventori Medis",
    shortTitle: "Inventori",
    category: "support_and_diagnostic",
    categoryLabel: "Penunjang Medis & Rantai Pasok",
    description:
      "Pencatatan mutasi stok gudang, kartu stok FEFO/FIFO, stok opname, dan pengadaan BMHP/alat kesehatan.",
    icon: Boxes,
    color: {
      bg: "bg-amber-50",
      text: "text-amber-700",
      border: "border-amber-200",
      hoverBg: "hover:border-amber-500 hover:shadow-amber-100",
      badge: "bg-amber-100 text-amber-800",
    },
    defaultNav: "inventory-stock",
    status: "in_progress",
    statusText: "Tahap Integrasi",
    permission: "item:read",
    features: ["Mutasi Stok Antar-Depo", "Kartu Stok FEFO", "Peringatan ED"],
  },

  // 3. Administrasi & Tata Kelola
  {
    id: "master",
    title: "Master Data & Organisasi",
    shortTitle: "Master RS",
    category: "admin_finance",
    categoryLabel: "Administrasi, Keuangan & Tata Kelola",
    description:
      "Katalog data pokok rumah sakit: pasien, tenaga medis/nakes, departemen, unit poli, ruangan/bed, rekanan penjamin, faskes rujukan, kelas tarif, dan katalog item universal.",
    icon: Database,
    color: {
      bg: "bg-sky-50",
      text: "text-sky-700",
      border: "border-sky-200",
      hoverBg: "hover:border-sky-500 hover:shadow-sky-100",
      badge: "bg-sky-100 text-sky-800",
    },
    defaultNav: "master-patient",
    status: "ready",
    statusText: "Aktif & Terintegrasi",
    permission: "patient:read",
    features: [
      "Demografi Pasien",
      "Katalog Dokter & Nakes",
      "Struktur Poli & Ruangan",
      "Debitur & Rekanan",
      "Faskes Rujukan",
      "Kelas Tarif",
      "Katalog Universal Item",
    ],
    quickStats: "9 Entitas Master",
  },
  {
    id: "accounting",
    title: "Chart of Account (COA)",
    shortTitle: "Akuntansi COA",
    category: "admin_finance",
    categoryLabel: "Administrasi, Keuangan & Tata Kelola",
    description:
      "Bagan akun standar rumah sakit (Aset, Kewajiban, Ekuitas, Pendapatan, Beban), konfigurasi akun kas & bank, serta hierarki pohon buku besar keuangan.",
    icon: BookOpen,
    color: {
      bg: "bg-indigo-50",
      text: "text-indigo-700",
      border: "border-indigo-200",
      hoverBg: "hover:border-indigo-500 hover:shadow-indigo-100",
      badge: "bg-indigo-100 text-indigo-800",
    },
    defaultNav: "accounting-coa",
    status: "ready",
    statusText: "COA Aktif & Terintegrasi",
    permission: "accounting:read",
    features: [
      "Bagan Akun (COA)",
      "Jurnal Umum Double-Entry",
      "Buku Kas & Bank",
      "Laporan Keuangan",
    ],
    quickStats: "Standar Akuntansi RS",
  },
  {
    id: "finance",
    title: "Keuangan",
    shortTitle: "Tarif & Finance",
    category: "admin_finance",
    categoryLabel: "Administrasi, Keuangan & Tata Kelola",
    description:
      "Pengaturan kelas tarif akomodasi (VVIP - Kelas 3), matriks komponen pendapatan jasa medis & sarana RS, serta pemetaan bagan akun (COA).",
    icon: Receipt,
    color: {
      bg: "bg-emerald-50",
      text: "text-emerald-700",
      border: "border-emerald-200",
      hoverBg: "hover:border-emerald-500 hover:shadow-emerald-100",
      badge: "bg-emerald-100 text-emerald-800",
    },
    defaultNav: "finance-tariff",
    status: "ready",
    statusText: "Aktif & Terintegrasi",
    permission: "tariff_class:read",
    features: [
      "Debitur & Penjamin (Customer)",
      "Buku Tarif (Price Plan)",
      "Kelas & Komponen Tarif",
    ],
    quickStats: "Master Keuangan RS",
  },
  {
    id: "billing",
    title: "Kasir & Billing Pasien",
    shortTitle: "Billing",
    category: "admin_finance",
    categoryLabel: "Administrasi, Keuangan & Tata Kelola",
    description:
      "Kalkulasi tarif tindakan, invoice penjaminan BPJS/Asuransi, kasir pembayaran, dan bukti setor.",
    icon: CreditCard,
    color: {
      bg: "bg-green-50",
      text: "text-green-700",
      border: "border-green-200",
      hoverBg: "hover:border-green-500 hover:shadow-green-100",
      badge: "bg-green-100 text-green-800",
    },
    defaultNav: "billing-workspace",
    status: "in_progress",
    statusText: "Tahap Integrasi",
    permission: "customer:read",
    features: [
      "Billing Tagihan Pasien",
      "Klaim BPJS & Asuransi",
      "Kasir Pembayaran",
    ],
  },
  {
    id: "auth",
    title: "Pengguna & Hak Akses (RBAC)",
    shortTitle: "Akses & User",
    category: "admin_finance",
    categoryLabel: "Administrasi, Keuangan & Tata Kelola",
    description:
      "Manajemen pengguna RS, pembuatan akun dokter/nakes, penugasan peran (role), dan konfigurasi matriks perizinan modul.",
    icon: ShieldCheck,
    color: {
      bg: "bg-violet-50",
      text: "text-violet-700",
      border: "border-violet-200",
      hoverBg: "hover:border-violet-500 hover:shadow-violet-100",
      badge: "bg-violet-100 text-violet-800",
    },
    defaultNav: "auth-users",
    status: "ready",
    statusText: "Aktif & Terintegrasi",
    permission: "user:read",
    features: ["Kelola Pengguna", "Peran & Matriks Izin", "Reset Password"],
    quickStats: "Terintegrasi Go Auth",
  },
  {
    id: "audit",
    title: "Audit Trail & Keamanan",
    shortTitle: "Audit Sistem",
    category: "admin_finance",
    categoryLabel: "Administrasi, Keuangan & Tata Kelola",
    description:
      "Append-only audit log aktivitas pengguna, riwayat mutasi data medis, dan rekam jejak kepatuhan.",
    icon: Shield,
    color: {
      bg: "bg-slate-50",
      text: "text-slate-700",
      border: "border-slate-200",
      hoverBg: "hover:border-slate-500 hover:shadow-slate-100",
      badge: "bg-slate-100 text-slate-800",
    },
    defaultNav: "audit-workspace",
    status: "in_progress",
    statusText: "Tahap Integrasi",
    permission: "patient:read",
    features: [
      "Immutable Audit Trail",
      "Pelacakan Mutasi RME",
      "Kepatuhan Regulasi",
    ],
  },
];
